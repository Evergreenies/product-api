package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/evergreenies/product-api/handlers"
	"github.com/gorilla/mux"
)

func main() {
	logger := log.New(os.Stdout, "product-api | ", log.LstdFlags)

	productsHandler := handlers.NewProducts(logger)

	// serveMux := http.NewServeMux()
	serveMux := mux.NewRouter()

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", productsHandler.GetProducts)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.Use(productsHandler.MiddlewareProductsValidation)
	putRouter.HandleFunc("/products/{id:[0-9]+}", productsHandler.UpdateProduct)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.Use(productsHandler.MiddlewareProductsValidation)
	postRouter.HandleFunc("/products", productsHandler.AddProduct)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		logger.Printf("Starting server at %d.", 8080)
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
