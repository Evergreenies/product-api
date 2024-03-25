package data

import (
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
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

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	regx := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := regx.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}

	return true
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
