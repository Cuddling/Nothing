package shopify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
)

// AddToCartResponse Struct when adding an item to cart
type AddToCartResponse struct {
	Id                           int64         `json:"id,omitempty"`
	Quantity                     int           `json:"quantity,omitempty"`
	VariantId                    int64         `json:"variant_id,omitempty"`
	Key                          string        `json:"key,omitempty"`
	Title                        string        `json:"title,omitempty"`
	Price                        int           `json:"price,omitempty"`
	OriginalPrice                int           `json:"original_price,omitempty"`
	DiscountedPrice              int           `json:"discounted_price,omitempty"`
	LinePrice                    int           `json:"line_price,omitempty"`
	TotalDiscount                int           `json:"total_discount,omitempty"`
	Discounts                    interface{}   `json:"discounts"`
	Sku                          string        `json:"sku,omitempty"`
	Grams                        int           `json:"grams,omitempty"`
	Vendor                       string        `json:"vendor,omitempty"`
	Taxable                      bool          `json:"taxable,omitempty"`
	ProductId                    int64         `json:"product_id,omitempty"`
	ProductHasOnlyDefaultVariant bool          `json:"product_has_only_default_variant,omitempty"`
	GiftCard                     bool          `json:"gift_card,omitempty"`
	FinalPrice                   int           `json:"final_price,omitempty"`
	FinalLinePrice               int           `json:"final_line_price,omitempty"`
	Url                          string        `json:"url,omitempty"`
	Image                        string        `json:"image,omitempty"`
	Handle                       string        `json:"handle,omitempty"`
	RequiresShipping             bool          `json:"requires_shipping,omitempty"`
	ProductType                  string        `json:"product_type,omitempty"`
	ProductTitle                 string        `json:"product_title,omitempty"`
	ProductDescription           string        `json:"product_description,omitempty"`
	VariantTitle                 string        `json:"variant_title,omitempty"`
	VariantOptions               []string      `json:"variant_options,omitempty"`
	OptionsWithValues            interface{}   `json:"options_with_values,omitempty"`
	LineLevelDiscountAllocations []interface{} `json:"line_level_discount_allocations,omitempty"`
	LineLevelTotalDiscount       int           `json:"line_level_total_discount,omitempty"`
}

// ClearCartResponse Struct when clearing the cart of items
type ClearCartResponse struct {
	Token            string      `json:"token,omitempty"`
	Note             interface{} `json:"note,omitempty"`
	Attributes       interface{} `json:"attributes,omitempty"`
	TotalPrice       int         `json:"total_price,omitempty"`
	TotalWeight      json.Number `json:"total_weight,omitempty"`
	ItemCount        int         `json:"item_count,omitempty"`
	Items            interface{} `json:"items,omitempty"`
	RequiresShipping bool        `json:"requires_shipping,omitempty"`
}

// AddToCart Adds a variant to the cart
func (t *Task) AddToCart(variant int64) (AddToCartResponse, *resty.Response, error) {

	resp, err := t.Client.R().
		SetHeader("Accept", "application/json, text/javascript, */*. q=0.01").
		SetHeader("Accept-Encoding", "gzip, deflate").
		SetHeader("Accept-Language", "en-US,en.q=0.9").
		SetHeader("Connection", "keep-alive").
		SetHeader("Content-Type", "application/x-www-form-urlencoded. charset=UTF-8").
		SetHeader("Host", t.Site.GetHostName()).
		SetHeader("Origin", t.Site.Url).
		SetHeader("Referer", fmt.Sprintf("%v/collections/all", t.Site.Url)).
		SetHeader("sec-ch-ua", "\" Not A.Brand\".v=\"99\", \"Chromium\".v=\"102\", \"Google Chrome\".v=\"102\"").
		SetHeader("sec-ch-ua-mobile", "?0").
		SetHeader("sec-ch-ua-platform", "\"Windows\"").
		SetHeader("Sec-Fetch-Dest", "empty").
		SetHeader("Sec-Fetch-Mode", "cors").
		SetHeader("Sec-Fetch-Site", "same-origin").
		// SetHeader("X-Requested-With", "XMLHttpRequest"). - Gives 422 OOS Error when uncommented
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0. Win64. x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.62 Safari/537.36").
		SetFormData(map[string]string{
			"form_type": "product",
			"utf8":      "✓",
			"id":        strconv.FormatInt(variant, 10),
			"quantity":  strconv.Itoa(t.Quantity),
		}).
		Post(fmt.Sprintf("%v/cart/add.js", t.Site.Url))

	if err != nil {
		return AddToCartResponse{}, resp, err
	}

	if resp.IsError() {
		return AddToCartResponse{}, resp, errors.New(fmt.Sprintf("request failed with status %v", resp.StatusCode()))
	}

	var data AddToCartResponse
	err = json.Unmarshal(resp.Body(), &data)

	if err != nil {
		return AddToCartResponse{}, resp, err
	}

	if data.VariantId != variant {
		return AddToCartResponse{}, resp, errors.New("variant was not added to the cart")
	}

	t.VariantInCart = variant
	return data, resp, nil
}

// ClearCart Removes all items from the cart
func (t *Task) ClearCart() (ClearCartResponse, *resty.Response, error) {
	resp, err := t.Client.R().
		SetHeader("Accept", "application/json, text/javascript, */*. q=0.01").
		SetHeader("Accept-Encoding", "gzip, deflate").
		SetHeader("Accept-Language", "en-US,en.q=0.9").
		SetHeader("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8").
		SetHeader("origin", t.Site.Url).
		SetHeader("referer", t.Site.Url).
		SetHeader("sec-ch-ua", "\" Not A.Brand\".v=\"99\", \"Chromium\".v=\"102\", \"Google Chrome\".v=\"102\"").
		SetHeader("sec-ch-ua-mobile", "?0").
		SetHeader("sec-ch-ua-platform", "\"Windows\"").
		SetHeader("Sec-Fetch-Dest", "empty").
		SetHeader("Sec-Fetch-Mode", "cors").
		SetHeader("Sec-Fetch-Site", "same-origin").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0. Win64. x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.62 Safari/537.36\"").
		SetHeader("X-Requested-With", "XMLHttpRequest").
		Post(fmt.Sprintf("%v/cart/clear.js", t.Site.Url))

	if err != nil {
		return ClearCartResponse{}, resp, err
	}

	if resp.IsError() {
		return ClearCartResponse{}, resp, errors.New(fmt.Sprintf("cart clear failed with status: %v", resp.StatusCode()))
	}

	var data ClearCartResponse
	err = json.Unmarshal(resp.Body(), &data)

	if err != nil {
		return ClearCartResponse{}, resp, err
	}

	if data.ItemCount != 0 {
		return data, resp, errors.New(fmt.Sprintf("expected 0 item_count but got %v", data.ItemCount))
	}

	t.VariantInCart = 0
	return data, resp, nil
}
