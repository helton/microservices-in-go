package main

import (
	"encoding/json"
	"flag"
	"github.com/nu7hatch/gouuid"
	"github.com/streadway/amqp"
	"log"
	"order/db"
	"order/queue"
	"os"
	"time"
)

type Product struct {
	Uuid    string  `json:"uuid"`
	Product string  `json:"product"`
	Price   float64 `json:"price,string"`
}

type Order struct {
	Uuid      string    `json:"uuid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	ProductId string    `json:"product_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at,string"`
}

func createOrder(payload []byte) Order {
	var order Order
	json.Unmarshal(payload, &order)

	uuid, _ := uuid.NewV4()
	order.Uuid = uuid.String()
	order.Status = "pending"
	order.CreatedAt = time.Now()
	saveOrder(order)

	return order
}

func saveOrder(order Order) {
	json, _ := json.Marshal(order)
	connection := db.Connect()
	err := connection.Set(order.Uuid, string(json), 0).Err()
	if err != nil {
		panic(err.Error())
	}

	log.Println("Order saved/updated: ", string(json))
}

func notifyOrderCreated(order Order, ch *amqp.Channel) {
	json, _ := json.Marshal(order)
	queue.Notify(json, os.Getenv("RABBITMQ_ORDER_EXCHANGE"), "", ch)
	log.Println("Order created: ", string(json))
}

func main() {
	var param string

	flag.StringVar(&param, "opt", "", "Usage")
	flag.Parse()

	in := make(chan []byte)
	connection := queue.Connect()

	switch param {
	case "checkout":
		log.Println("Consuming checkout queue...")

		queue.StartConsuming(os.Getenv("RABBITMQ_CHECKOUT_QUEUE"), connection, in)
		for payload := range in {
			order := createOrder(payload)
			notifyOrderCreated(order, connection)
		}
	case "payment":
		log.Println("Consuming payment queue...")

		queue.StartConsuming(os.Getenv("RABBITMQ_PAYMENT_QUEUE"), connection, in)
		var order Order
		for payload := range in {
			json.Unmarshal(payload, &order)
			saveOrder(order)
			log.Println("Order payment status updated: ", string(payload))
		}
	}

}
