package main

import (
	"github.com/mamedvedkov/BankBot/internals/bot"
	"github.com/mamedvedkov/BankBot/internals/repository"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	r := repository.NewRepo()
	b := bot.NewBot(r)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go b.BotWork()

	<-done

	log.Print("Бот отключился")
}
