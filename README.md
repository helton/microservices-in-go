# microservices-in-go

Microservices examples using Go

## How to run

Run all microservices with **RabbitMQ** and **Redis** with:
```
docker-compose up -d --build
```

## Structure

- `catalog`: serves web pages retrieving data from the `product` service
- `product`: stores and retrieves products
- `checkout`: creates orders and send them to a queue (using RabbitMQ) to be processed
- `order`: processes the order from the queue (using RabbitMQ) and save it to the database (currently on Redis). It also updates the saved order when the payment status was updated (notified from queue)
- `payment`: process the order payment and notify via queue

## Flow

1 - List all available products (using the `catalog` service)
```
http://localhost:8082
```

2 - View a product
```
http://localhost:8082/45688cd6-7a27-4a7b-89c5-a9b604eefe2f
```

3 - Checkout a product
```
http://localhost:8083/45688cd6-7a27-4a7b-89c5-a9b604eefe2f
```

4 - See the order saved on Redis

If you have `redis-cli` installed you can use ...
```
redis-cli
KEYS *
GET eba1fd9a-439b-4d20-60ff-2b6aac87a18c
```

... otherwise you might want to use `nc` (netcat):
```
nc -v localhost 6379
KEYS *
GET eba1fd9a-439b-4d20-60ff-2b6aac87a18c
```

## Todo

- [ ] Create `Dockerfile` for Redis to create a custom image with the exchanges and queues definitions