package main

import (
	"gracia-bot/internal/telegram"
	"log"
	"os"
	"gracia-bot/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Ошибка загрузки .env")
	}

	token := os.Getenv("BOT_TOKEN")
	dbPassword := os.Getenv("DB_PASSWORD")

	db := database.Connect(dbPassword)

	defer db.Close()

	bot, err := telegram.NewBot(token,db)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Бот запущен")

	bot.Start()
	
	}
