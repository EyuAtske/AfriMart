# AfriMart — Copilot Instructions

## 1. Project Overview

AfriMart is a marketplace application organized as a monorepo containing two independently deployed applications:

* `frontend/` — Nuxt.js frontend application
* `backend/` — Go backend application

The frontend communicates with the backend through the backend's API Gateway.

The backend contains several functional areas such as authentication, shops, products, orders, payments, AI, and analytics. These are **not independently deployed services**. They are components/modules of the single Go backend application.

The repository should therefore be treated as:

```text
AfriMart
├── frontend/              # Independently deployed Nuxt application
└── backend/               # Independently deployed Go application
    ├── cmd/               # Backend entry points
    └── internal/          # Backend application modules
```

Do not assume that every directory under `backend/cmd` or `backend/internal` represents an independently deployed microservice.

---

## 2. Current Architecture

The high-level architecture is:

```text
Nuxt Frontend
       │
       ▼
 API Gateway
       │
       ▼
 Go Backend
 ┌─────┼──────────────────────────────┐
 │     │       │       │       │      │
Auth  Shop  Product  Order  Payment  AI
 │                                      │
 └──────────────────────────────────────┤
                                        │
                                   Analytics
```

The backend is one deployable Go application.

The repository currently contains backend entry-point directories including:

```text
backend/cmd/
├── ai/
├── analytics/
├── auth/
├── gateway/
├── order/
├── payment/
├── product/
└── shop/
```

These should be understood as backend application entry points/components, not separate production deployments.

---

## 3. Technology Stack

### Frontend

The frontend uses:

* Nuxt.js
* TypeScript
* SSR
* Client-side state management
* API client communicating with the backend API Gateway

### Backend

The backend uses:

* Go 1.25+
* PostgreSQL
* JWT-based authentication
* Password hashing
* SQL migrations
* SQL queries
* `sqlc` for type-safe generated database access where applicable
* RabbitMQ for asynchronous messaging where required
* API Gateway for frontend-facing API access

Current Go module:

```text
github.com/EyuAtske/AfriMart/backend
```

---

## 4. Backend Structure

The backend currently follows this general structure:

```text
backend/
├── cmd/
│   ├── ai/
│   ├── analytics/
│   ├── auth/
│   ├── gateway/
│   ├── order/
│   ├── payment/
│   ├── product/
│   └── shop/
│
├── config/
│
├── internal/
│   ├── ai/
│   ├── analytics/
│   ├── auth/
│   ├── commErr/
│   ├── database/
│   ├── handlers/
│   ├── order/
│   ├── payment/
│   ├── product/
│   └── shop/
│
└── go.mod
```

Keep backend business logic inside appropriate `internal` packages.

Do not unnecessarily move business logic into `cmd`.

`cmd` should primarily contain application entry points and wiring.

---

## 5. Functional Backend Areas

The backend supports the following major areas.

### Authentication

Authentication includes:

* User registration
* User login
* Logout/token handling
* Password hashing
* JWT/session authentication
* Role/permission handling
* User profile functionality
* Input validation
* Authorization middleware
* Protected routes

Authentication-related code belongs primarily under:

```text
backend/internal/auth/
```

Do not bypass authentication or authorization middleware for protected endpoints.

---

### Shop

The marketplace allows registered users to act as sellers and create/manage shops.

Shop functionality includes:

* Creating shops
* Getting shop information
* Updating shops
* Deactivating/deleting shops
* Shop ownership authorization
* Seller permissions
* Shop status

Shop-related logic belongs under:

```text
backend/internal/shop/
```

A user must not be allowed to modify another user's shop unless the authorization rules explicitly permit it.

---

### Products

Product functionality includes:

* Product creation
* Product updates
* Product deletion
* Product images
* Pricing
* Categories
* Stock management
* Product status
* Product ownership
* Product search
* Product filtering
* Product pagination
* Product details

Product-related logic belongs under:

