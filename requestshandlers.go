package main

import (
	"database/sql"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"math"
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

func handleBackRequest(db *sql.DB, data string, user *UserInfo) (string, tgbotapi.InlineKeyboardMarkup) {
	if user.State == RatingQSState {
		text, rateQSMenu := handleRatingQSRequest(db, data)
		return text, rateQSMenu
	} else if user.State == FindUniState {
		text, findUniMenu := handleFindUniRequest(db, user.Query + "#" + data)
		return text, findUniMenu
	} else if user.State == UniState {
		text, searchUniMenu := handleSearchUniRequest(db, data, user)
		return text, searchUniMenu
	}

	text := "Добро пожаловать в бота для подбора университета!\n\n" +
		"Здесь вы можете узнать, какие университеты подходят вам, исходя из ваших баллов ЕГЭ и других запросов."

	return text, mainMenu
}

func handleUniRequest(db *sql.DB, data string) (string, tgbotapi.InlineKeyboardMarkup) {
	page := takePage(data)
	uniId := takeId(data)
	uni := getUniFromDb(db, uniId)

	ratingQS := getUniQSRateFromDb(db, uniId)
	text := makeTextUni(uni, ratingQS)
	uniMenu := makeUniMenu(uni, page)
	return text, uniMenu
}

func handleRatingQSRequest(db *sql.DB, data string) (string, tgbotapi.InlineKeyboardMarkup) {
	page := takePage(data)

	text := "*Международный рейтинг вузов QS*\n\n" +
		"Для более подробной информации посетите сайт QS, нажав на кнопку *Перейти на сайт QS*\n\n"

	unisQS := getUnisQSPageFromDb(db, makeOffset(page))
	text += makeTextUnisQS(unisQS)

	unisQSNum := getUnisQSNumFromDb(db)
	rateQSMenu := makeRatingQsMenu(unisQSNum, unisQS, page)

	return text, rateQSMenu
}

func handleFacsRequest(db *sql.DB, data string) (string, tgbotapi.InlineKeyboardMarkup) {
	pages := takePages(data)
	uniId := takeId(data)

	uni := getUniFromDb(db, uniId)
	text := "*" + uni.Name + "*\n\n" +
		"Факультеты:\n\n"

	facs := getFacsPageFromDb(db, uniId, makeOffset(pages[1]))
	text += makeTextFacs(facs)

	facsNum := getFacsNumFromDb(db, uniId)
	facsMenu := makeFacsMenu(facsNum, facs, pages)

	return text, facsMenu
}

func handleFacRequest(db *sql.DB, data string) (string, tgbotapi.InlineKeyboardMarkup) {
	pages := takePages(data)
	facId := takeId(data)
	fac := getFacFromDb(db, facId)

	text := makeTextFac(fac)
	facMenu := makeFacMenu(fac, pages)
	return text, facMenu
}

func handleFindUniRequest(db *sql.DB, data string) (string, tgbotapi.InlineKeyboardMarkup) {
	query := takeUniQuery(data)
	page := takePage(data)

	unisNum := getFindUnisNumFromDb(db, query)
	if unisNum == 0 {
		text := "По запросу *\"" + query + "\"* ничего не найдено " + makeEmoji(CryingEmoji) + "\n\n" +
			"Возможно нужно ввести полное название университета " + makeEmoji(WinkEmoji)
		return text, tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(mainButton))
	}

	text := "Результаты поиска по запросу *\"" + query + "\"*:\n\n"

	unis := findUnisInDb(db, query, makeOffset(page))
	text += makeTextUnis(unis)

	unisMenu := makeUnisMenu(unisNum, unis, "findUniPage", "", page)

	return text, unisMenu
}

func handleProfsRequest(db *sql.DB, data string) (string, tgbotapi.InlineKeyboardMarkup) {
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

		fac := getFacFromDb(db, uniOrFacId)
		text = "*" + fac.Name + "*\n\n" +
			"Профили:\n\n"

		profs = getFacProfsPageFromDb(db, uniOrFacId, makeOffset(curPage))

		profsNum = getFacProfsNumFromDb(db, uniOrFacId)

		pagesPattern += "#" + facsPage
		backPattern = "getFac&" + backPattern + "#" + unisPage + "#" + facsPage
	} else {
		curPage = pages[1]
		uni := getUniFromDb(db, uniOrFacId)
		text = "*" + uni.Name + "*\n\n" +
			"Профили:\n\n"

		profs = getUniProfsPageFromDb(db, uniOrFacId, makeOffset(curPage))

		profsNum = getUniProfsNumFromDb(db, uniOrFacId)
		backPattern = "getUni&" + backPattern + "#" + unisPage
	}

	text += makeTextProfs(profs)
	profsMenu := makeProfsMenu(profsNum, profs, pagesPattern, backPattern, curPage)

	return text, profsMenu
}

