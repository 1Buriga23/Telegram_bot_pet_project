package models

type Category struct {
	ID   int64
	Name string
}

type Item struct {
	ID         int64
	CategoryID int64
	Name       string
	// Description string
}

type Video struct {
	ID     int64
	ItemID int64
	Name   string
	FileID string
}
 