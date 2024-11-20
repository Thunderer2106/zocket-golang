package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"zocket/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateProductHandler(t *testing.T) {
    gin.SetMode(gin.TestMode)
    db.SetupTestDB()

    r := gin.Default()
    r.POST("/products", CreateProductHandler)

    body := `{
        "user_id": 1,
        "product_name": "Test Product",
        "product_description": "Test Description",
        "product_images": ["http://example.com/image1.jpg"],
        "product_price": 100.0
    }`

    req, _ := http.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    resp := httptest.NewRecorder()

    r.ServeHTTP(resp, req)

    assert.Equal(t, http.StatusOK, resp.Code)
    assert.Contains(t, resp.Body.String(), "Product added successfully")
}
