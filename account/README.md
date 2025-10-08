# 🧩 Account Microservice

This microservice provides account registration, authentication (login), and account retrieval functionalities using **gRPC** and **PostgreSQL**.
It is designed as a modular component for the `orderStream` system but can also run independently for testing and development.

---

## 🚀 Features

* Register a new account
* Login with email and password
* Fetch account details by ID
* List all accounts (with pagination support)
* PostgreSQL as database
* gRPC-based service interface with reflection enabled

---

## 🏗️ Project Structure

```
account/
├── internal/           # Database and repository logic
│   └── repo.go
├── models/             # Data models
│   └── account.go
├── proto/              # Protobuf definitions and generated Go code
│   └── pb/
├── client/             # gRPC client for interacting with the service
│   └── client.go
├── service/            # gRPC server implementation
│   └── account_service.go
├── main.go             # Entry point for the service
├── go.mod / go.sum
└── .env                # Environment variables
```

---

## ⚙️ Prerequisites

Before running the service, ensure the following are installed:

* **Go** ≥ 1.21
* **PostgreSQL**
* **grpcurl** (for manual testing)

---

## 🧾 Environment Setup

Create a `.env` file in the project root:

```bash
DATABASE_URL=postgres://postgres:postgres@localhost:5432/order_stream?sslmode=disable
GRPC_PORT=8080
SECRET_KEY="secret_key"
ISSUER="order-stream"
```

---

## 🗄️ Database Setup

Start PostgreSQL locally, then create the database and `accounts` table:

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

---

## ▶️ Run the Service

1. Install dependencies:

   ```bash
   go mod tidy
   ```

2. Start the gRPC server:

   ```bash
   go run account/cmd/account/main.go
   ```

3. You should see logs like:

   ```
   gRPC server listening on port 8080(something like that)
   Connected to PostgreSQL at localhost:5432
   ```

---

## 🧪 Test the API with grpcurl

### 1️⃣ Register a new account

```bash
grpcurl -plaintext -d '{"name":"Abhishek","email":"abhishek.work71@gmail.com","password":"123456"}' localhost:8080 pb.AccountService/Register
```

### 2️⃣ Login

```bash
grpcurl -plaintext -d '{"email":"abhishek.work71@gmail.com","password":"123456"}' localhost:8080 pb.AccountService/Login
```

### 3️⃣ Get account by ID

```bash
grpcurl -plaintext -d '1' localhost:8080 pb.AccountService/GetAccount
```

### 4️⃣ List all accounts

```bash
grpcurl -plaintext -d '{"skip":0,"take":10}' localhost:8080 pb.AccountService/GetAccounts
```

---

## 🧩 Notes

* Reflection is enabled — you can list all available services:

  ```bash
  grpcurl -plaintext localhost:8080 list
  ```
* JWT signing uses the provided EC private key (`SECRET_KEY`).
* For production, use environment variables securely and connect to a managed PostgreSQL instance.

---

## 📦 Example Response

**Register:**

```json
{
  "value": "Account registered successfully"
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

## 🧰 Tech Stack

* **Language:** Go
* **Framework:** gRPC
* **Database:** PostgreSQL
* **Protocol Buffers:** v3
* **Authentication:** JWT (HS256)

---
