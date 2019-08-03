package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

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