```text
backend/internal/product/
```

Product ownership and seller authorization must be enforced server-side.

---

### Orders and Checkout

Order functionality includes:

* Cart operations
* Quantity updates
* Stock validation
* Checkout
* Delivery information
* Order creation
* Order history
* Order details
* Buyer orders
* Seller orders
* Order status
* Order cancellation
* Basic tracking

Order-related logic belongs under:

```text
backend/internal/order/
```

Do not trust prices, quantities, ownership, or stock values supplied by the client without server-side validation.

---

### Payments

Payment functionality includes:

* Payment initialization
* Payment verification
* Transaction records
* Payment status
* Cash-on-delivery where supported
* Payment failure handling

Payment-related logic belongs under:

```text
backend/internal/payment/
```

Never expose payment secrets, credentials, or sensitive payment information in API responses, logs, commits, or error messages.

---

### AI

The AI component provides marketplace assistance by allowing users to describe what they are looking for and receiving relevant marketplace products.

AI functionality may include:

* AI provider/model integration
* Prompt handling
* Product search integration
* Product recommendation logic
* Error handling
* Rate limiting
* Safe logging

AI-related logic belongs under:

```text
backend/internal/ai/
```

Do not expose API keys or other AI provider credentials.

Do not log sensitive user information or credentials.

---

### Analytics

Analytics functionality includes platform statistics such as:

* User statistics
* Product statistics
* Order statistics
* Sales statistics
* Popular products

Analytics-related logic belongs under:

```text
backend/internal/analytics/
```

Analytics should not unnecessarily expose sensitive user information.

---

## 6. Database Architecture

The current project uses PostgreSQL.

The project has intentionally moved away from treating each backend functional area as having its own physical database.

Use the shared PostgreSQL database with clear table/data ownership boundaries.

Current database-related files include:

```text
sql/
├── queries/
└── schema/
```

Migrations are stored under:

```text
sql/schema/
```

Queries are stored under:

```text
sql/queries/
```

`sqlc.yaml` defines SQL code-generation configuration.

### Database Rules

* Use parameterized queries.
* Never construct SQL using unsafe string concatenation with user input.
* Keep schema changes in migrations.
* Do not silently modify existing migrations that may already have been applied.
* Add a new migration for schema changes.
* Keep SQL queries organized by domain.
* Do not store passwords in plaintext.
* Do not expose database credentials.
* Do not commit `.env` files or secrets.

---

## 7. Authentication and Security

Security is a core requirement.

When modifying authentication or protected functionality, always consider:

* Authentication
* Authorization
* Password hashing
* JWT validation
* Token expiration
* Refresh token handling
* User ownership
* Role/permission checks
* Input validation
* SQL injection
* Sensitive data exposure
* Error information leakage

Never:

* Hard-code secrets.
* Commit JWT secrets.
* Commit database passwords.
* Commit RabbitMQ credentials.
* Log passwords or tokens.
* Return password hashes to clients.
* Trust user IDs supplied by clients for authorization decisions.
* Disable authorization merely to make a test pass.

Environment-specific secrets must come from environment variables or an appropriate secret-management mechanism.

---

## 8. API Design

The backend provides APIs consumed by the Nuxt frontend through the API Gateway.

When creating or modifying an endpoint:

* Follow existing routing conventions.
* Validate request input.
* Authenticate protected endpoints.
* Authorize access to resources.
* Return appropriate HTTP status codes.
* Use the project's common error-response format.
* Avoid leaking implementation details in errors.
* Keep request and response structures predictable.

Do not introduce a new response format when an existing project convention already exists.

---

## 9. Synchronous vs Asynchronous Communication

Use synchronous communication when the client or another component requires an immediate response.

Use asynchronous communication through RabbitMQ when an operation is event-driven and does not require an immediate response.

Do not introduce RabbitMQ simply because an operation exists.

When introducing an event:

