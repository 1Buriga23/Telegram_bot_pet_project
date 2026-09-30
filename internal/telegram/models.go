package telegram

type Category struct {
	Name string
	ID string
	Type string
}

type Item struct {
	ID string
	Name string
	CategoryID string
	// Description string
}

type Video struct {
	ID string
	ItemID string
	Name string
	FileID string
}