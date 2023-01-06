package shopify

import (
	"Mystery/tasks"
	"fmt"
	"testing"
	"time"
)

func TestGetPaymentToken(t *testing.T) {
	site := tasks.Website{Name: "Oneness Boutique", Url: "https://www.onenessboutique.com"}
	variant := 39885244498004
	task := NewTaskShopify(site, &testProfile, nil, tasks.ModeShopifySafe, []string{string(rune(variant))}, []string{"L"})

	token, _, err := task.GetPaymentToken()

	if err != nil {
		t.Fatal(err)
	}

	t.Log(fmt.Sprintf("Payment Token: %v", token))
}

// Tests submitting an order on oneness
func TestSubmitPaymentOneness(t *testing.T) {
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
			time.Sleep(750 * time.Millisecond)
			continue
		case nil:
			break
		default:
			t.Fatal(err)
		}
	}

	page, _, _ := task.SubmitShippingRate(rate)

	for page.GetCheckoutStep() == CheckoutStepCalculatingTaxes {
		ok, _, err := task.CalculateTaxes()

		if err != nil {
			t.Fatal(err)
		}

		if !ok {
			time.Sleep(1 * time.Second)
			continue
		}

		t.Log("Finished calculating taxes!")
		break
	}

	token, _, err := task.GetPaymentToken()

	if err != nil {
		t.Fatal(err)
	}

	page, _, err = task.SubmitPayment(false, token, "")

	if err != nil {
		t.Fatal(err)
	}

	t.Log(fmt.Sprintf("Redirected to: %v", page.Url))
}
