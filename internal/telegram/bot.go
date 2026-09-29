package telegram

import (

	"log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	API *tgbotapi.BotAPI
}

func NewBot(token string) (*Bot,error){

	api, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Fatal(err)
	}

	return &Bot{
		API: api,
	},nil
}