
package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"zocket/internal/api/handlers"
	"zocket/redis"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
    // Initialize Redis
    redis.InitRedis()

    // Run tests
    os.Exit(m.Run())
}

func TestCacheEffectiveness(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.Default()
    r.GET("/products/:id", handlers.GetProductByIDHandler)
    client := redis.GetClient()
    client.Set(context.Background(), "product:1", `{"id":1,"name":"Test Product"}`, 0)

    req, _ := http.NewRequest(http.MethodGet, "/products/1", nil)
    resp := httptest.NewRecorder()
    r.ServeHTTP(resp, req)

    if resp.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", resp.Code)
    }
    if !contains(resp.Body.String(), "Test Product") {
        t.Errorf("Response body does not contain 'Test Product'")
    }

    client.FlushAll(context.Background()) 
    resp2 := httptest.NewRecorder()
    r.ServeHTTP(resp2, req)

    if resp2.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", resp2.Code)
    }
    if !contains(resp2.Body.String(), "Test Product") {
        t.Errorf("Response body does not contain 'Test Product'")
    }
}

func contains(body, substring string) bool {
    return strings.Contains(body, substring)
}
