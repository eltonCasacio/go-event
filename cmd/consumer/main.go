package main

import (
	"fmt"

	"github.com/eltoncasacio/go-event/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	ch := rabbitmq.OpenChannel()
	defer ch.Close()

	msgs := make(chan amqp.Delivery)

	go rabbitmq.Consume(ch, msgs, "minha_fila")

	for msg := range msgs {
		fmt.Println(string(msg.Body))
		msg.Ack(true)
	}
}
