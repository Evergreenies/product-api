package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/evergreenies/product-api/data"
	"github.com/gorilla/mux"
)

type Products struct {
	logger *log.Logger
}

func NewProducts(logger *log.Logger) *Products {
	return &Products{logger}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("GET /products or `/` - to get all products")
	allProducts := data.GetProducts()
	err := allProducts.ToJSON(rw)
	if err != nil {
		http.Error(rw, "unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("POST /products - to create products")
	product := r.Context().Value(KeyProduct{}).(*data.Product)
	p.logger.Printf("Product: %#v\n", product)
	data.AddProduct(product)
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(rw, "expected `id` in path", http.StatusBadRequest)
	}

	p.logger.Println("PUT /products - to update products")
	product := r.Context().Value(KeyProduct{}).(*data.Product)

	err = data.UpdateProduct(id, product)
	if err == data.ErrProductNotFoun {
		http.Error(rw, "product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "product not found", http.StatusInternalServerError)
	}
}

type KeyProduct struct{}

func (p *Products) MiddlewareProductsValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		product := &data.Product{}

		err := product.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "unable to unmarshal json", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}
