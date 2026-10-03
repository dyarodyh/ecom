package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"example.com/ecom/backend/internal/database"
	"example.com/ecom/backend/internal/httpapi"
	"example.com/ecom/backend/internal/products"

	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// 1. Baca pengaturan dari file .env yang sudah ada.
	values, err := godotenv.Read("../devops/.env")
	if err != nil {
		return errors.New(
			"gagal membaca ../devops/.env; periksa lokasi dan format file",
		)
	}

	for _, key := range []string{
		"POSTGRES_DB",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
	} {
		if values[key] == "" {
			return fmt.Errorf("%s belum diisi di devops/.env", key)
		}
	}

	// 2. Susun alamat koneksi PostgreSQL.
	databaseURL := url.URL{
		Scheme: "postgres",
		User: url.UserPassword(
			values["POSTGRES_USER"],
			values["POSTGRES_PASSWORD"],
		),
		Host:     "127.0.0.1:5432",
		Path:     "/" + values["POSTGRES_DB"],
		RawQuery: "sslmode=disable",
	}

	// 3. Hubungkan backend ke database.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	pool, err := database.Open(ctx, databaseURL.String())
	cancel()

	if err != nil {
		return err
	}
	defer pool.Close()

	log.Println("Koneksi PostgreSQL berhasil")

	// 4. Siapkan repository, handler, dan server.
	productRepository := products.NewRepository(pool)
	productHandler := products.NewHandler(productRepository)

	server := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           httpapi.NewRouter(productHandler),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Backend berjalan di http://localhost:8080")

	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("menjalankan server HTTP: %w", err)
	}

	return nil
}
