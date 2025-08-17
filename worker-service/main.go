package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/erikkvale/nornir/proto/workflowpb"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type server struct {
	workflowpb.UnimplementedWorkflowServiceServer
}

// StartWorkflow implements the gRPC method
func (s *server) StartWorkflow(ctx context.Context, req *workflowpb.StartWorkflowRequest) (*workflowpb.StartWorkflowResponse, error) {
	log.Printf("Received request to start workflow: name=%s", req.Name)

	workflowId := uuid.New().String()

	return &workflowpb.StartWorkflowResponse{
		Id:     workflowId,
		Status: "STARTED",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v")
	}

	grpcServer := grpc.NewServer()
	workflowpb.RegisterWorkflowServiceServer(grpcServer, &server{})

	fmt.Println("Worker service listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
