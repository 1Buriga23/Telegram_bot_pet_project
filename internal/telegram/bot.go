package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Bot struct {
	API *tgbotapi.BotAPI
	DB *pgxpool.Pool
}

func NewBot(token string,db *pgxpool.Pool) (*Bot,error){

	api, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Fatal(err)
	}

	return &Bot{
		API: api,
		DB: db,
	},nil
}