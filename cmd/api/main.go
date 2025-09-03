package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joshL1215/RemoteCodeSandbox/internal/docker"
	"github.com/joshL1215/RemoteCodeSandbox/internal/handler"
)

func main() {
	cli := docker.ConnectToDaemon()

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", handler.InputHandler(cli))

	server := &http.Server{
		Addr:    ":5001",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Listen and serve failed", err)
	}
}
