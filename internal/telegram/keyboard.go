package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

func ElementsMenu(categoryID string) tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{}
	for _, element := range Elements {

		if element.CategoryID != categoryID {
			continue
		}

		button := tgbotapi.NewInlineKeyboardButtonData(
			element.Name,
			"element:"+element.ID,
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

func PerformancesMenu() tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{}
	for _, performance := range Performances {

		button := tgbotapi.NewInlineKeyboardButtonData(
			performance.Name,
			"performances:"+performance.ID,
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
