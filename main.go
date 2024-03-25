package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/evergreenies/product-api/handlers"
)

func main() {
	logger := log.New(os.Stdout, "product-api", log.LstdFlags)

	helloHandler := handlers.NewHello(logger)
	byHandler := handlers.NewGoodBye(logger)

	serveMux := http.NewServeMux()

	serveMux.Handle("/", helloHandler)
	serveMux.Handle("/bye", byHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	logger.Println("Terminating..., \ngraceful shutdown", sig)

	contxt, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(contxt)
}
