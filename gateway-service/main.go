package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/erikkvale/nornir/proto/workflowpb"
)

type startRequest struct {
	Name string `json:"name"`
}

type startResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type WorkflowHandler struct {
	client pb.WorkflowServiceClient
}

func NewWorkflowHandler(client pb.WorkflowServiceClient) *WorkflowHandler {
	return &WorkflowHandler{client: client}
}

func (h *WorkflowHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handlePost(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *WorkflowHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	req, err := h.decodeRequest(r)
	if err != nil {
		http.Error(w, "invalid payload: need name", http.StatusBadRequest)
		return
	}

	resp, err := h.startWorkflow(r.Context(), req)
	if err != nil {
		log.Printf("gateway: StartWorkflow error: %v", err)
		http.Error(w, "upstream error", http.StatusBadGateway)
		return
	}

	h.writeResponse(w, resp)
}

func (h *WorkflowHandler) decodeRequest(r *http.Request) (*startRequest, error) {
	defer r.Body.Close()
	var req startRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	return &req, nil
}

func (h *WorkflowHandler) startWorkflow(ctx context.Context, req *startRequest) (*startResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	grpcResp, err := h.client.StartWorkflow(ctx, &pb.StartWorkflowRequest{Name: req.Name})
	if err != nil {
		return nil, err
	}

	return &startResponse{
		ID:     grpcResp.GetId(),
		Status: grpcResp.GetStatus(),
	}, nil
}

func (h *WorkflowHandler) writeResponse(w http.ResponseWriter, resp *startResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	// 1. Read environment variable for the worker service address (default to localhost:50051)
	workerServiceAddress := os.Getenv("WORKER_ADDRESS")
	if workerServiceAddress == "" {
		workerServiceAddress = "localhost:50051"
	}

	// 2. Create a gRPC client connection to the worker service
	conn, err := grpc.Dial(workerServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("gateway: failed to connect to worker: %v", err)
	}
	defer conn.Close()

	// 3. Initialize a WorkflowServiceClient from the generated gRPC code
	client := pb.NewWorkflowServiceClient(conn)

	// 4. Set up an HTTP mux/handler
	mux := http.NewServeMux()
	handler := NewWorkflowHandler(client)
	mux.Handle("/workflows", handler)

	addr := os.Getenv("HTTP_ADDRESS")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("gateway: listening on %s (worker=%s)", addr, workerServiceAddress)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("gateway: server error %v", err)
	}

}
