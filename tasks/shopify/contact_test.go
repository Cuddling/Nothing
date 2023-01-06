package shopify

import (
	"Mystery/tasks"
	"testing"
)

// Tests submitting contact information on oneness
func TestSubmitContactInfoOneness(t *testing.T) {
	site := tasks.Website{Name: "Oneness Boutique", Url: "https://www.onenessboutique.com"}
	variant := 39885244498004
	task := NewTaskShopify(site, &testProfile, nil, tasks.ModeShopifySafe, []string{string(rune(variant))}, []string{"L"})

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

	page, _, err := task.SubmitContactInfo()

	if err != nil {
		t.Fatal(err)
	}

	t.Log(page.Url)
}
