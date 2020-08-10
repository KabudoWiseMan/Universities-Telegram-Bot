package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"strconv"
	"strings"
)

func takeId(data string) int {
	splitted := strings.Split(data, "#")
	splitted2 := strings.Split(splitted[0], "&")
	id, _ := strconv.Atoi(splitted2[len(splitted2) - 1])
	return id
}

func takeIds(data string) []int {
	var ids []int
	splitted := strings.Split(data, "#")
	splitted2 := strings.Split(splitted[0], "&")
	for _, s := range splitted2[1:] {
		id, _ := strconv.Atoi(s)
		ids = append(ids, id)
	}

	return ids
}

func takePage(data string) int {
	splitted := strings.Split(data, "#")
	page, _ := strconv.Atoi(splitted[len(splitted) - 1])
	return page
}

func takePages(data string) []int {
	var pages []int
	splitted := strings.Split(data, "#")
	for _, s := range splitted[1:] {
		page, _ := strconv.Atoi(s)
		pages = append(pages, page)
	}

	return pages
}

func takeUniQuery(data string) string {
	splitted := strings.Split(data, "#")
	return splitted[0]
}

func handleBackRequest(data string, user *UserInfo) (string, tgbotapi.InlineKeyboardMarkup) {
	if strings.Contains(data, "Uni") {
		page := takePage(data)
		uniId := takeId(data)
		uni := getUniFromDb(uniId)

		text := makeTextUni(uni)
		uniMenu := makeUniMenu(uni, page)
		return text, uniMenu
	} else if strings.Contains(data, "Facs") {
		text, facsMenu := handleFacsRequest(data)
		return text, facsMenu
	} else if strings.Contains(data, "Fac") {
		text, facMenu := handleFacRequest(data)
		return text, facMenu
	} else if strings.Contains(data, "Profs") {
		text, facMenu := handleProfsRequest(data)
		return text, facMenu
	} else {
		if user.State == RatingQSState {
			text, rateQSMenu := handleRatingQSRequest(data)
			return text, rateQSMenu
		} else if user.State == FindUniState {
			text, findUniMenu := handleFindUniRequest(user.Query + "#" + data)
			return text, findUniMenu
		}
	}

	text := "Добро пожаловать в бота для подбора университета!\n\n" +
		"Здесь вы можете узнать, какие университеты подходят вам, исходя из ваших баллов ЕГЭ и других запросов."

	return text, mainMenu
}

func handleUniRequest(data string) (string, tgbotapi.InlineKeyboardMarkup) {
	page := takePage(data)
	uniId := takeId(data)
	uni := getUniFromDb(uniId)

	text := makeTextUni(uni)
	uniMenu := makeUniMenu(uni, page)
	return text, uniMenu
}

func handleRatingQSRequest(data string) (string, tgbotapi.InlineKeyboardMarkup) {
	page := takePage(data)

	text := "*Международный рейтинг вузов QS*\n\n" +
		"Для более подробной информации посетите сайт QS, нажав на кнопку *Перейти на сайт QS*\n\n"

	unisQS := getUnisQSPageFromDb((page - 1) * 5)
	text += makeTextUnisQS(unisQS)

	unisQSNum := getUnisQSNumFromDb()
	rateQSMenu := makeRatingQsMenu(unisQSNum, unisQS, page)

	return text, rateQSMenu
}

func handleFacsRequest(data string) (string, tgbotapi.InlineKeyboardMarkup) {
	pages := takePages(data)
	uniId := takeId(data)

	uni := getUniFromDb(uniId)
	text := "*" + uni.Name + "*\n\n" +
		"Факультеты:\n\n"

	facs := getFacsPageFromDb(uniId, (pages[1] - 1) * 5)
	text += makeTextFacs(facs)

	facsNum := getFacsNumFromDb(uniId)
	facsMenu := makeFacsMenu(facsNum, facs, pages)

	return text, facsMenu
}

func handleFacRequest(data string) (string, tgbotapi.InlineKeyboardMarkup) {
	pages := takePages(data)
	facId := takeId(data)
	fac := getFacFromDb(facId)

	text := makeTextFac(fac)
	facMenu := makeFacMenu(fac, pages)
	return text, facMenu
}

func handleFindUniRequest(data string) (string, tgbotapi.InlineKeyboardMarkup) {
	query := takeUniQuery(data)
	page := takePage(data)

	unisNum := getFindUnisNumFromDb(query)
	if unisNum == 0 {
		text := "По запросу *\"" + query + "\"* ничего не найдено " + makeEmoji(CryingEmoji) + "\n\n" +
			"Возможно нужно ввести полное название университета " + makeEmoji(WinkEmoji)
		return text, tgbotapi.NewInlineKeyboardMarkup()
	}

	text := "Результаты поиска по запросу *\"" + query + "\"*:\n\n"

	unis := findUnisInDb(query, (page - 1) * 5)
	text += makeTextUnis(unis)

	rateQSMenu := makeUnisMenu(unisNum, unis, page)

	return text, rateQSMenu
}

