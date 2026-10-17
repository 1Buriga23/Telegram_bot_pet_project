package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AdminMode string

const (
	AdminModeNone                AdminMode = ""
	AdminModeAddCategory         AdminMode = "add_category"
	AdminModeAddItem             AdminMode = "add_item"
	AdminModeAddVideo            AdminMode = "add_video"
	AdminModeSelectVideoCategory AdminMode = "select_video_category"
	AdminModeSelectVideoItem     AdminMode = "select_video_item"
)

type UserSession struct {
	MenuMessageID        int
	VideoMessages        []int
	SearchMode           bool
	FromSearch           bool
	WaitingAdminPassword bool
	AdminMode            AdminMode
	SelectedCategoryID   int64
	SelectedItemID       int64
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

func (b *Bot) DeleteChatMessage(chatID int64, messageID int) {
	msg := tgbotapi.NewDeleteMessage(chatID, messageID)

	if _, err := b.API.Request(msg); err != nil {
		log.Println("Не удалось удалить сообщение:", err)
	}
}
