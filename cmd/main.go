package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mamedvedkov/BankBot/internals"
)

// TODO: убрать в енв
const token string = "1010464119:AAHIAE-Z_XK3A3Zuawjm_2bWyfEwn1rExgg"
const botName string = "kjhgskdf743_bot"
const mainChatId int64 = 42

func main() {
	adapter := internals.NewRepo()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go botWork(adapter)

	<-done

	log.Print("Бот отключился")
}

func botWork(repo *internals.Repo) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic("Всё упало на создании бота")
	}

	log.Printf("Бот включился, имя: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// TODO: чтение из канала несколькими воркерами??
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic("Всё упало на создании канала апдейтов")
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		// TODO: роутер запросов
		if cmd := update.Message.Command(); cmd != "" {
			log.Printf("from [%s] id = %v cmd = %s", update.Message.From.UserName, update.Message.From.ID, cmd)
			response := internals.Process(repo, cmd, update.Message.Chat.ID == mainChatId, update)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
			msg.ParseMode = "Markdown"
			bot.Send(msg)
		}
	}
}
