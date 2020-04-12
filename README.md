# Promocao

This is a project developed in the Go programming language to search for promotions in any e-commerce, using as a base a JSON file, which will contain all the configurations so that the crawler can search the products in a list on any e-commerce.

This JSON with the settings can be found in the `main.go` file or be viewed just below this paragraph. In its structure, for example, it contains information from the URL where it should start searching for products (a list of products, which on the same page contains a link to go to the next page), the CSS so that the crawler can find the product, the CSS of the product price, the CSS with the reduced product value if is in the sale, the CSS that contains the product link, the CSS of the link that contains the URL for the next page, among many others.

```json
{
    "seller": "E-commerce",
    "url": "https://www.e-commerce.com/...",
    "css_products": "div.product",
    "css_products_price": "div.price",
    "css_products_price_attribute": "data-price",
    "css_products_final_price": "div.final_price",
    "css_products_final_price_attribute": "",
    "css_products_url": "a.link",
    "css_products_url_attribute": "href",
    "css_product_seller": "div.seller_name",
    "css_product_seller_attribute": "",
    "css_product_image": "img",
    "css_product_image_attribute": "src",
    "css_product_title": "h1.title",
    "css_product_title_attribute": "",
    "css_product_price": "p.price",
    "css_product_price_attribute": "",
    "css_next_page": "a.next_page",
    "css_next_page_attribute": "href",
    "regex_url_old": "(.*)",
    "regex_url_new": "https:$1",
    "regex_price_old": "",
    "regex_price_new": "",
    "minimum_porcentage_discount": 15,
    "maximum_value": 100000,
    "minimum_value": 1,
    "maximum_pages": 20,
    "output": "output.html",
    "sleep": 1000
}
```


It is possible to configure the number of pages that the crawler must go through, the maximum price of the product, the minimum price of the product, the minimum value of the percentage of the product if is in the sale, among many others, all of this in JSON.

This project is personal and has been used by me a few times. It may be that it doesn't work the way you need it, or even on the e-commerce that you want to use it. And it also doesn't have all the features I imagined, but it is 100% functional.

## Run the application

To run the application use the command:

`go run .`

If you prefer, build the application with the command below:

`go build .`

Then to run the application, use the following command:

`./promocao`