func handleSpecsRequest(db *sql.DB, data string) (string, tgbotapi.InlineKeyboardMarkup) {
	pages := takePages(data)
	unisPage := pages[0]
	ids := takeIds(data)
	profId := ids[0]
	uniOrFacId := ids[1]

	prof := getProfFromDb(db, profId)

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

		fac := getFacFromDb(db, uniOrFacId)
		text = "*" + fac.Name + "*\n\n"

		specs = getFacSpecsPageFromDb(db, uniOrFacId, profId, makeOffset(curPage))
		specsNum = getFacSpecsNumFromDb(db, uniOrFacId, profId)

		progsPattern += "#" + facsPage
		pagesPattern += "#" + facsPage
		backPattern += "#" + facsPage
	} else {
		profsPage = pages[1]
		curPage = pages[2]

		uni := getUniFromDb(db, uniOrFacId)
		text = "*" + uni.Name + "*\n\n"

		specs = getUniSpecsPageFromDb(db, uniOrFacId, profId, makeOffset(curPage))
		specsNum = getUniSpecsNumFromDb(db, uniOrFacId, profId)
	}

	pagesPattern += "#" + profsPage
	backPattern += "#" + profsPage
	progsPattern += "#" + profsPage

	text += "Специальности по профилю *" + makeProfOrSpecCode(prof.ProfileId) + "* " + prof.Name + ":\n\n"
	text += makeTextSpecs(specs)
	specsMenu := makeSpecsMenu(specsNum, specs, pagesPattern, backPattern, progsPattern, curPage)

	return text, specsMenu
}

func handleProgsRequest(db *sql.DB, data string) (string, tgbotapi.InlineKeyboardMarkup) {
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

			fac := getFacFromDb(db, uniOrFacId)
			text = "*" + fac.Name + "*\n\n"

			progs = getFacProgsPageFromDb(db, uniOrFacId, makeOffset(curPage))
			progsNum = getFacProgsNumFromDb(db, uniOrFacId)

			pagesPattern += "#" + facsPage
			backPattern = "getFac&" + backPattern + "#" + facsPage
		} else {
			curPage = pages[1]

			uni := getUniFromDb(db, uniOrFacId)
			text = "*" + uni.Name + "*\n\n"

			progs = getUniProgsPageFromDb(db, uniOrFacId, makeOffset(curPage))
			progsNum = getUniProgsNumFromDb(db, uniOrFacId)

			backPattern = "getUni&" + backPattern
		}

		text += "Программы обучения:\n\n"
	} else {
		uniOrFacId := ids[1]
		specId := ids[0]
		spec := getSpecFromDb(db, specId)

		var profsPage, specsPage string

		pagesPattern = "&" + strconv.Itoa(spec.SpecialityId) + "&" + uniOrFacId + "#" + unisPage
		backPattern = "specs&" + strconv.Itoa(spec.ProfileId) + "&" + uniOrFacId + "#" + unisPage

		if len(pages) == 5 {
			facsPage := pages[1]
			profsPage = pages[2]
			specsPage = pages[3]
			curPage = pages[4]

			fac := getFacFromDb(db, uniOrFacId)
			text = "*" + fac.Name + "*\n\n"

			progs = getFacSpecProgsPageFromDb(db, uniOrFacId, specId, makeOffset(curPage))
			progsNum = getFacSpecProgsNumFromDb(db, uniOrFacId, specId)

			pagesPattern += "#" + facsPage
			backPattern += "#" + facsPage
		} else {
			profsPage = pages[1]
			specsPage = pages[2]
			curPage = pages[3]

			uni := getUniFromDb(db, uniOrFacId)
			text = "*" + uni.Name + "*\n\n"

			progs = getUniSpecProgsPageFromDb(db, uniOrFacId, specId, makeOffset(curPage))
			progsNum = getUniSpecProgsNumFromDb(db, uniOrFacId, specId)
		}

		pagesPattern += "#" + profsPage + "#" + specsPage
		backPattern += "#" + profsPage + "#" + specsPage

		text += "Программы обучения по специальности *" + makeProfOrSpecCode(spec.SpecialityId) + "* " + spec.Name + ":\n\n"
	}

	text += makeTextProgs(progs)
	progsMenu := makeProgsMenu(progsNum, progs, pagesPattern, backPattern, "#" + strings.Join(pages, "#"), curPage)

	return text, progsMenu
}

