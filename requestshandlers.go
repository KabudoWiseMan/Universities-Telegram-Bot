package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"strconv"
	"strings"
)

func makeOffset(page string) string {
	pageNum, _ := strconv.Atoi(page)
	offset := (pageNum - 1) * 5
	return strconv.Itoa(offset)
}

func takeId(data string) string {
	splitted := strings.Split(data, "#")
	splitted2 := strings.Split(splitted[0], "&")
	return splitted2[len(splitted2) - 1]
}

func takeIds(data string) []string {
	var ids []string
	splitted := strings.Split(data, "#")
	splitted2 := strings.Split(splitted[0], "&")
	for _, s := range splitted2[1:] {
		ids = append(ids, s)
	}

	return ids
}

func takePage(data string) string {
	splitted := strings.Split(data, "#")
	return splitted[len(splitted) - 1]
}

func takePages(data string) []string {
	var pages []string
	splitted := strings.Split(data, "#")
	for _, s := range splitted[1:] {
		pages = append(pages, s)
	}

	return pages
}

func takeUniQuery(data string) string {
	splitted := strings.Split(data, "#")
	return splitted[0]
}

func handleBackRequest(data string, user *UserInfo) (string, tgbotapi.InlineKeyboardMarkup) {
	if user.State == RatingQSState {
		text, rateQSMenu := handleRatingQSRequest(data)
		return text, rateQSMenu
	} else if user.State == FindUniState {
		text, findUniMenu := handleFindUniRequest(user.Query + "#" + data)
		return text, findUniMenu
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

	unisQS := getUnisQSPageFromDb(makeOffset(page))
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

	facs := getFacsPageFromDb(uniId, makeOffset(pages[1]))
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

	unis := findUnisInDb(query, makeOffset(page))
	text += makeTextUnis(unis)

	rateQSMenu := makeUnisMenu(unisNum, unis, page)

	return text, rateQSMenu
}

func handleProfsRequest(data string) (string, tgbotapi.InlineKeyboardMarkup) {
	pages := takePages(data)
	unisPage := pages[0]
	uniOrFacId := takeId(data)

	var text, curPage string
	var profs []*Profile
	var profsNum int
	pagesPattern := "&" + uniOrFacId + "#" + unisPage
	backPattern := uniOrFacId

	if len(pages) == 3 {
		facsPage := pages[1]
		curPage = pages[2]

		fac := getFacFromDb(uniOrFacId)
		text = "*" + fac.Name + "*\n\n" +
			"Профили:\n\n"

		profs = getFacProfsPageFromDb(uniOrFacId, makeOffset(curPage))

		profsNum = getFacProfsNumFromDb(uniOrFacId)

		pagesPattern += "#" + facsPage
		backPattern = "getFac&" + backPattern + "#" + unisPage + "#" + facsPage
	} else {
		curPage = pages[1]
		uni := getUniFromDb(uniOrFacId)
		text = "*" + uni.Name + "*\n\n" +
			"Профили:\n\n"

		profs = getUniProfsPageFromDb(uniOrFacId, makeOffset(curPage))

		profsNum = getUniProfsNumFromDb(uniOrFacId)
		backPattern = "getUni&" + backPattern + "#" + unisPage
	}

	text += makeTextProfs(profs)
	profsMenu := makeProfsMenu(profsNum, profs, pagesPattern, backPattern, curPage)

	return text, profsMenu
}

func handleSpecsRequest(data string) (string, tgbotapi.InlineKeyboardMarkup) {
	pages := takePages(data)
	unisPage := pages[0]
	ids := takeIds(data)
	profId := ids[0]
	uniOrFacId := ids[1]

	prof := getProfFromDb(profId)

	var text, profsPage, curPage string
	var specs []*Speciality
	var specsNum int
	progsPattern := "&" + uniOrFacId + "#" + unisPage
	pagesPattern := "&" + profId + "&" + uniOrFacId + "#" + unisPage
	backPattern := "profs&" + uniOrFacId + "#" + unisPage

	if len(pages) == 4 {
		facsPage := pages[1]
		profsPage = pages[2]
		curPage = pages[3]

		fac := getFacFromDb(uniOrFacId)
		text = "*" + fac.Name + "*\n\n"

		specs = getFacSpecsPageFromDb(uniOrFacId, profId, makeOffset(curPage))
		specsNum = getFacSpecsNumFromDb(uniOrFacId, profId)

		progsPattern += "#" + facsPage
		pagesPattern += "#" + facsPage
		backPattern += "#" + facsPage
	} else {
		profsPage = pages[1]
		curPage = pages[2]

		uni := getUniFromDb(uniOrFacId)
		text = "*" + uni.Name + "*\n\n"

		specs = getUniSpecsPageFromDb(uniOrFacId, profId, makeOffset(curPage))
		specsNum = getUniSpecsNumFromDb(uniOrFacId, profId)
	}

	pagesPattern += "#" + profsPage
	backPattern += "#" + profsPage
	progsPattern += "#" + profsPage

	text += "Специальности по профилю *" + makeProfOrSpecCode(prof.ProfileId) + "* " + prof.Name + ":\n\n"
	text += makeTextSpecs(specs)
	specsMenu := makeSpecsMenu(specsNum, specs, pagesPattern, backPattern, progsPattern, curPage)

	return text, specsMenu
}

func handleProgsRequest(data string) (string, tgbotapi.InlineKeyboardMarkup) {
	pages := takePages(data)
	unisPage := pages[0]
	ids := takeIds(data)

	var text, curPage string
	var progs []*Program
	var progsNum int
	var pagesPattern, backPattern string

	if len(pages) == 2 || len(pages) == 3 {
		uniOrFacId := ids[0]
		pagesPattern = "&" + uniOrFacId + "#" + unisPage
		backPattern = uniOrFacId + "#" + unisPage

		if len(pages) == 3 {
			facsPage := pages[1]
			curPage = pages[2]

			fac := getFacFromDb(uniOrFacId)
			text = "*" + fac.Name + "*\n\n"

			progs = getFacProgsPageFromDb(uniOrFacId, makeOffset(curPage))
			progsNum = getFacProgsNumFromDb(uniOrFacId)

			pagesPattern += "#" + facsPage
			backPattern = "getFac&" + backPattern + "#" + facsPage
		} else {
			curPage = pages[1]

			uni := getUniFromDb(uniOrFacId)
			text = "*" + uni.Name + "*\n\n"

			progs = getUniProgsPageFromDb(uniOrFacId, makeOffset(curPage))
			progsNum = getUniProgsNumFromDb(uniOrFacId)

			backPattern = "getUni&" + backPattern
		}

		text += "Программы обучения:\n\n"
	} else {
		uniOrFacId := ids[1]
		specId := ids[0]
		spec := getSpecFromDb(specId)

		var profsPage, specsPage string

		pagesPattern = "&" + strconv.Itoa(spec.SpecialityId) + "&" + uniOrFacId + "#" + unisPage
		backPattern = "specs&" + strconv.Itoa(spec.ProfileId) + "&" + uniOrFacId + "#" + unisPage

		if len(pages) == 5 {
			facsPage := pages[1]
			profsPage = pages[2]
			specsPage = pages[3]
			curPage = pages[4]

			fac := getFacFromDb(uniOrFacId)
			text = "*" + fac.Name + "*\n\n"

			progs = getFacSpecProgsPageFromDb(uniOrFacId, specId, makeOffset(curPage))
			progsNum = getFacSpecProgsNumFromDb(uniOrFacId, specId)

			pagesPattern += "#" + facsPage
			backPattern += "#" + facsPage
		} else {
			profsPage = pages[1]
			specsPage = pages[2]
			curPage = pages[3]

			uni := getUniFromDb(uniOrFacId)
			text = "*" + uni.Name + "*\n\n"

			progs = getUniSpecProgsPageFromDb(uniOrFacId, specId, makeOffset(curPage))
			progsNum = getUniSpecProgsNumFromDb(uniOrFacId, specId)
		}

		pagesPattern += "#" + profsPage + "#" + specsPage
		backPattern += "#" + profsPage + "#" + specsPage

		text += "Программы обучения по специальности *" + makeProfOrSpecCode(spec.SpecialityId) + "* " + spec.Name + ":\n\n"
	}

	text += makeTextProgs(progs)
	progsMenu := makeProgsMenu(progsNum, progs, pagesPattern, backPattern, "#" + strings.Join(pages, "#"), curPage)

	return text, progsMenu
}

func handleProgRequest(data string) (string, tgbotapi.InlineKeyboardMarkup) {
	pages := takePages(data)
	unisPage := pages[0]
	progId := takeId(data)
	prog := getProgInfoFromDb(progId)

	facIdStr := strconv.Itoa(prog.FacultyId)

	var text, progsPage string
	backPattern := "#" + unisPage

	if len(pages) % 2 == 0 {
		uni := getUniOfFacFromDb(facIdStr)
		backPattern = "&" + strconv.Itoa(uni.UniversityId) + backPattern
		text = "*" + uni.Name + "*\n\n"

		if len(pages) == 4 {
			profsPage := pages[1]
			specsPage := pages[2]
			progsPage = pages[3]

			backPattern = "&" + strconv.Itoa(prog.SpecialityId) + backPattern + "#" + profsPage + "#" + specsPage
		} else {
			progsPage = pages[1]
		}
	} else {
		facsPage := pages[1]
		fac := getFacFromDb(facIdStr)
		backPattern = "&" + facIdStr + backPattern + "#" + facsPage
		text = "*" + fac.Name + "*\n\n"

		if len(pages) == 5 {
			profsPage := pages[2]
			specsPage := pages[3]
			progsPage = pages[4]

			backPattern = "&" + strconv.Itoa(prog.SpecialityId) + backPattern + "#" + profsPage + "#" + specsPage
		} else {
			progsPage = pages[2]
		}
	}

	backPattern += "#" + progsPage
	backPattern = "progs" + backPattern

	text += "Программа обучения: " + makeTextProg(prog)
	progMenu := makeProgMenu(backPattern)

	return text, progMenu
}
