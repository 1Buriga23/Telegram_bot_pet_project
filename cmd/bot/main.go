package main

import (
	"gracia-bot/internal/telegram"
	"log"
	"os"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Ошибка загрузки .env")
	}

	token := os.Getenv("BOT_TOKEN")

	bot, err := telegram.NewBot(token)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Бот запущен")

	bot.Start()
	
	}
