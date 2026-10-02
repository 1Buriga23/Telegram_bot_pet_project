package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"gracia-bot/internal/models"
)

func GetCategories(ctx context.Context, db *pgxpool.Pool) ([]models.Category, error) {
	rows, err := db.Query(
		ctx,
		`
	SELECT id, name
	FROM categories
	`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	categories := []models.Category{}
	for rows.Next(){
		var category models.Category
		err := rows.Scan(&category.ID,&category.Name)

		if err != nil {
			return nil, err
		}

		categories = append(categories, category)

	}
	return  categories, nil
}

func GetCategoryByID(
	ctx context.Context,
	db *pgxpool.Pool,
	id int64,
) (*models.Category, error) {


	var category models.Category


	err := db.QueryRow(
		ctx,
		`
		SELECT id, name
		FROM categories
		WHERE id = $1
		`,
		id,
	).Scan(
		&category.ID,
		&category.Name,
	)


	if err != nil {
		return nil, err
	}


	return &category, nil
}
