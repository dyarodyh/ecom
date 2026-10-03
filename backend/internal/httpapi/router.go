package httpapi

import (
	"net/http"

	"example.com/ecom/backend/internal/health"
	"example.com/ecom/backend/internal/products"
)

func NewRouter(productHandler *products.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", health.Handler)
	mux.HandleFunc("GET /api/products", productHandler.List)
	mux.HandleFunc("POST /api/products", productHandler.Create)

	return mux
}
