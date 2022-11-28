package upload

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/majorperfect/guardrails-test/code-scanner/internal"
)

func NewPostUploadHandler(fileScanner internal.ScannerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fileScanner(r.Context(), "./files"); err != nil {
			log.Printf("fileScanner error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(map[string]string{"error": "fileScanner"}); err != nil {
				log.Printf("json encode error: %s", err)
			}
			return
		}
	}
}
