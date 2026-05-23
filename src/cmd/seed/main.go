package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"go-crud/internal/db"
)

func main() {
	_ = godotenv.Load(".env.example")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/crud_db?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	q := db.New(pool)
	ctx := context.Background()

	if err := seed(ctx, q); err != nil {
		log.Fatalf("seed: %v", err)
	}
	fmt.Println("Seeding completed.")
}

func seed(ctx context.Context, q *db.Queries) error {
	users, err := seedUsers(ctx, q)
	if err != nil {
		return fmt.Errorf("users: %w", err)
	}
	products, err := seedProducts(ctx, q)
	if err != nil {
		return fmt.Errorf("products: %w", err)
	}
	if err := seedCarts(ctx, q, users, products); err != nil {
		return fmt.Errorf("carts: %w", err)
	}
	return nil
}

func seedUsers(ctx context.Context, q *db.Queries) ([]db.User, error) {
	rows := []db.CreateUserParams{
		{Name: "Алексей Иванов", Email: "alexey@example.com"},
		{Name: "Мария Петрова", Email: "maria@example.com"},
		{Name: "Дмитрий Сидоров", Email: "dmitry@example.com"},
	}

	var created []db.User
	for _, r := range rows {
		r.ID = uuid.Must(uuid.NewV7())
		u, err := q.CreateUser(ctx, r)
		if err != nil {
			log.Printf("skip user %s: %v", r.Email, err)
			continue
		}
		created = append(created, u)
		fmt.Printf("  user  %s  %s\n", u.ID, u.Email)
	}
	return created, nil
}

func ptr(s string) *string { return &s }

func seedProducts(ctx context.Context, q *db.Queries) ([]db.Product, error) {
	rows := []db.CreateProductParams{
		{Name: "Ноутбук Pro 15", Description: ptr("15-дюймовый ноутбук с процессором i7"), Price: 89999.99},
		{Name: "Беспроводные наушники", Description: ptr("Шумоподавление, 30 часов автономности"), Price: 12499.00},
		{Name: "Механическая клавиатура", Description: ptr("Переключатели Cherry MX Red, RGB"), Price: 7990.00},
		{Name: "Игровая мышь", Description: ptr("16000 DPI, 6 программируемых кнопок"), Price: 4590.00},
		{Name: "Монитор 27\" 4K", Description: ptr("IPS, 144 Гц, HDR400"), Price: 54990.00},
	}

	var created []db.Product
	for _, r := range rows {
		r.ID = uuid.Must(uuid.NewV7())
		p, err := q.CreateProduct(ctx, r)
		if err != nil {
			log.Printf("skip product %s: %v", r.Name, err)
			continue
		}
		created = append(created, p)
		fmt.Printf("  product  %s  %s — %.2f₽\n", p.ID, p.Name, p.Price)
	}
	return created, nil
}

func seedCarts(ctx context.Context, q *db.Queries, users []db.User, products []db.Product) error {
	if len(users) == 0 || len(products) == 0 {
		return nil
	}

	items := []db.AddToCartParams{
		{UserID: users[0].ID, ProductID: products[0].ID, Quantity: 1},
		{UserID: users[0].ID, ProductID: products[1].ID, Quantity: 2},
		{UserID: users[1].ID, ProductID: products[2].ID, Quantity: 1},
		{UserID: users[1].ID, ProductID: products[3].ID, Quantity: 3},
		{UserID: users[2].ID, ProductID: products[4].ID, Quantity: 1},
	}

	for _, item := range items {
		item.ID = uuid.Must(uuid.NewV7())
		c, err := q.AddToCart(ctx, item)
		if err != nil {
			log.Printf("skip cart item user=%s product=%s: %v", item.UserID, item.ProductID, err)
			continue
		}
		fmt.Printf("  cart  %s  user=%s product=%s qty=%d\n", c.ID, c.UserID, c.ProductID, c.Quantity)
	}
	return nil
}