* Define the event clearly.
* Define its payload.
* Consider retries.
* Consider duplicate delivery.
* Consider failure handling.
* Avoid putting sensitive information into event payloads unnecessarily.

---

## 10. Error Handling

The backend has a common error-handling area:

```text
backend/internal/commErr/
```

Follow existing error-handling conventions.

Errors returned to clients should:

* Be useful to the client.
* Avoid exposing internal implementation details.
* Avoid stack traces in production responses.
* Avoid database credentials or SQL internals.
* Use consistent status codes and response structures.

Do not silently ignore errors.

Handle errors explicitly and return or propagate them appropriately.

---

## 11. Go Coding Standards

Follow idiomatic Go.

Prefer:

* Small focused functions.
* Explicit error handling.
* Clear package responsibilities.
* Dependency injection where appropriate.
* Context propagation.
* Meaningful names.
* Minimal interfaces.
* Standard library functionality where sufficient.

Avoid:

* Global mutable state.
* Unnecessary abstractions.
* Over-engineering.
* Large functions with unrelated responsibilities.
* `panic` for normal application errors.
* Ignoring returned errors.

Run formatting on modified Go files.

Use:

```bash
gofmt
```

and verify the backend with:

```bash
go test ./...
go build ./...
```

---

## 12. Testing Requirements

Testing should happen continuously during development.

The backend should eventually contain:

* Unit tests
* API tests
* Database tests
* Authentication tests
* Authorization tests
* Integration tests
* Message broker tests
* Payment integration tests
* Failure tests
* Performance/load tests

Currently, not every area necessarily has tests.

When adding or modifying functionality, add tests where practical.

A bug fix should preferably include a regression test.

Do not remove or weaken an existing test merely to make CI pass.

---

## 13. CI Requirements

Backend CI must verify that the entire Go backend builds and tests successfully.

The current backend CI checks are:

```text
backend-test
backend-build
```

The expected commands are:

```bash
cd backend
go test ./...
go build ./...
```

The backend should be treated as one application.

Do not create separate CI requirements for `auth`, `product`, `order`, `payment`, `ai`, or `analytics` as though they were independently deployed services unless the architecture is explicitly changed in the future.

---

## 14. Frontend Integration

The frontend is a separate deployable application.

Architecture:

```text
Nuxt Frontend
      │
      ▼
API Gateway
      │
      ▼
Go Backend
```

Backend changes that affect API contracts should consider frontend compatibility.

When changing an API:

* Check existing consumers.
* Preserve backward compatibility where practical.
* Update API documentation when necessary.
* Clearly communicate breaking changes.

Do not place frontend-specific logic inside backend domain packages.

---

## 15. Environment Configuration

Environment-specific configuration belongs in environment variables.

Examples include:

* Database URL
* JWT secret
* RabbitMQ URL
* AI provider credentials
* Payment provider credentials
* API URLs

Never hard-code these values into application source code.

Never commit real secrets.

Use `.env.example` or equivalent documentation for required configuration names without exposing real values.

---

## 16. Docker and Local Development

The repository contains:

```text
docker-compose.yaml
```

Docker Compose is used for local infrastructure and development.

When modifying infrastructure:

* Preserve existing service configuration unless there is a clear reason to change it.
* Avoid committing credentials.
* Ensure environment variables are configurable.
* Keep local development reproducible.

Do not assume production infrastructure is identical to local Docker Compose.

---

## 17. Code Review Rules

When reviewing code, prioritize issues in this order:

1. Security vulnerabilities
2. Authentication/authorization problems
3. Data corruption or incorrect business logic
4. Race conditions/concurrency problems
5. API contract breakage
6. Database correctness
7. Error handling
8. Test failures or missing important tests
9. Performance problems
10. Maintainability/style

Do not raise style comments as blocking issues when the code is correct and consistent with the project.

Every review finding should explain:

* What is wrong.
* Why it matters.
* Where it occurs.
* How it can be fixed.

