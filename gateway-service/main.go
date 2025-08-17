package main

import (
	"context"
	"encoding/json"
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

	// 5. Add handler for POST /workflows
	mux.HandleFunc("/workflows", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		defer r.Body.Close()

		var req startRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
			http.Error(w, "invalid payload: need name", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		grpcResponse, err := client.StartWorkflow(ctx, &pb.StartWorkflowRequest{Name: req.Name})
		if err != nil {
			log.Printf("gateway: StarWorkflow error: %v", err)
			http.Error(w, "upstream error", http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(startResponse{
			ID:     grpcResponse.GetId(),
			Status: grpcResponse.GetStatus(),
		})
	})

	addr := os.Getenv("HTTP_ADDRESS")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("gateway: listening on %s (worker=%s)", addr, workerServiceAddress)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("gateway: server error %v", err)
	}

}
