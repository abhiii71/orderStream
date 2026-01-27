# OrderStream - Microservices E-Commerce Platform

A modern, event-driven e-commerce platform built with microservices architecture using Go, Python, GraphQL, gRPC, Kafka, and multiple databases.

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Client Applications                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                       │
                                       ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         GraphQL Gateway (Port 8080)                          │
│                    - Authentication (JWT)                                    │
│                    - Request Routing                                         │
│                    - API Aggregation                                         │
└─────────────────────────────────────────────────────────────────────────────┘
                                       │
                    ┌──────────────────┼──────────────────┐
                    │                  │                  │
                    ▼                  ▼                  ▼
┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│   Account   │ │   Product   │ │    Order    │ │   Payment   │ │ Recommender │
│   Service   │ │   Service   │ │   Service   │ │   Service   │ │   Service   │
│   (gRPC)    │ │   (gRPC)    │ │   (gRPC)    │ │   (gRPC)    │ │   (gRPC)    │
│             │ │             │ │             │ │             │ │   Python    │
└──────┬──────┘ └──────┬──────┘ └──────┬──────┘ └──────┬──────┘ └──────┬──────┘
       │               │               │               │               │
       ▼               ▼               ▼               ▼               ▼
┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│ PostgreSQL  │ │Elasticsearch│ │ PostgreSQL  │ │ PostgreSQL  │ │ PostgreSQL  │
│ account_db  │ │ product_db  │ │  order_db   │ │ payment_db  │ │recommender_db│
└─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘
                                       │
                                       ▼
                    ┌─────────────────────────────────────┐
                    │         Apache Kafka                │
                    │   (Event Streaming Platform)        │
                    │   - product_events                  │
                    │   - interaction_events              │
                    └─────────────────────────────────────┘
```

## 📦 Services Description

### 1. **Account Service** (Go)
- **Port**: 8080 (internal gRPC)
- **Database**: PostgreSQL
- **Responsibilities**:
  - User registration with password hashing
  - User login with JWT token generation
  - Account retrieval by ID or list

### 2. **Product Service** (Go)
- **Port**: 8080 (internal gRPC)
- **Database**: Elasticsearch
- **Responsibilities**:
  - Product CRUD operations
  - Full-text search for products
  - Publishes product events to Kafka

### 3. **Order Service** (Go)
- **Port**: 8080 (internal gRPC)
- **Database**: PostgreSQL
- **Responsibilities**:
  - Create orders with multiple products
  - Retrieve orders for an account
  - Update order payment status
  - Publishes purchase events to Kafka for recommendations

### 4. **Payment Service** (Go)
- **Ports**: 8080 (gRPC), 8081 (Webhook)
- **Database**: PostgreSQL
- **Responsibilities**:
  - Customer management
  - Checkout session creation (Dodo Payments integration)
  - Payment webhook handling
  - Consumes product events from Kafka

### 5. **Recommender Service** (Python)
- **Port**: 8080 (internal gRPC)
- **Database**: PostgreSQL
- **Responsibilities**:
  - ML-based product recommendations using SVD algorithm
  - Consumes interaction events from Kafka
  - Provides personalized recommendations based on user history

### 6. **GraphQL Gateway** (Go)
- **Port**: 8080 (exposed)
- **Responsibilities**:
  - Single entry point for all client requests
  - JWT authentication
  - Request routing to appropriate microservices
  - API aggregation and response formatting

## 🚀 Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.24+ (for local development)
- Python 3.11+ (for recommender service development)

### Quick Start

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd orderStream
   ```

2. **Build the base image first**
   ```bash
   docker-compose build base
   ```

3. **Start all services**
   ```bash
   docker-compose up -d
   ```

4. **Run database migrations**
   ```bash
   # Account DB
   docker exec -i account_db psql -U abhiii71 -d abhiii71 < account/db/migrations/000001_create_accounts_table.up.sql

   # Order DB
   docker exec -i order_db psql -U abhiii71 -d abhiii71 < order/db/migrations/000002_create_orders_table.up.sql
   docker exec -i order_db psql -U abhiii71 -d abhiii71 < order/db/migrations/000003_create_order_products_table.up.sql

   # Payment DB
   docker exec -i payment_db psql -U abhiii71 -d abhiii71 < payment/db/migrations/000004_create_customers_table.up.sql
   docker exec -i payment_db psql -U abhiii71 -d abhiii71 < payment/db/migrations/000005_create_transactions_table.up.sql
   ```

5. **Verify all services are running**
   ```bash
   docker-compose ps
   ```

