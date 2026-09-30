package telegram

import (
	"log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserSession struct {
	MenuMessageID int
	VideoMessages []int
}

func (b *Bot) DeleteVideos(chatID int64) {

	session, exists := Sessions[chatID]

	if !exists {
		return
	}


	for _, messageID := range session.VideoMessages {

		msg := tgbotapi.NewDeleteMessage(
			chatID,
			messageID,
		)


		_, err := b.API.Request(msg)

		if err != nil {
			log.Println(err)
		}
	}

	session.VideoMessages = nil
}

func (b *Bot) ClearAll(chatID int64) {

	session, exists := Sessions[chatID]

	if !exists {
		return
	}


	for _, messageID := range session.VideoMessages {

		msg := tgbotapi.NewDeleteMessage(
			chatID,
			messageID,
		)


		_, err := b.API.Request(msg)

		if err != nil {
			log.Println(err)
		}
	}

	if session.MenuMessageID != 0 {

		deleteMsg := tgbotapi.NewDeleteMessage(
			chatID,
			session.MenuMessageID,
		)

		b.API.Request(deleteMsg)
	}

	delete(
		Sessions,
		chatID,
	)
}

var Sessions = make(map[int64]*UserSession)