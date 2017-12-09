package main

import (
	"os"
	// "encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	// "time"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

func confirmOne(ack, nack chan uint64) {
	select {
	case tag := <-ack:
		log.Printf("acked: %d\n", tag)
	case tag := <-nack:
		log.Printf("nack alert: %d\n", tag)
	}
}

var connection *amqp.Connection
var channel *amqp.Channel
var ack, nack chan uint64
var q amqp.Queue

func restHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	vars := mux.Vars(r)

	// Stop here if its Preflighted OPTIONS request
	if r.Method == "POST" {

		ordersNum, err := strconv.Atoi(vars["ordersNum"])
		if err != nil {
			log.Println(err)
		}

		// log.Println("request received:", r)
		log.Printf("orders received: %d", ordersNum)

		min, _ := strconv.Atoi(os.Getenv("MIN"))
		max, _ := strconv.Atoi(os.Getenv("MAX"))

		for i := 0; i < ordersNum; i++ {
			message := strconv.Itoa(min + rand.Intn(max-min))

			err = channel.Publish("", q.Name, true, false, amqp.Publishing{
				Headers:         amqp.Table{},
				ContentType:     "text/plain",
				ContentEncoding: "UTF-8",
				Body:            []byte(message),
				DeliveryMode:    amqp.Transient,
				Priority:        0,
			},
			)
			if err != nil {
				log.Println(err)
				return
			}

			confirmOne(ack, nack)
		}
	}

	output := "{\"ordersNum\": \"" + vars["ordersNum"] + "\"}"
	fmt.Fprintln(w, string(output))
}

func main() {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	connection = conn

	ch, err := connection.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()
	channel = ch

	channel.Confirm(false)

	ack, nack = channel.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

	q, err = channel.QueueDeclare(os.Getenv("queueName"), true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/{ordersNum:[0-9]+}", restHandler).Methods("POST", "OPTIONS")

	fmt.Println("Demo Producer started, listening on port ", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))

}
