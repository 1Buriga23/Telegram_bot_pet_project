package main

import (
	"gracia-bot/internal/database"
	"gracia-bot/internal/telegram"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			log.Fatal("Ошибка загрузки .env:", err)
		}
	}

	token := os.Getenv("BOT_TOKEN")
	dbPassword := os.Getenv("DB_PASSWORD")

	if token == "" {
		log.Fatal("Переменная BOT_TOKEN не задана")
	}

	if dbPassword == "" {
		log.Fatal("Переменная DB_PASSWORD не задана")
	}

	db := database.Connect(dbPassword)

	defer db.Close()

	bot, err := telegram.NewBot(token, db)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Бот запущен")

	bot.Start()

}
