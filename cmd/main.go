package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/tapmahtec/TNL_bot/handlers"
	"github.com/tapmahtec/TNL_bot/repository"
	"github.com/tapmahtec/TNL_bot/service"
)

func main() {

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Failed to load environment: %v", err.Error())
	}

	db, err := repository.NewMySQLDB(repository.Config{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("PORT"),
		User:     os.Getenv("USER"),
		Password: os.Getenv("PASSWORD"),
		DBname:   os.Getenv("DBNAME"),
	})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	err = repository.InitDB(db)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	repo := repository.NewRepository(db)
	services := service.NewService(repo)

	bot, err := handlers.NewBot(os.Getenv("TOKEN"), os.Getenv("ALLOWED_CHANNEL_ID"), services)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	if err := bot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}

	fmt.Println("Bot has been stopped.")
}
