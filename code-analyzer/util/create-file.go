package util

import (
	"encoding/base64"
	"errors"
	"os"
	"strings"

	pb "github.com/majorperfect/guardrails-test/code-analyzer/proto/code-analyzer"
)

// create file to local disk and return path to file
type CreateFileFunc = func(req *pb.FileUploadRequest) (string, error)

func NewCreateFile(root string) CreateFileFunc {
	return func(req *pb.FileUploadRequest) (string, error) {
		dec, err := base64.StdEncoding.DecodeString(req.Base64)
		if err != nil {
			return "", err
		}

		p := strings.Split(req.Name, "/")
		path := root + "/" + strings.Join(p[1:], "/")
		dir := root + "/" + strings.Join(p[1:len(p)-1], "/")
		if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(dir, os.ModePerm)
			if err != nil {
				return "", err
			}
		}
		f, err := os.Create(path)
		if err != nil {
			return "", err

		}
		defer f.Close()

		if _, err := f.Write(dec); err != nil {
			return "", err

		}
		if err := f.Sync(); err != nil {
			return "", err
		}

		return path, nil
	}
}
