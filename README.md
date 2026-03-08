# Fullstack E‑Commerce Microservices Demo (Go + React + Kafka)

A production‑grade, microservices‑based e‑commerce demo with independent per‑service SQLite databases, KrakenD API gateway, Kafka events, RBAC security, and a modern React frontend.

## Architecture Overview

- Gateway: KrakenD fronts all REST (`/api/v1`), forwards `Authorization` header, preserves query, health at `/__health`
- Services (independent, per‑service SQLite): auth, product, cart, order, payment, notification, chat
- Events: Kafka topics per domain (`user.*`, `product.*`, `cart.*`, `order.*`, `payment.*`, `chat.*`)
- Frontend: Nginx serves UI; `/api/` proxied to KrakenD; `/ws` proxied directly to chat service

```
[Frontend] → [KrakenD]
               ├─ /auth → Auth Service
               ├─ /products → Product Service
               ├─ /cart → Cart Service
               ├─ /orders → Order Service
               ├─ /payments → Payment Service
               ├─ /notifications → Notification Service
               └─ /ws → Chat Service (direct WebSocket)

Kafka ← events (product/cart/order/payment/chat)
SQLite ← per service DB under /data
```

## Features

- Auth: register/login/refresh/logout, profile, change password; RBAC roles (admin, seller, customer, moderator, support)
- Addresses: manage shipping addresses (list/create/update/delete; default address support)
- Products: CRUD, pagination/sort, stock updates
- Cart: per‑user cart; add/remove/clear; persisted items
- Orders: checkout with address snapshot; status transitions (pending → paid → processing → shipped → delivered → canceled)
- Payments: Stripe/Razorpay intents and webhooks; order payment state updates; audit records
- Notifications: consumes order/payment topics; user notices (payment success/failure, shipped/delivered)
- Chat: WebSocket user↔admin messaging; publishes `chat.message_sent`
- COD: Cash on Delivery option; creates order with status `processing`

## Databases (SQLite)

- `auth-service` → `/data/auth.db`
- `product-service` → `/data/product.db`
- `cart-service` → `/data/cart.db`
- `order-service` → `/data/order.db`
- `payment-service` → `/data/payment.db`
- `notification-service` → `/data/notification.db`
- `chat-service` → `/data/chat.db`

Driver: pure‑Go `github.com/glebarez/sqlite` (no CGO).

## Kafka Topics

- `product.created|updated|deleted`
- `cart.item_added|item_removed|cleared`
- `order.created|status_changed`
- `payment.processed|failed|refunded`
- `chat.message_sent`

## API (via KrakenD `/api/v1`)

- Auth: `POST /auth/register`, `POST /auth/login`, `POST /auth/refresh`, `POST /auth/logout`, `GET /auth/profile`, `PUT /auth/profile`, `POST /auth/change-password`
- Addresses: `GET /auth/addresses`, `POST /auth/addresses`, `PUT /auth/addresses/:id`, `DELETE /auth/addresses/:id`
- Products: `GET /products`, `GET /products/:id`, `GET /products/sku/:sku`, `POST /products`, `PUT /products/:id`, `DELETE /products/:id`, `PATCH /products/:id/stock`
- Cart: `GET /cart/:user_id`, `DELETE /cart/:user_id`, `POST /cart/:user_id/items`, `DELETE /cart/:user_id/items/:sku`
- Orders: `POST /orders`, `GET /orders?user_id=...`, `GET /orders/:id`, `PATCH /orders/:id/status`, `POST /orders/checkout`
- Payments: `POST /payments/intent` (accepts `order_id`, calculates amount from order)
- Notifications: `GET /notifications?user_id=<uuid>`
- WebSocket: `GET /ws` (direct to chat service via Nginx)

## Health Endpoints

- Each service exposes `GET/HEAD /health` returning 200 when ready.

## Frontend

- React + Vite + TypeScript; Nginx serves `/` and proxies `/api/` to KrakenD and `/ws` to chat
- Pages: Login, Products, Cart, Checkout, Orders, Payments, Chat, Notifications
- Checkout:
  - Address selection/creation
  - Order summary
  - Payment: Pay Online (intent + provider widget) or COD (order becomes `processing`)

Environment (`frontend/.env`):

```
VITE_API_URL=http://krakend:8080/api/v1
VITE_CHAT_URL=http://chat-service:8086
```

For local dev, set `VITE_API_URL=http://localhost:8080/api/v1`.

## Quick Start (Docker Compose)

```
docker compose down
docker compose build --no-cache
docker compose up -d krakend product-service cart-service order-service auth-service notification-service payment-service chat-service kafka zookeeper frontend
```

- Open `http://localhost:3000`
- Gateway health: `http://localhost:8080/__health`

## Example Requests

Create address:

```
POST /api/v1/auth/addresses
{
  "first_name": "John",
  "last_name": "Doe",
  "phone": "1234567890",
  "address1": "123 Main St",
  "address2": "Apt 4",
  "city": "Springfield",
  "state": "IL",
  "country": "USA",
  "postal_code": "62704",
  "is_default": true,
  "type": "shipping"
}
```

Checkout:

```
POST /api/v1/orders/checkout
{
  "user_id": "<uuid>",
  "currency": "USD",
  "payment_method": "online", // or "cod"
  "items": [
    { "sku": "SKU-HEAD-001", "name": "Wireless Headphones", "unit_price": 8999, "quantity": 1 }
  ],
  "shipping": {
    "name": "John Doe",
    "phone": "1234567890",
    "address1": "123 Main St",
    "address2": "Apt 4",
    "city": "Springfield",
    "state": "IL",
    "country": "USA",
    "postal": "62704"
  }
}
```

Payment intent:

```
POST /api/v1/payments/intent
{
  "order_id": "<uuid>",
  "provider": "stripe"
}
```

## RBAC & Security

- JWT HMAC (`HS256`) with claims `{ user_id, email, roles }`
- Protected routes enforce auth; admin routes enforce role `admin`
- Avoid secrets in code; use environment variables

## Troubleshooting

- Ensure Docker has network access for base images and go modules
- Kafka/Zookeeper should be healthy before event‑dependent services fully start
- If a service fails to start, check its `/health` endpoint and logs

## Project Layout

```
fullstack-demo-app/
  frontend/
  microservices/
    krakend/
    auth-service/
    product-service/
    cart-service/
    order-service/
    payment-service/
    chat-service/
    notification-service/
  docker-compose.yml
```

## License

MIT (for demo purposes)