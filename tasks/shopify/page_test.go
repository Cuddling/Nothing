package shopify

import (
	"fmt"
	"testing"
)

// Tests retrieving an authenticity_token from a page
func TestGetAuthenticityToken(t *testing.T) {
	p := Page{
		Html: "<form class=\"edit_checkout\" action=\"/6269065/checkouts/b47d2e820ee2f10b267eca48abe743f0\" accept-charset=\"UTF-8\" method=\"post\"><input type=\"hidden\" name=\"_method\" value=\"patch\" autocomplete=\"off\" /><input type=\"hidden\" name=\"authenticity_token\" value=\"l9NDCVavRthdL3SDncpESnrrATdEdnk8MljYkXyglrKe5Tza5otp5kICfIEorx_YDwvInUBI3wlWYoYlOapxfg\" autocomplete=\"off\" />",
	}

	token, err := p.GetAuthenticityToken()

	if err != nil {
		t.Fatal(err)
	}

	t.Log(fmt.Sprintf("Auth Token: %v", token))
}

// Tests retrieving the total price the from the page
func TestGetTotalPrice(t *testing.T) {
	p := Page{
		Html: "<span class=\"payment-due__price skeleton-while-loading--lg\" data-checkout-payment-due-target=\"131739\">",
	}

	price, err := p.GetTotalPrice()

	if err != nil {
		t.Fatal(err)
	}

	t.Log(fmt.Sprintf("Total Price: %v", price))
}

// Tests getting the payment gateway from the page
func TestGetPaymentGateway(t *testing.T) {
	p := Page{
		Html: "data-select-gateway=\"26102467\"",
	}

	gateway, err := p.GetPaymentGateway()

	if err != nil {
		t.Fatal(err)
	}

	t.Log(fmt.Sprintf("Gateway: %v", gateway))
}

// Tests getting the notice text from the page
func TestGetNoticeText(t *testing.T) {
	p := Page{
		Html: "<p class=\"notice__text\">There was a problem processing the payment. Try refreshing this page or check your internet connection.</p>",
	}

	notice, err := p.GetNoticeText()

	if err != nil {
		t.Fatal(err)
	}

	t.Log(fmt.Sprintf("Notice Text: %v", notice))
}

// Tests getting the product title from the page
func TestGetProductTitle(t *testing.T) {
	p := Page{
		Html: "<span class=\"product__description__name order-summary__emphasis\">New Balance M2002RMB</span>",
	}

	title := p.GetProductTitle()

	if title != "New Balance M2002RMB" {
		t.Fatalf("Incorrect title. Got %v", title)
	}
}

// Tests getting the product size from the page
func TestGetProductSize(t *testing.T) {
	p := Page{
		Html: "<span class=\"product__description__variant order-summary__small-text\">11.5</span>",
	}

	title := p.GetProductSize()

	if title != "11.5" {
		t.Fatalf("Incorrect size. Got %v", title)
	}
}

// Tests getting the product image from the page
func TestGetProductImage(t *testing.T) {
	p := Page{
		Html: "<img alt=\"New Balance M2002RMB - 11.5\" class=\"product-thumbnail__image\" src=\"//cdn.shopify.com/s/files/1/0187/5180/products/new-balance-2002r-mule-m2002rmb_small.webp?v=1666364325\" />",
	}

	image := p.GetProductImage()

	if image != "https://cdn.shopify.com/s/files/1/0187/5180/products/new-balance-2002r-mule-m2002rmb_small.webp?v=1666364325" {
		t.Fatalf("Incorrect image. Got %v", image)
	}
}

func TestGetOrderNumber(t *testing.T) {
	p := Page{
		Html: "\t   \t<span class=\"os-order-number\">\n\t                    Order #SP456201\n\t                  </span>",
	}

	number := p.GetOrderNumber()

	if number == "" {
		t.Fatal("got empty order number")
	}

	fmt.Println(number)
}
