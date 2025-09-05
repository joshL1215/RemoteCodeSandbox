package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/joshL1215/RemoteCodeSandbox/internal/docker"
	"github.com/joshL1215/RemoteCodeSandbox/internal/models"
)

type RequestPayload struct {
	Language string        `json:"language"`
	Code     string        `json:"code"`
	Cases    []models.Case `json:"cases"`
}

func InputHandler(cli *client.Client) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var payload RequestPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Payload is invalid json: "+err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Printf("Received payload: %+v\n", payload)

		result, err := docker.RunJudgeJob(cli, payload.Language, payload.Code, payload.Cases)
		if err != nil {
			http.Error(w, "Issue with server: "+err.Error(), http.StatusInternalServerError)
		}

		fmt.Println(result)
	}
}
