package shopify

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

// CalculateTaxes Polls the calculating taxes page. Returns if taxes were successfully calculated
// or we were moved to a different page.
func (t *Task) CalculateTaxes() (bool, *resty.Response, error) {
	page, resp, err := t.GetCheckoutPage(fmt.Sprintf("%v?previous_step=shipping_method&step=payment_method", t.CheckoutUrl))

	if err != nil {
		return false, nil, err
	}

	if resp.IsError() {
		return false, resp, errors.New(fmt.Sprintf("calculating taxes failed (%v)", resp.StatusCode()))
	}
	
	return page.GetCheckoutStep() != CheckoutStepCalculatingTaxes, resp, nil
}
