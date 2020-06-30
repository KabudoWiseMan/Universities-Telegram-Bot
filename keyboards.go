package main

import tgbotapi "github.com/Syfaro/telegram-bot-api"

//var numericKeyboard = tgbotapi.NewReplyKeyboard(
//	tgbotapi.NewKeyboardButtonRow(
//		tgbotapi.NewKeyboardButton("1"),
//		tgbotapi.NewKeyboardButton("2"),
//		tgbotapi.NewKeyboardButton("3"),
//	),
//	tgbotapi.NewKeyboardButtonRow(
//		tgbotapi.NewKeyboardButton("4"),
//		tgbotapi.NewKeyboardButton("5"),
//		tgbotapi.NewKeyboardButton("6"),
//	),
//)

//var numericInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
//	tgbotapi.NewInlineKeyboardRow(
//		tgbotapi.NewInlineKeyboardButtonURL("1.com","http://1.com"),
//		tgbotapi.NewInlineKeyboardButtonSwitch("2sw","open 2"),
//		tgbotapi.NewInlineKeyboardButtonData("3","3"),
//	),
//	tgbotapi.NewInlineKeyboardRow(
//		tgbotapi.NewInlineKeyboardButtonData("4","4"),
//		tgbotapi.NewInlineKeyboardButtonData("5","5"),
//		tgbotapi.NewInlineKeyboardButtonData("6","6"),
//	),
//)

var blank = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Ещё не готово", "nil")),
)

var mainMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Подбор вуза","uni")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Найти вуз","fUni")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сравнить вузы","comp")),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Рейтинг вузов","rate")),
)
