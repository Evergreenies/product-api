package data

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

type Products []*Product

var ErrProductNotFoun = errors.New("product not found")

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error {
	_, index, err := findProduct(id)
	if err != nil {
		return err
	}

	p.ID = id
	productList[index] = p

	return nil
}

func getNextID() int {
	if len(productList) == 0 {
		return 1
	}
	lastProductID := productList[len(productList)-1].ID
	lastProductID++
	return lastProductID
}

func findProduct(id int) (*Product, int, error) {
	for index, product := range productList {
		if product.ID == id {
			return product, index, nil
		}
	}

	return nil, -1, ErrProductNotFoun
}

var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frooty milky coffee",
		Price:       2.45,
		SKU:         "abc234",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	}, {
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "somting 123",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
