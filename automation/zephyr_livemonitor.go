package automation

import (
	"encoding/json"
	"fmt"
	"time"
)

type ZephyrMonitorLive struct {
	Type string                `json:"type,omitempty"`
	Body ZephyrMonitorLiveBody `json:"body,omitempty"`
}

type ZephyrMonitorLiveBody struct {
	Payload ZephyrMonitorLivePayload `json:"payload,omitempty"`
	Channel string                   `json:"channel,omitempty"`
	Type    string                   `json:"type,omitempty"`
	Event   string                   `json:"event,omitempty"`
}

type ZephyrMonitorLivePayload struct {
	Product ZephyrMonitorLiveProduct `json:"product,omitempty"`
	Store   string                   `json:"store,omitempty"`
}

type ZephyrMonitorLiveProduct struct {
	Images        []interface{}                     `json:"images,omitempty"`
	Available     bool                              `json:"available,omitempty"`
	CreatedAt     time.Time                         `json:"created_at,omitempty"`
	Handle        string                            `json:"handle,omitempty"`
	Variants      []ZephyrMonitorLiveProductVariant `json:"variants,omitempty"`
	Title         string                            `json:"title,omitempty"`
	Uuid          string                            `json:"uuid,omitempty"`
	ProductType   string                            `json:"product_type,omitempty"`
	UpdatedAt     time.Time                         `json:"updated_at,omitempty"`
	Filtered      bool                              `json:"filtered,omitempty"`
	Vendor        string                            `json:"vendor,omitempty"`
	ProductId     int64                             `json:"product_id,omitempty"`
	Options       []interface{}                     `json:"options,omitempty"`
	Id            int64                             `json:"id,omitempty"`
	PublishedAt   time.Time                         `json:"published_at,omitempty"`
	Timestamp     int64                             `json:"timestamp,omitempty"`
	FeaturedImage string                            `json:"featured_image,omitempty""`
}

type ZephyrMonitorLiveProductImage struct {
	UpdatedAt  time.Time     `json:"updated_at,omitempty"`
	Src        string        `json:"src,omitempty"`
	ProductId  int64         `json:"product_id,omitempty"`
	Width      int           `json:"width,omitempty"`
	CreatedAt  time.Time     `json:"created_at,omitempty"`
	VariantIds []interface{} `json:"variant_ids,omitempty"`
	Id         int64         `json:"id,omitempty"`
	Position   int           `json:"position,omitempty"`
	Height     int           `json:"height,omitempty"`
}

type ZephyrMonitorLiveProductVariant struct {
	FormattedPrice   string                                  `json:"formatted_price,omitempty"`
	CompareAtPrice   json.Number                             `json:"compare_at_price,omitempty"`
	Taxable          bool                                    `json:"taxable,omitempty"`
	RequiresShipping bool                                    `json:"requires_shipping,omitempty"`
	OptionValues     []ZephyrMonitorLiveProductVariantOption `json:"option_values,omitempty"`
	Available        bool                                    `json:"available,omitempty"`
	CreatedAt        time.Time                               `json:"created_at,omitempty"`
	Title            string                                  `json:"title,omitempty"`
	UpdatedAt        time.Time                               `json:"updated_at,omitempty"`
	Price            json.Number                             `json:"price,omitempty"`
	Id               int64                                   `json:"id,omitempty"`
	Position         int                                     `json:"position,omitempty"`
	Grams            int                                     `json:"grams,omitempty"`
	Sku              string                                  `json:"sku,omitempty"`
	Barcode          string                                  `json:"barcode,omitempty"`
	Name             string                                  `json:"name,omitempty"`
}

type ZephyrMonitorLiveProductVariantOption struct {
	Name     string `json:"name,omitempty"`
	OptionId int64  `json:"option_id,omitempty"`
	Value    string `json:"value,omitempty"`
}

func (p *ZephyrMonitorLiveProduct) GetImage() string {
	if p.FeaturedImage != "" {
		return fmt.Sprintf("https:%v", p.FeaturedImage)
	}

	if len(p.Images) == 0 {
		return ""
	}

	switch i := p.Images[0].(type) {
	case string:
		return fmt.Sprintf("https:%v", i)
	case map[string]interface{}:
		return fmt.Sprintf("%v", i["src"])
	}

	return ""
}
