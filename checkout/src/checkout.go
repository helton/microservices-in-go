package main

import (
	"checkout/queue"
	"encoding/json"
	"fmt"
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

type Order struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	ProductId string `json:"product_id"`
}

var port string
var productsUrl string

func init() {
	port = os.Getenv("SERVICE_PORT")
	productsUrl = os.Getenv("PRODUCTS_URL")
}

func displayCheckout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response, err := http.Get(productsUrl + "/products/" + vars["id"])
	if err != nil {
		log.Printf("Could not connect to the product service to access /products/%s", vars["id"])
		log.Fatal(err.Error())
	}

	data, _ := ioutil.ReadAll(response.Body)
	var product Product
	json.Unmarshal(data, &product)

	t := template.Must(template.ParseFiles("templates/checkout.html"))
	t.Execute(w, product)
}

func finish(w http.ResponseWriter, r *http.Request) {
	var order Order
	order.Name = r.FormValue("name")
	order.Email = r.FormValue("email")
	order.Phone = r.FormValue("phone")
	order.ProductId = r.FormValue("product_id")

	data, _ := json.Marshal(order)
	fmt.Println(string(data))

	connection := queue.Connect()
	queue.Notify(data, os.Getenv("RABBITMQ_CONSUMER_EXCHANGE"), "", connection)

	w.Write([]byte("Processed"))
}

func main() {
	if len(port) == 0 {
		log.Fatal("SERVICE_PORT should be specified")
	}

	router := mux.NewRouter()
	router.HandleFunc("/finish", finish)
	router.HandleFunc("/{id}", displayCheckout)

	log.Printf("Checkout service running on port %v", port)
	log.Fatal(http.ListenAndServe(port, router))
}
