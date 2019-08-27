package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"sync"
)

type configuration struct {
	Seller                         string  `json:"seller"`
	URL                            string  `json:"url"`
	CSSProducts                    string  `json:"css_products"`
	CSSProductsPrice               string  `json:"css_products_price"`
	CSSProductsPriceAttribute      string  `json:"css_products_price_attribute"`
	CSSProductsFinalPrice          string  `json:"css_products_final_price"`
	CSSProductsFinalPriceAttribute string  `json:"css_products_final_price_attribute"`
	CSSProductsURL                 string  `json:"css_products_url"`
	CSSProductsURLAttribute        string  `json:"css_products_url_attribute"`
	CSSProductSeller               string  `json:"css_product_seller"`
	CSSProductSellerAttribute      string  `json:"css_product_seller_attribute"`
	CSSProductImage                string  `json:"css_product_image"`
	CSSProductImageAttribute       string  `json:"css_product_image_attribute"`
	CSSProductTitle                string  `json:"css_product_title"`
	CSSProductTitleAttribute       string  `json:"css_product_title_attribute"`
	CSSProductPrice                string  `json:"css_product_price"`
	CSSProductPriceAttribute       string  `json:"css_product_price_attribute"`
	CSSNextPage                    string  `json:"css_next_page"`
	CSSNextPageAttribute           string  `json:"css_next_page_attribute"`
	RegexURLOld                    string  `json:"regex_url_old"`
	RegexURLNew                    string  `json:"regex_url_new"`
	RegexPriceOld                  string  `json:"regex_price_old"`
	RegexPriceNew                  string  `json:"regex_price_new"`
	MinimumPorcentageDiscount      float64 `json:"minimum_porcentage_discount"`
	MaximumValue                   float64 `json:"maximum_value"`
	MinimumValue                   float64 `json:"minimum_value"`
	MaximumPages                   int64   `json:"maximum_pages"`
	Output                         string  `json:"output"`
}

type product struct {
	Image              string
	Title              string
	URL                string
	Price              float64
	PorcentageDiscount float64
}

var config configuration
var page int64 = 1
var products = make([]product, 0)
var waitGroup sync.WaitGroup
var finalHTML string

func main() {
	args := os.Args
	if len(args) != 2 {
		return
	}

	configFile, err := os.Open(args[1])
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)
	json.Unmarshal(byteValue, &config)

	loadProducts(config.URL)
	waitGroup.Wait()

	sort.Slice(products, func(i, j int) bool {
		return products[i].PorcentageDiscount > products[j].PorcentageDiscount
	})

	fmt.Println()

	deleteHTML()
	createHTML()
	for _, values := range products {
		appendHTML(values)
	}
	finishHTML()
	writeHTML()
	showHTML()
}
