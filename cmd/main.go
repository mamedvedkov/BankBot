package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mamedvedkov/BankBot/internals"
)

const token string = "1010464119:AAHIAE-Z_XK3A3Zuawjm_2bWyfEwn1rExgg"
const botName string = "kjhgskdf743_bot"
const mainChatId int64 = 42

func main() {
	adapter := internals.NewAdapter()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go botWork(adapter)

	//пример GetValues
	//adapter.GetValues("A1:D16")

	<-done

	log.Print("Бот отключился")
}

func botWork(adapter *internals.Adapter) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic("Всё упало на создании бота")
	}

	log.Printf("Бот включился, имя: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic("Всё упало на создании канала апдейтов")
	}
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		// todo: роутер запросов
		// принимает adapter *App, upd update
		// update.Message.Command() была ли команда
		if cmd := update.Message.Command(); cmd != "" {
			log.Printf("from [%s] id = %v cmd = %s", update.Message.From.UserName, update.Message.From.ID, cmd)
			response := internals.Process(adapter, cmd, update.Message.Chat.ID == mainChatId, update)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
			bot.Send(msg)
		}

		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID

		//bot.Send(msg)
	}
}
