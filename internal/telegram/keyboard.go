package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MainMenu() tgbotapi.InlineKeyboardMarkup {

	buttons := [][]tgbotapi.InlineKeyboardButton{}
	for _, category := range Categories {
		button := tgbotapi.NewInlineKeyboardButtonData(
			category.Name,
			"category:"+category.ID,
		)

		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{button})
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func ItemsMenu(categoryID string) tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{}
	for _, item := range Items {

		if item.CategoryID != categoryID {
			continue
		}

		button := tgbotapi.NewInlineKeyboardButtonData(
			item.Name,
			"item:"+item.ID,
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

func ItemMenu(categoryID string) tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData("⬅ Назад","category:"+categoryID,)},
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}