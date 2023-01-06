package shopify

import (
	"Mystery/profiles"
	"Mystery/tasks"
	"Mystery/utils"
	"strconv"
	"testing"
)

// Tests monitoring a product by variants
func TestMonitorProductsVariants(t *testing.T) {
	site := tasks.Website{Name: "Kith", Url: "https://kith.com"}
	variant := "39250265571456"
	task := NewTaskShopify(site, &profiles.Profile{}, nil, tasks.ModeShopifySafe, []string{variant}, []string{})

	_, err := task.MonitorProducts()

	if err != nil {
		t.Fatal(err)
	}

	if strconv.FormatInt(task.Variant.Id, 10) != variant {
		t.Fatal("task product variant does not match the given one")
	}
}

// Tests monitoring a product by URL
func TestMonitorProductsUrlRandom(t *testing.T) {
	site := tasks.Website{Name: "Kith", Url: "https://kith.com"}
	task := NewTaskShopify(site, &profiles.Profile{}, nil, tasks.ModeShopifySafe, []string{"https://kith.com/products/nkdd1399-300"}, []string{})

	_, err := task.MonitorProducts()

	if err != nil {
		t.Fatal(err)
	}

	if task.Product.Id == 0 {
		t.Fatal("got empty product")
	}

	if task.Variant.Id == 0 {
		t.Fatal("got empty product variant")
	}
}

// Tests monitoring a product by URL on Kith (Raw numbers)
func TestMonitorProductsUrlRangeKith(t *testing.T) {
	site := tasks.Website{Name: "Kith", Url: "https://kith.com"}
	task := NewTaskShopify(site, &profiles.Profile{}, nil, tasks.ModeShopifySafe, []string{"https://kith.com/products/nkdd1399-300"}, []string{"9"})

	_, err := task.MonitorProducts()

	if err != nil {
		t.Fatal(err)
	}

	if task.Product.Id == 0 {
		t.Fatal("got empty product")
	}

	if task.Variant.Id == 0 {
		t.Fatal("got empty product variant")
	}

	if task.Variant.Option1 != "9" {
		t.Fatal("got incorrect size")
	}
}

// Tests monitoring a product by URL on DSMNY (US 5)
func TestMonitorProductsUrlRangeDSMNY(t *testing.T) {
	site := tasks.Website{Name: "DSMNY", Url: "https://shop-us.doverstreetmarket.com"}
	url := "https://shop-us.doverstreetmarket.com/collections/sneaker-space-adidas-1/products/adidas-x-yeezy-yeezy-boost-350-v2-flax-fx9028-aw22"
	task := NewTaskShopify(site, &profiles.Profile{}, nil, tasks.ModeShopifySafe, []string{url}, []string{"7.5"})

	_, err := task.MonitorProducts()

	if err != nil {
		t.Fatal(err)
	}

	if task.Product.Id == 0 {
		t.Fatal("got empty product")
	}

	if task.Variant.Id == 0 {
		t.Fatal("got empty product variant")
	}

	if task.Variant.Option1 != "Size US 7.5" {
		t.Fatal("got incorrect size")
	}
}

// Tests monitoring a product by keywords on Kith
func TestMonitorProductsKeywords(t *testing.T) {
	site := tasks.Website{Name: "Kith", Url: "https://kith.com"}
	task := NewTaskShopify(site, &profiles.Profile{}, nil, tasks.ModeShopifySafe, []string{"+a"}, []string{})

	_, err := task.MonitorProducts()

	if err != nil {
		t.Fatal(err)
	}

	if task.Product.Id == 0 {
		t.Fatal("got empty product")
	}

	if task.Variant.Id == 0 {
		t.Fatal("got empty product variant")
	}
}

func TestKws(t *testing.T) {
	kw := "+dunk,+low,-necklace,-peace,-cream,-halloween,-womens,-(w),-women,-wmn,-scream,-arizona,-miami,-gorge,-women's,-wmns,-qs,-scream,-kumquat,-racer,-mujer,-keychain,-picture,-store,-raffle,-instore,-photo,-sweater,-molina,-bot,-image,-fake,-infant,-toddler,-td,-toy,-poster,-card,-bag,-cap,-hat,-zen,-defiant,-deviant,-zoom,-inf,-coupon,-couture,-flight,-fearless,-satin,-ps,-infant,-polo,-pants,-shirt,-hoodie,-crib,-disrupt,-beanie,-top,-under,-bottom,-iso,-jan,-truck,-book,-jacket,-ease,-mid,-disrupt,-ps,-td,-little,-toddler,-infant,-book,-candle,-tee,-hood,-hoodie,-jacket,-deck,-skateboard,-omni,-rebel,-up,-lemon,-cashmere,-sunset,-aged,-easter,-camo,-up,-rebel,-animal,-1985,-acid,-crater,-instinct,-gold,-metallic,-sequoia,-primal,-avocado,-avacado,-banana,-safari,-barbershop,-prism,-teal,-lisa,-leslie,-nh,-rough,-nature,-comme,-cider,-emerald,-velvet,-dunk low nn,-noir,-coconut,-milk,-quilt,-wheat,-flax,-oxford,-safety,-sail,-scream,-quilted"

	if !utils.IsKeywordMatch("Nike Dunk Low Retro Premium - Vast Grey/Summit White", kw) {
		t.Fatal("no match")
	}
}
