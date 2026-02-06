# RabbitMQ Event-Driven Microservices (Go)

This guide teaches you **from zero → production-style understanding** of RabbitMQ, AMQP, channels, exchanges, queues, publishers, consumers, routing keys, and how data actually flows between microservices.

---

# What is RabbitMQ?

**RabbitMQ is a message broker.**

It sits between services and moves data safely from one service to another.

Instead of this:

```
Order Service → Direct HTTP → Email Service
```

You get:

```
Order Service → RabbitMQ → Email Service
```

### Why this is powerful:

✅ Services don’t depend on each other
✅ Systems don’t crash if one service is down
✅ Easy to scale
✅ Handles retries
✅ Buffers traffic spikes

This is the foundation of **event-driven architecture** used in modern backend systems.

---

# What is AMQP?

**AMQP = Advanced Message Queuing Protocol**

It is the protocol RabbitMQ uses to communicate.

Think of it like:

```
HTTP → REST APIs
AMQP → Message Brokers
```

AMQP defines concepts like:

- Exchanges
- Queues
- Bindings
- Routing keys
- Channels

So when Go connects to RabbitMQ — it is speaking **AMQP**.

---

# Core Building Blocks

---

## Connection

A **TCP connection** between your microservice and RabbitMQ.

```
Service → TCP → RabbitMQ
```

### Important:

⚠️ Connections are **EXPENSIVE**
✔️ Create very few.

**Production pattern:**

```
1 microservice = 1 connection
```

---

## Channel

A **channel is a lightweight virtual connection inside a TCP connection.**

Think of it like:

```
Connection = Highway
Channel = Lane
```

You don’t build a new highway for every car.

### Rule:

✔️ Few connections
✔️ Many channels

---

# 🔥 Do Both Publisher and Consumer Need Channels?

## Short Answer:

👉 **YES. Always.**

### Producer Flow:

```
Connection
   ↓
Channel
   ↓
Publish → Exchange
```

### Consumer Flow:

```
Connection
   ↓
Channel
   ↓
Read ← Queue
```

---

# Exchange

An **exchange is a router.**

It decides:

> “Which queue should receive this message?”

### Without Exchange (Bad Design)

```
Producer → Queue
```

Producer must know queue names.

This creates **tight coupling**.

---

### With Exchange (Correct)

```
Producer → Exchange → Many Queues
```

Producer knows nothing about consumers.

This is called:

## ✅ Decoupling

A core principle of microservices.

---

# Queue

A queue stores messages until a consumer reads them.

Features:

✅ Durable storage
✅ Retry capability
✅ Ordering
✅ Backpressure handling

If a consumer is offline — messages wait safely.

---

# Binding

Binding connects a queue to an exchange.

```
Exchange → Queue
```

It tells RabbitMQ:

> “Send messages matching THIS pattern to THIS queue.”

---

# Routing Key

A routing key is basically the **event name.**

Examples:

```
order.placed
order.cancelled
payment.failed
user.created
```

The exchange uses this key to decide where messages go.

---

# Publisher

A publisher is a service that **emits events.**

Example:

Order service publishes:

```
order.placed
```

### Publisher MUST:

✅ Connect to RabbitMQ
✅ Open channel
✅ Declare exchange
✅ Publish events

❗ **Publisher NEVER cares about queues.**

Consumers own queues.

---

# Consumer

A consumer listens for events.

Example:

Email service listens for:

```
order.placed
```

### Consumer MUST:

✅ Connect
✅ Open channel
✅ Declare exchange (safe even if exists)
✅ Declare queue
✅ Bind queue
✅ Consume messages

---

# 🚨 Senior Rule Most Beginners Don’t Know

## ALWAYS declare topology on BOTH sides.

Why?

Services may start in any order.

- If consumer starts first → it creates queue
- If producer starts first → it creates exchange

System never breaks.

---

# Data Flow Example

Let’s say a new order is created.

```
Order Service
   |
   v
(order.placed)
   |
   v
Exchange
   |
   |----> email-queue → Email Service
   |
   |----> analytics-queue → Analytics Service
```

One event → multiple consumers.

This is how large systems scale effortlessly.

---

# Microservices Example (Go)

We will build TWO real services:

✅ Order Service → Publisher
✅ Email Service → Consumer

You can run them immediately.

---

---

# 🥇 Order Service (Publisher)

Sends event:

```
order.placed
```

### `order-service/main.go`

```go
package main

import (
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	// SAFE: creates exchange if missing
	err = ch.ExchangeDeclare(
		"order-exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	body := "NEW ORDER CREATED"

	err = ch.Publish(
		"order-exchange",
		"order.placed",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Order event published!")
}
```

---

# 🥈 Email Service (Consumer)

Listens for:

```
order.placed
```

### `email-service/main.go`

```go
package main

import (
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	ch.ExchangeDeclare(
		"order-exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	q, _ := ch.QueueDeclare(
		"email-queue",
		true,
		false,
		false,
		false,
		nil,
	)

	ch.QueueBind(
		q.Name,
		"order.placed",
		"order-exchange",
		false,
		nil,
	)

	msgs, _ := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	log.Println("📩 Waiting for order events...")

	for msg := range msgs {
		log.Println("📧 Sending email for:", string(msg.Body))
	}
}
```

---

# How To Test

## Step 1 — Start RabbitMQ

```bash
docker run -d \
-p 5672:5672 \
-p 15672:15672 \
rabbitmq:3-management
```

Dashboard:

```
http://localhost:15672
guest / guest
```

---

## Step 2 — Run Consumer FIRST

```
go run email-service/main.go
```

You should see:

```
Waiting for order events...
```

---

## Step 3 — Run Publisher

```
go run order-service/main.go
```

Consumer output:

```
Sending email for: NEW ORDER CREATED
```

You just built a real event-driven system.

---

# Architecture You Built

```
Order Service
   |
   v
order-exchange
   |
email-queue
   |
Email Service
```

### ZERO coupling.

Order service has **no idea** email exists.

That is why large distributed systems scale cleanly.

---
