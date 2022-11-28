package service

import (
	"io"

	pb "github.com/majorperfect/guardrails-test/code-analyzer/proto/code-analyzer"
	"github.com/majorperfect/guardrails-test/code-analyzer/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// read from buffer and scan line by line
type HandleUploadStreamFileFunc = func(stream pb.CodeAnalyzerService_AnalyzeUploaderServer) error

func NewHandleUploadStreamFile(scanBuffer util.BufferScannerFunc) HandleUploadStreamFileFunc {
	return func(stream pb.CodeAnalyzerService_AnalyzeUploaderServer) error {
		firstStream := true
		var newFile *util.File
		for {
			req, err := stream.Recv()
			if firstStream {
				newFile = util.NewFile(req.Name)
				// create report file queueing time
				firstStream = false
			}
			if err == io.EOF {
				// improve by saving file to storage
				// can do send a worker to analyze a file read from kafka or else
				// for a demo we save in local dir and read it again

				go scanBuffer(newFile)
				return stream.SendAndClose(&pb.UploadResponse{Name: newFile.GetName()})
			} else {
				newFile.Write(req.Chunk)
			}
			if err != nil {
				return status.Error(codes.Internal, err.Error())
			}

		}
	}
}
