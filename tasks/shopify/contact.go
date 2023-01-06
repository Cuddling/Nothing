package shopify

import (
	"Mystery/utils"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strings"
)

func (t *Task) SubmitContactInfo() (Page, *resty.Response, error) {
	token, _ := t.CurrentPage.GetAuthenticityToken()
	body := utils.FormBody{}

	body.Add("_method", "patch")
	body.Add("authenticity_token", token)
	body.Add("previous_step", "contact_information")
	body.Add("step", "shipping_method")

	body.Add("checkout[email]", t.Profile.ShippingAddress.Email)
	body.Add("checkout[buyer_accepts_marketing]", "0")
	body.Add("checkout[buyer_accepts_marketing]", "1")

	// Add the same information twice because of honeypot / anti-bot
	for i := 0; i < 2; i++ {
		body.Add("checkout[shipping_address][first_name]", t.Profile.ShippingAddress.GetFirstName())
		body.Add("checkout[shipping_address][last_name]", t.Profile.ShippingAddress.GetLastName())
		body.Add("checkout[shipping_address][address1]", t.Profile.ShippingAddress.Line1)
		body.Add("checkout[shipping_address][address2]", t.Profile.ShippingAddress.Line2)
		body.Add("checkout[shipping_address][city]", t.Profile.ShippingAddress.City)
		body.Add("checkout[shipping_address][country]", t.Profile.ShippingAddress.Country)
		body.Add("checkout[shipping_address][province]", t.Profile.ShippingAddress.State)
		body.Add("checkout[shipping_address][zip]", t.Profile.ShippingAddress.PostCode)
		body.Add("checkout[shipping_address][phone]", t.Profile.ShippingAddress.Phone)
	}

	// Kicks Lounge
	if strings.Contains(t.CurrentPage.Html, "checkout[buyer_accepts_sms]") {
		body.Add("checkout[buyer_accepts_sms]", "0")
	}

	// Kicks Lounge
	if strings.Contains(t.CurrentPage.Html, "checkout[sms_marketing_phone]") {
		body.Add("checkout[sms_marketing_phone]", "")
	}

	// Pickup
	if strings.Contains(t.CurrentPage.Html, "checkout[pick_up_in_store][selected]") {
		body.Add("checkout[pick_up_in_store][selected]", "false")
	}

	if strings.Contains(t.CurrentPage.Html, "checkout[id]") {
		body.Add("checkout[id]", "delivery-shipping")
	}

	// Slam Jam
	if strings.Contains(t.CurrentPage.Html, "checkout[buyer_accepts_privacy_policy]") {
		body.Add("checkout[buyer_accepts_privacy_policy]", "0")
		body.Add("checkout[buyer_accepts_privacy_policy]", "on")
	}

	// Slam Jam
	if strings.Contains(t.CurrentPage.Html, "checkout[buyer_accepts_sms]") {
		body.Add("checkout[buyer_accepts_sms]", "0")
	}

	// Slam Jam
	if strings.Contains(t.CurrentPage.Html, "checkout[sms_marketing_phone]") {
		body.Add("checkout[sms_marketing_phone]", "")
	}

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
		SetHeader("Connection", "Keep-Alive").
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Host", t.Site.GetHostName()).
		SetHeader("Origin", t.Site.Url).
		SetHeader("Referer", t.CheckoutUrl).
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
		return Page{}, nil, errors.New(fmt.Sprintf("submit contact info failed (%v)", resp.StatusCode()))
	}

	page, err := GetPageFromResponse(resp)

	if err != nil {
		return Page{}, nil, err
	}

	t.UpdateCurrentPage(page)
	return page, resp, nil
}
