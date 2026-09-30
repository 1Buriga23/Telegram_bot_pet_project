package telegram

func GetCategoryByID(id string) *Category {
	for _, category := range Categories {
		if category.ID == id {
			return &category
		}
	}
	return nil
}

func GetItemByID(id string) *Item {
	for _, item := range Items {
		if item.ID == id {
			return &item
		}
	}
	return nil
}

func GetVideoByItemID(itemID string) []Video {

	var videos []Video

	for _, video := range Videos {
		if video.ItemID == itemID {
			videos = append(videos, video)
		}
	}
	return videos
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
		ID:   "cruella_show",
		CategoryID: "performances",
	},
}

var Videos = []Video{

	{
		ID: "backflip_1",
		ItemID: "flips",
		Name: "▶ Выполнение",
		FileID: "BAACAgIAAxkBAAOmar1nmvEBrl8Kmv6jP7zyJgRpDrUAAmKtAALHEPBJWpLQrzfQFD89BA",
	},


	// {
	// 	ID: "backflip_2",
	// 	ItemID: "backflip",
	// 	Name: "▶ Разбор ошибок",
	// 	FileID: "ВАШ_FILE_ID",
	// },


	{
		ID: "show_1",
		ItemID: "cruella_show",
		Name: "▶ Выступление",
		FileID: "BAACAgIAAxkBAAOnar1oBpLXozvDCGYwss97vr_I3IkAAmatAALHEPBJzmauXCzliAI9BA",
	},

}
