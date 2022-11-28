package server

import (
	"context"
	"sync"

	"github.com/majorperfect/guardrails-test/code-analyzer/service"
	"github.com/majorperfect/guardrails-test/code-analyzer/util"

	pb "github.com/majorperfect/guardrails-test/code-analyzer/proto/code-analyzer"
)

type Backend struct {
	mu *sync.RWMutex
	pb.UnimplementedCodeAnalyzerServiceServer
	handleUploadFile       service.HandleUploadFileFunc
	handleUploadStreamFile service.HandleUploadStreamFileFunc
}

func New() pb.CodeAnalyzerServiceServer {
	return &Backend{
		mu:                     &sync.RWMutex{},
		handleUploadFile:       service.NewHandleUploadFile(util.NewCreateFile("./source"), util.NewScanCode(util.NewLineAnalyzer())),
		handleUploadStreamFile: service.NewHandleUploadStreamFile(util.NewBufferScanner(util.NewLineAnalyzer())),
	}
}

func (b *Backend) HealthCheck(ctx context.Context, p *pb.Empty) (*pb.HealthCheckResponse, error) {

	b.mu.Lock()
	defer b.mu.Unlock()
	return &pb.HealthCheckResponse{}, nil

}

func (b *Backend) AnalyzeUploader(stream pb.CodeAnalyzerService_AnalyzeUploaderServer) error {

	b.mu.Lock()
	defer b.mu.Unlock()
	return b.handleUploadStreamFile(stream)

}

func (b *Backend) FileUpload(ctx context.Context, req *pb.FileUploadRequest) (*pb.FileUploadResponse, error) {

	b.mu.Lock()
	defer b.mu.Unlock()

	if err := b.handleUploadFile(req); err != nil {
		return &pb.FileUploadResponse{}, err
	}
	return &pb.FileUploadResponse{
		Message: "",
		Status:  "Success",
		Code:    "200",
		Data: &pb.FileUploadResponse_Data{
			Name: req.Name,
		},
	}, nil

}
