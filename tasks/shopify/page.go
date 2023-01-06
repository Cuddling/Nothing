package shopify

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Page struct {
	Url      string
	Html     string
	Document *goquery.Document
}

type CheckoutStep int

const (
	CheckoutStepNone CheckoutStep = iota
	CheckoutStepContact
	CheckoutStepShippingMethod
	CheckoutStepCalculatingTaxes
	CheckoutStepPaymentMethod
	CheckoutStepProcessing
	CheckoutStepQueue
	CheckoutStepOrderConfirmation
	CheckoutStepCheckpoint
	CheckoutStepLogin
	CheckoutStepLoginChallenge
)

// ErrNoShippingRateAvailable No shipping rates have not been calculated yet
var ErrNoShippingRateAvailable = errors.New("shipping rate not found")
var ErrNoTotalPriceFound = errors.New("total price not found")
var ErrNoGatewayFound = errors.New("gateway not found")
var ErrNoNoticeTextFound = errors.New("notice__text not found")

// NewPage Creates a new page struct instance and parses the associated HTML
func NewPage(url string, html string) Page {
	return Page{url, html, nil}
}

// GetCheckoutStep Returns which step of checkout we're currently on
func (p *Page) GetCheckoutStep() CheckoutStep {
	urlLower := strings.ToLower(p.Url)

	if strings.Contains(urlLower, "/throttle") || strings.Contains(urlLower, "/queue") {
		return CheckoutStepQueue
	}

	if strings.Contains(urlLower, "/thank_you") {
		return CheckoutStepOrderConfirmation
	}

	if strings.Contains(urlLower, "/checkpoint") {
		return CheckoutStepCheckpoint
	}

	if strings.Contains(urlLower, "/challenge") {
		return CheckoutStepLoginChallenge
	}

	if strings.Contains(urlLower, "/account/login") {
		return CheckoutStepLogin
	}

	if strings.Contains(urlLower, "/checkouts/") {
		if strings.Contains(urlLower, "/processing") {
			return CheckoutStepProcessing
		}

		htmlLower := strings.ToLower(p.Html)

		if strings.Contains(htmlLower, "calculating taxes") {
			return CheckoutStepCalculatingTaxes
		}

		u, err := url.Parse(p.Url)

		if err != nil {
			return CheckoutStepNone
		}

		pageStep := p.GetCheckoutStepFromString(p.GetCheckoutStepStringFromPage(), "")

		if pageStep != CheckoutStepNone {
			return pageStep
		}

		urlStep := p.GetCheckoutStepFromString(u.Query().Get("step"), u.Query().Get("previous_step"))

		if urlStep != CheckoutStepNone {
			return urlStep
		}
	}

	return CheckoutStepNone
}

// GetCheckoutStepStringFromPage Retrieves the checkout step from the page
func (p *Page) GetCheckoutStepStringFromPage() string {
	rgx := regexp.MustCompile("Shopify\\.Checkout\\.step = \\\"(.+)\\\";")
	rs := rgx.FindStringSubmatch(p.Html)

	if len(rs) <= 1 {
		return ""
	}

	return rs[1]
}

func (p *Page) GetCheckoutStepFromString(step string, previous string) CheckoutStep {
	if previous == "payment_method" && step == "" || strings.Contains(p.Url, "validate=") {
		return CheckoutStepPaymentMethod
	}

	if step == "" || step == "contact_information" {
		return CheckoutStepContact
	}

	switch step {
	case "shipping_method":
		return CheckoutStepShippingMethod
	case "payment_method":
		return CheckoutStepPaymentMethod
	case "":
		return CheckoutStepPaymentMethod
	default:
		return CheckoutStepNone
	}
}

// GetAuthenticityToken Retrieves the authenticity_token from the page if one exists
// On error, it returns an empty space as a token. For some reason this makes the request
// go through.
func (p *Page) GetAuthenticityToken() (string, error) {
	rgx := regexp.MustCompile("name=\"authenticity_token\" value=\"(.+)\"")
	rs := rgx.FindStringSubmatch(p.Html)

	if len(rs) <= 1 {
		return "", ErrNoTotalPriceFound
	}

	return rs[1], nil
}

