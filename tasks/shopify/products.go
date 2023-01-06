package shopify

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

type Product struct {
	Id                   int64            `json:"id,omitempty"`
	Title                string           `json:"title,omitempty"`
	Handle               string           `json:"handle,omitempty"`
	BodyHtml             string           `json:"body_html,omitempty"`
	Description          string           `json:"description,omitempty"`
	PublishedAt          time.Time        `json:"published_at,omitempty"`
	CreatedAt            time.Time        `json:"created_at,omitempty"`
	UpdatedAt            time.Time        `json:"updated_at,omitempty"`
	Vendor               string           `json:"vendor,omitempty"`
	ProductType          string           `json:"product_type,omitempty"`
	Type                 string           `json:"type,omitempty"`
	Tags                 []string         `json:"tags,omitempty"`
	Price                json.Number      `json:"price,omitempty"`
	PriceMin             json.Number      `json:"price_min,omitempty"`
	PriceMax             json.Number      `json:"price_max,omitempty"`
	Available            bool             `json:"available,omitempty"`
	PriceVaries          bool             `json:"price_varies,omitempty"`
	CompareAtPrice       interface{}      `json:"compare_at_price,omitempty"`
	CompareAtPriceMin    json.Number      `json:"compare_at_price_min,omitempty"`
	CompareAtPriceMax    json.Number      `json:"compare_at_price_max,omitempty"`
	CompareAtPriceVaries interface{}      `json:"compare_at_price_varies,omitempty"`
	Variants             []ProductVariant `json:"variants,omitempty"`
}

type ProductVariant struct {
	Id               int64       `json:"id,omitempty"`
	Title            string      `json:"title,omitempty"`
	Option1          string      `json:"option1,omitempty"`
	Option2          string      `json:"option2,omitempty"`
	Option3          string      `json:"option3,omitempty"`
	Sku              string      `json:"sku,omitempty"`
	RequiresShipping bool        `json:"requires_shipping,omitempty"`
	Taxable          bool        `json:"taxable,omitempty"`
	FeaturedImage    interface{} `json:"featured_image,omitempty"`
	Available        bool        `json:"available,omitempty"`
	Price            interface{} `json:"price,omitempty"`
	Grams            json.Number `json:"grams,omitempty"`
	CompareAtPrice   interface{} `json:"compare_at_price,omitempty"`
	Position         json.Number `json:"position,omitempty"`
	ProductId        json.Number `json:"product_id,omitempty"`
	CreatedAt        time.Time   `json:"created_at,omitempty"`
	UpdatedAt        time.Time   `json:"updated_at,omitempty"`
}

// GetProducts Fetches the most recent loaded products
func (t *Task) GetProducts() (*resty.Response, []Product, error) {
	resp, err := t.Client.R().
		SetHeader("Accept", "text/html,application/xhtml+xml,application/xml.q=0.9,image/avif,image/webp,image/apng,*/*.q=0.8,application/signed-exchange.v=b3.q=0.9").
		SetHeader("Accept-Encoding", "gzip, deflate").
		SetHeader("Accept-Language", "en-US,en.q=0.9").
		SetHeader("Cache-Control", "max-age=0").
		SetHeader("Connection", "keep-alive").
		SetHeader("Host", t.Site.GetHostName()).
		SetHeader("sec-ch-ua", "\" Not A.Brand\".v=\"99\", \"Chromium\".v=\"102\", \"Google Chrome\".v=\"102\"").
		SetHeader("sec-ch-ua-mobile", "?0").
		SetHeader("sec-ch-ua-platform", "\"Windows\"").
		SetHeader("Sec-Fetch-Dest", "document").
		SetHeader("Sec-Fetch-Mode", "navigate").
		SetHeader("Sec-Fetch-Site", "none").
		SetHeader("Sec-Fetch-User", "?1").
		SetHeader("Upgrade-Insecure-Requests", "1").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0. Win64. x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.62 Safari/537.36").
		Get(fmt.Sprintf("%v/products.json", t.Site.Url))

	if err != nil {
		return resp, nil, err
	}

	if resp.IsError() {
		return resp, nil, nil
	}

	type respStruct struct {
		Products []Product `json:"products"`
	}

	var data respStruct
	err = json.Unmarshal(resp.Body(), &data)

	if err != nil {
		return resp, nil, err
	}

	return resp, data.Products, err
}

// GetProductSpecific Fetch a specific product's information
func (t *Task) GetProductSpecific(url string) (*resty.Response, Product, error) {
	resp, err := t.Client.R().
		SetHeader("Accept", "text/html,application/xhtml+xml,application/xml.q=0.9,image/avif,image/webp,image/apng,*/*.q=0.8,application/signed-exchange.v=b3.q=0.9").
		SetHeader("Accept-Encoding", "gzip, deflate").
		SetHeader("Accept-Language", "en-US,en.q=0.9").
		SetHeader("Cache-Control", "max-age=0").
		SetHeader("Connection", "keep-alive").
		SetHeader("Host", t.Site.GetHostName()).
		SetHeader("sec-ch-ua", "\" Not A.Brand\".v=\"99\", \"Chromium\".v=\"102\", \"Google Chrome\".v=\"102\"").
		SetHeader("sec-ch-ua-mobile", "?0").
		SetHeader("sec-ch-ua-platform", "\"Windows\"").
		SetHeader("Sec-Fetch-Dest", "document").
		SetHeader("Sec-Fetch-Mode", "navigate").
		SetHeader("Sec-Fetch-Site", "none").
		SetHeader("Sec-Fetch-User", "?1").
		SetHeader("Upgrade-Insecure-Requests", "1").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0. Win64. x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.62 Safari/537.36").
		Get(fmt.Sprintf("%v.js", url))

	if err != nil {
		return resp, Product{}, err
	}

	if resp.IsError() {
		return resp, Product{}, nil
	}

	var product Product
	err = json.Unmarshal(resp.Body(), &product)

	if err != nil {
		return resp, Product{}, err
	}

	return resp, product, nil
}
