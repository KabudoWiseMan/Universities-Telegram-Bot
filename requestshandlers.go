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

func makeDbErrorResponseData(err error, backPattern string) (string, tgbotapi.InlineKeyboardMarkup) {
	log.Println("data base error:", err)
	return "Простите, произошла ошибка" + makeEmoji(CryingEmoji), makeMainBackMenu(backPattern)
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
	uni, err := getUniFromDb(db, uniId)
	if err != nil {
		return makeDbErrorResponseData(err, "")
	}

	ratingQS, err := getUniQSRateFromDb(db, uniId)
	if err != nil {
		return makeDbErrorResponseData(err, "")
	}
	text := makeTextUni(uni, ratingQS)
	uniMenu := makeUniMenu(uni, page)
	return text, uniMenu
}

func handleRatingQSRequest(db *sql.DB, data string) (string, tgbotapi.InlineKeyboardMarkup) {
	page := takePage(data)

	text := "*Международный рейтинг вузов QS*\n\n" +
		"Для более подробной информации посетите сайт QS, нажав на кнопку *Перейти на сайт QS*\n\n"

	unisQS, err := getUnisQSPageFromDb(db, makeOffset(page))
	if err != nil || len(unisQS) == 0 {
		return makeDbErrorResponseData(err, "")
	}
	text += makeTextUnisQS(unisQS)

	unisQSNum, err := getUnisQSNumFromDb(db)
	if err != nil {
		return makeDbErrorResponseData(err, "")
	}
	rateQSMenu := makeRatingQsMenu(unisQSNum, unisQS, page)

	return text, rateQSMenu
}

func handleFacsRequest(db *sql.DB, data string) (string, tgbotapi.InlineKeyboardMarkup) {
	pages := takePages(data)
	uniId := takeId(data)
	backPattern := "getUni&" + uniId + "#" + pages[0]

	uni, err := getUniFromDb(db, uniId)
	if err != nil {
		return makeDbErrorResponseData(err, backPattern)
	}
	text := "*" + uni.Name + "*\n\n" +
		"Факультеты:\n\n"

	facs, err := getFacsPageFromDb(db, uniId, makeOffset(pages[1]))
	if err != nil {
		return makeDbErrorResponseData(err, backPattern)
	}
	text += makeTextFacs(facs)

	facsNum, err := getFacsNumFromDb(db, uniId)
	if err != nil {
		return makeDbErrorResponseData(err, backPattern)
	}
	facsMenu := makeFacsMenu(facsNum, facs, backPattern, pages)

	return text, facsMenu
}

func handleFacRequest(db *sql.DB, data string) (string, tgbotapi.InlineKeyboardMarkup) {
	pages := takePages(data)
	facId := takeId(data)
	fac, err := getFacFromDb(db, facId)
	if err != nil {
		return makeDbErrorResponseData(err, "")
	}

	text := makeTextFac(fac)
	facMenu := makeFacMenu(fac, pages)
	return text, facMenu
}

