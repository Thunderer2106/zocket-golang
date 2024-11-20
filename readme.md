# Zocket API

Zocket API is a RESTful API that allows users to interact with a PostgreSQL database for product-related operations. The API provides endpoints to manage products, caching, and other operations while leveraging a PostgreSQL connection pool and Redis cache for efficient data retrieval.

## Features

- **Product Management**: Endpoints to manage products (create, read).
- **RabbitMQ Integration**
   - When a new product is created, the API sends a message to RabbitMQ. This can be used to trigger other actions or notify services about the new product.
   - RabbitMQ is integrated into the backend, and any interaction with the product creation endpoint will result in a message being queued for further processing.
- **S3 Storage**
   - The images and compressed images associated with the product are uploaded to Amazon S3 for persistent storage.
   - The product creation request will process the image files and store them in the S3 bucket. The API will return the URLs of the stored images.
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
    REDIS_DB=0
    AWS_ACCESS_KEY=""
    AWS_SECRET_ACCESS_KEY=""
    AWS_REGION=""
    ```

5. Run the application:

    ```bash
    go run cmd/api/main.go
    ```

6. Your application will be running on `http://localhost:8080` by default.

## API Endpoints

### GET `/products/{id}`

Fetch a product by its ID. It first checks the cache and then queries the database if not found.

#### Example Request:

```bash
GET /products/1

{
  "id": 1,
  "user_id": 10,
  "product_name": "Test Product",
  "product_description": "This is a test product.",
  "product_images": ["image1.jpg", "image2.jpg"],
  "compressed_product_images": ["compressed_image1.jpg", "compressed_image2.jpg"],
  "product_price": 19.99,
  "created_at": "2024-11-20T00:00:00Z"
}
```

# POST /products

This endpoint allows you to create a new product by sending a POST request with the required product data. The server will process the request and return the created product's details.

## Request

### Method

- `POST /products`

### Request Body

The request body must contain a JSON object with the following fields:

| Field                      | Type        | Description                                       | Required  |
|----------------------------|-------------|---------------------------------------------------|-----------|
| `user_id`                  | `int`       | The ID of the user creating the product.          | Yes       |
| `product_name`             | `string`    | The name of the product.                          | Yes       |
| `product_description`      | `string`    | A brief description of the product.               | Yes       |
| `product_images`           | `[]string`     | An array of product image filenames (URLs or paths). | Yes       |
| `product_price`            | `float64`   | The price of the product.                         | Yes       |

### Example Request

```bash
POST /products
Content-Type: application/json

{
  "user_id": 1,
  "product_name": "New Product",
  "product_description": "Description of the new product.",
  "product_images": ["image1.jpg"],
  "product_price": 29.99
}
```
### Example Response
```bash
{
  "id": 2,
  "user_id": 1,
  "product_name": "New Product",
  "product_description": "Description of the new product.",
  "product_images": ["image1.jpg"],
  "compressed_product_images": ["compressed_image1.jpg"],
  "product_price": 29.99,
  "created_at": "2024-11-20T00:00:00Z"
}
```

