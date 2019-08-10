package main

import (
	"regexp"
	"strconv"
	"strings"
)

func setRegex(value string, regexOld string, regexNew string) string {
	var re = regexp.MustCompile(regexOld)
	newValue := re.ReplaceAllString(value, regexNew)
	return newValue
}

func formatPrice(price string) float64 {
	if price == "" {
		return 0
	}

	priceNew := setRegex(price, config.RegexPriceOld, config.RegexPriceNew)

	priceFloat, err := strconv.ParseFloat(priceNew, 64)
	if err != nil {
		panic(err)
	}

	return priceFloat
}

func formatURL(url string) string {
	if config.RegexURLNew == "" || url == "" {
		return url
	}

	urlNew := setRegex(url, config.RegexURLOld, config.RegexURLNew)
	return urlNew
}

func clearString(value string) string {
	value = strings.ReplaceAll(value, "\n", "")
	value = strings.ReplaceAll(value, "\r", "")
	value = strings.TrimSpace(value)
	return value
}

func compareSeller(seller1 string, seller2 string) bool {
	seller1 = strings.TrimSpace(seller1)
	seller1 = strings.ToLower(seller1)
	seller1 = strings.TrimSuffix(seller1, "\n")
	seller1 = strings.TrimSuffix(seller1, "\r")

	seller2 = strings.TrimSpace(seller2)
	seller2 = strings.ToLower(seller2)
	seller2 = strings.TrimSuffix(seller2, "\n")
	seller2 = strings.TrimSuffix(seller2, "\r")

	return seller1 == seller2
}
