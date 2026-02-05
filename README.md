# Microservices Architecture (HTTP-Based)

## Overview

This project implements a distributed microservices architecture written in **Go**, where services communicate synchronously using **HTTP + JSON**.

The system demonstrates core backend engineering patterns including:

- Service-to-service communication
- Centralized logging
- Email service integration
- Authentication with PostgreSQL
- NoSQL logging with MongoDB
- Dockerized infrastructure
- Health checks between services

This branch represents the **stable HTTP-based architecture** before migrating to an event-driven system using RabbitMQ.

---

## Architecture

```
Client (Frontend)
        |
        v
   Broker Service
        |
  -------------------------
  |          |           |
Auth     Logger      Mail
(Postgres) (MongoDB) (SMTP/MailHog)
```

### Communication Style

All services communicate via:

👉 **REST APIs (HTTP/JSON)**

Example:

- Auth → Logger (log authentication events)
- Auth → Mail (send welcome emails)
- Frontend → Mail (service health check)

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

**Why it exists:**

- Prevent frontend from calling multiple services directly
- Centralize request handling
- Improve security boundaries
- Simplify future scaling

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

## Known Limitations

Synchronous communication introduces:

- Tight coupling between services
- Retry complexity
- Latency stacking
- Reduced fault tolerance

---

## Next Evolution

The system will migrate toward:

👉 **Event-Driven Architecture using RabbitMQ**

Goals:

- Asynchronous workflows
- Improved resilience
- Better scalability
- Reduced service dependency

This HTTP branch is preserved as a **baseline architecture**.

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

## Future Improvements

- RabbitMQ integration
- Event-driven workflows
- gRPC for internal communication
- Kubernetes deployment
- Distributed tracing
- Circuit breakers
- Rate limiting

---

## Author Notes

This repository reflects an intentional progression from:

**Monolith Thinking → Microservices → Event-Driven Systems**

The goal is not just building services — but understanding **system design evolution**.
