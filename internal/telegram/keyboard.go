package telegram

import (
	"gracia-bot/internal/models"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BackToMenu() tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData(
                    "⬅️ Назад",
                    "menu:main",
                ),
		},
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func MainMenu(categories []models.Category) tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{}
	for _, category := range categories {
		button := tgbotapi.NewInlineKeyboardButtonData(
			category.Name,
			"category:"+strconv.FormatInt(category.ID, 10),
		)

		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{button})
	}
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			"🔎 Поиск",
			"search",
		),
	),
	)
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func ItemsMenu(items []models.Item) tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{}
	for _, item := range items {

		button := tgbotapi.NewInlineKeyboardButtonData(
			item.Name,
			"item:"+strconv.FormatInt(item.ID, 10),
		)

		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{button})
	}

	buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(
			"⬅ Назад",
			"menu:main",
		),
	})
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func ItemMenu(categoryID int64) tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData("⬅ Назад", "category:"+strconv.FormatInt(categoryID, 10))},
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func SearchResultsMenu(items []models.Item) tgbotapi.InlineKeyboardMarkup {

	buttons := [][]tgbotapi.InlineKeyboardButton{}

	for _, item := range items {

		button := tgbotapi.NewInlineKeyboardButtonData(
			item.Name,
			"item:"+strconv.FormatInt(item.ID, 10),
		)

		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{button})
	}

	buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(
			"⬅ Назад",
			"menu:main",
		),
	})

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func SearchItemMenu() tgbotapi.InlineKeyboardMarkup {

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"⬅ Назад",
				"menu:main",
			),
		),
	)
}

func AdminPanelMenu() tgbotapi.InlineKeyboardMarkup {

	buttons := [][]tgbotapi.InlineKeyboardButton{}

	buttons = append(buttons,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"➕ Добавить категорию",
				"admin:add_category",
			),
		),
	)

	buttons = append(buttons,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"➕ Добавить элемент/выступление",
				"admin:add_item",
			),
		),
	)

	buttons = append(buttons,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🎥 Добавить видео",
				"admin:add_video",
			),
		),
	)

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}