6. **Access GraphQL Playground**
   - Open: http://localhost:8080/playground

## 🧪 Testing the API

### GraphQL Playground
Access the interactive GraphQL playground at: `http://localhost:8080/playground`

### Authentication Flow

#### 1. Register a new user
```graphql
mutation {
  register(account: {
    name: "John Doe"
    email: "john@example.com"
    password: "securepassword123"
  }) {
    token
  }
}
```

#### 2. Login
```graphql
mutation {
  login(account: {
    email: "john@example.com"
    password: "securepassword123"
  }) {
    token
  }
}
```

#### 3. Use the token
Add the token to HTTP headers:
```json
{
  "Authorization": "Bearer <your-jwt-token>"
}
```

### Product Operations (Requires Authentication)

#### Create a Product
```graphql
mutation {
  createProduct(product: {
    name: "Wireless Headphones"
    description: "High-quality Bluetooth headphones with noise cancellation"
    price: 149.99
  }) {
    id
    name
    description
    price
  }
}
```

#### Query Products
```graphql
query {
  product(pagination: {skip: 0, take: 10}) {
    id
    name
    description
    price
  }
}
```

#### Search Products
```graphql
query {
  product(pagination: {skip: 0, take: 10}, query: "headphones") {
    id
    name
    description
    price
  }
}
```

#### Update a Product
```graphql
mutation {
  updateProduct(product: {
    id: "<product-id>"
    name: "Updated Headphones"
    description: "Updated description"
    price: 129.99
  }) {
    id
    name
    description
    price
  }
}
```

#### Delete a Product
```graphql
mutation {
  deleteProduct(id: "<product-id>")
}
```

### Order Operations (Requires Authentication)

#### Create an Order
```graphql
mutation {
  createOrder(order: {
    products: [
      {id: "<product-id-1>", quantity: 2}
      {id: "<product-id-2>", quantity: 1}
    ]
  }) {
    id
    totalPrice
    createdAt
    products {
      id
      name
      price
      quantity
    }
  }
}
```

### Account Operations (Requires Authentication)

#### Get All Accounts
```graphql
query {
  accounts(pagination: {skip: 0, take: 10}) {
    id
    name
    email
  }
}
```

#### Get Single Account with Orders
```graphql
query {
  accounts(id: 1) {
    id
    name
    email
    Orders {
      id
      totalPrice
      createdAt
      products {
        id
        name
        price
        quantity
      }
    }
  }
}
```

### Payment Operations

#### Create Customer Portal Session
```graphql
mutation {
  createCustomerPortalSession(credentials: {
    accounntId: 1
    email: "john@example.com"
    name: "John Doe"
  }) {
    url
  }
}
```

#### Create Checkout Session
```graphql
mutation {
  createCheckoutSession(details: {
    accounId: 1
    name: "John Doe"
    email: "john@example.com"
    redirectUrl: "http://localhost:3000/success"
    products: [
      {id: "<product-id>", quantity: 2}
    ]
    orderId: 1
  }) {
    url
  }
}
```

## 🔧 Testing Individual Microservices (gRPC)

You can test individual microservices using `grpcurl`:

### Install grpcurl
```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

### Test Account Service
```bash
# List services
docker run --rm --network orderstream_app-network fullstorydev/grpcurl -plaintext account:8080 list

# Register
docker run --rm --network orderstream_app-network fullstorydev/grpcurl -plaintext \
  -d '{"name":"Test User", "email":"test@example.com", "password":"password123"}' \
  account:8080 pb.AccountService/Register

# Login
docker run --rm --network orderstream_app-network fullstorydev/grpcurl -plaintext \
  -d '{"email":"test@example.com", "password":"password123"}' \
  account:8080 pb.AccountService/Login

# Get Account
docker run --rm --network orderstream_app-network fullstorydev/grpcurl -plaintext \
  -d '"1"' \
  account:8080 pb.AccountService/GetAccount
```

### Test Product Service
```bash
# Create Product
docker run --rm --network orderstream_app-network fullstorydev/grpcurl -plaintext \
  -d '{"name":"Test Product", "description":"A test product", "price":99.99, "accountId":1}' \
  product:8080 pb.ProductService/PostProduct

# Get Products
docker run --rm --network orderstream_app-network fullstorydev/grpcurl -plaintext \
  -d '{"skip":0, "take":10}' \
  product:8080 pb.ProductService/GetProducts
