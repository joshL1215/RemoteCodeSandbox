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

	cli, err := docker.ConnectToDaemon()
	if err != nil {
		fmt.Println(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Post("/", handler.InputHandler(cli))

	server := &http.Server{
		Addr:    ":5001",
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Listn and serve failed ", err)
	}
}
