package shopify

import (
	"Mystery/tasks"
	"fmt"
	"testing"
	"time"
)

func TestGetShippingRatesOneness(t *testing.T) {
	site := tasks.Website{Name: "Oneness Boutique", Url: "https://www.onenessboutique.com"}
	variant := 39885244498004
	task := NewTaskShopify(site, &testProfile, nil, tasks.ModeShopifySafe, []string{string(rune(variant))}, []string{"L"})

	_, _, _ = task.AddToCart(int64(variant))
	_, _, _ = task.CreateCheckout()
	_, _, _ = task.GetCheckoutPage(task.CheckoutUrl)
	_, _, _ = task.SubmitContactInfo()

	rate := ""

	for rate == "" {
		var err error
		rate, _, err = task.FetchShippingRate()

		switch err {
		case ErrNoShippingRateAvailable:
			time.Sleep(1 * time.Second)
			continue
		case nil:
			break
		default:
			t.Fatal(err)
		}
	}

	t.Log(fmt.Sprintf("Shipping Rate: %v", rate))
}

func TestSubmitShippingRateOneness(t *testing.T) {
	site := tasks.Website{Name: "Oneness Boutique", Url: "https://www.onenessboutique.com"}
	variant := 39885244498004
	task := NewTaskShopify(site, &testProfile, nil, tasks.ModeShopifySafe, []string{string(rune(variant))}, []string{"L"})

	_, _, _ = task.AddToCart(int64(variant))
	_, _, _ = task.CreateCheckout()
	_, _, _ = task.GetCheckoutPage(task.CheckoutUrl)
	_, _, _ = task.SubmitContactInfo()

	rate := ""

	for rate == "" {
		var err error
		rate, _, err = task.FetchShippingRate()

		switch err {
		case ErrNoShippingRateAvailable:
			time.Sleep(1 * time.Second)
			continue
		case nil:
			break
		default:
			t.Fatal(err)
		}
	}

	t.Log(fmt.Sprintf("Shipping Rate: %v", rate))

	page, _, err := task.SubmitShippingRate(rate)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(page.Url)
}

func TestGetShippingRatePrice(t *testing.T) {
	if GetShippingRatePrice("Advanced Shipping Rules-1-12.05") != "12.05" {
		t.Fatal("incorrect shipping rate price")
	}
}
