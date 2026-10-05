package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)


func IsAdmin(
	ctx context.Context,
	db *pgxpool.Pool,
	telegramID int64,
) (bool, error) {

	var exists bool

	err := db.QueryRow(
		ctx,
		`
		SELECT EXISTS(
			SELECT 1
			FROM admins
			WHERE telegram_id = $1
		)
		`,
		telegramID,
	).Scan(&exists)

	return exists, err
}


func AddAdmin(
	ctx context.Context,
	db *pgxpool.Pool,
	telegramID int64,
) error {

	_, err := db.Exec(
		ctx,
		`
		INSERT INTO admins(telegram_id)
		VALUES($1)
		ON CONFLICT DO NOTHING
		`,
		telegramID,
	)

	return err
}