package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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

func init() {
	port = os.Getenv("SERVICE_PORT")
}

func loadData() []byte {
	jsonFile, err := os.Open("products.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer jsonFile.Close()
	data, _ := ioutil.ReadAll(jsonFile)
	return data
}

func ListProducts(w http.ResponseWriter, r *http.Request) {
	products := loadData()
	w.Write(products)
}

func GetProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data := loadData()

	var products Products
	json.Unmarshal(data, &products)

	for _, v := range products.Products {
		if v.Uuid == vars["id"] {
			product, _ := json.Marshal(v)
			w.Write(product)
		}
	}
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	if len(port) == 0 {
		log.Fatal("SERVICE_PORT should be specified")
	}

	router := mux.NewRouter()
	router.Use(commonMiddleware)

	router.HandleFunc("/products", ListProducts)
	router.HandleFunc("/products/{id}", GetProductById)
	log.Printf("Product service running on port %v", port)
	log.Fatal(http.ListenAndServe(port, router))
}
