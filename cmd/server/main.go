package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"arca3/config"
	"arca3/handlers"
	"arca3/spreadsheet"
)

const (
	shutdownTimeout = 30 * time.Second
)

func main() {
	ctxMain := context.Background()
	ctxSignal, cancel := signal.NotifyContext(ctxMain, syscall.SIGINT, syscall.SIGTERM)

	defer cancel()

	env := config.LoadConfig()

	spreadsheet := spreadsheet.New(ctxSignal, env.ServiceCredentialsPath, env.SpreadsheetID)

	server := launchServer(env, spreadsheet)

	<-ctxSignal.Done()

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), shutdownTimeout)

	defer cancelShutdown()

	log.Println("Received shutdown signal, shutting down server...")
	server.Shutdown(ctxShutdown)
	log.Println("Server gracefully stopped")
}

func launchServer(env *config.Config, spreadsheet *spreadsheet.Spreadsheet) *http.Server {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	wallsHandlers := handlers.NewWallsHandler(spreadsheet)
	router.Get("/api/v1/read_all", wallsHandlers.ReadAll)

	server := &http.Server{
		Addr:    env.ServerAddress,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Error listening and serving: %v\n", err)

			return
		}
	}()

	return server
}
