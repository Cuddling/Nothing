package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var randomSeeded = false

// GenerateRandomString Generates a random string of a given length
func GenerateRandomString(length int) string {
	if !randomSeeded {
		rand.Seed(time.Now().UnixNano())
		randomSeeded = true
	}

	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

// IsKeywordMatch Returns if a string matches a given keyword set
func IsKeywordMatch(str string, keywords string) bool {
	if str == "" || keywords == "" {
		return false
	}

	str = strings.ToLower(str)
	keywords = strings.ToLower(strings.ReplaceAll(keywords, " ", ""))

	for _, keyword := range strings.Split(keywords, ",") {
		if keyword[0] == '+' {
			if !strings.Contains(str, strings.Replace(keyword, "+", "", 1)) {
				return false
			}
		} else if keyword[0] == '-' {
			if strings.Contains(str, strings.Replace(keyword, "-", "", 1)) {
				return false
			}
		}
	}

	return true
}
