package shopify

import (
	"Mystery/profiles"
	"Mystery/tasks"
	"fmt"
	"testing"
)

// Tests fetching the products.json on Kith
func TestGetProductsKith(t *testing.T) {
	site := tasks.Website{Name: "Kith", Url: "https://kith.com"}
	task := NewTaskShopify(site, &profiles.Profile{}, nil, tasks.ModeShopifySafe, []string{}, []string{})

	resp, products, err := task.GetProducts()

	if err != nil {
		t.Fatal(err)
	}

	if resp.IsError() {
		t.Fatal(fmt.Sprintf("request failed with error %v", resp.StatusCode()))
	}

	if len(products) != 30 {
		t.Fatal(fmt.Sprintf("expected 30 products. got %v", len(products)))
	}
}

// Tests fetching the products.json on DSMNY E-flash. Expecting a 401 because there's a password page always enabled.
func TestGetProductsPasswordPage(t *testing.T) {
	site := tasks.Website{Name: "DSMNY E-Flash", Url: "https://eflash-us.doverstreetmarket.com"}
	task := NewTaskShopify(site, &profiles.Profile{}, nil, tasks.ModeShopifySafe, []string{}, []string{})

	resp, products, err := task.GetProducts()

	if err != nil {
		t.Fatal(err)
	}

	if !resp.IsError() || resp.StatusCode() != 401 {
		t.Fatal(fmt.Sprintf("expected 401 status code. Got %v", resp.StatusCode()))
	}

	if len(products) != 0 {
		t.Fatal(fmt.Sprintf("expected 0 products. got %v", len(products)))
	}
}

// Tests fetching a specific product on Kith.
func TestGetProductSpecificKith(t *testing.T) {
	site := tasks.Website{Name: "Kith", Url: "https://kith.com"}
	task := NewTaskShopify(site, &profiles.Profile{}, nil, tasks.ModeShopifySafe, []string{}, []string{})

	resp, product, err := task.GetProductSpecific("https://kith.com/products/nkdd1399-300")

	if err != nil {
		t.Fatal(err)
	}

	if resp.IsError() {
		t.Fatal(fmt.Sprintf("request failed with error %v", resp.StatusCode()))
	}

	if product.Title != "Nike Dunk High Retro BTTYS - Noble Green / White" {
		t.Fatal(fmt.Sprintf("retrieved incorrect product title - '%v'", product.Title))
	}
}

// Tests fetching a specific product on DSMNY. Expecting an error due to there always being a password page enabled.
func TestGetProductSpecificDSMNY(t *testing.T) {
	site := tasks.Website{Name: "DSMNY", Url: "https://eflash-us.doverstreetmarket.com"}
	task := NewTaskShopify(site, &profiles.Profile{}, nil, tasks.ModeShopifySafe, []string{}, []string{})

	resp, _, err := task.GetProductSpecific("https://eflash-us.doverstreetmarket.com/products/testproduct123")

	if err != nil {
		t.Fatal(err)
	}

	if !resp.IsError() {
		t.Fatal("expected request to fail")
	}
}
