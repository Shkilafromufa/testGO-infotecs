# testGO-infotecs


This is a simple wallet and transaction service written in Go, designed for a test assignment. The application allows creating wallets, querying balances, and transferring funds between wallets, with transaction history stored in a SQLite database.

## Features

* **Create and initialize wallets**: Generates a set of wallets with random hashes and starting balances.
* **Get wallet balance**: Query the balance of a specific wallet by its address.
* **Send transactions**: Transfer funds from one wallet to another with validation (no self-transfers, sufficient balance).
* **Transaction history**: Retrieve the latest transactions with a configurable limit.
* **SQLite persistence**: Stores wallet and transaction data in a local SQLite database.

## Prerequisites

* Go 1.18+ installed
* SQLite3

## Project Structure

```text
testCaseGO/
├── internal
│   ├── handler
│   │   └── handler.go          # HTTP routers and handlers
│   ├── model
│   │   ├── transaction.go      # Transaction struct
│   │   └── wallet.go           # Wallet struct
│   └── service
│       ├── database.go         # DB initialization and helpers
│       ├── service.go          # Service functions (Send, GetLast, GetBalance)
│       └── storage.go          # CRUD operations and business logic
├── go.mod
└── main.go                     # Application entry point
```

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/shkilafromufa/testGO-infotecs.git
   cd testGO-infotecs
   ```
2. Download dependencies:

   ```bash
   go mod download
   ```

## Database Setup

The service uses SQLite. By default, the database file is located at `internal/service/database/database.db`. The application will automatically create the required tables and initialize wallets on startup.

If you need to reset the database, simply delete the `database.db` file:

```bash
rm internal/service/database/database.db
```

The next run will recreate the schema and generate 10 wallets with a starting balance of 100 units each.

## Running the Application

Build and run:

```bash
go build -o wallet-service main.go
./wallet-service
```

The service listens on port `8080` by default.

## API Endpoints

### Get Latest Transactions

* **URL**: `/api/transactions?count={n}`
* **Method**: `GET`
* **Query Parameters**:

  * `count` (int, required): Number of most recent transactions to return.
* **Response**: JSON array of transactions.

**Example**:

```bash
curl "http://localhost:8080/api/transactions?count=5"
```

```json
[
  {"from":"<hash1>","to":"<hash2>","amount":25},
  {"from":"<hash3>","to":"<hash4>","amount":10}
]
```

### Send Transaction

* **URL**: `/api/send`
* **Method**: `POST`
* **Headers**: `Content-Type: application/json`
* **Body**:

  ```json
  {
    "from": "<sender_hash>",
    "to": "<recipient_hash>",
    "amount": 50
  }
  ```
* **Responses**:

  * `200 OK`: Transaction successful.
  * `400 Bad Request`: Validation error (self-transfer, insufficient funds, unknown wallet).

**Example**:

```bash
curl -X POST http://localhost:8080/api/send \
  -H "Content-Type: application/json" \
  -d '{"from":"addr1","to":"addr2","amount":50}'
```

### Get Wallet Balance

* **URL**: `/api/wallet/{address}/balance`
* **Method**: `GET`
* **Path Parameters**:

  * `address` (string, required): Wallet hash.
* **Response**: JSON number indicating balance.

**Example**:

```bash
curl http://localhost:8080/api/wallet/<hash>/balance
```

```json
100
```
