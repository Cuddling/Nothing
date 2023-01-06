package utils

import (
	"net/url"
	"strings"
)

type MonitorInputType int

const (
	MonitorInputTypeVariant = iota
	MonitorInputTypeUrl
	MonitorInputTypeKeywords
)

// GetMonitorInputTypeFromString Returns the type of monitor input from the given string
func GetMonitorInputTypeFromString(s string) MonitorInputType {
	if strings.HasPrefix(s, "+") || strings.HasPrefix(s, "-") {
		return MonitorInputTypeKeywords
	}

	_, err := url.ParseRequestURI(s)

	if err == nil {
		return MonitorInputTypeUrl
	}

	return MonitorInputTypeVariant
}
