package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

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
var rabbitChanDeliveryTemp1 <-chan amqp.Delivery
var rabbitChanDeliveryTemp2 <-chan amqp.Delivery

var memcachedServer string = "192.168.99.100" //"memcached1"
var rabbitServer string = "192.168.99.100"    //"some-rabbit"

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
	if hostname != "" {
		fmt.Println("Running from Container")
		memcachedServer = "memcached1"
		rabbitServer = "some-rabbit"
	} else {
		fmt.Println("Running from Localhost")
	}

	//Esta IP es la IP del container, también se puede poner la de la máquina virtual y funciona
	//(que es más práctico porque es siempre la misma para todos los containers), la otra es poner el nomber del link

	//Create the memcached client with the container name linked
	mc = memcache.New(memcachedServer + ":11211")

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
	rabbitConn, rabbitErr = amqp.Dial("amqp://guest:guest@" + rabbitServer + ":5672/")
	failOnError(rabbitErr, "Failed to connect to RabbitMQ")

	//Creathe the channel
	rabbitChannel, rabbitErr = rabbitConn.Channel()
	failOnError(rabbitErr, "Failed to open a channel")

	//Set the channel to prefetch only 1, in order to avoid lots of un-aknowledged messages
	rabbitErr = rabbitChannel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(rabbitErr, "Failed to set QoS")

	//create the exchange of fanout type
	rabbitErr = rabbitChannel.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(rabbitErr, "Failed to declare an exchange")

	//Creathe the queue (the standar one)
	rabbitQueue, rabbitErr = rabbitChannel.QueueDeclare(
		"visits_durable", // name
		true,             // durable --> put in true, the queue will be stored in disk
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	failOnError(rabbitErr, "Failed to declare a queue")

	//create two queues, to use with the exchange which will send duplicate events to both of them
	rabbitQueueTemp1, err1 := rabbitChannel.QueueDeclare(
		"",    // name -->without name because we want to end it at the end
		false, // durable
		false, // delete when usused
		true,  // exclusive --> indicates that when the connection is close the queue is deleted
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err1, "Failed to declare a queue")

	rabbitQueueTemp2, err2 := rabbitChannel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err2, "Failed to declare a queue")

	//create the two bindings between the exchange and the temporal queues
	rabbitErr = rabbitChannel.QueueBind(
		rabbitQueueTemp1.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil)
	failOnError(rabbitErr, "Failed to bind a queue")

	rabbitErr = rabbitChannel.QueueBind(
		rabbitQueueTemp2.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil)
	failOnError(rabbitErr, "Failed to bind a queue")

	//Create the channel to consume the events
	rabbitChanDelivery, rabbitErr = rabbitChannel.Consume(
		rabbitQueue.Name, // queue
		"",               // consumer
		false,            // auto-ack --> have to ack(false) each message
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)

	//Create the channel to consume the events
	rabbitChanDeliveryTemp1, rabbitErr = rabbitChannel.Consume(
		rabbitQueueTemp1.Name, // queue
		"",    // consumer
		false, // auto-ack --> have to ack(false) each message
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	//Create the channel to consume the events
	rabbitChanDeliveryTemp2, rabbitErr = rabbitChannel.Consume(
		rabbitQueueTemp2.Name, // queue
		"",    // consumer
		false, // auto-ack --> have to ack(false) each message
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

}

func sendRabbit(data []byte) {

	rabbitErr = rabbitChannel.Publish(
		"",               // exchange --> Use the default exchange, thats why we use "", it forward the message to the queue specified in "routing_key"
		rabbitQueue.Name, // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, //enables the message as persistant, so it is saved to disk
			ContentType:  "text/plain",
			Body:         data,
		})

	rabbitErr = rabbitChannel.Publish(
		"logs", // exchange --> Use the logs exchage
		"",     // routing key
		false,  // mandatory
		false,  // immediate
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

	message := <-rabbitChanDelivery
	messageTemp1 := <-rabbitChanDeliveryTemp1
	messageTemp2 := <-rabbitChanDeliveryTemp2

	message.Ack(false)
	messageTemp1.Ack(false)
	messageTemp2.Ack(false)

	text := "Mensaje recibido: " + string(message.Body) + " - " + string(messageTemp1.Body) + " - " + string(messageTemp2.Body)
	fmt.Println(text)
	return text

}
