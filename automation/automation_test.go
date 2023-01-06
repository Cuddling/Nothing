package automation

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

// Tests if the product is a match via keywords
func TestAutomationProductMatchKeywords(t *testing.T) {
	auto := Automation{
		Name:          "Test",
		MonitorInputs: []string{"+dunk,+high,+bttys"},
		CheckUrl:      false,
	}

	zephyr := ZephyrMonitorLive{
		Body: ZephyrMonitorLiveBody{
			Payload: ZephyrMonitorLivePayload{
				Product: ZephyrMonitorLiveProduct{
					Handle: "nkdd1399-300",
					Title:  "Nike Dunk High Retro BTTYS - Noble Green / White",
				},
				Store: "https://kith.com/",
			},
		},
	}

	if !auto.IsProductMatch(&zephyr) {
		t.Fatalf("Product does not match")
	}
}

// Tests if the product is a match via keywords with check URL
func TestAutomationProductMatchCheckUrl(t *testing.T) {
	auto := Automation{
		Name:          "Test",
		MonitorInputs: []string{"+nkdd1399,+300"},
		CheckUrl:      true,
	}

	zephyr := ZephyrMonitorLive{
		Body: ZephyrMonitorLiveBody{
			Payload: ZephyrMonitorLivePayload{
				Product: ZephyrMonitorLiveProduct{
					Handle: "nkdd1399-300",
					Title:  "Nike Dunk High Retro BTTYS - Noble Green / White",
				},
				Store: "https://kith.com/",
			},
		},
	}

	if !auto.IsProductMatch(&zephyr) {
		t.Fatalf("Product does not match")
	}
}

// Tests if an automation website matches with whitelist only
func TestAutomationWebsiteMatchWhitelist(t *testing.T) {
	auto := Automation{
		Name:          "Test",
		SiteWhitelist: []string{"https://www.shoepalace.com"},
	}

	zephyr := ZephyrMonitorLive{
		Body: ZephyrMonitorLiveBody{
			Payload: ZephyrMonitorLivePayload{
				Store: "https://www.shoepalace.com/",
			},
		},
	}

	if !auto.IsWebsiteMatch(&zephyr) {
		t.Fatalf("Product website does not match")
	}
}

// Tests if an automation website matches with whitelist only
func TestAutomationWebsiteMatchBlacklist(t *testing.T) {
	auto := Automation{
		Name:          "Test",
		SiteBlacklist: []string{"https://kith.com"},
	}

	zephyr := ZephyrMonitorLive{
		Body: ZephyrMonitorLiveBody{
			Payload: ZephyrMonitorLivePayload{
				Store: "https://kith.com/",
			},
		},
	}

	if auto.IsWebsiteMatch(&zephyr) {
		t.Fatalf("Website is not being blacklisted")
	}
}

// Tests price matching for automation
func TestAutomationPriceMatch(t *testing.T) {
	auto := Automation{
		Name:         "Test",
		PriceMinimum: 90,
		PriceMaximum: 125,
	}

	zephyr := ZephyrMonitorLive{
		Body: ZephyrMonitorLiveBody{
			Payload: ZephyrMonitorLivePayload{
				Product: ZephyrMonitorLiveProduct{
					Variants: []ZephyrMonitorLiveProductVariant{
						{
							Price: "124.69",
						},
					},
				},
				Store: "",
			},
		},
	}

	if !auto.IsPriceMatch(&zephyr) {
		t.Fatalf("Price does not match")
	}
}

