package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"net/http"
	"os"
)

func MainHandler(resp http.ResponseWriter, _ *http.Request) {
	resp.Write([]byte("Hi there! I'm Choose University bot!"))
}

func isAdmin(chatID int64) bool {
	return chatID == CreatorID
}

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("1"),
		tgbotapi.NewKeyboardButton("2"),
		tgbotapi.NewKeyboardButton("3"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("4"),
		tgbotapi.NewKeyboardButton("5"),
		tgbotapi.NewKeyboardButton("6"),
	),
)

func main() {
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

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		userID := int64(update.Message.From.ID)

		log.Printf("[%s u: %d c: %d] %s\n", update.Message.From.UserName, userID, chatID, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		switch update.Message.Text {
		case "open":
			msg.ReplyMarkup = numericKeyboard
		case "close":
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		}

		bot.Send(msg)

		//if update.Message.IsCommand() {
		//	msg := tgbotapi.NewMessage(chatID, "")
		//	switch update.Message.Command() {
		//	case "start":
		//		msg.Text = "Добро пожаловать в бота для записи на лабораторные работы. " +
		//			"Ознакомьтесь с моими возможностями, набрав /help"
		//	case "help":
		//		msg.Text = "Мне доступны следующие команды:\n\n" +
		//			"/list — посмотреть текущий список\n" +
		//			"/insert l1 l2 ... ln — добавить себя и лабы, которые вы хотите сдать, в список; " +
		//			"где li — это номер лабы; пример: /insert 5 6 7\n" +
		//			"ВНИМАНИЕ: люди добавляются в список по приоритету наименьшей лабы, по которой ещё не прошёл дедлайн.\n" +
		//			"/delete — удалить себя из списка\n" +
		//			"/success l1 l2 ... ln — при успешной сдаче лаб; если отправлена без аргументов, то вы будете удалены из списка," +
		//			"если аргументы (номера лаб) есть, то вы будете добавлены в список заново только с " +
		//			"неуказанными (несданными) вами лабами\n" +
		//			"/fail — при неуспешной сдаче лаб; вы будете перемещены в конец списка, сохраняя приоритет, если он у вас был\n" +
		//			"/deadlinelab - посмотреть минимальную лабу, по которой сейчас дедлайн"
		//	case "list":
		//		msg.Text = listToString(list)
		//	case "insert":
		//		args := update.Message.CommandArguments()
		//
		//		if len(args) == 0 {
		//			msg.Text = "Недостаточно аргументов, пожалуйста, введите лабы, которые собираетесь сдавать"
		//		} else {
		//			labs, err := getLabsFromString(args)
		//			if err != nil {
		//				msg.Text = "Неверно введены лабы, пожалуйста, введите номера лаб через пробел, " +
		//					"каждая лаба должна быть не меньше 1 и не больше 15"
		//			} else {
		//				var resp int
		//				list, resp = insert(list, keys, userID, labs, deadlineLab)
		//				writeToSheet(srv, SpreadsheetID, listToInrerface(list, keys), QueueRows)
		//				msg.Text = processResponse(resp)
		//			}
		//		}
		//	case "delete":
		//		var resp int
		//		list, resp = deleteElem(list, keys, userID)
		//		if resp == Delete {
		//			writeToSheet(srv, SpreadsheetID, listToInrerface(list, keys), QueueRows)
		//		}
		//		msg.Text = processResponse(resp)
		//	case "success":
		//		args := update.Message.CommandArguments()
		//		var resp int
		//
		//		if len(args) == 0 {
		//			list, resp = deleteElem(list, keys, userID)
		//			if resp == Delete {
		//				writeToSheet(srv, SpreadsheetID, listToInrerface(list, keys), QueueRows)
		//			}
		//			msg.Text = processResponse(resp)
		//		} else {
		//			labs, err := getLabsFromString(args)
		//			if err != nil {
		//				msg.Text = "Неверно введены лабы, пожалуйста, введите номера лаб через пробел"
		//			} else {
		//				ok, remainingLabs := subOfArrays(list[userID].labsToPass, labs)
		//				if ok {
		//					list, resp = deleteElem(list, keys, userID)
		//					msg.Text = processResponse(resp)
		//					if resp == Delete {
		//						writeToSheet(srv, SpreadsheetID, listToInrerface(list, keys), QueueRows)
		//						if len(remainingLabs) != 0 {
		//							list, resp = insert(list, keys, userID, remainingLabs, deadlineLab)
		//							writeToSheet(srv, SpreadsheetID, listToInrerface(list, keys), QueueRows)
		//							msg.Text = processResponse(resp)
		//						}
		//					}
		//				} else {
		//					msg.Text = "Неверно введены лабы, вы ввели те лабы, которые не собирались сдавать"
		//				}
		//			}
		//
		//		}
		//	case "fail":
		//		var resp int
		//		labs := list[userID].labsToPass
		//		list, resp = deleteElem(list, keys, userID)
		//		msg.Text = processResponse(resp)
		//		if resp == Delete {
		//			list, resp = insert(list, keys, userID, labs, deadlineLab)
		//			writeToSheet(srv, SpreadsheetID, listToInrerface(list, keys), QueueRows)
		//			msg.Text = processResponse(resp)
		//		}
		//	case "deadlinelab":
		//		msg.Text = fmt.Sprint(deadlineLab)
		//	case "forceinsert":
		//		if isAdmin(chatID) {
		//			args := update.Message.CommandArguments()
		//			splittedArgs := strings.Split(args, "-")
		//
		//			studentUserID, _ := strconv.ParseInt(splittedArgs[0], 10, 64)
		//			_, exists := list[studentUserID]
		//			if !exists {
		//				msg := tgbotapi.NewMessage(chatID, "Такого студента нет в списке группы")
		//				bot.Send(msg)
		//				continue
		//			}
		//
		//			labs, _ := getLabsFromString(splittedArgs[1])
		//			var resp int
		//
		//			if len(splittedArgs) == 3 {
		//				pos, _ := strconv.Atoi(splittedArgs[2])
		//				list, resp = forceInsert(list, keys, studentUserID, labs, pos)
		//			} else {
		//				list, resp = insert(list, keys, studentUserID, labs, deadlineLab)
		//			}
		//
		//			writeToSheet(srv, SpreadsheetID, listToInrerface(list, keys), QueueRows)
		//
		//			if resp == Success {
		//				msg.Text = fmt.Sprintf("Студент %s успешно добавлен", list[studentUserID].name)
		//			} else {
		//				msg.Text = processResponse(resp)
		//			}
		//
		//			if chatID != CreatorID {
		//				bot.Send(tgbotapi.NewMessage(CreatorID, list[chatID].name + " executed: /" +
		//					update.Message.Command() + " " + args))
		//			}
		//		} else {
		//			msg.Text = "У меня нет такой команды"
		//		}
		//	case "forcedelete":
		//		if isAdmin(chatID) {
		//			args := update.Message.CommandArguments()
		//			studentUserID, _ := strconv.ParseInt(args, 10, 64)
		//			_, exists := list[studentUserID]
		//			if !exists {
		//				msg := tgbotapi.NewMessage(chatID, "Такого студента нет в списке группы")
		//				bot.Send(msg)
		//				continue
		//			}
		//
		//			var resp int
		//			list, resp = deleteElem(list, keys, studentUserID)
		//
		//			writeToSheet(srv, SpreadsheetID, listToInrerface(list, keys), QueueRows)
		//
		//			if resp == Delete {
		//				msg.Text = fmt.Sprintf("Студент %s успешно удалён", list[studentUserID].name)
		//			} else {
		//				msg.Text = processResponse(resp)
		//			}
		//
		//			if chatID != CreatorID {
		//				bot.Send(tgbotapi.NewMessage(CreatorID, list[chatID].name + " executed: /" +
		//					update.Message.Command() + " " + args))
		//			}
		//		} else {
		//			msg.Text = "У меня нет такой команды"
		//		}
		//	case "changedl":
		//		if isAdmin(chatID) {
		//			args := update.Message.CommandArguments()
		//			prevDeadlineLab := deadlineLab
		//			deadlineLab, _ = strconv.Atoi(args)
		//
		//			var rec [][]interface{}
		//			rec = append(rec, []interface{}{deadlineLab})
		//
		//			writeToSheet(srv, SpreadsheetID, rec, DeadlineLabRow)
		//
		//			list = updateList(list, keys, prevDeadlineLab, deadlineLab)
		//			writeToSheet(srv, SpreadsheetID, listToInrerface(list, keys), QueueRows)
		//
		//			msg.Text = "Минимальная лаба с дедлайном изменена"
		//
		//			if chatID != CreatorID {
		//				bot.Send(tgbotapi.NewMessage(CreatorID, list[chatID].name + " executed: /" +
		//					update.Message.Command() + " " + args))
		//			}
		//		} else {
		//			msg.Text = "У меня нет такой команды"
		//		}
		//	case "updatelist":
		//		if isAdmin(chatID) {
		//			list, keys = getCurrentListWithKeys(getSheetRecordsByRange(srv, SpreadsheetID, ListRows))
		//			msg.Text = "Cписок успешно обновлён"
		//
		//			if chatID != CreatorID {
		//				bot.Send(tgbotapi.NewMessage(CreatorID, list[chatID].name + " executed: /" + update.Message.Command()))
		//			}
		//		} else {
		//			msg.Text = "У меня нет такой команды"
		//		}
		//	case "wholelist":
		//		if isAdmin(chatID) {
		//			msg.Text = fullListToString(list, keys)
		//		} else {
		//			msg.Text = "У меня нет такой команды"
		//		}
		//	default:
		//		msg.Text = "У меня нет такой команды"
		//	}
		//	bot.Send(msg)
		//} else {
		//	bot.Send(tgbotapi.NewMessage(chatID, "Я не знаю, что вам на это ответить"))
		//}
	}

}