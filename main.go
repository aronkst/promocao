package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
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
	RegexPriceOld                  string  `json:"regex_price_old"`
	RegexPriceNew                  string  `json:"regex_price_new"`
	MinimumPorcentageDiscount      float64 `json:"minimum_porcentage_discount"`
	MaximumPages                   int     `json:"maximum_pages"`
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
var page = 1
var products = make([]product, 0)
var waitGroup sync.WaitGroup

func loadSite(url string) *goquery.Document {
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		err := fmt.Sprintf("status code error: %d %s", response.StatusCode, response.Status)
		panic(err)
	}

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		panic(err)
	}

	return document
}

func formatPrice(price string) float64 {
	if price == "" {
		return 0
	}

	var re = regexp.MustCompile(config.RegexPriceOld)
	priceNew := re.ReplaceAllString(price, config.RegexPriceNew)

	priceFloat, err := strconv.ParseFloat(priceNew, 64)
	if err != nil {
		panic(err)
	}

	return priceFloat
}

func clearString(value string) string {
	value = strings.ReplaceAll(value, "\n", "")
	value = strings.ReplaceAll(value, "\r", "")
	value = strings.TrimSpace(value)
	return value
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

func getProductsPrice(selector *goquery.Selection) float64 {
	return getPrice("ProductsPrice", selector)
}

func getProductsFinalPrice(selector *goquery.Selection) float64 {
	return getPrice("ProductsFinalPrice", selector)
}

func getProductsURL(selector *goquery.Selection) string {
	return getCSSValue("ProductsURL", selector)
}

func getProductSeller(selector *goquery.Selection) string {
	return getCSSValue("ProductSeller", selector)
}

func getProductImage(selector *goquery.Selection) string {
	return getCSSValue("ProductImage", selector)
}

func getProductTitle(selector *goquery.Selection) string {
	return getCSSValue("ProductTitle", selector)
}

func getProductPrice(selector *goquery.Selection) float64 {
	return getPrice("ProductPrice", selector)
}

func getNextPage(selector *goquery.Selection) string {
	return getCSSValue("NextPage", selector)
}

func loadProducts(url string) {
	document := loadSite(url)
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
		loadProducts(nextPage)
	}
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

func loadProduct(url string, porcentageDiscount float64) {
	defer waitGroup.Done()

	document := loadSite(url)
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

func deleteHTML() {
	if _, err := os.Stat(config.Output); os.IsNotExist(err) {
		return
	}

	err := os.Remove(config.Output)
	if err != nil {
		panic(err)
	}
}

func writeHTML(text string) {
	file, err := os.OpenFile(config.Output, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(text)
	if err != nil {
		panic(err)
	}

	err = file.Sync()
	if err != nil {
		panic(err)
	}
}

func createHTML() {
	deleteHTML()

	html := `<!doctype html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
		<title>Result</title>
		<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
	</head>
	<body>
		<main role="main">
			<div class="album py-5 bg-light">
				<div class="container">
					<div class="row">
`

	text := []byte(html)
	err := ioutil.WriteFile(config.Output, text, 0644)
	if err != nil {
		panic(err)
	}
}

func appendHTML(values product) {
	html := `
<div class="col-md-4">
	<div class="card mb-4 shadow-sm">
		<img src="{{ image }}" class="bd-placeholder-img card-img-top" width="100%" focusable="false">
		<div class="card-body">
			<p class="card-text">{{ title }}</p>
			<div class="d-flex justify-content-between align-items-center">
				<div class="btn-group">
					<a href="{{ url }}" class="btn btn-sm btn-outline-secondary" target="_blank">Visualizar</a>
				</div>
				<small class="text-muted">{{ porcentageDiscount }}</small>
				<small class="text-muted">
					<strong>{{ price }}</strong>
				</small>
			</div>
		</div>
	</div>
</div>
`

	porcentageDiscountString := fmt.Sprintf("%.2f", values.PorcentageDiscount)
	priceString := fmt.Sprintf("%.2f", values.Price)

	html = strings.ReplaceAll(html, "{{ image }}", values.Image)
	html = strings.ReplaceAll(html, "{{ title }}", values.Title)
	html = strings.ReplaceAll(html, "{{ url }}", values.URL)
	html = strings.ReplaceAll(html, "{{ porcentageDiscount }}", porcentageDiscountString)
	html = strings.ReplaceAll(html, "{{ price }}", priceString)

	writeHTML(html)
}

func finishHTML() {
	html := `					</div>
				</div>
			</div>
		</main>
		<script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js" integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1" crossorigin="anonymous"></script>
		<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js" integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM" crossorigin="anonymous"></script>
	</body>
</html>`

	writeHTML(html)
}

func showHTML() {
	path, err := filepath.Abs(filepath.Dir(config.Output))
	if err != nil {
		panic(err)
	}
	file := filepath.Join(path, config.Output)
	fmt.Println(file)
}

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
	fmt.Println()

	createHTML()
	for _, values := range products {
		appendHTML(values)
	}
	finishHTML()
	showHTML()
}
