package telegram

import (
	"context"
	"fmt"
	"gracia-bot/internal/models"
	"gracia-bot/internal/repository"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

	categories, err := repository.GetCategories(
		context.Background(),
		b.DB,
	)

	message := tgbotapi.NewMessage(
		chatID,
		text,
	)

	message.ReplyMarkup = MainMenu(categories)

	msg, err := b.API.Send(message)

	if err != nil {
		log.Println(err)
		return
	}

	session, exists := Sessions[chatID]

	if !exists {
		session = &UserSession{}
	}

	session.MenuMessageID = msg.MessageID

	Sessions[chatID] = session
}

func (b *Bot) HandleCallback(callback *tgbotapi.CallbackQuery) {
	data := callback.Data

	switch {

	case strings.HasPrefix(data, "category:"):
		b.DeleteVideos(callback.Message.Chat.ID)
		categoryID := strings.TrimPrefix(data, "category:")
		b.HandleCategory(callback, categoryID)

	case strings.HasPrefix(data, "item:"):

		itemID := strings.TrimPrefix(data, "item:")
		b.HandleItem(callback, itemID)

	case data == "menu:main":
		chatID := callback.Message.Chat.ID
		messageID := callback.Message.MessageID

		categories, _ := repository.GetCategories(
			context.Background(),
			b.DB,
		)

		b.EditMenu(
			chatID,
			messageID,
			"🤸 Gracia Bot 🤸\n\n\nВыбери каталог:",
			MainMenu(categories),
		)

	}

}

func (b *Bot) HandleMessage(message *tgbotapi.Message) {
	if message.IsCommand() {
		chatID := message.Chat.ID
		b.ClearAll(chatID)
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

	if message.Video != nil {

		fmt.Println("Video FileID:")
		fmt.Println(message.Video.FileID)

		return
	}
}

func (b *Bot) HandleCategory(callback *tgbotapi.CallbackQuery, categoryID string) {

	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID

	id, err := strconv.ParseInt(
		categoryID,
		10,
		64,
	)

	if err != nil {
		log.Println(err)
		return
	}

	items, err := repository.GetItemsByCategoryID(
		context.Background(),
		b.DB,
		id,
	)

	if err != nil {
		log.Println(err)
		return
	}

	category, err := repository.GetCategoryByID(
		context.Background(),
		b.DB,
		id,
	)

	if err != nil {
		log.Println(err)
		return
	}

	text := fmt.Sprintf("%s\n\nВыберите элемент:", category.Name)

	if category.Name == "Выступления" {
		text = fmt.Sprintf("%s\n\nВыберите номер:", category.Name)
	}

	b.EditMenu(
		chatID,
		messageID,
		text,
		ItemsMenu(items),
	)
}

func (b *Bot) HandleItem(callback *tgbotapi.CallbackQuery, itemID string) {
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID

	id, err := strconv.ParseInt(
		itemID,
		10,
		64,
	)

	if err != nil {
		log.Println(err)
		return
	}

	item, err := repository.GetItemByID(
		context.Background(),
		b.DB,
		id,
	)

	text := fmt.Sprintf("%s", item.Name)

	b.EditMenu(
		chatID,
		messageID,
		text,
		ItemMenu(item.ID),
	)

	videos, err := repository.GetVideosByItemID(
		context.Background(),
		b.DB,
		id,
	)

	if err != nil {
		log.Println(err)
		return
	}

	b.SendVideos(
		chatID,
		videos,
	)
}

func (b *Bot) SendVideos(chatID int64, videos []models.Video) {
	if len(videos) == 0 {
		return
	}

	var media []interface{}

	for _, video := range videos {
		item := tgbotapi.NewInputMediaVideo(tgbotapi.FileID(video.FileID))

		item.Caption = video.Name

		media = append(media, item)
	}

	message := tgbotapi.NewMediaGroup(
		chatID,
		media,
	)
	messages, err := b.API.SendMediaGroup(message)

	if err != nil {
		fmt.Println(err)
	}

	var messageIDs []int

	for _, message := range messages {
		messageIDs = append(messageIDs, message.MessageID)
	}

	session, exists := Sessions[chatID]

	if !exists {
		session = &UserSession{}
	}

	session.VideoMessages = messageIDs

	Sessions[chatID] = session
}

func (b *Bot) Start() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := b.API.GetUpdatesChan(updateConfig)

	for update := range updates {
		b.HandleUpdate(update)
	}
}