// GetShippingRate Retrieves the shipping rate on the page if one exists
func (p *Page) GetShippingRate() (string, error) {
	rgx := regexp.MustCompile("type=\"radio\" value=\"(.+)\" name=\"checkout\\[shipping_rate\\]\\[id\\]\"")
	rs := rgx.FindStringSubmatch(p.Html)

	if len(rs) <= 1 {
		return "", ErrNoShippingRateAvailable
	}

	return rs[1], nil
}

// GetTotalPrice Retrieves the total price on the page if one exists
func (p *Page) GetTotalPrice() (string, error) {
	rgx := regexp.MustCompile("data-checkout-payment-due-target=\\\"(\\d+)\\\"")
	rs := rgx.FindStringSubmatch(p.Html)

	if len(rs) <= 1 {
		return "", ErrNoTotalPriceFound
	}

	return rs[1], nil
}

// GetPaymentGateway Retrieves the payment gateway from the page if one exists. Otherwise, it returns 0 and an error.
func (p *Page) GetPaymentGateway() (string, error) {
	gateway := GetGateway(p.GetShopId())

	if gateway != 0 {
		return strconv.Itoa(gateway), nil
	}

	rgx := regexp.MustCompile("payment_gateway_(\\d+)")
	rs := rgx.FindStringSubmatch(p.Html)

	if len(rs) <= 1 {
		return "0", ErrNoGatewayFound
	}

	return rs[1], nil
}

// GetNoticeText Retrieves the notice text from the page if one exists
func (p *Page) GetNoticeText() (string, error) {
	rgx := regexp.MustCompile("<p class=\\\"notice__text\\\">(.+)<\\/p>")
	rs := rgx.FindStringSubmatch(p.Html)

	if len(rs) <= 1 {
		return "", ErrNoNoticeTextFound
	}

	return rs[1], nil
}

// GetProductTitle Retrieves the product title from the page if one exists
func (p *Page) GetProductTitle() string {
	rgx := regexp.MustCompile("<span class=\\\"product__description__name order-summary__emphasis\\\">(.+)<\\/span>")
	rs := rgx.FindStringSubmatch(p.Html)

	if len(rs) <= 1 {
		return "None"
	}

	return rs[1]
}

// GetProductSize Retrieves the product size from the page if one exists
func (p *Page) GetProductSize() string {
	rgx := regexp.MustCompile("<span class=\\\"product__description__variant order-summary__small-text\\\">(.+)<\\/span>")
	rs := rgx.FindStringSubmatch(p.Html)

	if len(rs) <= 1 {
		return "None"
	}

	return rs[1]
}

// GetProductImage Retrieves the product image from the page if one exists
func (p *Page) GetProductImage() string {
	rgx := regexp.MustCompile("class=\\\"product-thumbnail__image\\\" src=\\\"(.+)\\\"")
	rs := rgx.FindStringSubmatch(p.Html)

	if len(rs) <= 1 {
		return ""
	}

	img := rs[1]

	if !strings.HasPrefix(img, "http") {
		img = "https:" + img
	}

	return img
}

// GetShopId Gets the shop id from the URL
func (p *Page) GetShopId() int {
	rgx := regexp.MustCompile(".+\\/(\\d+)\\/checkouts\\/(.+)")
	rs := rgx.FindStringSubmatch(p.Url)

	if len(rs) <= 1 {
		return 0
	}

	id, err := strconv.Atoi(rs[1])

	if err != nil {
		fmt.Println(err)
		return 0
	}

	return id
}

// GetOrderNumber Returns the order number from the order confirmation page
func (p *Page) GetOrderNumber() string {
	var err error
	p.Document, err = goquery.NewDocumentFromReader(strings.NewReader(p.Html))

	if err != nil {
		return ""
	}

	orderNumber := ""

	p.Document.Find(".os-order-number").Each(func(i int, s *goquery.Selection) {
		orderNumber = strings.TrimSpace(s.Text())
		orderNumber = strings.Replace(orderNumber, "Order ", "", -1)
	})

	return orderNumber
}

// GetPageFromResponse Returns a parsed page struct instance from a response
func GetPageFromResponse(r *resty.Response) (Page, error) {
	redirect, err := GetResponseUrl(r)

	if err != nil {
		return Page{}, err
	}

	p := NewPage(redirect, string(r.Body()))
	return p, nil
}

// GetResponseUrl Gets a response's current page url
func GetResponseUrl(r *resty.Response) (string, error) {
	return r.RawResponse.Request.URL.String(), nil
}
