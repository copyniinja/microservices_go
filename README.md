# Microservices Architecture (Event-Driven with RabbitMQ)

## Overview

This project implements a distributed microservices architecture written in **Go** with hybrid communication patterns:

- **Asynchronous**: Event-driven messaging via **RabbitMQ**
- **Synchronous**: Service-to-service calls via **HTTP + JSON**

The system demonstrates core backend engineering patterns including:

- Event-driven architecture (pub/sub)
- Service-to-service communication
- Centralized logging
- Email service integration
- Authentication with PostgreSQL
- NoSQL logging with MongoDB
- Dockerized infrastructure
- Message queue consumer pattern
- Health checks between services

---

## Architecture

```
Client (Frontend)
        |
        v
   Broker Service ←──────────────────┐
        |                              |
        ├─→ Auth Service              |
        ├─→ Mail Service              |
        └─→ RabbitMQ (publish events) │
               |                       │
               v                       │
         Listener Service (consume)────┘
               |
               v
         Logger Service (HTTP call)
         (Postgres/MongoDB)
```

### Communication Patterns

**1. Event-Driven (Async)**

- Broker Service → publishes log events to RabbitMQ (`logs_topic` exchange)
- Listener Service → subscribes to `log.INFO` and `log.ERROR` topics
- Uses topic-based routing for flexible event distribution

**2. Synchronous (HTTP/JSON)**

- Listener Service → Logger Service (persists logs)
- Broker Service → Auth Service (user validation)
- Broker Service → Mail Service (email delivery)

---

## Services

### Auth Service

Responsible for user management and authentication.

**Features:**

- User registration
- Credential validation
- Password hashing
- PostgreSQL integration
- Logging trigger
- Welcome email trigger

**Database:** PostgreSQL

---

### Logger Service

Centralized logging service for capturing system events.

**Features:**

- Accepts logs from other services
- Stores structured logs
- Designed for horizontal scalability

**Database:** MongoDB
(Perfect for high-write workloads)

---

### Mail Service

Handles transactional email delivery.

**Features:**

- `/send-email` HTTP endpoint
- Template-based emails
- SMTP integration
- MailHog support for local testing
- Designed to support real providers (SendGrid, SES, etc.)

---

### Broker Service

Acts as the **entry point** for client requests and routes traffic to appropriate services.

**Features:**

- `/log` endpoint that publishes events to RabbitMQ
- Routes requests to Auth, Logger, and Mail services
- Publishes log events with severity levels (`log.INFO`, `log.ERROR`)
- HTTP clients for service-to-service calls

**Why it exists:**

- Prevent frontend from calling multiple services directly
- Centralize request handling & event publishing
- Decouple services via event queue
- Improve security boundaries

---

### Listener Service

Asynchronous message consumer that bridges RabbitMQ and the Logger service.

**Features:**

- Subscribes to RabbitMQ `logs_topic` exchange
- Listens for `log.INFO` and `log.ERROR` topics
- Automatically binds to dynamically created queues
- Consumes messages and forwards to Logger Service via HTTP
- Runs continuously with graceful error handling
- Horizontal scalability via competing consumers pattern

**Flow:**

1. Broker Service publishes event to RabbitMQ: `POST /log` → RabbitMQ push
2. Listener Service consumes from queue (topic-based routing)
3. Unmarshals JSON payload and calls Logger Service HTTP API
4. Logger Service persists the log to MongoDB

**RabbitMQ Configuration:**

- **Exchange:** `logs_topic` (Topic exchange for flexible routing)
- **Queue:** Auto-generated (exclusive, durable)
- **Routing Keys:** `log.INFO`, `log.ERROR`

This decouples the Broker from direct Logger dependency and allows future consumers (email alerts, monitoring services, etc.) to bind to the same topics without modifying existing code.

---

### Frontend

Used primarily for:

- System testing
- Service health validation
- Integration verification

(Not focused on UI complexity.)

---

## Infrastructure

### Dockerized Environment

All services run inside Docker containers and communicate over a shared Docker network using **service names** (NOT localhost).

Example:

```
http://logger-service:6001
http://mail-service:6002
```

Docker Compose automatically creates a bridge network enabling internal DNS resolution.

---

## Email Testing (MailHog)

MailHog is used as a fake SMTP server for development.

**SMTP:** `mailhog:1025`
**Web UI:** `http://localhost:8025`

No real emails are sent.

---

## RabbitMQ Integration

**Current Implementation:**

The system already uses RabbitMQ for event-driven logging workflows:

- **Publisher:** Broker Service publishes log events
- **Consumer:** Listener Service consumes and processes events
- **Pattern:** Topic exchange with selective subscription
- **Reliability:** Auto-acknowledgment with durable queues

**Benefits:**

- Decouples Broker from Logger (no direct HTTP dependency)
- Allows multiple consumers to subscribe to the same topics
- Better resilience: if Logger is down, events stay in queue
- Foundation for future event-driven features (alerts, webhooks, analytics)

---

## Known Characteristics

**Synchronous (HTTP):**

- Broker ↔ Auth Service
- Broker ↔ Mail Service
- Listener ↔ Logger Service

**Asynchronous (RabbitMQ):**

- Broker → Listener (via message queue)

---

## Running the Project

### Start all services

```bash
docker compose up --build
```

---

### Verify Containers

```bash
docker ps
```

---

### Access front-end UI

```
http://localhost:8080
```

### Access MailHog UI

```
http://localhost:8025
```

---

## Engineering Highlights

This project demonstrates understanding of:

- Distributed system fundamentals
- Microservice boundaries
- Database-per-service pattern
- Async triggers
- Infrastructure design
- Production-style repository organization

---

## Future Improvements & Enhancements

- Expand event types (auth events, mail events, system events)
- Dead-letter queues (DLQ) for failed message handling
- Message retry strategies and circuit breakers
- Multiple consumer groups for different listening patterns
- gRPC for internal service communication
- Distributed tracing (OpenTelemetry)
- Kubernetes deployment
- Consumer lag monitoring
- Event sourcing for audit trails

---

## Author Notes

This repository reflects an intentional progression from:

**Monolith Thinking → Microservices → Event-Driven Systems**

The goal is not just building services — but understanding **system design evolution**.
