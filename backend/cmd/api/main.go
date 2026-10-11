package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
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
	if err := godotenv.Load("../devops/.env"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.New(
				"gagal membaca ../devops/.env; periksa format dan akses file",
			)
		}
	}

	// konfigurasi env
	for _, key := range []string{
		"POSTGRES_DB",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
	} {
		if os.Getenv(key) == "" {
			return fmt.Errorf("%s belum diisi pada environment", key)
		}
	}

	//menggunakan alamat lokal jika di env belum diberikan
	databaseAddr := os.Getenv("DATABASE_ADDR")
	if databaseAddr == "" {
		databaseAddr = "127.0.0.1:5432"
	}

	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = "127.0.0.1:8080"
	}

	// alamat koneksi postgres
	databaseURL := url.URL{
		Scheme: "postgres",
		User: url.UserPassword(
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
		),

		Host:     databaseAddr,
		Path:     "/" + os.Getenv("POSTGRES_DB"),
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
		Addr:              httpAddr,
		Handler:           httpapi.NewRouter(productHandler),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Backend mendengarkan di %s", httpAddr)

	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("menjalankan server HTTP: %w", err)
	}

	return nil
}
