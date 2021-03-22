package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
)

var rabbitHost = os.Getenv("RABBIT_HOST")
var rabbitPort = os.Getenv("RABBIT_PORT")
var rabbitUser = os.Getenv("RABBIT_USERNAME")
var rabbitPassword = os.Getenv("RABBIT_PASSWORD")

func main() {
	ingest()
}

func ingest() {

	conn, err := amqp.Dial("amqp://" + rabbitUser + ":" + rabbitPassword + "@" + rabbitHost + ":" + rabbitPort + "/")

	if err != nil {
		log.Fatalf("%s: %s", "RabbitMQ Connection Failed", err)
	}

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalf("%s: %s", "Channel Creation Failed", err)
	}

	que, err := ch.QueueDeclare(
		"publisher", // name
		true,        // durable
		false,       // autoDelete
		false,       // exclusive
		false,       // noWait
		nil,         // args
	)

	if err != nil {
		log.Fatalf("%s: %s", "Queue Declaration Failed", err)
	}

	fmt.Println("Channel Creation and Queue Declaration Established")

	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		que.Name, // queue
		"",       // consumer
		false,    // autoAck
		false,    // exclusive
		false,    // noLocal
		false,    // noWait
		nil,      // args
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed To Register Consumer", err)
	}

	always := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			d.Ack(false)
		}
	}()

	fmt.Println("Running...")
	<-always
}
