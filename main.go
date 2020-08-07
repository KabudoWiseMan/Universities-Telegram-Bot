package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"net/http"
	"os"
	"strings"
)

func MainHandler(resp http.ResponseWriter, _ *http.Request) {
	resp.Write([]byte("Hi there! I'm Choose University bot!"))
}

func isAdmin(chatID int64) bool {
	return chatID == CreatorID
}

func main() {
	//db, err := connect()
	//if err != nil {
	//	log.Panic(err)
	//}

	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//updates := bot.ListenForWebhook("/" + bot.Token)
	updates, err := bot.GetUpdatesChan(u)

	http.HandleFunc("/", MainHandler)
	go http.ListenAndServe(":" + os.Getenv("PORT"), nil)

	users := InitUsers()

	for update := range updates {
		if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.Message.Chat.ID
			log.Printf("[%s u: %d c: %d] %s\n", update.CallbackQuery.From.UserName, update.CallbackQuery.From.ID, chatID, update.CallbackQuery.Data)

			msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "")
			msg.ParseMode = "markdown"

			switch update.CallbackQuery.Data {
			case "main":
				users.Delete(chatID)
				msg.Text = "Добро пожаловать в бота для подбора университета!\n\n" +
					"Здесь вы можете узнать, какие университеты подходят вам, исходя из ваших баллов ЕГЭ и других запросов."
				msg.ReplyMarkup = &mainMenu
			case "uni":
				msg.Text = "Введите один или несколько критериев для получения подборки университетов"
				msg.ReplyMarkup = &unisCompilationMenu
			case "fUni":
				msg.Text = "Введите название университета"
			case "rate":
				users.User(chatID).State = RatingQSState
				text, rateQSMenu := handleRatingQSRequest("rateQSPage#1")
				msg.Text = text
				msg.ReplyMarkup = &rateQSMenu
			default:
				data := update.CallbackQuery.Data
				if strings.Contains(data, "rateQSPage") {
					text, rateQSMenu := handleRatingQSRequest(data)
					msg.Text = text
					msg.ReplyMarkup = &rateQSMenu
				} else if strings.Contains(data, "getUni") {
					text, uniMenu := handleUniRequest(data)
					msg.Text = text
					msg.ReplyMarkup = &uniMenu
				} else if strings.Contains(data, "facs") {
					text, facsMenu := handleFacsRequest(data)
					msg.Text = text
					msg.ReplyMarkup = &facsMenu
				} else if strings.Contains(data, "back") {
					user := users.User(chatID)
					text, rateQSMenu := handleBackRequest(data, user)
					msg.Text = text
					msg.ReplyMarkup = &rateQSMenu
				}
				//else if strings.Contains(data, "getFac") {
				//	text, facMenu := handleFacRequest(data, chatID, users)
				//	msg.Text = text
				//	msg.ReplyMarkup = &facMenu
				//}
			}

			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Done"))
			bot.Send(msg)
		}

		if update.Message != nil {
			chatID := update.Message.Chat.ID
			userID := update.Message.From.ID

			log.Printf("[%s u: %d c: %d] %s\n", update.Message.From.UserName, userID, chatID, update.Message.Text)

			if update.Message.IsCommand() {
				msg := tgbotapi.NewMessage(chatID, "")
				switch update.Message.Command() {
				case "start", "help":
					users.Delete(chatID)
					msg.Text = "Добро пожаловать в бота для подбора университета!\n\n" +
						"Здесь вы можете узнать, какие университеты подходят вам, исходя из ваших баллов ЕГЭ и других запросов."
					msg.ReplyMarkup = mainMenu
				default:
					msg.Text = "У меня нет такой команды"
				}
				bot.Send(msg)
				continue
			}

			bot.Send(tgbotapi.NewMessage(chatID, "Я не знаю, что вам на это ответить"))
		}
	}

}