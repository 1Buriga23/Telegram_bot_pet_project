package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"gracia-bot/internal/models"
)


// Получение всех элементов определённой категории
func GetItemsByCategoryID(
	ctx context.Context,
	db *pgxpool.Pool,
	categoryID int64,
) ([]models.Item, error) {

	rows, err := db.Query(
		ctx,
		`
		SELECT id, name, category_id
		FROM items
		WHERE category_id = $1
		`,
		categoryID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()


	var items []models.Item


	for rows.Next() {

		var item models.Item


		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.CategoryID,
		)

		if err != nil {
			return nil, err
		}


		items = append(
			items,
			item,
		)
	}


	return items, nil
}



// Получение одного элемента по ID
func GetItemByID(
	ctx context.Context,
	db *pgxpool.Pool,
	id int64,
) (*models.Item, error) {


	var item models.Item


	err := db.QueryRow(
		ctx,
		`
		SELECT id, name, category_id
		FROM items
		WHERE id = $1
		`,
		id,
	).Scan(
		&item.ID,
		&item.Name,
		&item.CategoryID,
	)


	if err != nil {
		return nil, err
	}


	return &item, nil
}