func handleProfsRequest(data string) (string, tgbotapi.InlineKeyboardMarkup) {
	pages := takePages(data)
	uniPage := pages[0]
	uniOrFacId := takeId(data)

	var text string
	var profs []*Profile
	var profsNum, curPage int
	pagesPattern := "&" + strconv.Itoa(uniOrFacId) + "#" + strconv.Itoa(uniPage)
	backPattern := strconv.Itoa(uniOrFacId)

	if len(pages) == 3 {
		facsPage := pages[1]
		curPage = pages[2]

		fac := getFacFromDb(uniOrFacId)
		text = "*" + fac.Name + "*\n\n" +
			"Профили:\n\n"

		profs = getFacProfsPageFromDb(uniOrFacId, (curPage - 1) * 5)

		profsNum = getFacProfsNumFromDb(uniOrFacId)

		pagesPattern += "#" + strconv.Itoa(facsPage)
		backPattern = "backFac&" + backPattern + "#" + strconv.Itoa(uniPage) + "#" + strconv.Itoa(facsPage)
	} else {
		curPage = pages[1]
		uni := getUniFromDb(uniOrFacId)
		text = "*" + uni.Name + "*\n\n" +
			"Профили:\n\n"

		profs = getUniProfsPageFromDb(uniOrFacId, (curPage - 1) * 5)

		profsNum = getUniProfsNumFromDb(uniOrFacId)
		backPattern = "backUni&" + backPattern + "#" + strconv.Itoa(uniPage)
	}

	text += makeTextProfs(profs)
	profsMenu := makeProfsMenu(profsNum, profs, pagesPattern, backPattern, curPage)

	return text, profsMenu
}

func handleSpecsRequest(data string) (string, tgbotapi.InlineKeyboardMarkup) {
	pages := takePages(data)
	uniPage := pages[0]
	ids := takeIds(data)
	profId := ids[0]
	uniOrFacId := ids[1]

	prof := getProfFromDb(profId)

	var text string
	var specs []*Speciality
	var specsNum, profsPage, curPage int
	progsPattern := "&" + strconv.Itoa(uniOrFacId) + "#" + strconv.Itoa(uniPage)
	pagesPattern := "&" + strconv.Itoa(profId) + "&" + strconv.Itoa(uniOrFacId) + "#" + strconv.Itoa(uniPage)
	backPattern := "backProfs&" + strconv.Itoa(uniOrFacId) + "#" + strconv.Itoa(uniPage)

	if len(pages) == 4 {
		facsPage := pages[1]
		profsPage = pages[2]
		curPage = pages[3]

		fac := getFacFromDb(uniOrFacId)
		text = "*" + fac.Name + "*\n\n"

		specs = getFacSpecsPageFromDb(uniOrFacId, profId, (curPage - 1) * 5)
		specsNum = getFacSpecsNumFromDb(uniOrFacId, profId)

		progsPattern += "#" + strconv.Itoa(facsPage)
		pagesPattern += "#" + strconv.Itoa(facsPage)
		backPattern += "#" + strconv.Itoa(facsPage)
	} else {
		profsPage = pages[1]
		curPage = pages[2]

		uni := getUniFromDb(uniOrFacId)
		text = "*" + uni.Name + "*\n\n"

		specs = getUniSpecsPageFromDb(uniOrFacId, profId, (curPage - 1) * 5)
		specsNum = getUniSpecsNumFromDb(uniOrFacId, profId)
	}

	pagesPattern += "#" + strconv.Itoa(profsPage)
	backPattern += "#" + strconv.Itoa(profsPage)
	progsPattern += "#" + strconv.Itoa(profsPage)

	text += "Специальности по профилю *" + makeProfOrSpecCode(prof.ProfileId) + "* " + prof.Name + ":\n\n"
	text += makeTextSpecs(specs)
	specsMenu := makeSpecsMenu(specsNum, specs, pagesPattern, backPattern, progsPattern, curPage)

	return text, specsMenu
}

//func handleProgsRequest(data string) (string, tgbotapi.InlineKeyboardMarkup) {
//	pages := takePages(data)
//	uniPage := pages[0]
//	ids := takeIds(data)
//	profId := ids[0]
//	uniOrFacId := ids[1]
//
//	prof := getProfFromDb(profId)
//
//	var text string
//	var specs []*Speciality
//	var specsNum, profsPage, curPage int
//	pagesPattern := "&" + strconv.Itoa(profId) + "&" + strconv.Itoa(uniOrFacId) + "#" + strconv.Itoa(uniPage)
//	backPattern := "backProfs&" + strconv.Itoa(uniOrFacId) + "#" + strconv.Itoa(uniPage)
//
//	if len(pages) == 4 {
//		facsPage := pages[1]
//		profsPage = pages[2]
//		curPage = pages[3]
//
//		fac := getFacFromDb(uniOrFacId)
//		text = "*" + fac.Name + "*\n\n"
//
//		specs = getFacSpecsPageFromDb(uniOrFacId, profId, (curPage - 1) * 5)
//		specsNum = getFacSpecsNumFromDb(uniOrFacId, profId)
//
//		pagesPattern += "#" + strconv.Itoa(facsPage)
//		backPattern += "#" + strconv.Itoa(facsPage)
//	} else {
//		profsPage = pages[1]
//		curPage = pages[2]
//
//		uni := getUniFromDb(uniOrFacId)
//		text = "*" + uni.Name + "*\n\n"
//
//		specs = getUniSpecsPageFromDb(uniOrFacId, profId, (curPage - 1) * 5)
//		specsNum = getUniSpecsNumFromDb(uniOrFacId, profId)
//	}
//
//	pagesPattern += "#" + strconv.Itoa(profsPage)
//	backPattern += "#" + strconv.Itoa(profsPage)
//
//	text += "Специальности по профилю *" + makeProfOrSpecCode(prof.ProfileId) + "* " + prof.Name + ":\n\n"
//	text += makeTextSpecs(specs)
//	specsMenu := makeSpecsMenu(specsNum, specs, pagesPattern, backPattern, curPage)
//
//	return text, specsMenu
//}
