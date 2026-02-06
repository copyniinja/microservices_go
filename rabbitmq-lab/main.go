package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Handle Failure
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s:%s", msg, err)
	}
}

// Connect to RabbitMQ Server //url : "amqp://guest:guest@localhost:5672/"
func connectRabbitMQ(url string) *amqp.Connection {
	// amqp.Dial
	conn, err := amqp.Dial(url)

	failOnError(err, "Failed to connect to Rabbitmq")

	defer conn.Close()

	return conn
}

/*
Channels are lightweight connections inside the TCP connection.
Use: Every producer and consumer usually has one channel.
*/
// Create Channel
func createChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	return ch
}

/*
Exchange:
-> The post office clerk. Decides which queue(s) get the message.
-> Producers never send messages directly to queues, only to exchanges.
Types of exchanges:

1. Direct	   Route by exact routing key to matching queue
2. Fanout	   Route to all queues bound to this exchange (broadcast)
3. Topic	   Route to queues matching a pattern (like order.*)
4. Headers	 Route based on message headers instead of key
*/
func declareExchange(ch *amqp.Channel, name, kind string) {
	err := ch.ExchangeDeclare(
		name, // "order-exchange", "payment-exchange", "user-exchange"
		kind,
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare Exchange")
}

/*
	this function does two critical things:

// 1. Creates a queue (message storage)
//    -> Messages sit here until consumed.
// 2. Binds the queue to an exchange
//    -> Tells RabbitMQ:
//       "If a message arrives at THIS exchange
//        with THIS routing key, (<domain>.<action> =order.placed)
//        deliver it to THIS queue."
// Without binding, messages will be dropped because the exchange won't know where to route them.
*/
func declareQueueAndBind(ch *amqp.Channel, queueName string, exchangeName string, routingKey string) {

	// 1. Declare a Queue
	q, err := ch.QueueDeclare(
		queueName,
		true,  //durable :Keeps the queue alive even if RabbitMQ restarts.
		false, // auto delete
		false, // exclusive :If true, only this connection can use the queue.
		false, // no wait :  Wait for RabbitMQ to confirm creation.
		nil,   // arguments : // Optional advanced configs like: dead-letter queues,message TTL,max length...
	)
	failOnError(err, "Failed to declare a queue")

	// 2. Bind the queue to an exchange
	err = ch.QueueBind(
		q.Name,
		routingKey,   // "order.email"   -> only email events land here.
		exchangeName, // The exchange that will route messages to this queue.
		false,        // no wait
		nil,          // arguments
	)
	failOnError(err, "Failed to bind queue")

}

// Producer
func publishMessage(ch *amqp.Channel, exchangeName, routingKey, body string) {
	err := ch.Publish(
		exchangeName,
		routingKey,
		false, // mandatory
		/*
		 mandatory = true → "Don't silently drop my message if it can't be routed to at least one queue"
		 mandatory = false (default) → "If it can't be routed to any queue → just drop it quietly
		*/
		false, // Immediate
		/*
			immediate = true → "Don't accept this message unless there's at least one active consumer ready to take it right now (no queueing allowed)"
			immediate = false (default) → "Queue it normally even if no consumer is connected yet
		*/
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)

	failOnError(err, "Failed to publish message")
	log.Printf("Publisher Sent message: %s (routing key: %s)", body, routingKey)
}

// Consumer
func consumeMessages(ch *amqp.Channel, queueName string) {
	msgs, err := ch.Consume(
		queueName,
		"",    // consumer tag (empty = auto)
		false, // autoAck
		false, // not exclusive
		false, // receive own messages too
		false, // wait for confirmation
		nil,   // no special args
	)
	failOnError(err, "Failed to register consumer")

	go func() {
		for d := range msgs {
			log.Printf("Consumer^^ %s received message: %s", queueName, d.Body)
		}
	}()

	log.Printf("Consumer waiting for messages on %s", queueName)

	select {} // block forever

}

// -----------------
func main() {

}
