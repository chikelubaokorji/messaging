package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"net/http"
	"os"
)

var rabbitHost = os.Getenv("RABBIT_HOST")
var rabbitPort = os.Getenv("RABBIT_PORT")
var rabbitUser = os.Getenv("RABBIT_USERNAME")
var rabbitPassword = os.Getenv("RABBIT_PASSWORD")

func main() {

	router := httprouter.New()

	handler := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		submit(p)
	}

	router.POST("/publish/:message", handler)

	fmt.Println("Running...")

	log.Fatal(http.ListenAndServe(":80", router))
}

func submit(p httprouter.Params) {

	message := p.ByName("message")

	fmt.Println("Received message: " + message)

	conn, err := amqp.Dial("amqp://" + rabbitUser + ":" + rabbitPassword + "@" + rabbitHost + ":" + rabbitPort + "/")

	if err != nil {
		log.Fatalf("%s: %s", "RabbitMQ Connection Failed", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalf("%s: %s", "Channel Creation Failed", err)
	}

	defer ch.Close()

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

	err = ch.Publish(
		"",       // exchange
		que.Name, // key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

	if err != nil {
		log.Fatalf("%s: %s", "Message Publishing Failed", err)
	}

	fmt.Println("Publishing Successful")
}
