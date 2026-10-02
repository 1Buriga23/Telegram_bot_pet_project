package telegram

import (
	"gracia-bot/internal/models"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MainMenu(categories []models.Category) tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{}
	for _, category := range categories {
		button := tgbotapi.NewInlineKeyboardButtonData(
			category.Name,
			"category:"+strconv.FormatInt(category.ID, 10),
		)

		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{button})
	}
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
