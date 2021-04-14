package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mamedvedkov/BankBot/internals"
	"github.com/mamedvedkov/BankBot/internals/repository"
	"log"
	"time"
)

// TODO: в енв
const (
	token      string = "1010464119:AAHIAE-Z_XK3A3Zuawjm_2bWyfEwn1rExgg"
	mainChatId int64  = 42
	adminId    int64  = 193893846
)

type bankBot struct {
	repo *repository.Repo
	bot  *tgbotapi.BotAPI
}

func NewBot(repo *repository.Repo) *bankBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic("Всё упало на создании бота")
	}

	log.Printf("Бот включился, имя: %s", bot.Self.UserName)
	return &bankBot{
		repo: repository.NewRepo(),
		bot:  bot,
	}
}

func (b *bankBot) BotWork() {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// TODO: чтение из канала несколькими воркерами??
	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic("Всё упало на создании канала апдейтов")
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		// TODO: роутер запросов
		// TODO: регистрация через админов
		cmd := update.Message.Command()
		if cmd == "" {
			continue
		}

		if i, _ := internals.GetRowByTgId(b.repo, update.Message.From.ID);
			i != 0 && update.Message.Chat.ID != mainChatId {
			b.sendMsgToAdmin(update)
			continue
		}

		log.Printf("cmd from [%s] id = %v cmd = %s", update.Message.From.UserName, update.Message.From.ID, cmd)

		t := time.Now()
		response := internals.Process(b.repo, cmd, update.Message.Chat.ID == mainChatId, update)

		log.Printf("time taken for respone: %v", time.Since(t))

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		msg.ParseMode = "Markdown"

		_, err = b.bot.Send(msg)
		if err != nil {
			log.Printf("error in send to admin: %w", err)
		}
	}
}

func (b *bankBot) sendMsgToAdmin(update tgbotapi.Update) bool {
	usr := update.Message.From

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Мы вас запомнили, информация отправлена админам")

	 _, err := b.bot.Send(msg)
	if err != nil {
		log.Printf("error in send to admin: %w", err)

		return false
	}

	msg = tgbotapi.NewMessage(adminId, fmt.Sprintf("aboutme: id=%v, firstName=%v, lastName=%v",
		usr.ID,
		usr.FirstName,
		usr.LastName,
	))

	_, err = b.bot.Send(msg)
	if err != nil {
		log.Printf("error in send to admin: %w", err)

		return false
	}

	return true
}
