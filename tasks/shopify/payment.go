package shopify

import (
	"Mystery/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"math"
	"strconv"
	"strings"
)

// LoadPaymentPage Retrieves the payment page
func (t *Task) LoadPaymentPage() (*resty.Response, error) {
	_, resp, err := t.GetCheckoutPage(fmt.Sprintf("%v?step=payment_method", t.CheckoutUrl))

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return resp, errors.New(fmt.Sprintf("loading payment page failed (%v)", resp.StatusCode()))
	}

	return resp, nil
}

// GetPaymentToken Retrieves the payment session id from deposits
func (t *Task) GetPaymentToken() (string, *resty.Response, error) {
	resp, err := t.Client.R().
		SetHeader("Accept", "application/json").
		SetHeader("Accept-Encoding", "gzip, deflate").
		SetHeader("Accept-Language", "en-US,en.q=0.9").
		SetHeader("Connection", "keep-alive").
		SetHeader("Content-Type", "application/json").
		SetHeader("Host", "deposit.us.shopifycs.com").
		SetHeader("Origin", "https://checkout.shopifycs.com").
		SetHeader("Referer", "https://checkout.shopifycs.com/").
		SetHeader("Sec-Fetch-Dest", "empty").
		SetHeader("Sec-Fetch-Mode", "cors").
		SetHeader("Sec-Fetch-Site", "same-site").
		SetHeader("sec-ch-ua", "\" Not A.Brand\".v=\"99\", \"Chromium\".v=\"102\", \"Google Chrome\".v=\"102\"").
		SetHeader("sec-ch-ua-mobile", "?0").
		SetHeader("sec-ch-ua-platform", "\"Windows\"").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0. Win64. x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.62 Safari/537.36").
		SetBody(map[string]interface{}{
			"credit_card": map[string]interface{}{
				"number":             t.Profile.CreditCard.Number,
				"name":               t.Profile.GetBillingAddress().Name,
				"month":              t.Profile.CreditCard.ExpiryMonth,
				"year":               t.Profile.CreditCard.ExpiryYear,
				"verification_value": t.Profile.CreditCard.CVV,
			},
			"payment_session_scope": t.Site.GetHostName(),
		}).
		Post("https://deposit.us.shopifycs.com/sessions")

	if err != nil {
		return "", resp, err
	}

	if resp.IsError() {
		return "", resp, errors.New(fmt.Sprintf("failed to get payment token (%v)", resp.StatusCode()))
	}

	type response struct {
		Id string `json:"id"`
	}

	var data response
	err = json.Unmarshal(resp.Body(), &data)

	if err != nil {
		return "", resp, err
	}

	return data.Id, resp, nil
}

// SubmitPayment Submits the payment
func (t *Task) SubmitPayment(fast bool, sessionId string, rate string) (Page, *resty.Response, error) {
	token, _ := t.CurrentPage.GetAuthenticityToken()
	gateway, _ := t.CurrentPage.GetPaymentGateway()
	price := t.GetTotalPaymentPrice()

	body := utils.FormBody{}
	body.Add("_method", "patch")
	body.Add("authenticity_token", token)
	body.Add("previous_step", "payment_method")
	body.Add("step", "")
	body.Add("s", sessionId)

	if fast && rate != "" {
		body.Add("checkout[shipping_rate][id]", rate)
	}

	body.Add("checkout[payment_gateway]", gateway)
	body.Add("checkout[credit_card][vault]", "false")
	body.Add("checkout[different_billing_address]", strconv.FormatBool(!t.Profile.SameBillingAddressAsShipping))

	if !t.Profile.SameBillingAddressAsShipping {
		// Honey pot
		body.Add("checkout[billing_address][first_name]", "")
		body.Add("checkout[billing_address][last_name]", "")
		body.Add("checkout[billing_address][company]", "")
		body.Add("checkout[billing_address][address1]", "")
		body.Add("checkout[billing_address][address2]", "")
		body.Add("checkout[billing_address][city]", "")
		body.Add("checkout[billing_address][country]", "")
		body.Add("checkout[billing_address][province]", "")
		body.Add("checkout[billing_address][zip]", "")
		body.Add("checkout[billing_address][phone]", "")

		address := t.Profile.GetBillingAddress()
		body.Add("checkout[billing_address][country]", address.Country)
		body.Add("checkout[billing_address][first_name]", address.GetFirstName())
		body.Add("checkout[billing_address][last_name]", address.GetLastName())
		body.Add("checkout[billing_address][company]", "")
		body.Add("checkout[billing_address][address1]", address.Line1)
		body.Add("checkout[billing_address][address2]", address.Line2)
		body.Add("checkout[billing_address][city]", address.City)
		body.Add("checkout[billing_address][province]", address.State)
		body.Add("checkout[billing_address][zip]", address.PostCode)
		body.Add("checkout[billing_address][phone]", address.Phone)
	}

	body.Add("checkout[remember_me]", "false")
	body.Add("checkout[remember_me]", "0")
	body.Add("checkout[vault_phone]", fmt.Sprintf("+1%v", t.Profile.ShippingAddress.Phone))
	body.Add("checkout[total_price]", price)
	body.Add("complete", "1")

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
		SetHeader("Referer", fmt.Sprintf("%v?previous_step=shipping_method&step=payment_method", t.CheckoutUrl)).
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
		return Page{}, resp, errors.New(fmt.Sprintf("submit payment failed (%v)", resp.StatusCode()))
	}

	page, err := GetPageFromResponse(resp)

	if err != nil {
		return Page{}, nil, err
	}

	t.UpdateCurrentPage(page)
	return page, resp, nil
}

// GetTotalPaymentPrice Attempts to retrieve the total payment price on the page.
func (t *Task) GetTotalPaymentPrice() string {
	price, err := t.CurrentPage.GetTotalPrice()

	if err != nil {
		t.Log(err)
		return "0"
	}

	// Given that we're on the payment step, it is a 100% confirmation that the total price is correct.
	if t.CurrentPage.GetCheckoutStep() == CheckoutStepPaymentMethod {
		return price
	}

	// Below attempts to calculate the total price from the shipping rate. This helps speed up the checkout process,
	// so we can skip the submitting shipping rate step and do it in the payment step.

	ratePrice, err := strconv.ParseFloat(GetShippingRatePrice(t.ShippingRate), 32)

	if err != nil {
		return price
	}

	ratePrice = math.Round(ratePrice*100) / 100
	priceInt, err := strconv.Atoi(price)

	if err != nil {
		return price
	}

	// Slam Jam has a shipping rate of 28.82 that gets rounded up to $30 by Shopify. This fixes that error so payment can still
	// be submitted fast
	if strings.Contains(strings.ToLower(t.Site.Url), "slamjam") &&
		(strings.Contains(t.ShippingRate, "28") || strings.Contains(t.ShippingRate, "29")) {
		ratePrice = 30
	}

	priceFinal := float64(priceInt) + ratePrice*100
	return fmt.Sprintf("%v", int(priceFinal))
}
