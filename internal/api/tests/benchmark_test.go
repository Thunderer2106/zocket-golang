package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"zocket/internal/api/handlers"
	"zocket/redis"

	"github.com/gin-gonic/gin"
)

func BenchmarkGetProductByIDWithCache(b *testing.B) {
    gin.SetMode(gin.TestMode)
    r := gin.Default()
    r.GET("/products/:id", handlers.GetProductByIDHandler)

    req, _ := http.NewRequest(http.MethodGet, "/products/1", nil)

    for i := 0; i < b.N; i++ {
        resp := httptest.NewRecorder()
        r.ServeHTTP(resp, req)
    }
}

func BenchmarkGetProductByIDWithoutCache(b *testing.B) {
    gin.SetMode(gin.TestMode)
    r := gin.Default()
    r.GET("/products/:id", handlers.GetProductByIDHandler)
	redis.InitRedis()
	client := redis.GetClient()
    client.FlushAll(context.Background())

    req, _ := http.NewRequest(http.MethodGet, "/products/1", nil)

    for i := 0; i < b.N; i++ {
        resp := httptest.NewRecorder()
        r.ServeHTTP(resp, req)
    }
}

