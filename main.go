package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func MainHandler(resp http.ResponseWriter, _ *http.Request) {
	resp.Write([]byte("Hi there! I'm Choose University bot!"))
}

func isAdmin(chatID int64) bool {
	return chatID == CreatorID
}

func main() {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	bot.Debug = true

	db, err := connectToDb()
	if err != nil {
		log.Println("couldn't connected to data base", err)
	}
	log.Println("Successfully connected to data base!")
	defer closeDb(db)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//updates := bot.ListenForWebhook("/" + bot.Token)
	updates, err := bot.GetUpdatesChan(u)

	//http.HandleFunc("/", MainHandler)
	//go func() {
	//	panic(http.ListenAndServe(":" + os.Getenv("PORT"), nil))
	//}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGTERM)

	users := InitUsers()

	for {
		select {
		case update := <- updates:
			if update.CallbackQuery != nil {
				chatID := update.CallbackQuery.Message.Chat.ID
				log.Printf("[%s u: %d c: %d] %s\n", update.CallbackQuery.From.UserName, update.CallbackQuery.From.ID, chatID, update.CallbackQuery.Data)

				msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "")
				msg.ParseMode = "markdown"

				user := users.User(chatID)

				switch update.CallbackQuery.Data {
				case "main":
					//users.Delete(chatID)
					user.State = NoState
					msg.Text = "Добро пожаловать в бота для подбора университета!\n\n" +
						"Здесь вы можете узнать, какие университеты подходят вам, исходя из ваших баллов ЕГЭ и других запросов."
					msg.ReplyMarkup = &mainMenu
				case "uni":
					user.State = UniState
					text, unisCompilationMenu := handleUnisCompRequest(db, user)
					msg.Text = text
					msg.ReplyMarkup = &unisCompilationMenu
				case "fUni":
					user.State = FindUniState
					msg.Text = "Введите название университета"
					mainBackMenu := makeMainBackMenu("")
					msg.ReplyMarkup = &mainBackMenu
				case "rate":
					user.State = RatingQSState
					text, rateQSMenu := handleRatingQSRequest(db, "rateQSPage#1")
					msg.Text = text
					msg.ReplyMarkup = &rateQSMenu
				case "dorm":
					user.Dormatary = !user.Dormatary
					text, unisCompilationMenu := handleUnisCompRequest(db, user)
					msg.Text = text
					msg.ReplyMarkup = &unisCompilationMenu
				case "army":
					user.MilitaryDep = !user.MilitaryDep
					text, unisCompilationMenu := handleUnisCompRequest(db, user)
					msg.Text = text
					msg.ReplyMarkup = &unisCompilationMenu
				case "entry":
					user.EntryTest = !user.EntryTest
					text, unisCompilationMenu := handleUnisCompRequest(db, user)
					msg.Text = text
					msg.ReplyMarkup = &unisCompilationMenu
				case "fee":
					user.State = FeeState
					msg.Text = "Введите максимальную цену за год обучения"
					var backPattern string
					if user.Fee == 0 {
						backPattern = "uni"
					} else {
						backPattern = "chOrCl&" + strconv.Itoa(FeeState)
					}
					mainBackMenu := makeMainBackMenu(backPattern)
					msg.ReplyMarkup = &mainBackMenu
				default:
					data := update.CallbackQuery.Data
					if strings.Contains(data, "rateQSPage") {
						text, rateQSMenu := handleRatingQSRequest(db, data)
						msg.Text = text
						msg.ReplyMarkup = &rateQSMenu
					} else if strings.Contains(data, "getUni") {
						text, uniMenu := handleUniRequest(db, data)
						msg.Text = text
						msg.ReplyMarkup = &uniMenu
					} else if strings.Contains(data, "facs") {
						text, facsMenu := handleFacsRequest(db, data)
						msg.Text = text
						msg.ReplyMarkup = &facsMenu
					} else if strings.Contains(data, "back") {
						text, rateQSMenu := handleBackRequest(db, data, user)
						msg.Text = text
						msg.ReplyMarkup = &rateQSMenu
					} else if strings.Contains(data, "getFac") {
						text, facMenu := handleFacRequest(db, data)
						msg.Text = text
						msg.ReplyMarkup = &facMenu
					} else if strings.Contains(data, "findUniPage") {
						text, findUniMenu := handleFindUniRequest(db, user.Query + "#" + data)
						msg.Text = text
						if len(findUniMenu.InlineKeyboard) != 0 {
							msg.ReplyMarkup = &findUniMenu
						}
					} else if strings.Contains(data, "profs") {
						text, profsMenu := handleProfsRequest(db, data)
						msg.Text = text
						msg.ReplyMarkup = &profsMenu
					} else if strings.Contains(data, "specs") {
						text, specsMenu := handleSpecsRequest(db, data)
						msg.Text = text
						msg.ReplyMarkup = &specsMenu
					} else if strings.Contains(data, "progs") {
						text, progsMenu := handleProgsRequest(db, data)
						msg.Text = text
						msg.ReplyMarkup = &progsMenu
					} else if strings.Contains(data, "getProg") {
						text, progMenu := handleProgRequest(db, data)
						msg.Text = text
						msg.ReplyMarkup = &progMenu
					}  else if strings.Contains(data, "setCity") {
						cityId, _ := strconv.Atoi(takeId(data))
						user.City = cityId
						text, unisCompilationMenu := handleUnisCompRequest(db, user)
						msg.Text = text
						msg.ReplyMarkup = &unisCompilationMenu
					} else if strings.Contains(data, "city") {
						text, citiesMenu := handleCitiesRequest(db, data, user)
						msg.Text = text
						msg.ReplyMarkup = &citiesMenu
					} else if strings.Contains(data, "proOrSpe") {
						text, specOrNotMenu := handleSpecOrNotRequest(data)
						msg.Text = text
						msg.ReplyMarkup = &specOrNotMenu
					} else if strings.Contains(data, "setPro") {
						profId, _ := strconv.Atoi(takeId(data))
						user.ProfileId = profId
						user.SpecialityId = 0
						text, unisCompilationMenu := handleUnisCompRequest(db, user)
						msg.Text = text
						msg.ReplyMarkup = &unisCompilationMenu
					} else if strings.Contains(data, "pro") {
						text, profilesMenu := handleProfilesRequest(db, data, user)
						msg.Text = text
						msg.ReplyMarkup = &profilesMenu
					} else if strings.Contains(data, "setSpe") {
						ids := takeIds(data)
						profId, _ := strconv.Atoi(ids[0])
						specId, _ := strconv.Atoi(ids[1])
						user.ProfileId = profId
						user.SpecialityId = specId
						text, unisCompilationMenu := handleUnisCompRequest(db, user)
						msg.Text = text
						msg.ReplyMarkup = &unisCompilationMenu
					} else if strings.Contains(data, "spe") {
						text, specialitiesMenu := handleSpecialitiesRequest(db, data)
						msg.Text = text
						msg.ReplyMarkup = &specialitiesMenu
					} else if strings.Contains(data, "chOrCl") {
						user.State = UniState
						text, changeOrClearMenu := handleChangeOrClearRequest(db, data, user)
						msg.Text = text
						msg.ReplyMarkup = &changeOrClearMenu
					} else if strings.Contains(data, "setEge") {
						subjId, _ := strconv.Atoi(takeId(data))
						user.Eges = append(user.Eges, Ege{SubjId: subjId, MinPoints: 100})
						text, egesMenu := handleEgesRequest(db, "ege#1", user)
						msg.Text = text
						msg.ReplyMarkup = &egesMenu
					} else if strings.Contains(data, "ege") {
						user.State = UniState
						text, egesMenu := handleEgesRequest(db, data, user)
						msg.Text = text
						msg.ReplyMarkup = &egesMenu
					} else if strings.Contains(data, "subj") {
						user.State = EgeState
						text, subjMenu := handleSubjRequest(data, user)
						msg.Text = text
						msg.ReplyMarkup = &subjMenu
					} else if strings.Contains(data, "chPoints") {
						page := takePage(data)
						subjName := getSubjNameFromDb(db, user.LastSubj)
						user.State = EgeState
						msg.Text = "Введите баллы ЕГЭ по предмету *" + subjName + "*"
						mainBackMenu := makeMainBackMenu("chOrCl&" + strconv.Itoa(SubjState) + "&" + strconv.Itoa(user.LastSubj) + "#" + page)
						msg.ReplyMarkup = &mainBackMenu
					} else if strings.Contains(data, "clear") {
						state, _ := strconv.Atoi(takeId(data))
						var text string
						var menu tgbotapi.InlineKeyboardMarkup
						switch state {
						case EgeState:
							user.Eges = nil
							text, menu = handleEgesRequest(db, "ege#1", user)
						case FeeState:
							user.Fee = math.MaxUint64
							text, menu = handleUnisCompRequest(db, user)
						case CityState:
							user.City = 0
							text, menu = handleUnisCompRequest(db, user)
						case ProfileState:
							user.ProfileId = 0
							text, menu = handleUnisCompRequest(db, user)
						case SpecialityState:
							user.SpecialityId = 0
							text, menu = handleUnisCompRequest(db, user)
						case SubjState:
							user.DeleteEge()
							text, menu = handleChangeOrClearRequest(db, "chOrCl&" + strconv.Itoa(EgeState) + "#1", user)
						case UniState:
							user.Clear()
							text, menu = handleUnisCompRequest(db, user)
						}
						msg.Text = text
						msg.ReplyMarkup = &menu
					} else if strings.Contains(data, "search") {
						user.State = UniState
						text, searchUniMenu := handleSearchUniRequest(db, data, user)
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

				if user.State == UniState {
					msg.Text = "Что-то здесь не так, проверьте, что вы нажали нужную кнопку" + makeEmoji(WinkEmoji)
				} else if user.State == RatingQSState {
					msg.Text = "Поиск здесь недоступен"
				} else if user.State == FeeState {
					feeStr := update.Message.Text
					fee, err := strconv.ParseUint(feeStr, 10, 64)
					if err != nil {
						msg.Text = "Пожалуйста, введите корректную сумму"
					} else {
						user.State = UniState
						user.Fee = fee
						text, unisCompilationMenu := handleUnisCompRequest(db, user)
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
							text, menu = handleChangeOrClearRequest(db, "chOrCl&" + strconv.Itoa(EgeState) + "#1", user)
						} else {
							text, menu = handleEgesRequest(db, "ege#1", user)
						}

						msg.Text = text
						msg.ReplyMarkup = &menu
					}
				} else {
					user.State = FindUniState
					user.Query = update.Message.Text
					text, findUniMenu := handleFindUniRequest(db, update.Message.Text + "#1")
					msg.Text = text
					if len(findUniMenu.InlineKeyboard) != 0 {
						msg.ReplyMarkup = &findUniMenu
					}
				}

				bot.Send(msg)
			}
		case <-stop:
			log.Println("Got interrupt signal. Aborting...")
			return
		}
	}
}