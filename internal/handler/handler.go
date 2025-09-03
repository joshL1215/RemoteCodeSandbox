package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/docker/docker/client"
)

type Case struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expectedOutput"`
}

type RequestPayload struct {
	Language    string `json:"language"`
	Code        string `json:"code"`
	PrelimCases []Case `json:"prelimCases"`
	TestCases   []Case `json:"testCases"`
}

func InputHandler(cli *client.Client) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var payload RequestPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Payload is invalid json: "+err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("Received payload: %+v\n", payload)
	}
}
