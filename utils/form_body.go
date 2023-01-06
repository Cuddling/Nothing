package utils

import (
	"fmt"
	"net/url"
	"strings"
)

type FormBody struct {
	Properties []url.Values
}

func (b *FormBody) Add(key string, value string) {
	params := url.Values{}
	params.Add(key, value)

	b.Properties = append(b.Properties, params)
}

func (b *FormBody) ToString() string {
	str := ""

	for _, prop := range b.Properties {
		str += fmt.Sprintf("%v&", prop.Encode())
	}

	// Remove last ampersand
	return strings.TrimSuffix(str, "&")
}
