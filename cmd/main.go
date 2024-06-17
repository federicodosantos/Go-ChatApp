package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	// "github.com/federicodosantos/Go-ChatApp/pkg/db/postgres"
	"github.com/federicodosantos/Go-ChatApp/pkg/log"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout",
	 time.Second * 15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// load environment variable
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("failed to load env: %s", err)
		os.Exit(1)
	}
	
	// init zap logger
	logger := log.NewLogger(os.Getenv("APP_ENV"))
	defer logger.Sync()

	// init database
	// db := postgres.DBInit(logger)
	
	// init mux
	mux := mux.NewRouter()

	server := &http.Server{
		Handler: mux,
		Addr: fmt.Sprintf("%s:%s", os.Getenv("URL"), os.Getenv("PORT")),
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	// Start the server
	logger.Info("Starting server",
		zap.String("address", server.Addr),
	)

	go func ()  {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start the server", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// block until we receive our signal
	<-c

	// create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("shutting down")
	os.Exit(0)
}