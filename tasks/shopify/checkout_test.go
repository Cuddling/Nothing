package shopify

import (
	"Mystery/profiles"
	"Mystery/tasks"
	"testing"
)

// Tests Creating a checkout session with a valid cart on kith
func TestCreateCheckoutKithWithCart(t *testing.T) {
	site := tasks.Website{Name: "Kith", Url: "https://kith.com"}
	variant := 39250265571456
	task := NewTaskShopify(site, &profiles.Profile{}, nil, tasks.ModeShopifySafe, []string{string(rune(variant))}, []string{})

	_, _, err := task.AddToCart(int64(variant))

	if err != nil {
		t.Fatal(err)
	}

	_, _, err = task.CreateCheckout()

	if err != nil {
		t.Fatal(err)
	}
}

// Tests Creating a checkout session with a valid cart on kith
func TestCreateCheckoutKithNoCart(t *testing.T) {
	site := tasks.Website{Name: "Kith", Url: "https://kith.com"}
	task := NewTaskShopify(site, &profiles.Profile{}, nil, tasks.ModeShopifySafe, []string{}, []string{})

	_, _, err := task.CreateCheckout()

	if err != nil {
		t.Fatal(err)
	}
}

// Tests Creating and retrieving a checkout page on Kith
func TestGetCheckoutPageKith(t *testing.T) {
	site := tasks.Website{Name: "Kith", Url: "https://kith.com"}
	variant := 39250265571456
	task := NewTaskShopify(site, &profiles.Profile{}, nil, tasks.ModeShopifySafe, []string{string(rune(variant))}, []string{})

	_, _, err := task.AddToCart(int64(variant))

	if err != nil {
		t.Fatal(err)
	}

	_, _, err = task.CreateCheckout()

	if err != nil {
		t.Fatal(err)
	}

	_, _, err = task.GetCheckoutPage(task.CheckoutUrl)

	if err != nil {
		t.Fatal(err)
	}
}
