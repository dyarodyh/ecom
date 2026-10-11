package products

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	repository *Repository
}

func NewHandler(repository *Repository) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)

	defer cancel()

	items, err := h.repository.List(ctx)
	if err != nil {
		log.Printf("Gagal mengambil produk: %v", err)

		http.Error(
			w,
			"Gagal mengambil produk",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(items); err != nil {
		log.Printf("Gagal mengirim respons produks: %v", err)
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string `json:"name"`
		Price *int64 `json:"price"`
		Stock *int   `json:"stock"`
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		http.Error(
			w,
			"Body harus berupa JSON produk yang valid, masksimal 1 MB",
			http.StatusBadRequest,
		)
		return
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		http.Error(
			w,
			"Body hanya boleh berisi satu objek JSON",
			http.StatusBadRequest,
		)
		return
	}

	input.Name = strings.TrimSpace(input.Name)

	if input.Name == "" {
		http.Error(w, "Nama Produk wajib diisi.", http.StatusBadRequest)
		return
	}
	if input.Price == nil || *input.Price < 0 {
		http.Error(
			w,
			"Harga wajib diisi dan tidak boleh negatif.",
			http.StatusBadRequest,
		)
		return
	}

	if input.Stock == nil || *input.Stock < 0 || *input.Stock > 2_147_483_647 {
		http.Error(
			w,
			"Stok wajib diisi antara 0 dan 2147483647.",
			http.StatusBadRequest,
		)
		return
	}

	product := Product{
		Name:  input.Name,
		Price: *input.Price,
		Stock: *input.Stock,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	createdProduct, err := h.repository.Create(ctx, product)
	if err != nil {
		log.Printf("Gagal menambahkan produk: %v", err)
		http.Error(
			w,
			"Gagal menambahkan produk.",
			http.StatusInternalServerError,
		)
		return
	}

	// 5. Kirim produk yang berhasil dibuat.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdProduct); err != nil {
		log.Printf("Gagal mengirim respons produk baru: %v", err)
	}
}
