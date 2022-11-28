package util

import (
	"bufio"
	"os"

	"github.com/majorperfect/guardrails-test/code-analyzer/service/models"
)

// open file from path and scan line by line
type ScanLineFunc = func(path string, projectPath string) ([]*models.Finding, error)

func NewScanCode(lineAnalyzer LineAnalyzerFunc) ScanLineFunc {
	return func(path string, projectPath string) ([]*models.Finding, error) {
		findings := []*models.Finding{}

		file, err := os.Open(path)
		if err != nil {
			return findings, err
		}
		fileScanner := bufio.NewScanner(file)
		fileScanner.Split(bufio.ScanLines)

		var count uint32 = 0
		for fileScanner.Scan() {
			//analyzing code line by line
			count++
			f, ok := lineAnalyzer(count, fileScanner.Text(), projectPath)
			if ok {
				findings = append(findings, f)
			}
		}

		return findings, nil
	}
}
