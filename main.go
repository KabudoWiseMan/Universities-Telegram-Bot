package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// ***FOR WEBHOOK***
//func MainHandler(resp http.ResponseWriter, _ *http.Request) {
//	resp.Write([]byte("Hi there! I'm Choose University bot!"))
//}

func isAdmin(chatID int64) bool {
	return chatID == CreatorID
}

func monitorUsers(ticker *time.Ticker, users *Users) {
	for {
		select {
		case <-ticker.C:
			var toDelete []int64
			for key, user := range users.Users {
				if time.Since(user.LastSeen).Hours() > 1 {
					toDelete = append(toDelete, key)
				}
			}

			for _, toDelId := range toDelete {
				users.Delete(toDelId)
			}
		}
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Panic("couldn't connect to bot", err)
	}
	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	bot.Debug = true

	db, err := connectToDb()
	if err != nil {
		log.Panic("couldn't connected to data base", err)
	}
	log.Println("Successfully connected to data base")
	defer closeDb(db)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	// ***FOR WEBHOOK***
	//updates := bot.ListenForWebhook("/" + bot.Token)
	//http.HandleFunc("/", MainHandler)
	//go func() {
	//	panic(http.ListenAndServe(":" + os.Getenv("PORT"), nil))
	//}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGTERM)

	users := InitUsers()

	usersTicker := time.NewTicker(30 * time.Minute)
	go monitorUsers(usersTicker, users)

	for {
		select {
		case update := <- updates:
			if update.CallbackQuery != nil {
				chatID := update.CallbackQuery.Message.Chat.ID
				log.Printf("[%s u: %d c: %d] %s\n", update.CallbackQuery.From.UserName, update.CallbackQuery.From.ID, chatID, update.CallbackQuery.Data)

				msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "")

				user := users.User(chatID)
				user.LastSeen = time.Now()

				var menu tgbotapi.InlineKeyboardMarkup
				var text string

				switch update.CallbackQuery.Data {
				case "main":
					user.State = NoState
					text = "Добро пожаловать в бота для подбора университета!\n\n" +
						"Здесь вы можете узнать, какие университеты подходят вам, исходя из ваших баллов ЕГЭ и других запросов."
					menu = mainMenu
				case "uni":
					user.State = UniState
					text, menu = handleUnisCompRequest(db, user)
				case "fUni":
					user.State = FindUniState
					text = "Введите название университета"
					menu = makeMainBackMenu("")
				case "rate":
					user.State = RatingQSState
					text, menu = handleRatingQSRequest(db, "rateQSPage#1")
				case "dorm":
					user.Dormatary = !user.Dormatary
					text, menu = handleUnisCompRequest(db, user)
				case "army":
					user.MilitaryDep = !user.MilitaryDep
					text, menu = handleUnisCompRequest(db, user)
				case "entry":
					user.EntryTest = !user.EntryTest
					text, menu = handleUnisCompRequest(db, user)
				case "fee":
					user.State = FeeState
					text = "Введите максимальную цену за год обучения"
					var backPattern string
					if user.Fee == math.MaxUint64 {
						backPattern = "uni"
					} else {
						backPattern = "chOrCl&" + strconv.Itoa(FeeState)
					}
					menu = makeMainBackMenu(backPattern)
				default:
					data := update.CallbackQuery.Data
					if strings.Contains(data, "rateQSPage") {
						text, menu = handleRatingQSRequest(db, data)
					} else if strings.Contains(data, "getUni") {
						text, menu = handleUniRequest(db, data)
					} else if strings.Contains(data, "facs") {
						text, menu = handleFacsRequest(db, data)
					} else if strings.Contains(data, "back") {
						text, menu = handleBackRequest(db, data, user)
					} else if strings.Contains(data, "getFac") {
						text, menu = handleFacRequest(db, data)
					} else if strings.Contains(data, "findUniPage") {
						text, menu = handleFindUniRequest(db, user.Query + "#" + data)
					} else if strings.Contains(data, "profs") {
						text, menu = handleProfsRequest(db, data)
					} else if strings.Contains(data, "specs") {
						text, menu = handleSpecsRequest(db, data)
					} else if strings.Contains(data, "progs") {
						text, menu = handleProgsRequest(db, data)
					} else if strings.Contains(data, "getProg") {
						text, menu = handleProgRequest(db, data)
					}  else if strings.Contains(data, "setCity") {
						cityId, _ := strconv.Atoi(takeId(data))
						user.City = cityId
						text, menu = handleUnisCompRequest(db, user)
					} else if strings.Contains(data, "city") {
						text, menu = handleCitiesRequest(db, data, user)
					} else if strings.Contains(data, "proOrSpe") {
						text, menu = handleSpecOrNotRequest(data)
					} else if strings.Contains(data, "setPro") {
						profId, _ := strconv.Atoi(takeId(data))
						user.ProfileId = profId
						user.SpecialityId = 0
						text, menu = handleUnisCompRequest(db, user)
					} else if strings.Contains(data, "pro") {
						text, menu = handleProfilesRequest(db, data, user)
					} else if strings.Contains(data, "setSpe") {
						ids := takeIds(data)
						profId, _ := strconv.Atoi(ids[0])
						specId, _ := strconv.Atoi(ids[1])
						user.ProfileId = profId
						user.SpecialityId = specId
						text, menu = handleUnisCompRequest(db, user)
					} else if strings.Contains(data, "spe") {
						text, menu = handleSpecialitiesRequest(db, data)
					} else if strings.Contains(data, "chOrCl") {
						user.State = UniState
						text, menu = handleChangeOrClearRequest(db, data, user)
					} else if strings.Contains(data, "setEge") {
						subjId, _ := strconv.Atoi(takeId(data))
						user.Eges = append(user.Eges, Ege{SubjId: subjId, MinPoints: 100})
						text, menu = handleEgesRequest(db, "ege#1", user)
					} else if strings.Contains(data, "ege") {
						user.State = UniState
						text, menu = handleEgesRequest(db, data, user)
					} else if strings.Contains(data, "subj") {
						user.State = EgeState
						text, menu = handleSubjRequest(data, user)
					} else if strings.Contains(data, "chPoints") {
						user.State = EgeState
						text, menu = handleChangePointsRequest(db, data, user)
					} else if strings.Contains(data, "clear") {
						text, menu = handleClearMenu(db, data, user)
					} else if strings.Contains(data, "search") {
						user.State = UniState
						text, menu = handleSearchUniRequest(db, data, user)
					}
				}

				msg.Text = text
				msg.ParseMode = "markdown"
				msg.ReplyMarkup = &menu

				bot.Send(msg)
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Done"))
			}

			if update.Message != nil {
				chatID := update.Message.Chat.ID
				userID := update.Message.From.ID

				log.Printf("[%s u: %d c: %d] %s\n", update.Message.From.UserName, userID, chatID, update.Message.Text)

				user := users.User(chatID)

				msg := tgbotapi.NewMessage(chatID, "")
				msg.ParseMode = "markdown"
				if update.Message.IsCommand() {
					switch update.Message.Command() {
					case "start", "help":
						msg.Text = "Добро пожаловать в бота для подбора университета!\n\n" +
							"Здесь вы можете узнать, какие университеты подходят вам, исходя из ваших баллов ЕГЭ и других запросов."
						msg.ReplyMarkup = mainMenu
					case "update":
						if isAdmin(chatID) {
							msg.Text = "Обновление началось"
							go updateDb()
						} else {
							msg.Text = "У меня нет такой команды"
						}
					case "updateU":
						if isAdmin(chatID) {
							msg.Text = "Обновление университетов началось"
							go updateUnis()
						} else {
							msg.Text = "У меня нет такой команды"
						}
					case "updateF":
						if isAdmin(chatID) {
							msg.Text = "Обновление факультетов началось"
							go updateFacs()
						} else {
							msg.Text = "У меня нет такой команды"
						}
					case "updateP":
						if isAdmin(chatID) {
							msg.Text = "Обновление программ началось"
							go updateProgsNInfo()
						} else {
							msg.Text = "У меня нет такой команды"
						}
					case "updateC":
						if isAdmin(chatID) {
							msg.Text = "Обновление городов началось"
							go updateCities()
						} else {
							msg.Text = "У меня нет такой команды"
						}
					case "updateS":
						if isAdmin(chatID) {
							msg.Text = "Обновление предметов началось"
							go updateSubjs()
						} else {
							msg.Text = "У меня нет такой команды"
						}
					case "updatePS":
						if isAdmin(chatID) {
							msg.Text = "Обновление профилей и специальностей началось"
							go updateProfsNSpecs()
						} else {
							msg.Text = "У меня нет такой команды"
						}
					case "updateR":
						if isAdmin(chatID) {
							msg.Text = "Обновление рейтинга QS началось"
							go updateRatingQS()
						} else {
							msg.Text = "У меня нет такой команды"
						}
					default:
						msg.Text = "У меня нет такой команды"
					}
					bot.Send(msg)
					continue
				}

				user.LastSeen = time.Now()

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