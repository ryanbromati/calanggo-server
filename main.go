package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	httpAdapter "calanggo-server/internal/adapters/http"
	"calanggo-server/internal/adapters/repository"
	"calanggo-server/internal/core/services"
)

func main() {
	// 1. Adapters (Driven)
	// repo := repository.NewMemoryRepository()

	repo, err := repository.NewSQLiteRepository("calanggo.db")
	if err != nil {
		log.Fatalf("Falha ao inicializar banco de dados: %v", err)
	}

	// 2. Core (Service)
	service := services.NewLinkService(repo)
	handler := httpAdapter.NewLinkHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/shorten", handler.CreateShortLink)
	mux.HandleFunc("GET /{code}", handler.Redirect)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serverAddr := ":" + port
	fmt.Printf("Calanggo Server iniciado na porta %s\n", port)

	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatal(err)
	}
}
