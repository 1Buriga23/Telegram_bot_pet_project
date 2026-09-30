package telegram

func GetCategoryByID(id string) *Category {
	for _, category := range Categories {
		if category.ID == id {
			return &category
		}
	}
	return nil
}

var Categories = []Category{
	{
		Name: "👦👩👦 Акробатика",
		ID:   "acrobatics",
		Type: "items",
	},

	{
		Name: "👦👩 Гимнастика",
		ID:   "gymnastics",
		Type: "items",
	},

	{
		Name: "🏆 Выступления",
		ID:   "performances",
		Type: "items",
	},
}

var Items = []Item{
	{
		Name:       "🤸 Сальто",
		ID:         "flips",
		CategoryID: "acrobatics",
	},

	{
		Name:       "🔄 Перевороты",
		ID:         "twists",
		CategoryID: "acrobatics",
	},

	{
		Name:       "⚖️ Балансы",
		ID:         "balances",
		CategoryID: "acrobatics",
	},

	{
		Name: "Пираты",
		ID:   "pirates",
		CategoryID: "performances",
	},

	{
		Name: "Обезьяны",
		ID:   "monkey",
		CategoryID: "performances",
	},

	{
		Name: "Круэлла",
		ID:   "cruela",
		CategoryID: "performances",
	},
}
