package discord

import (
	"testing"
)

func TestSuccessWebhook(t *testing.T) {
	for i := 0; i < 1; i++ {
		data := CheckoutWebhookData{
			Success:       true,
			FailureReason: "",
			Site:          "Sneaker Politics",
			Mode:          "Fast (Auto)",
			ProductTitle:  "Nike Dunk Low - Test",
			ProductSize:   "12",
			ProductImage:  "https://cdn.shopify.com/s/files/1/0214/7974/products/DD8338-001-PHSLH000-2000.jpg",
			Profile:       "Test Profile",
			ProxyList:     "Test Proxies",
			Email:         "test@gmail.com",
			OrderNumber:   "#123456",
			OrderLink:     "https://sneakerpolitics.com",
		}

		SendCheckoutWebhook(data)
	}
}

func TestFailureWebhook(t *testing.T) {
	data := CheckoutWebhookData{
		Success:       false,
		FailureReason: "Card was declined",
		Site:          "Kith",
		Mode:          "Fast",
		ProductTitle:  "Nike Dunk High Retro Bttys Noble Green / White",
		ProductSize:   "5.5",
		ProductImage:  "https://cdn.shopify.com/s/files/1/0094/2252/products/DD1399-300-PHSRH000-2000_1080x.jpg?v=1659621681",
		Profile:       "Test US",
		ProxyList:     "Live",
		Email:         "test@gmail.com",
		OrderNumber:   "12345",
		OrderLink:     "https://kith.com",
	}

	SendCheckoutWebhook(data)
}
