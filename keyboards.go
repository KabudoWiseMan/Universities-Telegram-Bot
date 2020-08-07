package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"math"
	"strconv"
	"strings"
)

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
	qsButton = tgbotapi.NewInlineKeyboardButtonURL("Перейти на сайт QS", RatingQsSite)

	blankMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Ещё не готово", "nil")),
		tgbotapi.NewInlineKeyboardRow(mainButton),
	)

	mainMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Подбор ВУЗа","uni")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Найти ВУЗ","fUni")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Рейтинг ВУЗов QS","rate")),
	)

	unisCompilationMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Баллы ЕГЭ","ege")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Профиль","pro")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Специальность","spec")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Город","city")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Важно наличие военной кафедры","army")),
		tgbotapi.NewInlineKeyboardRow(mainButton),
	)
)

func makePaginator(elementsNum int, elementsOnPage int, curPage int, dataPattern string) []tgbotapi.InlineKeyboardButton {
	pagesNum := int(math.Ceil(float64(elementsNum) / float64(elementsOnPage)))

	if pagesNum == 1 {
		return []tgbotapi.InlineKeyboardButton{}
	}

	var paginatorButtons []tgbotapi.InlineKeyboardButton

	if pagesNum <= 5 {
		for i := 1; i <= pagesNum; i++ {
			text := strconv.Itoa(i)
			if i == curPage {
				text = "•" + text + "•"
			}

			paginatorButtons = append(paginatorButtons, tgbotapi.NewInlineKeyboardButtonData(text, dataPattern + "#" + strconv.Itoa(i)))
		}

		return paginatorButtons
	}

	if curPage <= 3 {
		for i := 1; i < 4; i++ {
			text := strconv.Itoa(i)
			if i == curPage {
				text = "•" + text + "•"
			}

			paginatorButtons = append(paginatorButtons, tgbotapi.NewInlineKeyboardButtonData(text, dataPattern + "#" + strconv.Itoa(i)))
		}

		paginatorButtons = append(paginatorButtons, tgbotapi.NewInlineKeyboardButtonData("4 ›", dataPattern + "#4"))
		paginatorButtons = append(paginatorButtons, tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(pagesNum) + " »", dataPattern + "#" + strconv.Itoa(pagesNum)))

		return paginatorButtons
	}

	if curPage > pagesNum - 3 {
		paginatorButtons = append(paginatorButtons, tgbotapi.NewInlineKeyboardButtonData("« 1", dataPattern + "#1"))
		paginatorButtons = append(paginatorButtons, tgbotapi.NewInlineKeyboardButtonData("‹ " + strconv.Itoa(pagesNum - 3), dataPattern + "#" + strconv.Itoa(pagesNum - 3)))

		for i := pagesNum - 2; i <= pagesNum; i++ {
			text := strconv.Itoa(i)
			if i == curPage {
				text = "•" + text + "•"
			}

			paginatorButtons = append(paginatorButtons, tgbotapi.NewInlineKeyboardButtonData(text, dataPattern + "#" + strconv.Itoa(i)))
		}

		return paginatorButtons
	}

	paginatorButtons = append(paginatorButtons, tgbotapi.NewInlineKeyboardButtonData("« 1", dataPattern + "#1"))
	paginatorButtons = append(paginatorButtons, tgbotapi.NewInlineKeyboardButtonData("‹ " + strconv.Itoa(curPage - 1), dataPattern + "#" + strconv.Itoa(curPage - 1)))
	paginatorButtons = append(paginatorButtons, tgbotapi.NewInlineKeyboardButtonData("•" + strconv.Itoa(curPage) + "•", dataPattern + "#" + strconv.Itoa(curPage)))
	paginatorButtons = append(paginatorButtons, tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(curPage + 1) + " ›", dataPattern + "#" + strconv.Itoa(curPage + 1)))
	paginatorButtons = append(paginatorButtons, tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(pagesNum) + " »", dataPattern + "#" + strconv.Itoa(pagesNum)))

	return paginatorButtons
}

func makeRatingQsMenu(unisQSNum int, unisQS []*UniversityQS, curPage int) tgbotapi.InlineKeyboardMarkup {
	var unisQSButtons [][]tgbotapi.InlineKeyboardButton
	for _, uniQS := range unisQS {
		unisQSButtons = append(unisQSButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(uniQS.Name, "getUni&" + strconv.Itoa(uniQS.UniversityId) + "#" + strconv.Itoa(curPage))))
	}

	paginator := makePaginator(unisQSNum, 5, curPage, "rateQSPage")

	var fullButtons [][]tgbotapi.InlineKeyboardButton
	fullButtons = append(fullButtons, unisQSButtons...)
	fullButtons = append(fullButtons, paginator)
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(qsButton), tgbotapi.NewInlineKeyboardRow(mainButton))
	
	ratingQSFullMenu := tgbotapi.NewInlineKeyboardMarkup(
		fullButtons...
	)

	return ratingQSFullMenu
}

func makeUniMenu(uni University, page int) tgbotapi.InlineKeyboardMarkup {
	var fullButtons [][]tgbotapi.InlineKeyboardButton
	if !strings.Contains(uni.Site, " ") {
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("Перейти на сайт ВУЗа", uni.Site)))
	}
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Факультеты", "facs&" + strconv.Itoa(uni.UniversityId) + "#" + strconv.Itoa(page) + "#1")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Профили", "profs#" + strconv.Itoa(uni.UniversityId))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Программы обучения", "progs#" + strconv.Itoa(uni.UniversityId))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Подобрать программу обучения", "findProg#" + strconv.Itoa(uni.UniversityId))),
		tgbotapi.NewInlineKeyboardRow(makeBackButton("back#" + strconv.Itoa(page))),
		tgbotapi.NewInlineKeyboardRow(mainButton),
	)

	uniFullMenu := tgbotapi.NewInlineKeyboardMarkup(
		fullButtons...
	)

	return uniFullMenu
}

func makeBackButton(data string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData("<< Назад", data)
}

func makeFacsMenu(facsNum int, uniId int, facs []*Faculty, pages []int) tgbotapi.InlineKeyboardMarkup {
	unisPage := pages[0]
	facsPage := pages[1]

	var facsButtons [][]tgbotapi.InlineKeyboardButton
	for _, fac := range facs {
		facsButtons = append(facsButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(fac.Name, "getFac&" + strconv.Itoa(fac.FacultyId) + "#" + strconv.Itoa(unisPage) + "#" + strconv.Itoa(facsPage))))
	}

	paginator := makePaginator(facsNum, 5, facsPage, "facs&" + strconv.Itoa(uniId) + "#" + strconv.Itoa(unisPage))

	var fullButtons [][]tgbotapi.InlineKeyboardButton
	fullButtons = append(fullButtons, facsButtons...)
	fullButtons = append(fullButtons, paginator)
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(makeBackButton("backUni&" + strconv.Itoa(uniId) + "#" + strconv.Itoa(unisPage))), tgbotapi.NewInlineKeyboardRow(mainButton))

	ratingQSFullMenu := tgbotapi.NewInlineKeyboardMarkup(
		fullButtons...
	)

	return ratingQSFullMenu
}
