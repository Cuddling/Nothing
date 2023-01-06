package shopify

import (
	"Mystery/utils"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"regexp"
)

// FetchShippingRate Attempts to fetch the shipping rate
func (t *Task) FetchShippingRate() (string, *resty.Response, error) {
	resp, err := t.Client.R().
		SetHeader("Accept", "*/*").
		SetHeader("Accept-Encoding", "gzip, deflate").
		SetHeader("Accept-Language", "en-US,en.q=0.9").
		SetHeader("Connection", "keep-alive").
		SetHeader("Host", t.Site.GetHostName()).
		SetHeader("Referer", t.Site.Url).
		SetHeader("sec-ch-ua", "\" Not A.Brand\".v=\"99\", \"Chromium\".v=\"102\", \"Google Chrome\".v=\"102\"").
		SetHeader("sec-ch-ua-mobile", "?0").
		SetHeader("sec-ch-ua-platform", "\"Windows\"").
		SetHeader("Sec-Fetch-Dest", "empty").
		SetHeader("Sec-Fetch-Mode", "cors").
		SetHeader("Sec-Fetch-Site", "same-origin").
		SetHeader("X-Requested-With", "XMLHttpRequest").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0. Win64. x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.62 Safari/537.36").
		Get(fmt.Sprintf("%v/shipping_rates?step=shipping_method", t.CheckoutUrl))

	if err != nil {
		return "", resp, err
	}

	if resp.IsError() {
		return "", nil, errors.New(fmt.Sprintf("polling shipping rate failed (%v)", resp.StatusCode()))
	}

	page := NewPage(resp.RawResponse.Request.URL.String(), string(resp.Body()))

	rate, err := page.GetShippingRate()

	if err != nil {
		return "", resp, err
	}

	return rate, nil, nil
}

// SubmitShippingRate Submits the shipping rate token to the server
func (t *Task) SubmitShippingRate(rate string) (Page, *resty.Response, error) {
	token, _ := t.CurrentPage.GetAuthenticityToken()
	body := utils.FormBody{}

	body.Add("_method", "patch")
	body.Add("authenticity_token", token)
	body.Add("previous_step", "shipping_method")
	body.Add("step", "payment_method")

	body.Add("checkout[shipping_rate][id]", rate)
	body.Add("checkout[client_details][browser_width]", "1903")
	body.Add("checkout[client_details][browser_height]", "979")
	body.Add("checkout[client_details][javascript_enabled]", "1")
	body.Add("checkout[client_details][color_depth]", "24")
	body.Add("checkout[client_details][java_enabled]", "false")
	body.Add("checkout[client_details][browser_tz]", "240")

	resp, err := t.Client.R().
		SetHeader("Accept", "text/html,application/xhtml+xml,application/xml.q=0.9,image/avif,image/webp,image/apng,*/*.q=0.8,application/signed-exchange.v=b3.q=0.9").
		SetHeader("Accept-Encoding", "gzip, deflate").
		SetHeader("Accept-Language", "en-US,en.q=0.9").
		SetHeader("Cache-Control", "max-age=0").
		SetHeader("Connection", "keep-alive").
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Host", t.Site.GetHostName()).
		SetHeader("Origin", t.Site.Url).
		SetHeader("Referer", fmt.Sprintf("%v?previous_step=contact_information&step=shipping_method", t.CheckoutUrl)).
		SetHeader("Sec-Fetch-Dest", "document").
		SetHeader("Sec-Fetch-Mode", "navigate").
		SetHeader("Sec-Fetch-Site", "same-origin").
		SetHeader("Sec-Fetch-User", "?1").
		SetHeader("Upgrade-Insecure-Requests", "1").
		SetHeader("sec-ch-ua", "\" Not A.Brand\".v=\"99\", \"Chromium\".v=\"102\", \"Google Chrome\".v=\"102\"").
		SetHeader("sec-ch-ua-mobile", "?0").
		SetHeader("sec-ch-ua-platform", "\"Windows\"").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0. Win64. x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.62 Safari/537.36").
		SetBody(body.ToString()).
		Post(t.CheckoutUrl)

	if err != nil {
		return Page{}, resp, err
	}

	if resp.IsError() {
		return Page{}, resp, errors.New(fmt.Sprintf("error submitting shipping rate (%v)", resp.StatusCode()))
	}

	page, err := GetPageFromResponse(resp)

	if err != nil {
		return Page{}, nil, err
	}

	t.UpdateCurrentPage(page)
	return page, resp, nil
}

// GetShippingRatePrice Returns the price of a shipping rate
// Example: Advanced Shipping Rules-1-12.05
// CARRIER-SOMETHING-PRICE
func GetShippingRatePrice(rate string) string {
	rgx := regexp.MustCompile(".+-.+-(.+)")
	rs := rgx.FindStringSubmatch(rate)

	if len(rs) <= 1 {
		return "0"
	}

	return rs[1]
}
