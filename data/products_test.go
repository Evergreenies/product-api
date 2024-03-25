package data

import "testing"

func TestStructValidation(t *testing.T) {
	product := &Product{
		Name:  "nics",
		Price: 1,
		SKU:   "a-b-s",
	}
	err := product.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
