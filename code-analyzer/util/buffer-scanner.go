package util

import (
	"bufio"
	"encoding/json"
	"fmt"

	"github.com/majorperfect/guardrails-test/code-analyzer/service/models"
)

type BufferScannerFunc = func(f *File) error

func NewBufferScanner(analyzer LineAnalyzerFunc) BufferScannerFunc {
	return func(f *File) error {
		scanner := bufio.NewScanner(f.buffer)
		var lineNum uint32 = 0
		var findings []*models.Finding
		for scanner.Scan() {
			lineNum++
			finding, found := analyzer(lineNum, scanner.Text(), f.name)
			if found {
				findings = append(findings[:], finding)
			}
		}

		b, err := json.Marshal(findings)
		if err != nil {
			return err

		}
		fmt.Println(string(b))
		if err := scanner.Err(); err != nil {
			return err
		}
		return nil

	}

}
