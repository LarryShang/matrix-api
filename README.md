# Matrix Operations API

A high-performance Go microservice for processing large matrix data via HTTP file uploads. This project demonstrates production-grade Go patterns, specifically focusing on handling large datasets efficiently using stream processing and concurrency.

![Go](https://img.shields.io/badge/go-1.21-blue.svg?style=for-the-badge&logo=go&logoColor=white)
![Architecture](https://img.shields.io/badge/architecture-clean-green.svg?style=for-the-badge)
![Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=for-the-badge)

## 📖 Overview

This service accepts `.csv` formatted matrices and performs mathematical operations (Sum, Multiply, Invert, Flatten).

### Key Features
* **Stream Processing:** Uses `io.Reader` to process files line-by-line, ensuring low memory footprint (`O(1)` space complexity for Sum/Multiply).
* **Concurrency:** Implements the **Producer-Consumer** pattern to decouple I/O (reading CSV) from CPU (BigInt calculations).
* **Clean Architecture:** Strictly separates `Domain`, `Service`, and `Handler` layers, utilizing interfaces for dependency injection.
* **Robustness:** Handles integer overflow using `math/big` and validates matrix structure (squareness) strictly.

---

## 🏗️ Architecture

The project follows the **Standard Go Layout**:

```text
matrix-api/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── domain/
│   │   └── errors.go         # Domain errors and constants
│   ├── handler/
│   │   └── matrix_handler.go # HTTP layer (Parsing multipart forms)
│   └── service/
│       ├── service.go        # Interface definition
│       ├── serial.go         # Serial implementation (Callback pattern)
│       └── concurrent.go     # Concurrent implementation (Producer-Consumer)
├── go.mod
└── README.md
```
The Producer-Consumer Pattern
For operations like Sum and Multiply, the service uses a concurrent pipeline to maximize throughput:

Producer (Goroutine): Reads the CSV file, parses strings to integers, and pushes them to a buffered channel. It handles I/O latency.

Consumer (Main Routine): Reads integers from the channel and updates the BigInt result. It handles CPU processing.

🚀 Getting Started
Prerequisites
Go 1.21+

Curl (for testing endpoints)

Running the Server
```Bash

# Clone the repository
git clone [https://github.com/](https://github.com/)[YourUsername]/matrix-api.git
cd matrix-api

# Run the server
go run cmd/server/main.go
```
Server will start on http://localhost:8080

## ⚡ API Reference
All endpoints expect a `POST` request with `multipart/form-data` containing a file field named `file`.

### 1. Echo Matrix
   Returns the matrix in string format.

Endpoint: `/echo`

Command:

```Bash

curl -F "file=@./matrix.csv" localhost:8080/echo
```
### 2. Invert Matrix
   Returns the matrix with columns and rows transposed.

Endpoint: `/invert`

Command:

```Bash

curl -F "file=@./matrix.csv" localhost:8080/invert
```
### 3. Flatten Matrix
   Returns the matrix as a single comma-separated line.

Endpoint: `/flatten`

Command:

```Bash

curl -F "file=@./matrix.csv" localhost:8080/flatten
```
### 4. Sum
   Returns the sum of all integers in the matrix.

Endpoint: `/sum`

Command:

```Bash

curl -F "file=@./matrix.csv" localhost:8080/sum
```
### 5. Multiply
   Returns the product of all integers in the matrix.

Endpoint: `/multiply`

Command:

```Bash

curl -F "file=@./matrix.csv" localhost:8080/multiply
```
**Note for Windows Users:** If using PowerShell, use curl.exe instead of curl to avoid alias conflicts with Invoke-WebRequest.

## 🧪 Running Tests
The project includes a comprehensive test suite covering both Unit (Service) and Integration (Handler) layers. The service tests verify both Serial and Concurrent implementations against the same dataset.

```Bash

# Run all tests with verbose output
go test ./... -v
```
## 💡 Design Decisions
### Q: Why use `math/big`?

Matrix multiplication grows exponentially. Standard int64 overflows very quickly. math/big ensures correctness even for large datasets.

### Q: Why separate `Serial` and `Concurrent` services?

This allows for easy A/B testing and benchmarking. For small files, the overhead of spinning up goroutines (context switching) might make the Serial version faster. The Concurrent version shines with large files where I/O blocking is the bottleneck.

### Q: How is the `Multiply` endpoint optimized?

It implements a Lazy Zero Check. If a `0` is encountered, the multiplication logic stops (result is guaranteed to be 0), but the stream continues draining to ensure the matrix is valid (square and properly formatted) before returning.