func TestAutomationGetSizeMatch(t *testing.T) {
	const data = "{\"body\":{\"payload\":{\"product\":{\"requires_selling_plan\":false,\"available\":true,\"created_at\":\"2022-12-02T15:14:17-06:00\",\"variants\":[{\"inventory_quantity\":2,\"requires_selling_plan\":false,\"compare_at_price\":20000,\"taxable\":true,\"inventory_management\":\"shopify\",\"requires_shipping\":true,\"available\":true,\"weight\":1361,\"title\":\"8\",\"featured_image\":null,\"public_title\":\"8\",\"inventory_policy\":\"deny\",\"price\":\"200.0\",\"option3\":null,\"name\":\"12 RETRO - 8\",\"options\":[\"8\"],\"selling_plan_allocations\":[],\"option1\":\"8\",\"id\":40845477576748,\"option2\":null,\"sku\":\"CT8013-071-8\",\"barcode\":\"\"},{\"inventory_quantity\":2,\"requires_selling_plan\":false,\"compare_at_price\":20000,\"taxable\":true,\"inventory_management\":\"shopify\",\"requires_shipping\":true,\"available\":true,\"weight\":1361,\"title\":\"8.5\",\"featured_image\":null,\"public_title\":\"8.5\",\"inventory_policy\":\"deny\",\"price\":20000,\"option3\":null,\"name\":\"12 RETRO - 8.5\",\"options\":[\"8.5\"],\"selling_plan_allocations\":[],\"option1\":\"8.5\",\"id\":40845477609516,\"option2\":null,\"sku\":\"CT8013-071-8.5\",\"barcode\":\"\"},{\"inventory_quantity\":1,\"requires_selling_plan\":false,\"compare_at_price\":20000,\"taxable\":true,\"inventory_management\":\"shopify\",\"requires_shipping\":true,\"available\":true,\"weight\":1361,\"title\":\"9\",\"featured_image\":null,\"public_title\":\"9\",\"inventory_policy\":\"deny\",\"price\":20000,\"option3\":null,\"name\":\"12 RETRO - 9\",\"options\":[\"9\"],\"selling_plan_allocations\":[],\"option1\":\"9\",\"id\":40845477642284,\"option2\":null,\"sku\":\"CT8013-071-9\",\"barcode\":\"\"},{\"inventory_quantity\":1,\"requires_selling_plan\":false,\"compare_at_price\":20000,\"taxable\":true,\"inventory_management\":\"shopify\",\"requires_shipping\":true,\"available\":true,\"weight\":1361,\"title\":\"9.5\",\"featured_image\":null,\"public_title\":\"9.5\",\"inventory_policy\":\"deny\",\"price\":20000,\"option3\":null,\"name\":\"12 RETRO - 9.5\",\"options\":[\"9.5\"],\"selling_plan_allocations\":[],\"option1\":\"9.5\",\"id\":40845477675052,\"option2\":null,\"sku\":\"CT8013-071-9.5\",\"barcode\":\"\"},{\"inventory_quantity\":0,\"requires_selling_plan\":false,\"compare_at_price\":20000,\"taxable\":true,\"inventory_management\":\"shopify\",\"requires_shipping\":true,\"available\":false,\"weight\":1361,\"title\":\"10\",\"featured_image\":null,\"public_title\":\"10\",\"inventory_policy\":\"deny\",\"price\":20000,\"option3\":null,\"name\":\"12 RETRO - 10\",\"options\":[\"10\"],\"selling_plan_allocations\":[],\"option1\":\"10\",\"id\":40845477707820,\"option2\":null,\"sku\":\"CT8013-071-10\",\"barcode\":\"\"},{\"inventory_quantity\":0,\"requires_selling_plan\":false,\"compare_at_price\":20000,\"taxable\":true,\"inventory_management\":\"shopify\",\"requires_shipping\":true,\"available\":false,\"weight\":1361,\"title\":\"10.5\",\"featured_image\":null,\"public_title\":\"10.5\",\"inventory_policy\":\"deny\",\"price\":20000,\"option3\":null,\"name\":\"12 RETRO - 10.5\",\"options\":[\"10.5\"],\"selling_plan_allocations\":[],\"option1\":\"10.5\",\"id\":40845477740588,\"option2\":null,\"sku\":\"CT8013-071-10.5\",\"barcode\":\"\"},{\"inventory_quantity\":0,\"requires_selling_plan\":false,\"compare_at_price\":20000,\"taxable\":true,\"inventory_management\":\"shopify\",\"requires_shipping\":true,\"available\":false,\"weight\":1361,\"title\":\"11\",\"featured_image\":null,\"public_title\":\"11\",\"inventory_policy\":\"deny\",\"price\":20000,\"option3\":null,\"name\":\"12 RETRO - 11\",\"options\":[\"11\"],\"selling_plan_allocations\":[],\"option1\":\"11\",\"id\":40845477773356,\"option2\":null,\"sku\":\"CT8013-071-11\",\"barcode\":\"\"},{\"inventory_quantity\":0,\"requires_selling_plan\":false,\"compare_at_price\":20000,\"taxable\":true,\"inventory_management\":\"shopify\",\"requires_shipping\":true,\"available\":false,\"weight\":1361,\"title\":\"11.5\",\"featured_image\":null,\"public_title\":\"11.5\",\"inventory_policy\":\"deny\",\"price\":20000,\"option3\":null,\"name\":\"12 RETRO - 11.5\",\"options\":[\"11.5\"],\"selling_plan_allocations\":[],\"option1\":\"11.5\",\"id\":40845477806124,\"option2\":null,\"sku\":\"CT8013-071-11.5\",\"barcode\":\"\"},{\"inventory_quantity\":0,\"requires_selling_plan\":false,\"compare_at_price\":20000,\"taxable\":true,\"inventory_management\":\"shopify\",\"requires_shipping\":true,\"available\":false,\"weight\":1361,\"title\":\"12\",\"featured_image\":null,\"public_title\":\"12\",\"inventory_policy\":\"deny\",\"price\":20000,\"option3\":null,\"name\":\"12 RETRO - 12\",\"options\":[\"12\"],\"selling_plan_allocations\":[],\"option1\":\"12\",\"id\":40845477838892,\"option2\":null,\"sku\":\"CT8013-071-12\",\"barcode\":\"\"}],\"title\":\"12 RETRO\",\"type\":\"Footwear\",\"price_min\":20000,\"uuid\":\"f08ae64b-f30b-4afc-bd93-afcb2e6656dc\",\"price_varies\":false,\"filtered\":true,\"vendor\":\"JORDAN\",\"price\":20000,\"options\":[{\"values\":[\"8\",\"8.5\",\"9\",\"9.5\",\"10\",\"10.5\",\"11\",\"11.5\",\"12\"],\"name\":\"Size\",\"position\":1}],\"selling_plan_groups\":[],\"id\":7009993883692,\"published_at\":\"2022-12-03T09:00:02-06:00\",\"timestamp\":1670255885034,\"images\":[\"//cdn.shopify.com/s/files/1/0268/1785/products/Air-Jordan-12-Retro-Black-Taxi-CT8013-071-1.jpg?v=1670019145\",\"//cdn.shopify.com/s/files/1/0268/1785/products/Air-Jordan-12-Retro-Black-Taxi-CT8013-071-2.jpg?v=1670019145\",\"//cdn.shopify.com/s/files/1/0268/1785/products/Air-Jordan-12-Retro-Black-Taxi-CT8013-071-3.jpg?v=1670019145\",\"//cdn.shopify.com/s/files/1/0268/1785/products/Air-Jordan-12-Retro-Black-Taxi-CT8013-071-4.jpg?v=1670019145\"],\"compare_at_price\":20000,\"compare_at_price_varies\":false,\"handle\":\"12-retro-6\",\"compare_at_price_min\":20000,\"featured_image\":\"//cdn.shopify.com/s/files/1/0268/1785/products/Air-Jordan-12-Retro-Black-Taxi-CT8013-071-1.jpg?v=1670019145\",\"url\":\"/products/12-retro-6\",\"compare_at_price_max\":20000,\"price_max\":20000},\"store\":\"https://www.saintalfred.com/\"},\"channel\":\"store\",\"type\":\"shopify\",\"event\":\"newProduct\"},\"type\":\"livemonitor\"}"

	var liveProduct ZephyrMonitorLive

	if err := json.Unmarshal([]byte(data), &liveProduct); err != nil {
		log.Printf("Failed to unmarshal monitor message: %v\n", err)
		return
	}

	autoRandom := Automation{
		Sizes: []string{},
	}

	if len(autoRandom.GetMatchingSizeVariants(&liveProduct)) != 4 {
		t.Fatal("expected 9 matching variants for random")
	}

	autoTwo := Automation{
		Sizes: []string{"9", "8"},
	}

	if len(autoTwo.GetMatchingSizeVariants(&liveProduct)) != 2 {
		t.Fatal("expected 2 matching variants")
	}
}

