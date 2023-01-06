package shopify

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

// ProcessOrder Fetches the process order page
func (t *Task) ProcessOrder() (*resty.Response, error) {
	_, resp, err := t.GetCheckoutPage(fmt.Sprintf("%v/processing?from_processing_page=1", t.CheckoutUrl))

	if err != nil {
		return resp, err
	}

	if resp.IsError() {
		return resp, errors.New(fmt.Sprintf("processing order failed (%v)", resp.StatusCode()))
	}

	return resp, nil
}
