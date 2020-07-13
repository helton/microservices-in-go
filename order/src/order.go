package main

import (
	"encoding/json"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"order/db"
	"order/queue"
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

func createOrder(payload []byte) {
	var order Order
	json.Unmarshal(payload, &order)

	uuid, _ := uuid.NewV4()
	order.Uuid = uuid.String()
	order.Status = "pending"
	order.CreatedAt = time.Now()
	saveOrder(order)
}

func saveOrder(order Order) {
	json, _ := json.Marshal(order)
	connection := db.Connect()
	err := connection.Set(order.Uuid, string(json), 0).Err()
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Order %s saved!\n", order.Uuid)
}

func main() {
	in := make(chan []byte)

	connection := queue.Connect()
	queue.StartConsuming(connection, in)

	for payload := range in {
		createOrder(payload)
	}
}
