package shopify

import (
	"Mystery/utils"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"regexp"
	"strconv"
)

var ErrProductNotFound = errors.New("product not found")

// MonitorProducts Monitors for a product to check out with all different monitoring input types
func (t *Task) MonitorProducts() (*resty.Response, error) {
	var keywords []string

	for _, input := range t.MonitorInputs {
		switch utils.GetMonitorInputTypeFromString(input) {
		case utils.MonitorInputTypeVariant:
			i, err := strconv.ParseInt(input, 10, 64)

			if err != nil {
				return nil, err
			}

			t.Variant = ProductVariant{Id: i}
			return nil, nil
		case utils.MonitorInputTypeUrl:
			resp, product, variant, err := t.MonitorProductUrl(input)

			if err != nil {
				return resp, err
			}

			t.Product = product
			t.Variant = variant
			return resp, nil
		case utils.MonitorInputTypeKeywords:
			keywords = append(keywords, input)
		}
	}

	resp, product, variant, err := t.MonitorProductKeywords(keywords)

	if err != nil {
		return resp, err
	}

	t.Product = product
	t.Variant = variant
	return resp, nil
}

// MonitorProductUrl Monitors for a product by URL
func (t *Task) MonitorProductUrl(url string) (*resty.Response, Product, ProductVariant, error) {
	resp, product, err := t.GetProductSpecific(url)

	if err != nil {
		return resp, Product{}, ProductVariant{}, err
	}

	if resp.IsError() {
		return resp, Product{}, ProductVariant{}, errors.New(fmt.Sprintf("request failed with status code: %v", resp.StatusCode()))
	}

	variant, err := getCheckoutVariant(product.Variants, t.Sizes)

	if err != nil {
		return resp, Product{}, ProductVariant{}, err
	}

	return resp, product, variant, nil
}

// MonitorProductKeywords Monitors for a product by keywords
func (t *Task) MonitorProductKeywords(keywordSets []string) (*resty.Response, Product, ProductVariant, error) {
	resp, products, err := t.GetProducts()

	if err != nil {
		return resp, Product{}, ProductVariant{}, err
	}

	if resp.IsError() {
		return resp, Product{}, ProductVariant{}, errors.New(fmt.Sprintf("request failed with status code: %v", resp.StatusCode()))
	}

	checkKeywords := func(usingHandle bool) (Product, ProductVariant, error) {
		for _, product := range products {
			foundProduct := false

			// Check if any of the keyword sets match the products title or handle
			for _, keywordSet := range keywordSets {
				var value string

				if usingHandle {
					value = product.Handle
				} else {
					value = product.Title
				}

				if utils.IsKeywordMatch(value, keywordSet) {
					foundProduct = true
					break
				}
			}

			if foundProduct {
				variant, err := getCheckoutVariant(product.Variants, t.Sizes)

				if err != nil {
					return Product{}, ProductVariant{}, err
				}

				return product, variant, nil
			}
		}

		return Product{}, ProductVariant{}, ErrProductNotFound
	}

	// Check all product titles first. This has precedence over handles
	product, variant, err := checkKeywords(false)

	if err == nil {
		return resp, product, variant, nil
	}

	// Now check the handles since no product was found
	product, variant, err = checkKeywords(true)

	if err == nil {
		return resp, product, variant, nil
	}

	// No product was found at all
	return nil, Product{}, ProductVariant{}, err
}

// Retrieves the variant that is in the size range.
// Available / in-stock variants take precedence. If there are no in-stock items, it selects a random one within the range.
func getCheckoutVariant(variants []ProductVariant, sizeRange []string) (ProductVariant, error) {
	var filtered []ProductVariant

	// No size range specified.
	if len(sizeRange) == 0 {
		available := getAvailableVariants(variants)

		if len(available) > 0 {
			filtered = available
		} else {
			filtered = variants
		}
		// Size range specified. Find all the variants that are within the size range.
	} else {
		var inSizeRange []ProductVariant

		for _, variant := range variants {
			for _, size := range sizeRange {
				// Shoe sizes are determined if the variant option contains ANY number
				if regexp.MustCompile(`\d`).MatchString(variant.Option1) {
					// Isolating the numbers from the size, so we can check against that
					regex := regexp.MustCompile("[^0-9.]+")
					rawSize := regex.ReplaceAllString(variant.Option1, "")

					if rawSize == size {
						inSizeRange = append(inSizeRange, variant)
					}
				} else {
					// Otherwise, we can safely assume it's clothing and match the exact size.
					// Not using contains because "L" can pick up XL, XXL, etc.
					if variant.Option1 == size {
						inSizeRange = append(inSizeRange, variant)
					}
				}
			}
		}

		available := getAvailableVariants(inSizeRange)

		if len(available) > 0 {
			filtered = available
		} else {
			filtered = inSizeRange
		}
	}

	if len(filtered) == 0 {
		return ProductVariant{}, errors.New("no variant found")
	}

	// Select a variant at random
	n := rand.Int() % len(filtered)
	return filtered[n], nil
}

// Gets only the available variants from a slice of them
func getAvailableVariants(variants []ProductVariant) []ProductVariant {
	var available []ProductVariant

	for _, v := range variants {
		if v.Available {
			available = append(available, v)
		}
	}

	return available
}
