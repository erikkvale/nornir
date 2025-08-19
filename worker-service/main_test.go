package main

import (
	"context"
	"net"
	"testing"

	pb "github.com/erikkvale/nornir/proto/workflowpb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestStartWorkflow(t *testing.T) {
	// Create a test server
	s := &server{}

	// Test cases
	tests := []struct {
		name    string
		req     *pb.StartWorkflowRequest
		wantErr bool
	}{
		{
			name: "successful workflow start",
			req: &pb.StartWorkflowRequest{
				Name: "test-workflow",
			},
			wantErr: false,
		},
		{
			name: "empty workflow name",
			req: &pb.StartWorkflowRequest{
				Name: "",
			},
			wantErr: false, // Currently no validation, could be changed to true if validation is added
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := s.StartWorkflow(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, resp.Id)
			assert.Equal(t, "STARTED", resp.Status)
		})
	}
}

func TestGRPCServer(t *testing.T) {
	// Start a test server on a random port
	lis, err := net.Listen("tcp", ":0")
	assert.NoError(t, err)

	s := grpc.NewServer()
	pb.RegisterWorkflowServiceServer(s, &server{})

	// Start server in a goroutine
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Errorf("failed to serve: %v", err)
		}
	}()

	// Create a client connection
	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewWorkflowServiceClient(conn)

	// Test the connection
	resp, err := client.StartWorkflow(context.Background(), &pb.StartWorkflowRequest{
		Name: "integration-test",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Id)
	assert.Equal(t, "STARTED", resp.Status)

	// Clean up
	s.GracefulStop()
}
