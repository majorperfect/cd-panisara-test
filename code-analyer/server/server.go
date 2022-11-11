package server

import (
	"context"
	"sync"

	pb "github.com/majorperfect/guardrails-test/code-analyzer/proto/code-analyzer"
)

type Backend struct {
	mu *sync.RWMutex
	pb.UnimplementedCodeAnalyzerServiceServer
}

func New() pb.CodeAnalyzerServiceServer {
	return &Backend{
		mu: &sync.RWMutex{},
	}
}

func (b *Backend) HealthCheck(ctx context.Context, p *pb.Empty) (*pb.HealthCheckResponse, error) {

	b.mu.Lock()
	defer b.mu.Unlock()
	return &pb.HealthCheckResponse{}, nil

}
