package main

import (
	"fmt"
	"net/http"
	"strconv"
	"os"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/streadway/amqp"
)

//Memcached client
var mc *memcache.Client

//Rabbit Connection
var rabbitConn *amqp.Connection
var rabbitChannel *amqp.Channel
var rabbitQueue amqp.Queue
var rabbitErr error
var rabbitChanDelivery <-chan amqp.Delivery

var memcachedServer string = "192.168.99.100" //"memcached1"
var rabbitServer string = "192.168.99.100" //"some-rabbit"

//Handle the incomming request
func handler(w http.ResponseWriter, r *http.Request) {

	consume := r.URL.Query().Get("consume")

	if consume == "" {
		//Get the stored visits from memcached
		it, err := mc.Get("visita")

		//Check the quantity
		quantity := getQuantity(it, err)

		fmt.Fprintf(w, "Visita %s!", strconv.Itoa(quantity))
		fmt.Println("Visita: ", strconv.Itoa(quantity))

		visits := []byte(strconv.Itoa(quantity))

		//Save the new value in memcached
		mc.Set(&memcache.Item{Key: "visita", Value: visits})

		//Publish the visit number to rabbit
		sendRabbit(visits)
	} else {
		//Consume from the queue
		strMessage := consumeRabbit()
		fmt.Fprintf(w, "Mensaje consumido: %s!", strMessage)
	}

}

func main() {

	fmt.Println("Iniciando Webserver....")
	hostname := os.Getenv("HOSTNAME")
	if hostname!="" {
		fmt.Println("Running from Container")
		memcachedServer = "memcached1"
		rabbitServer = "some-rabbit"		
	} else {
		fmt.Println("Running from Localhost")
	}

	//Esta IP es la IP del container, también se puede poner la de la máquina virtual y funciona
	//(que es más práctico porque es siempre la misma para todos los containers), la otra es poner el nomber del link

	//Create the memcached client with the container name linked
	mc = memcache.New(memcachedServer+":11211")

	//Create the Rabbit connection
	initRabbitMQ()

	//Creathe the webserver
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Println("Msg: ", msg, "Error: ", err)
	}
}

func initRabbitMQ() {

	//Create the connection with the rabbit cluster
	rabbitConn, rabbitErr = amqp.Dial("amqp://guest:guest@"+rabbitServer+":5672/")
	failOnError(rabbitErr, "Failed to connect to RabbitMQ")

	//Creathe the channel
	rabbitChannel, rabbitErr = rabbitConn.Channel()
	failOnError(rabbitErr, "Failed to open a channel")

	//Creathe the queue
	rabbitQueue, rabbitErr = rabbitChannel.QueueDeclare(
		"visits", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(rabbitErr, "Failed to declare a queue")

	//Create the channel to consume the events
	rabbitChanDelivery, rabbitErr = rabbitChannel.Consume(
		rabbitQueue.Name, // queue
		"",               // consumer
		false,            // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)

}

func sendRabbit(data []byte) {

	rabbitErr = rabbitChannel.Publish(
		"",               // exchange
		rabbitQueue.Name, // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		})

	fmt.Println("Sent to rabbit! bytes: ", len(data))
	failOnError(rabbitErr, "Failed to publish a message")
}

func getQuantity(it *memcache.Item, err error) int {

	quantity := 0
	if err == nil && it != nil {
		xs := string(it.Value)
		x, _ := strconv.Atoi(xs)
		quantity = x + 1
	} else {
		fmt.Println("Error", err)
	}

	return quantity
}

func consumeRabbit() string {

	failOnError(rabbitErr, "Failed to register a consumer")

	fmt.Println("Esperando Mensaje... ")
	
	message := <- rabbitChanDelivery

	message.Ack(false)

	fmt.Println("Mensaje recibido: ", string(message.Body))
	return string(message.Body)
			
}
