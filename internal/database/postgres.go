package database

import (
	"context"
	"log"
	"net"
	"net/url"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(password string) *pgxpool.Pool {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "gracia_bot"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(dbUser, password),
		Host:   net.JoinHostPort(host, port),
		Path:   dbName,
	}

	db, err := pgxpool.New(
		context.Background(),
		dsn.String(),
	)
	if err != nil {
		log.Fatal("Ошибка создания подключения к БД:", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		db.Close()
		log.Fatal("Ошибка подключения к PostgreSQL:", err)
	}

	return db
}
