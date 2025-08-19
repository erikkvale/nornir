package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	pb "github.com/erikkvale/nornir/proto/workflowpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockWorkflowServiceClient is a mock implementation of WorkflowServiceClient
type MockWorkflowServiceClient struct {
	mock.Mock
}

func (m *MockWorkflowServiceClient) StartWorkflow(ctx context.Context, in *pb.StartWorkflowRequest, opts ...grpc.CallOption) (*pb.StartWorkflowResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.StartWorkflowResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

// Add the missing GetStatus method
func (m *MockWorkflowServiceClient) GetStatus(ctx context.Context, in *pb.GetStatusRequest, opts ...grpc.CallOption) (*pb.GetStatusResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.GetStatusResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestWorkflowEndpoint(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		requestBody    interface{}
		setupMock      func(*MockWorkflowServiceClient)
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:   "successful workflow start",
			method: http.MethodPost,
			requestBody: startRequest{
				Name: "test-workflow",
			},
			setupMock: func(m *MockWorkflowServiceClient) {
				m.On("StartWorkflow", mock.Anything, &pb.StartWorkflowRequest{
					Name: "test-workflow",
				}, mock.Anything).Return(&pb.StartWorkflowResponse{
					Id:     "123",
					Status: "running",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]string{
				"id":     "123",
				"status": "running",
			},
		},
		{
			name:           "invalid method",
			method:         http.MethodGet,
			requestBody:    nil,
			setupMock:      func(m *MockWorkflowServiceClient) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "invalid payload",
			method:         http.MethodPost,
			requestBody:    struct{}{},
			setupMock:      func(m *MockWorkflowServiceClient) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client and set up expectations
			mockClient := new(MockWorkflowServiceClient)
			tt.setupMock(mockClient)

			// Create request
			var body bytes.Buffer
			if tt.requestBody != nil {
				json.NewEncoder(&body).Encode(tt.requestBody)
			}

			req := httptest.NewRequest(tt.method, "/workflows", &body)
			rec := httptest.NewRecorder()

			// Use the WorkflowHandler directly
			handler := NewWorkflowHandler(mockClient)
			handler.ServeHTTP(rec, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedBody != nil {
				var response map[string]string
				err := json.NewDecoder(rec.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestEnvironmentVariables(t *testing.T) {
	// Test WORKER_ADDRESS
	os.Setenv("WORKER_ADDRESS", "test:50051")
	defer os.Unsetenv("WORKER_ADDRESS")
	addr := os.Getenv("WORKER_ADDRESS")
	assert.Equal(t, "test:50051", addr)

	// Test HTTP_ADDRESS
	os.Setenv("HTTP_ADDRESS", ":9090")
	defer os.Unsetenv("HTTP_ADDRESS")
	httpAddr := os.Getenv("HTTP_ADDRESS")
	assert.Equal(t, ":9090", httpAddr)
}