```

### Test Order Service
```bash
# Create Order
docker run --rm --network orderstream_app-network fullstorydev/grpcurl -plaintext \
  -d '{"accountId":1, "products":[{"id":"<product-id>", "quantity":2}]}' \
  order:8080 pb.OrderService/PostOrder

# Get Orders for Account
docker run --rm --network orderstream_app-network fullstorydev/grpcurl -plaintext \
  -d '"1"' \
  order:8080 pb.OrderService/GetOrdersForAccount
```

## 📊 Database Schemas

### Account Service (PostgreSQL)
```sql
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Order Service (PostgreSQL)
```sql
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL,
    total_price DOUBLE PRECISION NOT NULL,
    payment_status TEXT DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE order_products (
    id SERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id VARCHAR(255) NOT NULL,
    quantity INT NOT NULL
);
```

### Payment Service (PostgreSQL)
```sql
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    customer_id VARCHAR(255) UNIQUE NOT NULL,
    billing_email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    customer_id VARCHAR(255) NOT NULL,
    payment_id VARCHAR(255),
    total_price BIGINT NOT NULL,
    settled_price BIGINT DEFAULT 0,
    currency VARCHAR(10) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## 🔄 Event Flow (Kafka)

### Product Events
When a product is created/updated/deleted:
```
Product Service → Kafka (product_events) → Payment Service
```

### Interaction Events
When an order is placed:
```
Order Service → Kafka (interaction_events) → Recommender Service
```

## 🛠️ Development

### Project Structure
```
orderStream/
├── account/                 # Account microservice
│   ├── client/             # gRPC client
│   ├── cmd/account/        # Main entry point
│   ├── config/             # Configuration
│   ├── db/migrations/      # Database migrations
│   ├── internal/           # Business logic
│   ├── models/             # Data models
│   └── proto/              # Protobuf definitions
├── product/                 # Product microservice (similar structure)
├── order/                   # Order microservice (similar structure)
├── payment/                 # Payment microservice (similar structure)
├── recommender/             # Recommender microservice (Python)
│   ├── app/
│   │   ├── db/             # Database models
│   │   ├── entry/          # Entry points
│   │   └── services/       # Business logic
│   ├── config/             # Configuration
│   └── generated/          # Generated protobuf files
├── graphql/                 # GraphQL Gateway
│   ├── cmd/graphql/        # Main entry point
│   ├── config/             # Configuration
│   ├── generated/          # Generated GraphQL code
│   ├── graph/              # Resolvers
│   └── schema.graphql      # GraphQL schema
├── pkg/                     # Shared packages
│   ├── auth/               # JWT authentication
│   ├── contextkeys/        # Context keys
│   ├── crypt/              # Password hashing
│   ├── kafka/              # Kafka producer/consumer
│   └── middleware/         # HTTP middleware
├── docker/                  # Dockerfiles
└── docker-compose.yml       # Docker Compose configuration
```

### Regenerating Protobuf Files

For Go services:
```bash
protoc --go_out=. --go-grpc_out=. proto/*.proto
```

For Python (Recommender):
```bash
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. recommender.proto
```

### Regenerating GraphQL Code
```bash
cd graphql
go run github.com/99designs/gqlgen generate
```

## 🐳 Docker Commands

```bash
# Build all services
docker-compose build

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f graphql

# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v

# Rebuild specific service
docker-compose up -d --build <service-name>
```

## 🔐 Environment Variables

### GraphQL Gateway
| Variable | Description |
|----------|-------------|
| ACCOUNT_URL | Account service URL |
| PRODUCT_URL | Product service URL |
| ORDER_URL | Order service URL |
| PAYMENT_URL | Payment service URL |
| RECOMMENDER_URL | Recommender service URL |

### Account/Order/Payment Services
| Variable | Description |
|----------|-------------|
| DATABASE_URL | PostgreSQL connection string |

### Product Service
| Variable | Description |
|----------|-------------|
| ELASTICSEARCH_URL | Elasticsearch URL |
| KAFKA_BOOTSTRAP_SERVERS | Kafka broker address |

### Payment Service
| Variable | Description |
|----------|-------------|
| DATABASE_URL | PostgreSQL connection string |
| ORDER_SERVICE_URL | Order service URL |
| KAFKA_BOOTSTRAP_SERVERS | Kafka broker address |
| DODO_API_KEY | Dodo Payments API key |
| DODO_WEBHOOK_SECRET | Webhook secret |
| DODO_TEST_MODE | Enable test mode |

## 📝 API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check |
| `/graphql` | POST | GraphQL API |
| `/playground` | GET | GraphQL Playground |

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📄 License

This project is licensed under the MIT License.

