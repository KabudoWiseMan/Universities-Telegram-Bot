package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"math"
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
				user := users.User(chatID)
				user.State = UniState
				text, unisCompilationMenu := handleUnisCompRequest(user)
				msg.Text = text
				msg.ReplyMarkup = &unisCompilationMenu
			case "fUni":
				users.User(chatID).State = FindUniState
				msg.Text = "Введите название университета"
				mainBackMenu := makeMainBackMenu("")
				msg.ReplyMarkup = &mainBackMenu
			case "rate":
				users.User(chatID).State = RatingQSState
				text, rateQSMenu := handleRatingQSRequest("rateQSPage#1")
				msg.Text = text
				msg.ReplyMarkup = &rateQSMenu
			case "dorm":
				user := users.User(chatID)
				user.Dormatary = !user.Dormatary
				text, unisCompilationMenu := handleUnisCompRequest(user)
				msg.Text = text
				msg.ReplyMarkup = &unisCompilationMenu
			case "army":
				user := users.User(chatID)
				user.MilitaryDep = !user.MilitaryDep
				text, unisCompilationMenu := handleUnisCompRequest(user)
				msg.Text = text
				msg.ReplyMarkup = &unisCompilationMenu
			case "entry":
				user := users.User(chatID)
				user.EntryTest = !user.EntryTest
				text, unisCompilationMenu := handleUnisCompRequest(user)
				msg.Text = text
				msg.ReplyMarkup = &unisCompilationMenu
			case "fee":
				user := users.User(chatID)
				user.State = FeeState
				msg.Text = "Введите максимальную цену за год обучения"
				mainBackMenu := makeMainBackMenu("uni")
				msg.ReplyMarkup = &mainBackMenu
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
				} else if strings.Contains(data, "getFac") {
					text, facMenu := handleFacRequest(data)
					msg.Text = text
					msg.ReplyMarkup = &facMenu
				} else if strings.Contains(data, "findUniPage") {
					text, findUniMenu := handleFindUniRequest(users.User(chatID).Query + "#" + data)
					msg.Text = text
					if len(findUniMenu.InlineKeyboard) != 0 {
						msg.ReplyMarkup = &findUniMenu
					}
				} else if strings.Contains(data, "profs") {
					text, profsMenu := handleProfsRequest(data)
					msg.Text = text
					msg.ReplyMarkup = &profsMenu
				} else if strings.Contains(data, "specs") {
					text, specsMenu := handleSpecsRequest(data)
					msg.Text = text
					msg.ReplyMarkup = &specsMenu
				} else if strings.Contains(data, "progs") {
					text, progsMenu := handleProgsRequest(data)
					msg.Text = text
					msg.ReplyMarkup = &progsMenu
				} else if strings.Contains(data, "getProg") {
					text, progMenu := handleProgRequest(data)
					msg.Text = text
					msg.ReplyMarkup = &progMenu
				}  else if strings.Contains(data, "setCity") {
					cityId, _ := strconv.Atoi(takeId(data))
					user := users.User(chatID)
					user.City = cityId
					text, unisCompilationMenu := handleUnisCompRequest(user)
					msg.Text = text
					msg.ReplyMarkup = &unisCompilationMenu
				} else if strings.Contains(data, "city") {
					user := users.User(chatID)
					//user.State = CityState
					text, citiesMenu := handleCitiesRequest(data, user)
					msg.Text = text
					msg.ReplyMarkup = &citiesMenu
				} else if strings.Contains(data, "proOrSpe") {
					user := users.User(chatID)
					text, specOrNotMenu := handleSpecOrNotRequest(data, user)
					msg.Text = text
					msg.ReplyMarkup = &specOrNotMenu
				} else if strings.Contains(data, "setPro") {
					user := users.User(chatID)
					profId, _ := strconv.Atoi(takeId(data))
					user.ProfileId = profId
					user.SpecialityId = 0
					text, unisCompilationMenu := handleUnisCompRequest(user)
					msg.Text = text
					msg.ReplyMarkup = &unisCompilationMenu
				} else if strings.Contains(data, "pro") {
					user := users.User(chatID)
					//user.State = ProfileState
					text, profilesMenu := handleProfilesRequest(data, user)
					msg.Text = text
					msg.ReplyMarkup = &profilesMenu
				} else if strings.Contains(data, "setSpe") {
					ids := takeIds(data)
					profId, _ := strconv.Atoi(ids[0])
					specId, _ := strconv.Atoi(ids[1])
					user := users.User(chatID)
					user.ProfileId = profId
					user.SpecialityId = specId
					text, unisCompilationMenu := handleUnisCompRequest(user)
					msg.Text = text
					msg.ReplyMarkup = &unisCompilationMenu
				} else if strings.Contains(data, "spe") {
					//user.State = SpecialityState
					text, specialitiesMenu := handleSpecialitiesRequest(data)
					msg.Text = text
					msg.ReplyMarkup = &specialitiesMenu
				} else if strings.Contains(data, "chOrCl") {
					user := users.User(chatID)
					text, changeOrClearMenu := handleChangeOrClearRequest(data, user)
					msg.Text = text
					msg.ReplyMarkup = &changeOrClearMenu
				} else if strings.Contains(data, "setEge") {
					subjId, _ := strconv.Atoi(takeId(data))
					user := users.User(chatID)
					user.Eges = append(user.Eges, Ege{SubjId: subjId, MinPoints: 100})
					text, egesMenu := handleEgesRequest("ege#1", user)
					msg.Text = text
					msg.ReplyMarkup = &egesMenu
				} else if strings.Contains(data, "ege") {
					user := users.User(chatID)
					text, egesMenu := handleEgesRequest(data, user)
					msg.Text = text
					msg.ReplyMarkup = &egesMenu
				} else if strings.Contains(data, "subj") {
					user := users.User(chatID)
					text, subjMenu := handleSubjRequest(data, user)
					msg.Text = text
					msg.ReplyMarkup = &subjMenu
				} else if strings.Contains(data, "chPoints") {
					user := users.User(chatID)
					page := takePage(data)
					subjName := getSubjNameFromDb(user.LastSubj)
					user.State = EgeState
					msg.Text = "Введите баллы ЕГЭ по предмету *" + subjName + "*"
					mainBackMenu := makeMainBackMenu("chOrCl&" + strconv.Itoa(SubjState) + "&" + strconv.Itoa(user.LastSubj) + "#" + page)
					msg.ReplyMarkup = &mainBackMenu
				} else if strings.Contains(data, "points") {
					subjId, _ := strconv.Atoi(takeId(data))
					page := takePage(data)
					subjName := getSubjNameFromDb(subjId)
					user := users.User(chatID)
					user.LastSubj = subjId
					user.State = EgeState
					msg.Text = "Введите баллы ЕГЭ по предмету *" + subjName + "*"
					mainBackMenu := makeMainBackMenu("subj&" + strconv.Itoa(subjId) + "#" + page)
					msg.ReplyMarkup = &mainBackMenu
				} else if strings.Contains(data, "clear") {
					state, _ := strconv.Atoi(takeId(data))
					user := users.User(chatID)
					var text string
					var menu tgbotapi.InlineKeyboardMarkup
					switch state {
					case EgeState:
						user.Eges = nil
						text, menu = handleEgesRequest("ege#1", user)
					case FeeState:
						user.Fee = math.MaxUint64
						text, menu = handleUnisCompRequest(user)
					case CityState:
						user.City = 0
						text, menu = handleUnisCompRequest(user)
					case ProfileState:
						user.ProfileId = 0
						text, menu = handleUnisCompRequest(user)
					case SpecialityState:
						user.SpecialityId = 0
						text, menu = handleUnisCompRequest(user)
					case SubjState:
						user.DeleteEge()
						text, menu = handleChangeOrClearRequest("chOrCl&" + strconv.Itoa(EgeState) + "#1", user)
					case UniState:
						user.Clear()
						text, menu = handleUnisCompRequest(user)
					}
					msg.Text = text
					msg.ReplyMarkup = &menu
				} else if strings.Contains(data, "search") {
					user := users.User(chatID)
					user.State = UniState
					text, searchUniMenu := handleSearchUniRequest(data, user)
					msg.Text = text
					if len(searchUniMenu.InlineKeyboard) != 0 {
						msg.ReplyMarkup = &searchUniMenu
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

			msg := tgbotapi.NewMessage(chatID, "")
			msg.ParseMode = "markdown"
			if update.Message.IsCommand() {
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

			user := users.User(chatID)

			if user.State == FeeState {
				feeStr := update.Message.Text
				fee, err := strconv.ParseUint(feeStr, 10, 64)
				if err != nil {
					msg.Text = "Пожалуйста, введите корректную сумму"
				} else {
					user.State = UniState
					user.Fee = fee
					text, unisCompilationMenu := handleUnisCompRequest(user)
					msg.Text = text
					msg.ReplyMarkup = &unisCompilationMenu
				}
			} else if user.State == EgeState {
				pointsStr := update.Message.Text
				points, err := strconv.ParseUint(pointsStr, 10, 64)
				if err != nil || points < 0 || points > 100 {
					msg.Text = "Пожалуйста, введите корректные баллы"
				} else {
					user.State = UniState
					found := user.AddEge(points)
					user.LastSubj = 0
					var text string
					var menu tgbotapi.InlineKeyboardMarkup
					if found {
						text, menu = handleChangeOrClearRequest("chOrCl&" + strconv.Itoa(EgeState) + "#1", user)
					} else {
						text, menu = handleEgesRequest("ege#1", user)
					}

					msg.Text = text
					msg.ReplyMarkup = &menu
				}
			} else {
				user.State = FindUniState
				user.Query = update.Message.Text
				text, findUniMenu := handleFindUniRequest(update.Message.Text + "#1")
				msg.Text = text
				if len(findUniMenu.InlineKeyboard) != 0 {
					msg.ReplyMarkup = &findUniMenu
				}
			}

			bot.Send(msg)
		}
	}

}