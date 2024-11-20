
package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
	"zocket/internal/db"
	"zocket/redis" 

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type Product struct {
	UserID          int      `json:"user_id" binding:"required"`
	ProductName     string   `json:"product_name" binding:"required"`
	ProductDesc     string   `json:"product_description" binding:"required"`
	ProductImages   []string `json:"product_images" binding:"required"`
	ProductPrice    float64  `json:"product_price" binding:"required"`
}

func enqueueImageURLs(images []string) error {
	for _, image := range images {
		body, err := json.Marshal(map[string]string{"image_url": image})
		if err != nil {
			log.Printf("Failed to marshal image URL %s: %v", image, err)
			return err
		}

		err = RabbitMQChannel.Publish(
			"",
			"imageProcessingQueue",
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		if err != nil {
			log.Printf("Failed to publish message to queue: %v", err)
			return err
		}
		log.Printf("Published message for image URL: %s", image)
	}
	return nil
}

func CreateProductHandler(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		INSERT INTO products (user_id, product_name, product_description, product_images, product_price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	var productID int
	err := db.DB.QueryRow(context.Background(), query,
		product.UserID, product.ProductName, product.ProductDesc, product.ProductImages, product.ProductPrice).Scan(&productID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert product into database"})
		return
	}


	err = redis.GetClient().Del(context.Background(), "product:"+strconv.Itoa(productID)).Err()
	if err != nil {
		log.Printf("Failed to invalidate cache for product %d: %v", productID, err)
	}

	err = enqueueImageURLs(product.ProductImages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue image URLs for processing"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "\033[1mProduct added successfully!\033[0m",
		"product_id": productID,
	})
}

func GetProductByIDHandler(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	cachedProduct, err := redis.GetClient().Get(context.Background(), "product:"+strconv.Itoa(productID)).Result()
	if err == nil {

		log.Printf("\033[1mCache hit for product %d\033[0m", productID)
		c.JSON(http.StatusOK, cachedProduct)
		return
	}
	log.Printf("\033[1mCache miss for product %d, querying database\033[0m", productID)

	query := `
		SELECT id, user_id, product_name, product_description, product_images, compressed_product_images, product_price, created_at
		FROM products
		WHERE id = $1;
	`

	var product struct {
		ID                     int       `json:"id"`
		UserID                 int       `json:"user_id"`
		ProductName            string    `json:"product_name"`
		ProductDesc            string    `json:"product_description"`
		ProductImages          []string  `json:"product_images"`
		CompressedProductImages []string `json:"compressed_product_images"`
		ProductPrice           float64   `json:"product_price"`
		CreatedAt              time.Time `json:"created_at"`
	}

	err = db.DB.QueryRow(context.Background(), query, productID).Scan(
		&product.ID,
		&product.UserID,
		&product.ProductName,
		&product.ProductDesc,
		&product.ProductImages,
		&product.CompressedProductImages,
		&product.ProductPrice,
		&product.CreatedAt,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	productJSON, err := json.Marshal(product)
	if err != nil {
		log.Printf("Failed to marshal product %d for caching: %v", productID, err)
	}

	err = redis.GetClient().Set(context.Background(), "product:"+strconv.Itoa(productID), productJSON, time.Hour).Err()
	if err != nil {
		log.Printf(" \033[1mFailed to cache product %d: %v \033[0m", productID, err)
	}

	c.JSON(http.StatusOK, product)
}
