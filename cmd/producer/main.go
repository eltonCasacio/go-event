package main

import "github.com/eltoncasacio/go-event/pkg/rabbitmq"

func main() {
	ch := rabbitmq.OpenChannel()
	defer ch.Close()

	if err := rabbitmq.Publish(ch, "???", "amq.direct"); err != nil {
		panic(err)
	}

}
