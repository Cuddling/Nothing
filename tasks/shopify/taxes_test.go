package shopify

import (
	"Mystery/tasks"
	"testing"
	"time"
)

// Tests calculating taxes on A-Ma-Maniere
// https://www.a-ma-maniere.com/collections/amiri/products/amiri-army-stencil-hoodie
func TestCalculateTaxesManiere(t *testing.T) {
	site := tasks.Website{Name: "A-Ma-Maniere", Url: "https://www.a-ma-maniere.com"}
	variant := 42319416819893
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
}