func TestSendAutomationWebhook(t *testing.T) {
	a := Automation{
		Name:             "TEST",
		MonitorInputs:    nil,
		Sizes:            []string{"7.5", "12.5"},
		Profiles:         nil,
		ProxyList:        "",
		PriceMinimum:     0,
		PriceMaximum:     0,
		Quantity:         1,
		TotalTaskCount:   100,
		SiteWhitelist:    nil,
		SiteBlacklist:    nil,
		PaymentRetries:   1,
		StopAfterMinutes: 5,
	}

	const data string = "{\"body\":{\"payload\":{\"product\":{\"images\":[{\"updated_at\":\"2022-11-15T10:08:54-06:00\",\"src\":\"https://cdn.shopify.com/s/files/1/0018/4506/7865/products/HP6586_5_FOOTWEAR_Photography_SideMedialCenterView_white.jpg?v=1668528534\",\"product_id\":7732447314107,\"width\":1200,\"created_at\":\"2022-11-15T10:08:48-06:00\",\"variant_ids\":[],\"id\":33878767763643,\"position\":1,\"height\":1200},{\"updated_at\":\"2022-11-15T10:08:54-06:00\",\"src\":\"https://cdn.shopify.com/s/files/1/0018/4506/7865/products/HP6586_1_FOOTWEAR_Photography_SideLateralCenterView_white.jpg?v=1668528534\",\"product_id\":7732447314107,\"width\":1200,\"created_at\":\"2022-11-15T10:08:48-06:00\",\"variant_ids\":[],\"id\":33878767534267,\"position\":2,\"height\":1200},{\"updated_at\":\"2022-11-15T10:08:54-06:00\",\"src\":\"https://cdn.shopify.com/s/files/1/0018/4506/7865/products/HP6586_8_FOOTWEAR_Photography_DetailView1_white.jpg?v=1668528534\",\"product_id\":7732447314107,\"width\":1200,\"created_at\":\"2022-11-15T10:08:48-06:00\",\"variant_ids\":[],\"id\":33878767599803,\"position\":3,\"height\":1200},{\"updated_at\":\"2022-11-15T10:08:54-06:00\",\"src\":\"https://cdn.shopify.com/s/files/1/0018/4506/7865/products/HP6586_9_FOOTWEAR_Photography_DetailView2_white.jpg?v=1668528534\",\"product_id\":7732447314107,\"width\":1200,\"created_at\":\"2022-11-15T10:08:48-06:00\",\"variant_ids\":[],\"id\":33878767665339,\"position\":4,\"height\":1200},{\"updated_at\":\"2022-11-15T10:08:54-06:00\",\"src\":\"https://cdn.shopify.com/s/files/1/0018/4506/7865/products/HP6586_3_FOOTWEAR_Photography_TopPortraitView_white.jpg?v=1668528534\",\"product_id\":7732447314107,\"width\":1200,\"created_at\":\"2022-11-15T10:08:48-06:00\",\"variant_ids\":[],\"id\":33878767501499,\"position\":5,\"height\":1200}],\"available\":true,\"created_at\":\"2022-11-15T10:03:43-06:00\",\"handle\":\"adifom-q-cblack-carbon-gresix-laces-mexico\",\"variants\":[{\"formatted_price\":\"$ 2,599\",\"compare_at_price\":\"2599.00\",\"taxable\":true,\"requires_shipping\":true,\"option_values\":[{\"name\":\"Size\",\"option_id\":9806870839483,\"value\":\"6.5US-4.5MX\"}],\"available\":false,\"created_at\":\"2022-11-15T10:03:43-06:00\",\"title\":\"6.5US-4.5MX\",\"updated_at\":\"2022-11-15T10:09:28-06:00\",\"price\":\"2599.0\",\"id\":43896303911099,\"position\":1,\"grams\":0,\"sku\":\"210000186284\",\"barcode\":\"4066747048275\"},{\"formatted_price\":\"$ 2,599\",\"compare_at_price\":\"2599.00\",\"taxable\":true,\"requires_shipping\":true,\"option_values\":[{\"name\":\"Size\",\"option_id\":9806870839483,\"value\":\"7US-5MX\"}],\"available\":true,\"created_at\":\"2022-11-15T10:03:43-06:00\",\"title\":\"7US-5MX\",\"updated_at\":\"2022-11-26T20:13:16-06:00\",\"price\":\"2599.00\",\"id\":43896303583419,\"position\":2,\"grams\":0,\"sku\":\"210000185357\",\"barcode\":\"4066747044611\"},{\"formatted_price\":\"$ 2,599\",\"compare_at_price\":\"2599.00\",\"taxable\":true,\"requires_shipping\":true,\"option_values\":[{\"name\":\"Size\",\"option_id\":9806870839483,\"value\":\"7.5US-5.5MX\"}],\"available\":true,\"created_at\":\"2022-11-15T10:03:43-06:00\",\"title\":\"7.5US-5.5MX\",\"updated_at\":\"2022-11-15T10:07:56-06:00\",\"price\":\"2599.00\",\"id\":43896303616187,\"position\":3,\"grams\":0,\"sku\":\"210000185358\",\"barcode\":\"4066747048299\"},{\"formatted_price\":\"$ 2,599\",\"compare_at_price\":\"2599.00\",\"taxable\":true,\"requires_shipping\":true,\"option_values\":[{\"name\":\"Size\",\"option_id\":9806870839483,\"value\":\"8US-6MX\"}],\"available\":true,\"created_at\":\"2022-11-15T10:03:43-06:00\",\"title\":\"8US-6MX\",\"updated_at\":\"2022-11-15T10:07:58-06:00\",\"price\":\"2599.00\",\"id\":43896303648955,\"position\":4,\"grams\":0,\"sku\":\"210000185359\",\"barcode\":\"4066747044543\"},{\"formatted_price\":\"$ 2,599\",\"compare_at_price\":\"2599.00\",\"taxable\":true,\"requires_shipping\":true,\"option_values\":[{\"name\":\"Size\",\"option_id\":9806870839483,\"value\":\"8.5US-6.5MX\"}],\"available\":true,\"created_at\":\"2022-11-15T10:03:43-06:00\",\"title\":\"8.5US-6.5MX\",\"updated_at\":\"2022-11-21T16:01:05-06:00\",\"price\":\"2599.00\",\"id\":43896303681723,\"position\":5,\"grams\":0,\"sku\":\"210000185360\",\"barcode\":\"4066747044581\"},{\"formatted_price\":\"$ 2,599\",\"compare_at_price\":\"2599.00\",\"taxable\":true,\"requires_shipping\":true,\"option_values\":[{\"name\":\"Size\",\"option_id\":9806870839483,\"value\":\"9US-7MX\"}],\"available\":true,\"created_at\":\"2022-11-15T10:03:43-06:00\",\"title\":\"9US-7MX\",\"updated_at\":\"2022-11-15T10:07:56-06:00\",\"price\":\"2599.00\",\"id\":43896303714491,\"position\":6,\"grams\":0,\"sku\":\"210000185361\",\"barcode\":\"4066747048336\"},{\"formatted_price\":\"$ 2,599\",\"compare_at_price\":\"2599.00\",\"taxable\":true,\"requires_shipping\":true,\"option_values\":[{\"name\":\"Size\",\"option_id\":9806870839483,\"value\":\"9.5US-7.5MX\"}],\"available\":true,\"created_at\":\"2022-11-15T10:03:43-06:00\",\"title\":\"9.5US-7.5MX\",\"updated_at\":\"2022-11-15T10:07:56-06:00\",\"price\":\"2599.00\",\"id\":43896303747259,\"position\":7,\"grams\":0,\"sku\":\"210000185362\",\"barcode\":\"4066747044642\"},{\"formatted_price\":\"$ 2,599\",\"compare_at_price\":\"2599.00\",\"taxable\":true,\"requires_shipping\":true,\"option_values\":[{\"name\":\"Size\",\"option_id\":9806870839483,\"value\":\"10US-8MX\"}],\"available\":true,\"created_at\":\"2022-11-15T10:03:43-06:00\",\"title\":\"10US-8MX\",\"updated_at\":\"2022-11-24T19:11:16-06:00\",\"price\":\"2599.00\",\"id\":43896303780027,\"position\":8,\"grams\":0,\"sku\":\"210000185363\",\"barcode\":\"4066747048343\"},{\"formatted_price\":\"$ 2,599\",\"compare_at_price\":\"2599.00\",\"taxable\":true,\"requires_shipping\":true,\"option_values\":[{\"name\":\"Size\",\"option_id\":9806870839483,\"value\":\"10.5US-8.5MX\"}],\"available\":false,\"created_at\":\"2022-11-15T10:03:43-06:00\",\"title\":\"10.5US-8.5MX\",\"updated_at\":\"2022-11-22T14:36:21-06:00\",\"price\":\"2599.00\",\"id\":43896303812795,\"position\":9,\"grams\":0,\"sku\":\"210000185364\",\"barcode\":\"4066747044598\"},{\"formatted_price\":\"$ 2,599\",\"compare_at_price\":\"2599.00\",\"taxable\":true,\"requires_shipping\":true,\"option_values\":[{\"name\":\"Size\",\"option_id\":9806870839483,\"value\":\"11US-9MX\"}],\"available\":false,\"created_at\":\"2022-11-15T10:03:43-06:00\",\"title\":\"11US-9MX\",\"updated_at\":\"2022-11-30T19:31:17-06:00\",\"price\":\"2599.00\",\"id\":43896303845563,\"position\":10,\"grams\":0,\"sku\":\"210000185365\",\"barcode\":\"4066747044628\"},{\"formatted_price\":\"$ 2,599\",\"compare_at_price\":\"2599.00\",\"taxable\":true,\"requires_shipping\":true,\"option_values\":[{\"name\":\"Size\",\"option_id\":9806870839483,\"value\":\"11.5US-9.5MX\"}],\"available\":false,\"created_at\":\"2022-11-15T10:03:43-06:00\",\"title\":\"11.5US-9.5MX\",\"updated_at\":\"2022-11-15T10:09:29-06:00\",\"price\":\"2599.00\",\"id\":43896303943867,\"position\":11,\"grams\":0,\"sku\":\"210000186285\",\"barcode\":\"4066747044567\"},{\"formatted_price\":\"$ 2,599\",\"compare_at_price\":\"2599.00\",\"taxable\":true,\"requires_shipping\":true,\"option_values\":[{\"name\":\"Size\",\"option_id\":9806870839483,\"value\":\"12US-10MX\"}],\"available\":false,\"created_at\":\"2022-11-15T10:03:43-06:00\",\"title\":\"12US-10MX\",\"updated_at\":\"2022-11-23T12:32:16-06:00\",\"price\":\"2599.00\",\"id\":43896303550651,\"position\":12,\"grams\":0,\"sku\":\"210000185356\",\"barcode\":\"4066747048251\"},{\"formatted_price\":\"$ 2,599\",\"compare_at_price\":\"2599.00\",\"taxable\":true,\"requires_shipping\":true,\"option_values\":[{\"name\":\"Size\",\"option_id\":9806870839483,\"value\":\"12.5US-10.5MX\"}],\"available\":false,\"created_at\":\"2022-11-15T10:03:43-06:00\",\"title\":\"12.5US-10.5MX\",\"updated_at\":\"2022-11-15T10:09:29-06:00\",\"price\":\"2599.00\",\"id\":43896303878331,\"position\":13,\"grams\":0,\"sku\":\"210000186283\",\"barcode\":\"4066747044635\"}],\"title\":\"adidas adiFOM Q Gresix\",\"uuid\":\"6a81ad13-c165-4a47-a70a-885919d73c06\",\"product_type\":\"FTW - SNEAKERS - MEN - ADULT - ONE SHOT - WI - 22\",\"updated_at\":\"2022-12-01T09:43:30-06:00\",\"filtered\":true,\"vendor\":\"ADIDAS\",\"product_id\":7732447314107,\"options\":[{\"product_id\":7732447314107,\"values\":[\"6.5US-4.5MX\",\"7US-5MX\",\"7.5US-5.5MX\",\"8US-6MX\",\"8.5US-6.5MX\",\"9US-7MX\",\"9.5US-7.5MX\",\"10US-8MX\",\"10.5US-8.5MX\",\"11US-9MX\",\"11.5US-9.5MX\",\"12US-10MX\",\"12.5US-10.5MX\"],\"name\":\"Size\",\"id\":9806870839483,\"position\":1}],\"id\":7732447314107,\"published_at\":\"2022-11-15T10:10:27-06:00\",\"timestamp\":1669909412989},\"store\":\"https://www.laces.mx/\"},\"channel\":\"store\",\"type\":\"shopify\",\"event\":\"newProduct\"},\"type\":\"livemonitor\"}"

	var product ZephyrMonitorLive

	if err := json.Unmarshal([]byte(data), &product); err != nil {
		log.Printf("Failed to unmarshal monitor message: %v\n", err)
		return
	}

	fmt.Println(product.Body.Payload.Product.GetImage())
	a.SendWebhook(true, &product)
}