Avoid speculative criticisms that are not supported by the code or project requirements.

---

## 18. Do Not Make Architectural Assumptions

Before introducing a major architectural change, verify the existing repository structure and documentation.

In particular, do not assume:

* `backend/cmd/*` are independent microservices.
* Each backend domain requires its own database.
* Every operation requires RabbitMQ.
* The frontend communicates directly with internal backend packages.
* Authentication can be bypassed for convenience.
* A new dependency is necessary when the standard library or existing project dependency is sufficient.

The current architecture consists of:

```text
One frontend application
        +
One backend application
        +
Shared PostgreSQL database
        +
RabbitMQ for asynchronous communication where appropriate
```

---

## 19. Project Development Priorities

The project development plan prioritizes the marketplace functionality before later platform capabilities.

Major functionality includes:

### Phase 1 — Core Marketplace

* Authentication and users
* Shop management
* Product management
* Cart and checkout
* Payments
* Orders

### Phase 2 — Backend Organization and Communication

* API Gateway
* Backend functional separation
* Message broker
* Event definitions
* Service/component communication
* Failure handling

### Phase 3 — Administration, Analytics and Observability

* Admin functionality
* Analytics
* Centralized logging
* Error tracking
* Metrics
* Performance monitoring
* Distributed tracing

### Phase 4 — AI Shopping Assistant

* AI service/component
* Product search integration
* Recommendations
* Rate limiting
* Error handling

### Phase 5 — Testing and CI/CD

* Unit testing
* Integration testing
* API testing
* E2E testing
* Failure testing
* Automated CI/CD

Do not prematurely implement future features when they are not required by the current task.

---

## 20. Observability

The project plans to support:

* Centralized logging
* Error tracking
* Metrics
* Performance monitoring
* Distributed tracing

When observability is implemented, avoid logging:

* Passwords
* JWTs
* Refresh tokens
* Database credentials
* Payment credentials
* AI provider credentials
* Other sensitive user information

Prefer structured logs and useful request/service context.

---

## 21. AI-Specific Development Rules

When working on AI functionality:

* Keep provider credentials outside source code.
* Do not expose API keys to the frontend.
* Validate user input.
* Apply appropriate rate limiting.
* Avoid unnecessarily sending sensitive information to external AI providers.
* Do not rely on AI output as authoritative authorization or payment logic.
* Product availability, price, stock, permissions, and payment state must be validated by normal backend logic.
* AI recommendations should ultimately resolve against real marketplace product data.

AI should assist the marketplace; it must not become the source of truth for transactional data.

---

## 22. General Copilot Behavior

When generating or modifying code:

1. Inspect the existing implementation first.
2. Follow existing project conventions.
3. Make the smallest reasonable change.
4. Avoid unnecessary dependencies.
5. Do not rewrite unrelated code.
6. Preserve existing behavior unless the task explicitly requires changing it.
7. Add or update tests for changed behavior.
8. Check error handling.
9. Check authentication and authorization for protected functionality.
10. Check for sensitive information exposure.
11. Run relevant tests/build commands when possible.
12. Explain architectural changes when they are necessary.

When uncertain about the intended architecture, prefer the existing repository structure and documented project decisions over assumptions.

Do not invent APIs, database tables, services, environment variables, or infrastructure components that do not exist.

---

## 23. Definition of Done

A backend change is considered complete when, where applicable:

* The implementation follows the existing architecture.
* Business logic is placed in the appropriate backend package.
* Authentication and authorization are enforced.
* Input is validated.
* Errors are handled consistently.
* Database access is safe.
* Tests cover the important behavior.
* Existing tests still pass.
* The backend builds successfully.
* No secrets are introduced.
* API contracts remain compatible or changes are documented.
* CI checks pass.

The minimum backend CI requirements are:

```bash
cd backend
go test ./...
go build ./...
```

These correspond to the required GitHub status checks:

```text
backend-test
backend-build
```
