package shopify

import (
	"Mystery/profiles"
	"Mystery/tasks"
	"testing"
)

var testProfile = profiles.Profile{
	Name: "Test US",
	ShippingAddress: profiles.Address{
		Name:     "Johnny Appleseeds",
		Email:    "jappleseed124@gmail.com",
		Phone:    "1234567891",
		Line1:    "542 6th Ave",
		Line2:    "Apartment 2",
		City:     "New York",
		State:    "New York",
		PostCode: "10011",
		Country:  "United States",
	},
	BillingAddress: profiles.Address{
		Name:     "J. Appleseed",
		Email:    "jappleseed124@gmail.com",
		Phone:    "1234567891",
		Line1:    "1234 Parkway Ave",
		Line2:    "Unit 2",
		City:     "New York",
		State:    "New York",
		PostCode: "10011",
		Country:  "United States",
	},
	CreditCard: profiles.Card{
		Number:      "5555555555555555",
		CVV:         "123",
		ExpiryMonth: 04,
		ExpiryYear:  2028,
	},
	SameBillingAddressAsShipping: true,
}

func TestRunFlowOnenessHuman(t *testing.T) {
	site := tasks.Website{Name: "Oneness Boutique", Url: "https://www.onenessboutique.com"}
	task := NewTaskShopify(site, &testProfile, nil, tasks.ModeShopifySafe, []string{"39913832677460"}, []string{})
	task.Run()
}

func TestRunFlowOnenessFast(t *testing.T) {
	site := tasks.Website{Name: "Oneness Boutique", Url: "https://www.onenessboutique.com"}
	task := NewTaskShopify(site, &testProfile, nil, tasks.ModeShopifyFast, []string{"39913832677460"}, []string{})
	task.Run()
}

func TestRunFastFlowManiere(t *testing.T) {
	site := tasks.Website{Name: "A-Ma-Maniere", Url: "https://www.a-ma-maniere.com"}
	task := NewTaskShopify(site, &testProfile, nil, tasks.ModeShopifyFast, []string{"42257116233909"}, []string{})
	task.Run()
}

func TestRunFastFlowSlamJam(t *testing.T) {
	site := tasks.Website{Name: "IT Slam Jam", Url: "https://slamjam.com"}
	task := NewTaskShopify(site, &testProfile, nil, tasks.ModeShopifyFast, []string{"https://slamjam.com/collections/sneakers/products/nike-jordan-footwear-air-jordan-1-low-se-white-j259432"}, []string{})
	task.Run()
}

func TestRunFastFlowSlamJamQuantity3(t *testing.T) {
	site := tasks.Website{Name: "IT Slam Jam", Url: "https://slamjam.com"}
	task := NewTaskShopify(site, &testProfile, nil, tasks.ModeShopifyFast, []string{"https://slamjam.com/collections/sneakers/products/nike-jordan-footwear-air-jordan-1-low-se-white-j259432"}, []string{})
	task.Quantity = 3
	task.Run()
}

func TestRunFastFlowManor(t *testing.T) {
	site := tasks.Website{Name: "Manor.", Url: "https://www.manorphx.com"}
	task := NewTaskShopify(site, &testProfile, nil, tasks.ModeShopifyFast, []string{"https://www.manorphx.com/products/suplmnt-x-manor-24-ounce-wattle-bottle"}, []string{})
	task.Run()
}

func TestRunFastFlowKicksLounge(t *testing.T) {
	site := tasks.Website{Name: "KicksLounge", Url: "https://www.kickslounge.com"}
	task := NewTaskShopify(site, &testProfile, nil, tasks.ModeShopifyFast, []string{"https://www.kickslounge.com/products/copy-of-jordan-women-1-low-og-white-dark-powder-blue-cz0775-104"}, []string{})
	task.Run()
}

func TestRunFastFlowSneakerPolitics(t *testing.T) {
	site := tasks.Website{Name: "Sneaker Politics", Url: "https://sneakerpolitics.com"}
	task := NewTaskShopify(site, &testProfile, nil, tasks.ModeShopifyFast, []string{"42366366187708"}, []string{})
	task.Run()
}

func TestRunFastFlowConcepts(t *testing.T) {
	site := tasks.Website{Name: "Concepts", Url: "https://cncpts.com"}
	task := NewTaskShopify(site, &testProfile, nil, tasks.ModeShopifyFast, []string{"https://cncpts.com/products/concepts-concepts-rival-logo-ny-hoodie-cnho22-106-558-black"}, []string{})
	task.Run()
}

func TestRunFastFlowDTLR(t *testing.T) {
	site := tasks.Website{Name: "DTLR Villa", Url: "https://www.dtlr.com"}
	task := NewTaskShopify(site, &testProfile, nil, tasks.ModeShopifyFast, []string{"https://www.dtlr.com/collections/men-footwear-basketball/products/reebok-question-mid-denver-nuggets-gw8854-white-blue"}, []string{})
	task.Run()
}

func TestRunFastFlowKith(t *testing.T) {
	site := tasks.Website{Name: "ShoePalace", Url: "https://www.shoepalace.com"}
	task := NewTaskShopify(site, &testProfile, nil, tasks.ModeShopifyFast, []string{"https://www.shoepalace.com/products/nike-dd1503-101-dunk-low-womens-lifestyle-shoes-black-white"}, []string{})
	task.Run()
}

func TestRunFastFlowHatClub(t *testing.T) {
	site := tasks.Website{Name: "Hat Club", Url: "https://www.hatclub.com"}
	task := NewTaskShopify(site, &testProfile, nil, tasks.ModeShopifySafe, []string{"+a"}, []string{})
	task.Run()
}
