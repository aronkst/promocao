package main

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"reflect"

	"github.com/PuerkitoBio/goquery"
)

func loadSite(url string) (*goquery.Document, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		message := fmt.Sprintf("status code error: %d %s", response.StatusCode, response.Status)
		err := errors.New(message)
		return nil, err
	}

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	return document, nil
}

func getCSSValue(field string, selector *goquery.Selection) string {
	var css, cssAttribute, value string

	reflector := reflect.ValueOf(config)
	reflectorCSS := reflect.Indirect(reflector).FieldByName("CSS" + field)
	reflectorCSSAttribute := reflect.Indirect(reflector).FieldByName("CSS" + field + "Attribute")

	css = reflectorCSS.String()
	cssAttribute = reflectorCSSAttribute.String()

	if cssAttribute == "" {
		value = selector.Find(css).Text()
	} else {
		value, _ = selector.Find(css).Attr(cssAttribute)
	}

	value = clearString(value)
	return value
}

func getPrice(price string, selector *goquery.Selection) float64 {
	priceString := getCSSValue(price, selector)
	priceFloat := formatPrice(priceString)
	return priceFloat
}

func getURL(url string, selector *goquery.Selection) string {
	urlCSS := getCSSValue(url, selector)
	urlNew := formatURL(urlCSS)
	return urlNew
}

func getProductsPrice(selector *goquery.Selection) float64 {
	return getPrice("ProductsPrice", selector)
}

func getProductsFinalPrice(selector *goquery.Selection) float64 {
	return getPrice("ProductsFinalPrice", selector)
}

func getProductsURL(selector *goquery.Selection) string {
	return getURL("ProductsURL", selector)
}

func getProductSeller(selector *goquery.Selection) string {
	return getCSSValue("ProductSeller", selector)
}

func getProductImage(selector *goquery.Selection) string {
	return getURL("ProductImage", selector)
}

func getProductTitle(selector *goquery.Selection) string {
	return getCSSValue("ProductTitle", selector)
}

func getProductPrice(selector *goquery.Selection) float64 {
	return getPrice("ProductPrice", selector)
}

func getNextPage(selector *goquery.Selection) string {
	return getURL("NextPage", selector)
}

func loadProducts(url string) {
	document, err := loadSite(url)
	if err != nil {
		return
	}

	fmt.Print("#")

	document.Find(config.CSSProducts).Each(func(i int, s *goquery.Selection) {
		fmt.Print(".")

		price := getProductsPrice(s)
		finalPrice := getProductsFinalPrice(s)
		if finalPrice <= 0 {
			return
		}

		porcentageDiscount := math.Abs(((finalPrice * 100) / price) - 100)
		if porcentageDiscount < config.MinimumPorcentageDiscount {
			return
		}

		urlProduct := getProductsURL(s)
		waitGroup.Add(1)
		go loadProduct(urlProduct, porcentageDiscount)
	})

	if page < config.MaximumPages {
		page++
		nextPage := getNextPage(document.Selection)
		if nextPage != "" {
			loadProducts(nextPage)
		}
	}
}

func loadProduct(url string, porcentageDiscount float64) {
	defer waitGroup.Done()

	if url == "" {
		return
	}

	document, err := loadSite(url)
	if err != nil {
		return
	}

	fmt.Print("+")

	if config.CSSProductSeller != "" {
		seller := getProductSeller(document.Selection)
		if compareSeller(seller, config.Seller) == false {
			return
		}
	}

	image := getProductImage(document.Selection)
	title := getProductTitle(document.Selection)
	price := getProductPrice(document.Selection)

	prod := product{
		Image:              image,
		Title:              title,
		URL:                url,
		Price:              price,
		PorcentageDiscount: porcentageDiscount,
	}
	products = append(products, prod)
}
