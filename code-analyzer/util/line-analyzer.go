package util

import (
	"strings"

	"github.com/majorperfect/guardrails-test/code-analyzer/service/models"
)

type LineAnalyzerFunc = func(lineNo uint32, line string, path string) (*models.Finding, bool)

func NewLineAnalyzer() LineAnalyzerFunc {
	return func(lineNo uint32, line string, path string) (*models.Finding, bool) {
		if strings.Contains(line, "public_key") || strings.Contains(line, "private_key") {
			// the logic could be improve to precisely detect the variable or just a comment or else
			return &models.Finding{
				Type:   models.SAST,
				RuleId: "G001",
				Location: models.Location{
					Path: path,
					Position: models.Position{
						Begin: models.Begin{
							Line: lineNo,
						},
					},
				},
				Metadata: models.Metadata{
					Description: "secret key was found, not secure",
					Severity:    "HIGH",
				},
			}, true
		}
		return &models.Finding{}, false

	}

}
