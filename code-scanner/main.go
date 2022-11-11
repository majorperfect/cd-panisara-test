package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/majorperfect/guardrails-test/code-scanner/config"
	"github.com/majorperfect/guardrails-test/code-scanner/internal"
)

func main() {
	cfg := config.SetupConfig()
	r := mux.NewRouter()
	r.HandleFunc("/", internal.Scanner()).Methods(http.MethodPost)

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	idleConnectionClosed := make(chan struct{}, 1)
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := httpServer.Shutdown(context.Background()); err != nil {
			// Error starting or closing listener:
			log.Printf("HTTP server shutdown failed: %v", err)
		}
		close(idleConnectionClosed)
	}()

	fmt.Printf("HTTP Server listen at port: %s", cfg.Port)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("Cannot start server: %s", err.Error())
		close(idleConnectionClosed)
	}
	<-idleConnectionClosed
}
