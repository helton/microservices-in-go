package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Product struct {
	Uuid    string  `json:"uuid"`
	Product string  `json:"product"`
	Price   float64 `json:"price,string"`
}

type Products struct {
	Products []Product
}

var port string
var productsUrl string

func init() {
	port = os.Getenv("SERVICE_PORT")
	productsUrl = os.Getenv("PRODUCTS_URL")
}

func loadProducts() []Product {
	response, err := http.Get(productsUrl + "/products")
	if err != nil {
		log.Println("Could not connect to the product service")
		log.Fatal(err.Error())
	}
	data, _ := ioutil.ReadAll(response.Body)

	var products Products
	json.Unmarshal(data, &products)

	return products.Products
}

func ListProducts(w http.ResponseWriter, r *http.Request) {
	products := loadProducts()
	t := template.Must(template.ParseFiles("templates/catalog.html"))
	t.Execute(w, products)
}

func ShowProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response, err := http.Get(productsUrl + "/products/" + vars["id"])
	if err != nil {
		log.Printf("Could not connect to the product service to access /products/%s", vars["id"])
		log.Fatal(err.Error())
	}

	data, _ := ioutil.ReadAll(response.Body)
	var product Product
	json.Unmarshal(data, &product)

	t := template.Must(template.ParseFiles("templates/view.html"))
	t.Execute(w, product)
}

func main() {
	if len(port) == 0 {
		log.Fatal("SERVICE_PORT should be specified")
	}

	router := mux.NewRouter()
	router.HandleFunc("/", ListProducts)
	router.HandleFunc("/{id}", ShowProduct)

	log.Printf("Catalog service running on port %v", port)
	log.Fatal(http.ListenAndServe(port, router))
}
