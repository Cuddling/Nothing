package profiles

type Profile struct {
	Name                         string
	ShippingAddress              Address
	BillingAddress               Address
	CreditCard                   Card
	SameBillingAddressAsShipping bool
}

// GetBillingAddress Returns the address to be used for billing information
func (p *Profile) GetBillingAddress() Address {
	if p.SameBillingAddressAsShipping {
		return p.ShippingAddress
	}

	return p.BillingAddress
}
