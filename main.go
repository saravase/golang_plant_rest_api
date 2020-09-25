package main

import (
	"context"
	"golang_microservice/plant-api/handlers"
	"os/signal"

	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	// New creates a new plant-api Logger.
	logger := log.New(os.Stdout, "product-plant-api", log.LstdFlags)

	// Initialize the plant struct properties
	plantHandler := handlers.NewPlant(logger)

	// NewServeMux allocates and returns a plant-api ServeMux.
	serveMux := http.NewServeMux()

	// Handle registers the handler for the given pattern
	serveMux.Handle("/plant/", plantHandler)

	// Initialize the plant-api server properties
	server := http.Server{
		Addr:         ":9090",
		Handler:      serveMux,
		IdleTimeout:  100 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// Initialize the go-routine function
	go func() {

		// ListenAndServe listens on the TCP network address specified in the server property
		listenAndServeError := server.ListenAndServe()

		if listenAndServeError != nil {
			logger.Fatal(listenAndServeError)
		}
	}()

	// Make the channel with type os.Signal
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	// Read the channel value
	sig := <-signalChannel

	logger.Println("Received os signal, graceful timeout", sig)

	//Canceling this context releases resources associated with it
	terminateContext, terminateContextError := context.WithTimeout(context.Background(), 30*time.Second)

	if terminateContextError != nil {
		logger.Fatal(terminateContextError)
	}

	// Shutdown gracefully shuts down the server without interrupting any active connections
	server.Shutdown(terminateContext)

}