func handleFindUniRequest(db *sql.DB, data string) (string, tgbotapi.InlineKeyboardMarkup) {
	query := takeUniQuery(data)
	page := takePage(data)

	unisNum, err := getFindUnisNumFromDb(db, query)
	if err != nil {
		return makeDbErrorResponseData(err, "")
	}
	if unisNum == 0 {
		text := "По запросу *\"" + query + "\"* ничего не найдено " + makeEmoji(CryingEmoji) + "\n\n" +
			"Возможно нужно ввести полное название университета " + makeEmoji(WinkEmoji)
		return text, tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(mainButton))
	}

	text := "Результаты поиска по запросу *\"" + query + "\"*:\n\n"

	unis, err := findUnisInDb(db, query, makeOffset(page))
	if err != nil {
		return makeDbErrorResponseData(err, "")
	}
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

		fac, err := getFacFromDb(db, uniOrFacId)
		if err != nil {
			return makeDbErrorResponseData(err, "")
		}
		text = "*" + fac.Name + "*\n\n" +
			"Профили:\n\n"

		profs, err = getFacProfsPageFromDb(db, uniOrFacId, makeOffset(curPage))
		if err != nil {
			return makeDbErrorResponseData(err, "")
		}

		profsNum, err = getFacProfsNumFromDb(db, uniOrFacId)
		if err != nil {
			return makeDbErrorResponseData(err, "")
		}

		pagesPattern += "#" + facsPage
		backPattern = "getFac&" + backPattern + "#" + unisPage + "#" + facsPage
	} else {
		curPage = pages[1]
		uni, err := getUniFromDb(db, uniOrFacId)
		if err != nil {
			return makeDbErrorResponseData(err, "")
		}
		text = "*" + uni.Name + "*\n\n" +
			"Профили:\n\n"

		profs, err = getUniProfsPageFromDb(db, uniOrFacId, makeOffset(curPage))
		if err != nil {
			return makeDbErrorResponseData(err, "")
		}

		profsNum, err = getUniProfsNumFromDb(db, uniOrFacId)
		if err != nil {
			return makeDbErrorResponseData(err, "")
		}
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

	prof, err := getProfFromDb(db, profId)
	if err != nil {
		return makeDbErrorResponseData(err, "")
	}

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

		fac, err := getFacFromDb(db, uniOrFacId)
		if err != nil {
			return makeDbErrorResponseData(err, "")
		}
		text = "*" + fac.Name + "*\n\n"

		specs, err = getFacSpecsPageFromDb(db, uniOrFacId, profId, makeOffset(curPage))
		if err != nil {
			return makeDbErrorResponseData(err, "")
		}
		specsNum, err = getFacSpecsNumFromDb(db, uniOrFacId, profId)
		if err != nil {
			return makeDbErrorResponseData(err, "")
		}

		progsPattern += "#" + facsPage
		pagesPattern += "#" + facsPage
		backPattern += "#" + facsPage
	} else {
		profsPage = pages[1]
		curPage = pages[2]

		uni, err := getUniFromDb(db, uniOrFacId)
		if err != nil {
			return makeDbErrorResponseData(err, "")
		}
		text = "*" + uni.Name + "*\n\n"

		specs, err = getUniSpecsPageFromDb(db, uniOrFacId, profId, makeOffset(curPage))
		if err != nil {
			return makeDbErrorResponseData(err, "")
		}
		specsNum, err = getUniSpecsNumFromDb(db, uniOrFacId, profId)
		if err != nil {
			return makeDbErrorResponseData(err, "")
		}
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

			fac, err := getFacFromDb(db, uniOrFacId)
			if err != nil {
				return makeDbErrorResponseData(err, "")
			}
			text = "*" + fac.Name + "*\n\n"

			progs, err = getFacProgsPageFromDb(db, uniOrFacId, makeOffset(curPage))
			if err != nil {
				return makeDbErrorResponseData(err, "")
			}
			progsNum, err = getFacProgsNumFromDb(db, uniOrFacId)
			if err != nil {
				return makeDbErrorResponseData(err, "")
			}

			pagesPattern += "#" + facsPage
			backPattern = "getFac&" + backPattern + "#" + facsPage
		} else {
			curPage = pages[1]

			uni, err := getUniFromDb(db, uniOrFacId)
			if err != nil {
				return makeDbErrorResponseData(err, "")
			}
			text = "*" + uni.Name + "*\n\n"

			progs, err = getUniProgsPageFromDb(db, uniOrFacId, makeOffset(curPage))
			if err != nil {
				return makeDbErrorResponseData(err, "")
			}
			progsNum, err = getUniProgsNumFromDb(db, uniOrFacId)
			if err != nil {
				return makeDbErrorResponseData(err, "")
			}

			backPattern = "getUni&" + backPattern
		}

		text += "Программы обучения:\n\n"
	} else {
		uniOrFacId := ids[1]
		specId := ids[0]
		spec, err := getSpecFromDb(db, specId)
		if err != nil {
			return makeDbErrorResponseData(err, "")
		}

		var profsPage, specsPage string

		pagesPattern = "&" + strconv.Itoa(spec.SpecialityId) + "&" + uniOrFacId + "#" + unisPage
		backPattern = "specs&" + strconv.Itoa(spec.ProfileId) + "&" + uniOrFacId + "#" + unisPage

		if len(pages) == 5 {
			facsPage := pages[1]
			profsPage = pages[2]
			specsPage = pages[3]
			curPage = pages[4]

			fac, err := getFacFromDb(db, uniOrFacId)
			if err != nil {
				return makeDbErrorResponseData(err, "")
			}
			text = "*" + fac.Name + "*\n\n"

			progs, err = getFacSpecProgsPageFromDb(db, uniOrFacId, specId, makeOffset(curPage))
			if err != nil {
				return makeDbErrorResponseData(err, "")
			}
			progsNum, err = getFacSpecProgsNumFromDb(db, uniOrFacId, specId)
			if err != nil {
				return makeDbErrorResponseData(err, "")
			}

			pagesPattern += "#" + facsPage
			backPattern += "#" + facsPage
		} else {
			profsPage = pages[1]
			specsPage = pages[2]
			curPage = pages[3]

			uni, err := getUniFromDb(db, uniOrFacId)
			if err != nil {
				return makeDbErrorResponseData(err, "")
			}
			text = "*" + uni.Name + "*\n\n"

			progs, err = getUniSpecProgsPageFromDb(db, uniOrFacId, specId, makeOffset(curPage))
			if err != nil {
				return makeDbErrorResponseData(err, "")
			}
			progsNum, err = getUniSpecProgsNumFromDb(db, uniOrFacId, specId)
			if err != nil {
				return makeDbErrorResponseData(err, "")
			}
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
	prog, err := getProgInfoFromDb(db, progId)
	if err != nil {
		return makeDbErrorResponseData(err, "")
	}

	facIdStr := strconv.Itoa(prog.FacultyId)

	var text, progsPage string
	backPattern := "#" + unisPage

	if len(pages) % 2 == 0 {
		uni, err := getUniOfFacFromDb(db, facIdStr)
		if err != nil {
			return makeDbErrorResponseData(err, "")
		}
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
		fac, err := getFacFromDb(db, facIdStr)
		if err != nil {
			return makeDbErrorResponseData(err, "")
		}
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
		subjs, err := getSubjsMapFromDb(db)
		if err != nil {
			return makeDbErrorResponseData(err, "uni")
		}
		text += "\n*ЕГЭ:*\n" + makeTextEges(user.Eges, subjs, "    ")
	}

	if user.EntryTest {
		text += "\n*Готов ко вступительным* " + makeEmoji(CheckEmoji)
	}

	if user.City != 0 {
		city, err := getCityNameFromDb(db, user.City)
		if err != nil {
			return makeDbErrorResponseData(err, "uni")
		}
		text += "\n*Город:* " + city
	}

	if user.ProfileId != 0 {
		if user.SpecialityId != 0 {
			spec, err := getSpecFromDb(db, strconv.Itoa(user.SpecialityId))
			if err != nil {
				return makeDbErrorResponseData(err, "uni")
			}
			text += "\n*Специальность:* " + spec.Name + " (" + makeProfOrSpecCode(user.SpecialityId) + ")"
		} else {
			prof, err := getProfFromDb(db, strconv.Itoa(user.ProfileId))
			if err != nil {
				return makeDbErrorResponseData(err, "uni")
			}
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

	var backPattern string
	var subjs map[int]string
	switch state {
	case FeeState:
		backPattern = "uni"
		text += " цену"
	case CityState:
		backPattern = "uni"
		text += " город"
	case ProfileState:
		backPattern = "uni"
		if user.SpecialityId != 0 {
			text += " профиль/специальность"
		} else {
			text += " профиль"
		}
	case EgeState:
		var err error
		backPattern = "ege#" + page
		subjs, err = getSubjsMapFromDb(db)
		if err != nil {
			return makeDbErrorResponseData(err, backPattern)
		}
		text += " ваши ЕГЭ\n\nВыбрано:\n" + makeTextEges(user.Eges, subjs,  "")
	case SubjState:
		backPattern = "chOrCl&" + strconv.Itoa(EgeState) + "#" + page
		var err error
		subjs, err = getSubjsMapFromDb(db)
		if err != nil {
			return makeDbErrorResponseData(err, backPattern)
		}
		subjId, _ := strconv.Atoi(ids[1])
		user.LastSubj = subjId
		subjName, err := getSubjNameFromDb(db, subjId)
		if err != nil {
			return makeDbErrorResponseData(err, backPattern)
		}
		text += " ЕГЭ по предмету *" + subjName + "*"
	}

	changeMenu := makeChangeOrClearMenu(state, user, subjs, backPattern, page)

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
	cities, err := getCitiesFromDb(db, makeOffset(page))
	if err != nil {
		return makeDbErrorResponseData(err, backPattern)
	}
	citiesNum, err := getCitiesNumFromDb(db)
	if err != nil {
		return makeDbErrorResponseData(err, backPattern)
	}
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
	profs, err := getProfsPageFromDb(db, makeOffset(page))
	if err != nil {
		return makeDbErrorResponseData(err, backPattern)
	}
	text += makeTextProfs(profs)
	profsNum, err := getProfsNumFromDb(db)
	if err != nil {
		return makeDbErrorResponseData(err, backPattern)
	}
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

	specs, err := getSpecsPageFromDb(db, makeOffset(curPage), profId)
	if err != nil {
		return makeDbErrorResponseData(err, backPattern)
	}
	text += makeTextSpecs(specs)
	specsNum, err := getSpecsNumFromDb(db, profId)
	if err != nil {
		return makeDbErrorResponseData(err, backPattern)
	}
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
		subjs, err := getSubjsMapFromDb(db)
		if err != nil {
			return makeDbErrorResponseData(err, "uni")
		}
		text = "*Выберите или измените ваши ЕГЭ*\n\nУже выбрано:\n"
		text += makeTextEges(user.Eges, subjs,  "")
	}

	subjs, err := getSubjsFromDb(db, makeOffset(page), user)
	if err != nil {
		return makeDbErrorResponseData(err, "uni")
	}
	subjsNum, err := getSubjsNumFromDb(db, user)
	if err != nil {
		return makeDbErrorResponseData(err, "uni")
	}
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

	innerQuery := makeSearchInnerQueryForDb(user)

	unisNum, err := getSearchUnisNumFromDb(db, innerQuery)
	if err != nil {
		makeDbErrorResponseData(err, "uni")
	}
	log.Println("Num:", unisNum)
	if unisNum == 0 {
		text := "Не удалось подобрать вузы по выбранным критериям " + makeEmoji(CryingEmoji)
		return text, makeMainBackMenu("uni")
	}

	text := "Результаты подбора вузов:\n\n"

	unis, err := searchUnisInDb(db, innerQuery, makeOffset(page))
	if err != nil {
		return makeDbErrorResponseData(err, "uni")
	}
	text += makeTextUnis(unis)

	unisMenu := makeUnisMenu(unisNum, unis, "search", "uni", page)

	return text, unisMenu
}

func handleChangePointsRequest(db *sql.DB, data string, user *UserInfo) (string, tgbotapi.InlineKeyboardMarkup) {
	page := takePage(data)
	subjName, err := getSubjNameFromDb(db, user.LastSubj)
	backPatter := "chOrCl&" + strconv.Itoa(SubjState) + "&" + strconv.Itoa(user.LastSubj) + "#" + page
	if err != nil {
		return makeDbErrorResponseData(err, backPatter)
	}
	text := "Введите баллы ЕГЭ *" + subjName + "*"
	mainBackMenu := makeMainBackMenu(backPatter)
	return text, mainBackMenu
}

func handleClearMenu(db *sql.DB, data string, user *UserInfo) (string, tgbotapi.InlineKeyboardMarkup) {
	state, _ := strconv.Atoi(takeId(data))
	var text string
	var changeMenu tgbotapi.InlineKeyboardMarkup
	switch state {
	case EgeState:
		user.Eges = nil
		text, changeMenu = handleEgesRequest(db, "ege#1", user)
	case FeeState:
		user.Fee = math.MaxUint64
		text, changeMenu = handleUnisCompRequest(db, user)
	case CityState:
		user.City = 0
		text, changeMenu = handleUnisCompRequest(db, user)
	case ProfileState:
		user.ProfileId = 0
		text, changeMenu = handleUnisCompRequest(db, user)
	case SpecialityState:
		user.SpecialityId = 0
		text, changeMenu = handleUnisCompRequest(db, user)
	case SubjState:
		user.DeleteEge()
		text, changeMenu = handleChangeOrClearRequest(db, "chOrCl&" + strconv.Itoa(EgeState) + "#1", user)
	case UniState:
		user.Clear()
		text, changeMenu = handleUnisCompRequest(db, user)
	}

	return text, changeMenu
}