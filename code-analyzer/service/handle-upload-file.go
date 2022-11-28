package service

import (
	"encoding/json"
	"fmt"

	pb "github.com/majorperfect/guardrails-test/code-analyzer/proto/code-analyzer"
	"github.com/majorperfect/guardrails-test/code-analyzer/util"
)

// open file from path and scan line by line
type HandleUploadFileFunc = func(req *pb.FileUploadRequest) error

func NewHandleUploadFile(createFile util.CreateFileFunc, scanLine util.ScanLineFunc) HandleUploadFileFunc {
	return func(req *pb.FileUploadRequest) error {
		path, err := createFile(req)
		if err != nil {
			return err
		}

		go func() {
			report, err := scanLine(path, req.Name)
			if err != nil {
				fmt.Println(err)
			}
			// save report
			for _, r := range report {
				b, err := json.Marshal(r)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println(string(b))
			}
		}()

		return err
	}
}
