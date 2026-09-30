package telegram

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) HandleUpdate(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		if update.CallbackQuery.Message != nil {
			b.HandleCallback(update.CallbackQuery)
		}
		return
	}

	if update.Message != nil {
		b.HandleMessage(update.Message)
		return
	}
}

func (b *Bot) EditMenu(
	chatID int64,
	messageID int,
	text string,
	keyboard tgbotapi.InlineKeyboardMarkup,
) {
	edit := tgbotapi.NewEditMessageText(
		chatID,
		messageID,
		text,
	)

	edit.ReplyMarkup = &keyboard

	_, err := b.API.Send(edit)
	if err != nil {
		log.Println(err)
	}
}

func (b *Bot) ShowMainMenu(
	chatID int64,
	text string,
) {
	message := tgbotapi.NewMessage(
		chatID,
		text,
	)

	message.ReplyMarkup = MainMenu()

	_, err := b.API.Send(message)
	if err != nil {
		log.Println(err)
	}
}

func (b *Bot) HandleCallback(callback *tgbotapi.CallbackQuery) {
	data := callback.Data

	switch {

	case strings.HasPrefix(data, "category:"):

		categoryID := strings.TrimPrefix(data, "category:")
		b.HandleCategory(callback, categoryID)

	case strings.HasPrefix(data, "item:"):

		elementID := strings.TrimPrefix(data, "item:")
		b.HandleItem(callback, elementID)

	case data == "menu:main":
		chatID := callback.Message.Chat.ID
		messageID := callback.Message.MessageID
		b.EditMenu(
			chatID,
			messageID,
			"🤸 Gracia Bot 🤸\n\n\nВыбери каталог:",
			MainMenu(),
		)

	}

}

func (b *Bot) HandleMessage(message *tgbotapi.Message) {
	if message.IsCommand() {
		chatID := message.Chat.ID

		switch message.Command() {

		case "start":

			textHello := fmt.Sprintf("Привет, %s!\n\n\n🤸 Gracia Bot 🤸\n\n\nВыбери каталог:", message.From.FirstName)

			if message.From.FirstName == "" {
				textHello = "Привет!\n\n\n🤸 Gracia Bot 🤸\n\n\nВыбери каталог:"
			}

			b.ShowMainMenu(
				chatID,
				textHello,
			)
		default:
			return
		}
	}
}

func (b *Bot) HandleCategory(callback *tgbotapi.CallbackQuery, categoryID string) {

	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID
	category := GetCategoryByID(categoryID)

	if category == nil {
	return
}
	
	text := fmt.Sprintf("%s\n\nВыберите элемент:", category.Name)

	if categoryID == "performances"{
		text = fmt.Sprintf("%s\n\nВыберите номер:", category.Name)
	}

	b.EditMenu(
		chatID,
		messageID,
		text,
		ItemsMenu(categoryID),
	)
}

func (b *Bot) HandleItem(callback *tgbotapi.CallbackQuery, categoryID string) {

}

func (b *Bot) Start() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := b.API.GetUpdatesChan(updateConfig)

	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		b.HandleUpdate(update)
	}
}
