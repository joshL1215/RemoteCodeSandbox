package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joshL1215/RemoteCodeSandbox/controller"
)

func main() {
	router := chi.NewRouter()
	router.Get("/", InputHandler)

	server := &http.Server{
		Addr: ":5001",
		Handler: router
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Listen and serve failed", err)
	}

}