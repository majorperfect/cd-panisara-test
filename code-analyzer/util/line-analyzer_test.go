package util_test

import (
	"testing"

	"github.com/majorperfect/guardrails-test/code-analyzer/service/models"
	"github.com/majorperfect/guardrails-test/code-analyzer/util"
	"github.com/stretchr/testify/assert"
)

func TestNewLineAnalyzer(t *testing.T) {
	var lineNo uint32 = 1
	var path string = "repo/dirname/me.go"

	testCases := []struct {
		Test   string
		Line   string
		Path   string
		LineNo uint32
		Output *models.Finding
		Detect bool
	}{
		{
			Test:   "1",
			Line:   "const private_key",
			Path:   path,
			LineNo: lineNo,
			Output: &models.Finding{
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
			},
			Detect: true,
		},
		{
			Test:   "2",
			Line:   "constfdggdfgdfgdfgd public_keyewretet end//",
			Path:   path,
			LineNo: lineNo,
			Output: &models.Finding{
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
			},
			Detect: true,
		},
		{
			Test:   "3",
			Line:   "conprivate_keydggdfgdfgdfgdpublic_keyewretet end//",
			Path:   path,
			LineNo: lineNo,
			Output: &models.Finding{
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
			},
			Detect: true,
		},
		{
			Test:   "4",
			Line:   "func Login() error{",
			Path:   path,
			LineNo: lineNo,
			Output: &models.Finding{},
			Detect: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.Test, func(t *testing.T) {
			testFunc := util.NewLineAnalyzer()

			findings, detect := testFunc(test.LineNo, test.Line, test.Path)
			assert.Equal(t, findings, test.Output)
			assert.Equal(t, detect, test.Detect)
		})

	}

}
