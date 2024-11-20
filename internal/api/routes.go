package api

import (
	"net/http"
	"zocket/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Test route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Product routes
	router.POST("/products", handlers.CreateProductHandler)
	router.GET("/products/:id", handlers.GetProductByIDHandler)
}
