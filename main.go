package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"net/http"
	"os"
	"strconv"
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
				text, rateQSMenu := handleRatingQSRequest(1)
				msg.Text = text
				msg.ReplyMarkup = &rateQSMenu
			case "back":
				user := users.User(chatID)
				if user.State == RatingQSState {
					text, rateQSMenu := handleRatingQSRequest(user.Page)
					msg.Text = text
					msg.ReplyMarkup = &rateQSMenu
				}
			default:
				data := update.CallbackQuery.Data
				if strings.Contains(data, "rateQSPage") {
					splitted := strings.Split(data, "#")
					page, _ := strconv.Atoi(splitted[len(splitted) - 1])
					users.User(chatID).Page = page
					text, rateQSMenu := handleRatingQSRequest(page)
					msg.Text = text
					msg.ReplyMarkup = &rateQSMenu
				} else if strings.Contains(data, "getUni") {
					splitted := strings.Split(data, "&")
					uniId, _ := strconv.Atoi(splitted[len(splitted) - 1])
					text, uniMenu := handleUniRequest(uniId)
					msg.Text = text
					msg.ParseMode = "markdown"
					msg.ReplyMarkup = &uniMenu
				} else if strings.Contains(data, "facs") {
					splitted := strings.Split(data, "#")
					page, _ := strconv.Atoi(splitted[len(splitted) - 1])
					splitted2 := strings.Split(splitted[0], "&")
					uniId, _ := strconv.Atoi(splitted2[len(splitted2) - 1])
					text, facsMenu := handleFacsRequest(uniId, page)
					msg.Text = text
					msg.ReplyMarkup = &facsMenu
				} else if strings.Contains(data, "back") {
					if strings.Contains(data, "Uni") {
						splitted := strings.Split(data, "#")
						uniId, _ := strconv.Atoi(splitted[len(splitted) - 1])
						text, uniMenu := handleUniRequest(uniId)
						msg.Text = text
						msg.ReplyMarkup = &uniMenu
					}
				}
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

func handleUniRequest(uniId int) (string, tgbotapi.InlineKeyboardMarkup) {
	uni := getUniFromDb(uniId)
	text := makeTextUni(uni)
	uniMenu := makeUniMenu(uni)
	return text, uniMenu
}

func handleRatingQSRequest(page int) (string, tgbotapi.InlineKeyboardMarkup) {
	text := "*Международный рейтинг вузов QS*\n\n" +
		"Для более подробной информации посетите сайт QS, нажав на кнопку *Перейти на сайт QS*\n\n"

	unisQS := getUnisQSPageFromDb((page - 1) * 5)
	text += makeTextUnis(unisQS)

	unisQSNum := getUnisQSNumFromDb()
	rateQSMenu := makeRatingQsMenu(unisQSNum, unisQS, page)

	return text, rateQSMenu
}

func handleFacsRequest(uniId int, page int) (string, tgbotapi.InlineKeyboardMarkup) {
	uni := getUniFromDb(uniId)
	text := "*" + uni.Name + "*\n\n" +
		"Факультеты:\n\n"

	facs := getFacsPageFromDb(uniId, (page - 1) * 5)
	text += makeTextFacs(facs)

	facsNum := getFacsNumFromDb(uniId)
	facsMenu := makeFacsMenu(facsNum, uniId, facs, page)

	return text, facsMenu
}