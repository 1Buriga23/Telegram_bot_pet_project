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
		Type: "elements",
	},

	{
		Name: "👦👩 Гимнастика",
		ID:   "gymnastics",
		Type: "elements",
	},

	{
		Name: "🏆 Выступления",
		ID:   "performances",
		Type: "performance",
	},
}

var Elements = []Element{
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
}

var Performances = []Performance{
	{
		Name: "Пираты",
		ID:   "pirates",
	},

	{
		Name: "Обезьяны",
		ID:   "monkey",
	},

	{
		Name: "Круэлла",
		ID:   "cruela",
	},
}
