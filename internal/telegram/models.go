package telegram

type Category struct {
	Name string
	ID string
	Type string
}

type Element struct {
	ID string
	Name string
	CategoryID string
}

type Performance struct {
	ID string
	Name string
}