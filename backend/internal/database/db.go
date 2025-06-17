package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool
var ctx = context.Background()

func Connect() {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		user, password, host, port, dbname,
	)

	// Create a context with a timeout for the connection
	// This is to ensure that the connection attempt does not hang indefinitely
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var err error
	DB, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("❌ Failed to create DB pool: %v", err)
	}

	if err := DB.Ping(ctx); err != nil {
		log.Fatalf("❌ Cannot connect to DB: %v", err)
	}

	log.Println("🐉 Connected to PostgreSQL using pgxpool!")
}

func Close() {
	if DB != nil {
		DB.Close()
		log.Println("🐉 PostgreSQL pool closed successfully!")
	} else {
		log.Println("⚠️ No DB connection pool to close.")
	}
}
