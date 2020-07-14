package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"os"
	"payment/queue"
	"time"
)

type Order struct {
	Uuid      string    `json:"uuid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	ProductId string    `json:"product_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at,string"`
}

func notifyPaymentProcessed(order Order, ch *amqp.Channel) {
	json, _ := json.Marshal(order)
	queue.Notify(json, os.Getenv("RABBITMQ_PAYMENT_EXCHANGE"), "", ch)
	log.Println("Order payment processed: ", string(json))
}

func main() {
	in := make(chan []byte)

	connection := queue.Connect()
	queue.StartConsuming(os.Getenv("RABBITMQ_ORDER_QUEUE"), connection, in)

	var order Order
	for payload := range in {
		json.Unmarshal(payload, &order)
		order.Status = "approved"
		notifyPaymentProcessed(order, connection)
	}
}
