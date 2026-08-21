# AfriMart — System Overview

## 1. Purpose

AfriMart is a learning-focused marketplace project designed to apply and demonstrate concepts in system design and system architecture.

The primary purpose of the project is not to build a production-ready marketplace. Instead, the project provides a practical environment for learning and implementing architectural concepts such as:

* Monorepo architecture
* Microservice architecture
* API Gateway patterns
* Service boundaries
* Independent data ownership
* Synchronous and asynchronous communication
* Message-driven architecture
* Authentication and authorization
* Observability
* Automated testing
* CI/CD
* AI-assisted development and code review

The marketplace functionality provides the business context through which these architectural concepts can be implemented and evaluated.

---

## 2. System Goals

The system is designed around the following goals:

1. Provide a working marketplace application.
2. Separate major business capabilities into independent backend services.
3. Provide a single API entry point through an API Gateway.
4. Establish clear ownership of data between services.
5. Support both synchronous API communication and asynchronous event-based communication.
6. Allow services to be developed and deployed independently.
7. Provide observability through logs, metrics, performance monitoring, and distributed tracing.
8. Establish automated testing and CI/CD practices.
9. Provide an AI service capable of assisting users with product discovery.
10. Provide an architecture that can be used to study the trade-offs of distributed systems.

---

## 3. High-Level Architecture

The system consists of a Nuxt frontend, an API Gateway, multiple Go backend services, independent service data stores, and a message broker.

The high-level architecture is:

```text
                         ┌─────────────────────┐
                         │    Nuxt Frontend    │
                         └──────────┬──────────┘
                                    │
                                    │ HTTP
                                    ▼
                         ┌─────────────────────┐
                         │     API Gateway     │
                         └──────────┬──────────┘
                                    │
              ┌─────────────────────┼─────────────────────┐
              │          │          │          │           │
              ▼          ▼          ▼          ▼           ▼
          ┌───────┐  ┌───────┐  ┌───────┐  ┌───────┐  ┌─────────┐
          │ Auth  │  │ Shop  │  │Product│  │ Order │  │ Payment │
          │Service│  │Service│  │Service│  │Service│  │ Service │
          └───┬───┘  └───┬───┘  └───┬───┘  └───┬───┘  └────┬────┘
              │          │          │          │             │
              ▼          ▼          ▼          ▼             ▼
           ┌────┐     ┌────┐     ┌────┐     ┌────┐       ┌────┐
           │ DB │     │ DB │     │ DB │     │ DB │       │ DB │
           └────┘     └────┘     └────┘     └────┘       └────┘

                         ┌─────────────────────┐
                         │    Message Broker   │
                         └──────────┬──────────┘
                                    │
                         ┌──────────┴──────────┐
                         ▼                     ▼
                    ┌─────────┐          ┌───────────┐
                    │   AI    │          │ Analytics │
                    │ Service │          │  Service  │
                    └─────────┘          └───────────┘
```

The project specification defines the initial backend services as Authentication, Shop, Product, Order, Payment, AI, and Analytics.

---

## 4. Major Components

### 4.1 Nuxt Frontend

The frontend is implemented using Nuxt.js with TypeScript.

Its responsibilities include:

* User interface
* Routing
* State management
* API integration
* Authentication state
* Marketplace pages
* Loading and error states
* AI assistant interface
* Administrative interfaces

The frontend communicates with the backend through the API Gateway rather than directly depending on individual backend services.

---

### 4.2 API Gateway

The API Gateway provides the main entry point into the backend system.

The intended request flow is:

```text
Nuxt Frontend
      │
      ▼
 API Gateway
      │
      ▼
Backend Services
```

The Gateway is responsible for routing external API requests to the appropriate backend services.

It also provides a boundary between the frontend and the internal service architecture.

---

### 4.3 Authentication Service

The Authentication Service is responsible for authentication and user-related access control.

Its responsibilities include:

