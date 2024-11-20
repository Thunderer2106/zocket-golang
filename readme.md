# Zocket API

Zocket API is a RESTful API that allows users to interact with a PostgreSQL database for product-related operations. The API provides endpoints to manage products, caching, and other operations while leveraging a PostgreSQL connection pool and Redis cache for efficient data retrieval.

## Features

- **Product Management**: Endpoints to manage products (create, read, update, delete).
- **Caching**: Redis caching for faster data retrieval.
- **Database Connection Pooling**: Optimized PostgreSQL connection pool for improved performance.

## Installation

To set up the project locally, follow the steps below:

### Prerequisites

- [Go](https://golang.org/dl/) 1.18 or higher
- [PostgreSQL](https://www.postgresql.org/download/)
- [Redis](https://redis.io/download)
- [Go Modules](https://blog.golang.org/using-go-modules)

### Steps to Install

1. Clone the repository:

    ```bash
    git clone https://github.com/your-username/zocket.git
    ```

2. Change directory to the project:

    ```bash
    cd zocket
    ```

3. Install dependencies using Go modules:

    ```bash
    go mod tidy
    ```

4. Set up your environment variables by creating a `.env` file in the root directory. Example:

    ```
    DB_USER=postgres
    DB_PASSWORD=your_password
    DB_HOST=localhost
    DB_PORT=5432
    DB_NAME=pgdb
    REDIS_HOST=localhost
    REDIS_PORT=6379
    ```

5. Run the application:

    ```bash
    go run main.go
    ```

6. Your application will be running on `http://localhost:8080` by default.

## API Endpoints

### GET `/products/{id}`

Fetch a product by its ID. It first checks the cache and then queries the database if not found.

#### Example Request:

```bash
GET /products/1
