package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(password string) *pgxpool.Pool {
	
	db, err := pgxpool.New(
		context.Background(),
		fmt.Sprintf("postgres://postgres:%s@localhost:5432/gracia_bot",password),
	)

	if err != nil{
		log.Fatal(err)
	}

	err = db.Ping(context.Background())

	if err != nil {
		log.Fatal(err)
	}


	return db
}