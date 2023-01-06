package shopify

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strings"
)

// CreateCheckout Starts a checkout session and returns the link to it
func (t *Task) CreateCheckout() (string, *resty.Response, error) {
	// Disabling redirects here because the first redirect is always the checkout URL.
	t.Client.SetRedirectPolicy(resty.NoRedirectPolicy())

	resp, err := t.Client.R().
		SetHeader("Accept", "text/html,application/xhtml+xml,application/xml.q=0.9,image/avif,image/webp,image/apng,*/*.q=0.8,application/signed-exchange.v=b3.q=0.9").
		SetHeader("Accept-Encoding", "gzip, deflate").
		SetHeader("Accept-Language", "en-US,en.q=0.9").
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
		Get(fmt.Sprintf("%v/checkout", t.Site.Url))

	t.Client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(1000))

	if err != nil && !strings.Contains(err.Error(), "auto redirect is disabled") {
		return "", resp, err
	}

	if resp.IsError() {
		return "", resp, errors.New(fmt.Sprintf("create checkout failed (%v)", resp.StatusCode()))
	}

	url, err := resp.RawResponse.Location()

	if err != nil {
		return "", resp, err
	}

	if strings.Contains(url.String(), "/checkouts/") {
		t.CheckoutUrl = fmt.Sprintf("https://%v%v", url.Host, url.Path)
	} else {
		t.CurrentPage = Page{Url: url.String()}
	}

	return t.CheckoutUrl, resp, nil
}

// CreateCheckoutFast Adds to cart + creates checkout fast
func (t *Task) CreateCheckoutFast(variant int64) (string, *resty.Response, error) {
	t.CurrentPage = Page{Url: fmt.Sprintf("%v/cart/%v:1", t.Site.Url, variant)}
	_, resp, err := t.GetCheckoutPage(t.CurrentPage.Url)

	if err != nil {
		return "", nil, err
	}

	redirectUrl := resp.RawResponse.Request.URL

	if strings.Contains(redirectUrl.String(), "/checkouts/") {
		t.CheckoutUrl = fmt.Sprintf("https://%v%v", redirectUrl.Host, redirectUrl.Path)
	}

	// If the variant is invalid, Shopify will redirect to an empty cart.
	if !strings.Contains(redirectUrl.String(), "/cart") {
		t.VariantInCart = variant
	}

	return t.CheckoutUrl, resp, err
}

// GetCheckoutPage Fetches a given URL's checkout page
func (t *Task) GetCheckoutPage(url string) (Page, *resty.Response, error) {
	resp, err := t.Client.R().
		SetHeader("Accept", "text/html,application/xhtml+xml,application/xml.q=0.9,image/avif,image/webp,image/apng,*/*.q=0.8,application/signed-exchange.v=b3.q=0.9").
		SetHeader("Accept-Encoding", "gzip, deflate").
		SetHeader("Accept-Language", "en-US,en.q=0.9").
		SetHeader("Connection", "keep-alive").
		SetHeader("Host", t.Site.GetHostName()).
		SetHeader("sec-ch-ua", "\"Chromium\".v=\"104\", \" Not A.Brand\".v=\"99\", \"Google Chrome\".v=\"104\"").
		SetHeader("sec-ch-ua-mobile", "?0").
		SetHeader("sec-ch-ua-platform", "\"Windows\"").
		SetHeader("Sec-Fetch-Dest", "document").
		SetHeader("Sec-Fetch-Mode", "navigate").
		SetHeader("Sec-Fetch-Site", "none").
		SetHeader("Sec-Fetch-User", "?1").
		SetHeader("Upgrade-Insecure-Requests", "1").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0. Win64. x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.62 Safari/537.36").
		Get(url)

	if err != nil {
		return Page{}, resp, err
	}

	if resp.IsError() {
		return Page{}, resp, errors.New(fmt.Sprintf("get checkout page failed (%v)", resp.StatusCode()))
	}

	page, err := GetPageFromResponse(resp)

	if err != nil {
		return Page{}, nil, err
	}

	t.UpdateCurrentPage(page)
	return page, resp, nil
}
