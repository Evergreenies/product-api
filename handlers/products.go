package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/evergreenies/product-api/data"
)

type Products struct {
	logger *log.Logger
}

func NewProducts(logger *log.Logger) *Products {
	return &Products{logger}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		p.logger.Println("PUT")
		regx := regexp.MustCompile(`/([0-9]+)`)
		grp := regx.FindAllStringSubmatch(r.URL.Path, -1)
		if len(grp) != 1 {
			http.Error(rw, "invalid url", http.StatusBadRequest)
			return
		}

		if len(grp[0]) != 2 {
			http.Error(rw, "invalid url, more than one group", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(grp[0][1])
		if err != nil {
			p.logger.Fatal("unable to convert to int")
		}

		p.logger.Println("got id: ", id)

		p.updateProduct(id, rw, r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("GET /products or `/` - to get all products")
	allProducts := data.GetProducts()
	err := allProducts.ToJSON(rw)
	if err != nil {
		http.Error(rw, "unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("POST /products - to create products")
	product := &data.Product{}
	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "unable to unmarshal json", http.StatusBadRequest)
	}

	p.logger.Printf("Product: %#v\n", product)
	data.AddProduct(product)
}

func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("PUT /products - to update products")
	product := &data.Product{}
	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "unable to unmarshal json", http.StatusBadRequest)
	}
	p.logger.Printf("Product: %#v\n", product)

	err = data.UpdateProduct(id, product)
	if err == data.ErrProductNotFoun {
		http.Error(rw, "product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "product not found", http.StatusInternalServerError)
	}
}
