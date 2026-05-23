// @title           Go CRUD API
// @version         1.0
// @description     REST API для управления пользователями, продуктами и корзинами.
// @host            localhost:5000
// @BasePath        /

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "go-crud/docs"
	"go-crud/internal/db"
	"go-crud/internal/handler"
)

func main() {

	_ = godotenv.Load(".env.example")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/crud?sslmode=disable"
	}

	if err := runMigrations(dsn); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}
	defer pool.Close()

	queries := db.New(pool)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	users := handler.NewUserHandler(queries)
	products := handler.NewProductHandler(queries)
	carts := handler.NewCartHandler(queries)

	r.Route("/users", func(r chi.Router) {
		r.Get("/", users.List)
		r.Post("/", users.Create)
		r.Get("/{id}", users.Get)
		r.Put("/{id}", users.Update)
		r.Delete("/{id}", users.Delete)
	})

	r.Route("/products", func(r chi.Router) {
		r.Get("/", products.List)
		r.Post("/", products.Create)
		r.Get("/{id}", products.Get)
		r.Put("/{id}", products.Update)
		r.Delete("/{id}", products.Delete)
	})

	r.Route("/carts", func(r chi.Router) {
		r.Get("/user/{userID}", carts.ListByUser)
		r.Post("/", carts.Add)
		r.Put("/{id}", carts.Update)
		r.Delete("/{id}", carts.Remove)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	addr := ":5000"
	fmt.Printf("Server listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

func runMigrations(dsn string) error {
	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("parse dsn: %w", err)
	}

	sqlDB := stdlib.OpenDB(*connConfig)
	defer sqlDB.Close()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run migrations: %w", err)
	}
	return nil
}
