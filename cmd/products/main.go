package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"universe_group_test_task/internal/products"
	"universe_group_test_task/pkg/kafka"
)

func main() {
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:root@localhost:5432/products_db?sslmode=disable")
	httpAddr := getEnv("HTTP_ADDR", ":8080")
	migrationsPath := getEnv("MIGRATIONS_PATH", "file://migrations")
	kafkaBroker := getEnv("KAFKA_BROKER", "localhost:9092")
	kafkaTopic := getEnv("KAFKA_TOPIC", "product-events")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("connected to database")

	runMigrations(dbURL, migrationsPath)

	producer := kafka.NewProducer([]string{kafkaBroker}, kafkaTopic)
	defer producer.Close()
	log.Println("kafka producer initialized")

	repo := products.NewRepository(pool)
	service := products.NewService(repo, producer)
	handler := products.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Mount("/products", handler.Routes())
	r.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:    httpAddr,
		Handler: r,
	}

	go func() {
		log.Printf("products service listening on %s", httpAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server stopped")
}

func runMigrations(dbURL, migrationsPath string) {
	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("migrations applied successfully")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
