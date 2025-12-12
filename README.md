# Matrix Operations Web Service

This is a web service written in Go that accepts an uploaded CSV file representing a square matrix and performs several mathematical and formatting operations on it.

The service is designed to be memory-efficient, processing large files as a stream to avoid loading them entirely into memory.

---

## Project Architecture

The project follows a layered architecture to ensure a clean separation of concerns.

* **`main.go`**: The entry point of the application. It is responsible for:
    * **Dependency Injection**: Instantiating the service, domain, and handler layers.
    * **Routing**: Mapping URL endpoints to their corresponding handler logic.

* **`handler` Package**: This layer is responsible for all HTTP-related tasks.
    * It parses incoming HTTP requests and extracts the uploaded file.
    * It calls the appropriate service method to perform the business logic.
    * It formats the result from the service into an HTTP response.
    * It uses a higher-order function (`FileHandler`) to eliminate repetitive boilerplate code for file handling and error checking.

* **`service` Package**: This is the core business logic layer.
    * It contains the functions that perform matrix operations (`Sum`, `Multiply`, `Flatten`, etc.).
    * It uses a generic streaming function (`processMatrixStream`) to read and validate the CSV file row by row, ensuring the application can handle very large files without crashing.
    * This layer knows nothing about HTTP; its methods could be reused in a command-line tool or another application.

* **`domain` Package**: A shared package that holds common data types and error definitions.
    * Its primary purpose is to break potential cyclic dependencies. For example, both `handler` and `service` need access to the same custom error types, and this package provides a central location for them.

---

## API Endpoints

The service exposes the following `POST` endpoints. Each endpoint expects a `multipart/form-data` request with a file field named `file`.

* `/echo`
  Returns the uploaded matrix as a string, provided it passes validation. **Note**: This operation loads the entire matrix into memory.

* `/invert`
  Returns the matrix with its rows and columns transposed. **Note**: This operation must load the entire matrix into memory.

* `/flatten`
  Returns the matrix as a single line of comma-separated values.

* `/sum`
  Returns the sum of all integers in the matrix.

* `/multiply`
  Returns the product of all integers in the matrix.

---

## How to Use

### 1. Run the Server
Navigate to the project's root directory and run:
```bash
go run .
```
The server will start on `localhost:8080`.

### 2. Send Requests
Use a tool like `curl` to send a CSV file to an endpoint. Assume you have a file named `matrix.csv` in current directory with the content `1,2,3\n4,5,6\n7,8,9`.

```bash
# Echo the matrix
curl -F 'file=@./matrix.csv' "localhost:8080/echo"

# Invert the matrix
curl -F 'file=@./matrix.csv' "localhost:8080/invert"

# Flatten the matrix
curl -F 'file=@./matrix.csv' "localhost:8080/flatten"

# Get the sum
curl -F 'file=@./matrix.csv' "localhost:8080/sum"

# Get the product
curl -F 'file=@./matrix.csv' "localhost:8080/multiply"
```

### 3. Run Tests
To run the full suite of unit tests for all packages, use the following command:
```bash
go test ./...
```