func handleProgRequest(db *sql.DB, data string) (string, tgbotapi.InlineKeyboardMarkup) {
	pages := takePages(data)
	unisPage := pages[0]
	progId := takeId(data)
	prog := getProgInfoFromDb(db, progId)

	facIdStr := strconv.Itoa(prog.FacultyId)

	var text, progsPage string
	backPattern := "#" + unisPage

	if len(pages) % 2 == 0 {
		uni := getUniOfFacFromDb(db, facIdStr)
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
		fac := getFacFromDb(db, facIdStr)
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

func handleUnisCompRequest(db *sql.DB, user *UserInfo) (string, tgbotapi.InlineKeyboardMarkup) {
	text := "Введите один или несколько критериев для получения подборки университетов\n"

	if len(user.Eges) != 0 {
		subjs := getSubjsMapFromDb(db)
		text += "\n*ЕГЭ:*\n" + makeTextEges(user.Eges, subjs, "    ")
	}

	if user.EntryTest {
		text += "\n*Готов ко вступительным* " + makeEmoji(CheckEmoji)
	}

	if user.City != 0 {
		city := getCityNameFromDb(db, user.City)
		text += "\n*Город:* " + city
	}

	if user.ProfileId != 0 {
		if user.SpecialityId != 0 {
			spec := getSpecFromDb(db, strconv.Itoa(user.SpecialityId))
			text += "\n*Специальность:* " + spec.Name + " (" + makeProfOrSpecCode(user.SpecialityId) + ")"
		} else {
			prof := getProfFromDb(db, strconv.Itoa(user.ProfileId))
			text += "\n*Профиль:* " + prof.Name + " (" + makeProfOrSpecCode(user.ProfileId) + ")"
		}
	}

	if user.Dormatary {
		text += "\n*Общежитие* " + makeEmoji(CheckEmoji)
	}

	if user.MilitaryDep {
		text += "\n*Военная кафедра* " + makeEmoji(CheckEmoji)
	}

	if user.Fee != math.MaxUint64 {
		text += "\n*Максимальная цена обучения:* " + strconv.Itoa(int(user.Fee))
	}

	unisCompilationMenu := makeUnisCompilationMenu(user)

	return text, unisCompilationMenu
}

func handleChangeOrClearRequest(db *sql.DB, data string, user *UserInfo) (string, tgbotapi.InlineKeyboardMarkup) {
	ids := takeIds(data)
	state, _ := strconv.Atoi(ids[0])
	page := takePage(data)
	text := "Измените или сбросьте"

	var subjs map[int]string
	switch state {
	case FeeState:
		text += " цену"
	case CityState:
		text += " город"
	case ProfileState:
		if user.SpecialityId != 0 {
			text += " профиль/специальность"
		} else {
			text += " профиль"
		}
	case EgeState:
		subjs := getSubjsMapFromDb(db)
		text += " ваши ЕГЭ\n\nВыбрано:\n" + makeTextEges(user.Eges, subjs,  "")
	case SubjState:
		subjs = getSubjsMapFromDb(db)
		subjId, _ := strconv.Atoi(ids[1])
		user.LastSubj = subjId
		text += " ЕГЭ по предмету *" + getSubjNameFromDb(db, subjId) + "*"
	}
	
	changeMenu := makeChangeOrClearMenu(state, user, subjs, page)

	return text, changeMenu
}

func handleCitiesRequest(db *sql.DB, data string, user *UserInfo) (string, tgbotapi.InlineKeyboardMarkup) {
	text := "Выберите город обучения"

	var backPattern string
	page := takePage(data)
	if user.City == 0 {
		backPattern = "uni"
	} else {
		backPattern = "chOrCl&" + strconv.Itoa(CityState)
	}
	cities := getCitiesFromDb(db, makeOffset(page))
	citiesNum := getCitiesNumFromDb(db)
	citiesMenu := makeCitiesMenu(citiesNum, cities, backPattern, page)

	return text, citiesMenu
}

func handleProfilesRequest(db *sql.DB, data string, user *UserInfo) (string, tgbotapi.InlineKeyboardMarkup) {
	text := "*Выберите профиль обучения*\n\n"

	var backPattern string
	page := takePage(data)
	if user.ProfileId == 0 {
		backPattern = "uni"
	} else {
		backPattern = "chOrCl&" + strconv.Itoa(ProfileState)
	}
	profs := getProfsPageFromDb(db, makeOffset(page))
	text += makeTextProfs(profs)
	profsNum := getProfsNumFromDb(db)
	profsMenu := makeProfsPageMenu(profsNum, profs, backPattern, page)

	return text, profsMenu
}

func handleSpecOrNotRequest(data string) (string, tgbotapi.InlineKeyboardMarkup) {
	text := "Выберите специальность или оставьте поиск только по профилю"
	page := takePage(data)
	profId := takeId(data)
	specOrNotMenu := makeSpecOrNotMenu(page, profId)

	return text, specOrNotMenu
}

func handleSpecialitiesRequest(db *sql.DB, data string) (string, tgbotapi.InlineKeyboardMarkup) {
	text := "*Выберите специальность обучения*\n\n"

	profId := takeId(data)
	pages := takePages(data)
	var pagesPattern, backPattern string
	var curPage string
	if len(pages) == 2 {
		profsPage := pages[0]
		pagesPattern = "#" + profsPage
		backPattern = "proOrSpe&" + profId + "#" + profsPage
		curPage = pages[1]
	} else {
		curPage = pages[0]
		backPattern = "chOrCl&" + strconv.Itoa(ProfileState)
	}

	specs := getSpecsPageFromDb(db, makeOffset(curPage), profId)
	text += makeTextSpecs(specs)
	specsNum := getSpecsNumFromDb(db, profId)
	specsMenu := makeSpecsPageMenu(specsNum, specs, profId, pagesPattern, backPattern, curPage)

	return text, specsMenu
}

func handleEgesRequest(db *sql.DB, data string, user *UserInfo) (string, tgbotapi.InlineKeyboardMarkup) {
	page := takePage(data)

	var text string
	isEges := false
	if len(user.Eges) == 0 {
		text = "*Выберите ваши ЕГЭ*\n\n"
	} else {
		isEges = true
		subjs := getSubjsMapFromDb(db)
		text = "*Выберите или измените ваши ЕГЭ*\n\nУже выбрано:\n"
		text += makeTextEges(user.Eges, subjs,  "")
	}

	subjs := getSubjsFromDb(db, makeOffset(page), user)
	subjsNum := getSubjsNumFromDb(db, user)
	egesMenu := makeEgesMenu(subjsNum, subjs, isEges, page)

	return text, egesMenu
}

func handleSubjRequest(data string, user *UserInfo) (string, tgbotapi.InlineKeyboardMarkup) {
	text := "Введите баллы или оставьте поиск только предмету"
	page := takePage(data)
	subjId := takeId(data)
	user.LastSubj, _ = strconv.Atoi(subjId)
	pointsOrNotMenu := makePointsOrNotMenu(page, subjId)

	return text, pointsOrNotMenu
}

func handleSearchUniRequest(db *sql.DB, data string, user *UserInfo) (string, tgbotapi.InlineKeyboardMarkup) {
	page := takePage(data)

	innerQuery := makeSearchInnerQueryForDb(db, user)

	unisNum := getSearchUnisNumFromDb(db, innerQuery)
	log.Println("Num:", unisNum)
	if unisNum == 0 {
		text := "Не удалось подобрать вузы по выбранным критериям " + makeEmoji(CryingEmoji)
		return text, makeMainBackMenu("uni")
	}

	text := "Результаты подбора вузов:\n\n"

	unis := searchUnisInDb(db, innerQuery, makeOffset(page))
	text += makeTextUnis(unis)

	unisMenu := makeUnisMenu(unisNum, unis, "search", "uni", page)

	return text, unisMenu
}