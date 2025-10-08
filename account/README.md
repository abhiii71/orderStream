---

# Account Microservice

This microservice provides account registration, authentication (login), and account retrieval functionalities using **gRPC** and **PostgreSQL**.
It is designed as a modular component for the `orderStream` system but can also run independently for testing and development.

---

## Features

* Register a new account
* Login with email and password
* Fetch account details by ID
* List all accounts with pagination support
* PostgreSQL as the database
* gRPC-based service interface with reflection enabled

---

## Project Structure

```
account/
├── client/             # gRPC client for interacting with the service
│   └── client.go
├── cmd/                # Main entry point
│   └── account/main.go
├── config/             # Environment variables and config
│   └── config.go
├── constants.go        # Shared constants
├── db/
│   └── migrations/     # SQL migration files
│       ├── 000001_create_accounts_table.up.sql
│       └── 000001_create_accounts_table.down.sql
├── internal/           # Service, server, and repository logic
│   ├── repository.go
│   ├── server.go
│   └── service.go
├── models/             # Data models
│   └── account.go
├── proto/              # Protobuf definitions
│   ├── account.proto
│   └── pb/
│       ├── account.pb.go
│       └── account_grpc.pb.go
├── README.md           # This file
└── tests/              # Unit tests
    ├── repo_test.go
    └── svc_test.go
```

---

## Prerequisites

* Go ≥ 1.21
* PostgreSQL
* grpcurl (for manual testing)

---

## Environment Setup

Create a `.env` file in the project root:

```bash
DATABASE_URL=postgres://postgres:postgres@localhost:5432/order_stream?sslmode=disable
GRPC_PORT=8080
SECRET_KEY="your_secret_key_here"
ISSUER="order-stream"
```

---

## Database Setup

### Option 1: Run SQL manually

```sql
CREATE DATABASE order_stream;

\c order_stream;

CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Option 2: Use migrations (recommended)

If you have migration files in the `migrations/` directory:

```bash
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/order_stream?sslmode=disable" up
```

---

## Run the Service

1. Install dependencies:

```bash
go mod tidy
```

2. Start the gRPC server:

```bash
go run account/cmd/account/main.go
```

3. Expected logs:

```
gRPC server listening on port 8080
Connected to PostgreSQL at localhost:5432
```

---

## Test the API with grpcurl

### Register a new account

```bash
grpcurl -plaintext -d '{"name":"Abhishek","email":"abhishek.work71@gmail.com","password":"123456"}' localhost:8080 pb.AccountService/Register
```

### Login

```bash
grpcurl -plaintext -d '{"email":"abhishek.work71@gmail.com","password":"123456"}' localhost:8080 pb.AccountService/Login
```

### Get account by ID

```bash
grpcurl -plaintext -d '1' localhost:8080 pb.AccountService/GetAccount
```

### List all accounts

```bash
grpcurl -plaintext -d '{"skip":0,"take":10}' localhost:8080 pb.AccountService/GetAccounts
```

---

## Notes

* Reflection is enabled, so you can list all available services:

```bash
grpcurl -plaintext localhost:8080 list
```

* JWT signing uses the `SECRET_KEY` from the environment.
* For production, use environment variables securely and connect to a managed PostgreSQL instance.

---

## Example Responses

**Register:**

```json
{
  "value": "1"
}
```

**GetAccount:**

```json
{
  "account": {
    "id": 1,
    "name": "Abhishek",
    "email": "abhishek.work71@gmail.com"
  }
}
```

---

## Tech Stack

* **Language:** Go
* **Framework:** gRPC
* **Database:** PostgreSQL
* **Protocol Buffers:** v3
* **Authentication:** JWT (HS256)

---