* User registration
* User login
* Logout/token handling
* Password hashing
* Authentication
* Roles and permissions
* User profile operations
* Authorization
* Protected access

The service owns the authentication-related user data.

---

### 4.4 Shop Service

The Shop Service manages marketplace shops and seller-related functionality.

Its responsibilities include:

* Creating shops
* Retrieving shops
* Updating shops
* Deactivating or deleting shops
* Shop ownership
* Seller permissions
* Shop status

---

### 4.5 Product Service

The Product Service manages marketplace products.

Its responsibilities include:

* Product creation
* Product updates
* Product deletion
* Categories
* Product images
* Pricing
* Stock management
* Product status
* Product ownership
* Product search
* Product filtering
* Product pagination
* Product details

---

### 4.6 Order Service

The Order Service manages the purchasing and order lifecycle.

Its responsibilities include:

* Order creation
* Order history
* Order details
* Buyer orders
* Seller orders
* Order status
* Order status updates
* Basic tracking
* Cancellation rules

The Order Service is also responsible for the order portion of the checkout flow.

---

### 4.7 Payment Service

The Payment Service handles payment-related operations.

Its responsibilities include:

* Payment initialization
* Payment verification
* Transaction records
* Payment status
* Cash-on-delivery logic
* Payment failure handling

The exact payment provider is not yet defined and will be selected later.

---

### 4.8 AI Service

The AI Service provides the marketplace's AI shopping assistant.

Its responsibilities include:

* AI provider/model integration
* AI API
* Prompt design
* Product search integration
* Product recommendation logic
* AI error handling
* Rate limiting
* Logging without exposing sensitive information

The intended user experience is that a user describes what they want in natural language and the system returns relevant marketplace products.

---

### 4.9 Analytics Service

The Analytics Service is responsible for collecting and providing marketplace statistics.

The planned analytics include:

* User statistics
* Product statistics
* Order statistics
* Sales statistics
* Popular products

The exact implementation and data-processing approach will be determined during the architecture and implementation phases.

---

### 4.10 Message Broker

The system will use a message broker to support asynchronous communication between services.

Asynchronous communication will be used for operations that do not require an immediate response.

For example:

```text
Order Service
      │
      │ OrderCreated
      ▼
Message Broker
      │
      ├──────────────► Payment Service
      │
      └──────────────► Analytics Service
```

The specific message broker has not yet been finalized.

---

## 5. Communication Model

The system supports two communication models.

### 5.1 Synchronous Communication

Synchronous communication is used when an immediate response is required.

For example:

```text
Frontend
   │
   ▼
API Gateway
   │
   ▼
Product Service
   │
   ▼
Product Database
   │
   ▼
Response
```

The caller waits for the service to process the request and return a response.

---

### 5.2 Asynchronous Communication

Asynchronous communication is used for operations that do not require an immediate response.

The message broker acts as the communication mechanism between event producers and consumers.

```text
Service A
    │
    │ Event
    ▼
Message Broker
    │
    ├────────► Service B
    │
    └────────► Service C
```

This allows services to react to events without requiring direct synchronous communication with every other service.

---

## 6. Data Ownership

Each service is responsible for its own data.

The architecture follows the principle that a service should not directly modify another service's database.

Conceptually:

```text
Auth Service      → Auth/User Data
Shop Service      → Shop Data
Product Service   → Product Data
Order Service     → Order Data
Payment Service   → Payment Data
Analytics Service → Analytics Data
```

If a service requires information owned by another service, communication should occur through an appropriate API or event rather than through direct database access.

This separation is an important part of the microservice architecture defined for the project.

---

## 7. Deployment Model

The project uses a monorepo containing the frontend, backend services, and supporting infrastructure.

The frontend and backend are deployed separately.

The intended deployment structure is:

```text
                    Monorepo
                       │
            ┌──────────┴──────────┐
            │                     │
            ▼                     ▼
       Frontend Deployment   Backend Deployment
            │                     │
          Nuxt              API Gateway
                                  │
                         ┌────────┼────────┐
                         ▼        ▼        ▼
                       Services  Services  ...
```

