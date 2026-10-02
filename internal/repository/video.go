package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"gracia-bot/internal/models"
)


func GetVideosByItemID(
	ctx context.Context,
	db *pgxpool.Pool,
	itemID int64,
) ([]models.Video, error) {


	rows, err := db.Query(
		ctx,
		`
		SELECT id, name, file_id, item_id
		FROM videos
		WHERE item_id = $1
		`,
		itemID,
	)


	if err != nil {
		return nil, err
	}


	defer rows.Close()


	var videos []models.Video


	for rows.Next() {

		var video models.Video


		err := rows.Scan(
			&video.ID,
			&video.Name,
			&video.FileID,
			&video.ItemID,
		)


		if err != nil {
			return nil, err
		}


		videos = append(
			videos,
			video,
		)
	}


	return videos, nil
}