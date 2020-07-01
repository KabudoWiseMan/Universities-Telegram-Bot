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

var (
	mainButton = tgbotapi.NewInlineKeyboardButtonData("<< Главное меню >>","main")
	qsButton = tgbotapi.NewInlineKeyboardButtonURL("Перейти на сайт QS","https://www.topuniversities.com/qs-world-university-rankings")
	ranhigsButton = tgbotapi.NewInlineKeyboardButtonURL("Перейти на сайт рейтинга РАНХиГС","https://russiaedu.ru/rating")

	blankMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Ещё не готово", "nil")),
		tgbotapi.NewInlineKeyboardRow(mainButton),
	)

	mainMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Подбор ВУЗа","uni")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Найти ВУЗ","fUni")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сравнить ВУЗы","comp")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Рейтинг ВУЗов","rate")),
		tgbotapi.NewInlineKeyboardRow(mainButton),
	)

	ratingMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Международный рейтинг QS","qs")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Российский рейтинг РАНХиГС","ranhigs")),
		tgbotapi.NewInlineKeyboardRow(mainButton),
	)

	unisCompilationMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Баллы ЕГЭ","ege")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Профиль","pro")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Специальность","spec")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Город","city")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Важно наличие военной кафедры","army")),
		tgbotapi.NewInlineKeyboardRow(mainButton),
	)

	compareUnisMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Добавить вуз или специальность","addSpecs")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Удалить вуз или специальность", "delSpecs")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сравнить","compare")),
		tgbotapi.NewInlineKeyboardRow(mainButton),
	)

	qsMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(qsButton),
		tgbotapi.NewInlineKeyboardRow(mainButton),
	)

	ranhigsMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(ranhigsButton),
		tgbotapi.NewInlineKeyboardRow(mainButton),
	)

	uniMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Специальности","specs")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Подобрать специальность","findSpecs")),
		tgbotapi.NewInlineKeyboardRow(mainButton),
	)
)
