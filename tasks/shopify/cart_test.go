package shopify

import (
	"Mystery/profiles"
	"Mystery/tasks"
	"testing"
)

// Tests adding to cart on Kith
func TestAddToCartKith(t *testing.T) {
	site := tasks.Website{Name: "Kith", Url: "https://kith.com"}
	variant := 39250265571456
	task := NewTaskShopify(site, &profiles.Profile{}, nil, tasks.ModeShopifySafe, []string{string(rune(variant))}, []string{})

	_, _, err := task.AddToCart(int64(variant))

	if err != nil {
		t.Fatal(err)
	}
}

// Tests adding to cart on Jimmy Jazz
func TestAddToCartJimmyJazz(t *testing.T) {
	site := tasks.Website{Name: "Jimmy Jazz", Url: "https://www.jimmyjazz.com"}
	variant := 42210012987599
	task := NewTaskShopify(site, &profiles.Profile{}, nil, tasks.ModeShopifySafe, []string{string(rune(variant))}, []string{})

	_, _, err := task.AddToCart(int64(variant))

	if err != nil {
		t.Fatal(err)
	}
}

// Tests clearing an item from the cart on Oneness
func TestClearCartOnenessBoutique(t *testing.T) {
	site := tasks.Website{Name: "Oneness Boutique", Url: "https://www.onenessboutique.com"}
	variant := 39880708456532
	task := NewTaskShopify(site, &profiles.Profile{}, nil, tasks.ModeShopifySafe, []string{string(rune(variant))}, []string{})

	_, _, err := task.AddToCart(int64(variant))

	if err != nil {
		t.Fatal(err)
	}

	_, _, err = task.ClearCart()

	if err != nil {
		t.Fatal(err)
	}
}