The Phase 0 specification explicitly defines one monorepo for the project and separate deployment for the Nuxt frontend and Go backend.

---

## 8. Observability

Because the system consists of multiple services, understanding system behavior requires more than examining individual service logs.

The project will therefore include observability capabilities covering:

* Centralized logging
* Error tracking
* Metrics
* Performance monitoring
* Distributed tracing

The goal is to make it possible to understand a request as it travels through multiple components.

For example:

```text
Frontend
   │
   ▼
API Gateway
   │
   ▼
Order Service
   │
   ▼
Message Broker
   │
   ▼
Payment Service
```

Observability should allow developers to identify where failures, latency, and other problems occur across this flow.

The project plan explicitly includes centralized logging, error tracking, metrics, performance monitoring, and distributed tracing.

---

## 9. Architectural Principles

The system will follow these principles:

### Service Independence

Each service should have a clearly defined responsibility and should avoid unnecessary coupling to other services.

### Data Ownership

A service owns and manages its own data. Other services should not directly access or modify its database.

### API Boundary

External clients communicate with the backend through the API Gateway.

### Appropriate Communication

Synchronous communication should be used when an immediate response is required. Asynchronous messaging should be used for operations that can be processed independently.

### Loose Coupling

Services should communicate through well-defined APIs and events rather than relying on internal implementation details.

### Observability

The system should provide sufficient logs, metrics, and traces to understand its behavior.

### Testability

Services and important user flows should be testable independently and together.

### Automation

Testing and deployment should eventually be integrated into CI/CD pipelines.

---

## 10. Out of Scope for the Initial Architecture

The following areas are intentionally not fully defined at this stage:

* Exact cloud provider
* Exact payment provider
* Exact message broker
* Exact database technology for every service
* Production-scale infrastructure
* Advanced recommendation algorithms
* Virtual try-on implementation
* Advanced seller analytics

These decisions can be made later as the project progresses and as their architectural requirements become clearer.

The project plan identifies virtual try-on, advanced recommendations, additional payment methods, and advanced seller analytics as future features rather than requirements for the main marketplace release.

---

## 11. Target Architecture

The overall target architecture can therefore be summarized as:

```text
                         ┌─────────────────┐
                         │  Nuxt Frontend  │
                         └────────┬────────┘
                                  │
                                  ▼
                         ┌─────────────────┐
                         │   API Gateway   │
                         └────────┬────────┘
                                  │
              ┌───────────────────┼───────────────────┐
              │                   │                   │
              ▼                   ▼                   ▼
        ┌──────────┐        ┌──────────┐        ┌──────────┐
        │   Auth   │        │ Product  │        │  Shop    │
        │ Service  │        │ Service  │        │ Service  │
        └────┬─────┘        └────┬─────┘        └────┬─────┘
             │                   │                   │
             ▼                   ▼                   ▼
           Auth DB            Product DB            Shop DB


              ┌───────────────────┼───────────────────┐
              │                   │                   │
              ▼                   ▼                   ▼
        ┌──────────┐        ┌──────────┐        ┌──────────┐
        │  Order   │        │ Payment  │        │    AI    │
        │ Service  │        │ Service  │        │ Service  │
        └────┬─────┘        └────┬─────┘        └──────────┘
             │                   │
             ▼                   ▼
          Order DB           Payment DB


                         ┌─────────────────┐
                         │ Message Broker  │
                         └────────┬────────┘
                                  │
                                  ▼
                         ┌─────────────────┐
                         │    Analytics    │
                         │     Service     │
                         └────────┬────────┘
                                  │
                                  ▼
                             Analytics DB
```

This architecture represents the **target state**, not necessarily the implementation state at the beginning of the project.

The system should evolve toward this architecture while the team uses each stage to understand the underlying system-design concepts.
