package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"math"
	"strconv"
	"strings"
)

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

func makeRatingQsMenu(unisQSNum int, unisQS []*UniversityQS, curPage string) tgbotapi.InlineKeyboardMarkup {
	var unisQSButtons [][]tgbotapi.InlineKeyboardButton
	for _, uniQS := range unisQS {
		unisQSButtons = append(unisQSButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(uniQS.Name, "getUni&" + strconv.Itoa(uniQS.UniversityId) + "#" + curPage)))
	}

	curPageNum, _ := strconv.Atoi(curPage)
	paginator := makePaginator(unisQSNum, 5, curPageNum, "rateQSPage")

	var fullButtons [][]tgbotapi.InlineKeyboardButton
	fullButtons = append(fullButtons, unisQSButtons...)
	fullButtons = append(fullButtons, paginator)
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(qsButton), tgbotapi.NewInlineKeyboardRow(mainButton))
	
	ratingQSFullMenu := tgbotapi.NewInlineKeyboardMarkup(
		fullButtons...
	)

	return ratingQSFullMenu
}

func makeUniMenu(uni University, page string) tgbotapi.InlineKeyboardMarkup {
	var fullButtons [][]tgbotapi.InlineKeyboardButton
	if uni.Site != "" && !strings.Contains(uni.Site, " ") {
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("Перейти на сайт ВУЗа", uni.Site)))
	}
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Факультеты", "facs&" + strconv.Itoa(uni.UniversityId) + "#" + page + "#1")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Профили", "profs&" + strconv.Itoa(uni.UniversityId) + "#" + page + "#1")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Программы обучения", "progs&" + strconv.Itoa(uni.UniversityId) + "#" + page + "#1")),
		tgbotapi.NewInlineKeyboardRow(makeBackButton("back#" + page)),
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

func makeFacsMenu(facsNum int, facs []*Faculty, pages []string) tgbotapi.InlineKeyboardMarkup {
	uniId := facs[0].UniversityId
	unisPage := pages[0]
	facsPage := pages[1]

	var facsButtons [][]tgbotapi.InlineKeyboardButton
	for _, fac := range facs {
		facsButtons = append(facsButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(fac.Name, "getFac&" + strconv.Itoa(fac.FacultyId) + "#" + unisPage + "#" + facsPage)))
	}

	facsPageNum, _ := strconv.Atoi(facsPage)
	paginator := makePaginator(facsNum, 5, facsPageNum, "facs&" + strconv.Itoa(uniId) + "#" + unisPage)

	var fullButtons [][]tgbotapi.InlineKeyboardButton
	fullButtons = append(fullButtons, facsButtons...)
	fullButtons = append(fullButtons, paginator)
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(makeBackButton("getUni&" + strconv.Itoa(uniId) + "#" + unisPage)), tgbotapi.NewInlineKeyboardRow(mainButton))

	facsFullMenu := tgbotapi.NewInlineKeyboardMarkup(
		fullButtons...
	)

	return facsFullMenu
}

func makeFacMenu(fac Faculty, pages []string) tgbotapi.InlineKeyboardMarkup {
	uniId := fac.UniversityId
	unisPage := pages[0]
	facsPage := pages[1]

	var fullButtons [][]tgbotapi.InlineKeyboardButton
	if fac.Site != "" && !strings.Contains(fac.Site, " ") {
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("Перейти на сайт факультета", fac.Site)))
	}
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Профили", "profs&" + strconv.Itoa(fac.FacultyId) + "#" + unisPage + "#" + facsPage + "#1")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Программы обучения", "progs&" + strconv.Itoa(fac.FacultyId) + "#" + unisPage + "#" + facsPage + "#1")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Подобрать программу обучения", "findProg&" + strconv.Itoa(fac.FacultyId) + "#" + unisPage + "#" + facsPage + "#1")),
		tgbotapi.NewInlineKeyboardRow(makeBackButton("facs&" + strconv.Itoa(uniId) + "#" + unisPage + "#" + facsPage)),
		tgbotapi.NewInlineKeyboardRow(mainButton),
	)

	uniFullMenu := tgbotapi.NewInlineKeyboardMarkup(
		fullButtons...
	)

	return uniFullMenu
}

func makeUnisMenu(unisNum int, unis []*University, pagesPattern string, backPattern string, curPage string) tgbotapi.InlineKeyboardMarkup {
	var unisButtons [][]tgbotapi.InlineKeyboardButton
	for _, uni := range unis {
		unisButtons = append(unisButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(uni.Name, "getUni&" + strconv.Itoa(uni.UniversityId) + "#" + curPage)))
	}

	curPageNum, _ := strconv.Atoi(curPage)
	paginator := makePaginator(unisNum, 5, curPageNum, pagesPattern)

	var fullButtons [][]tgbotapi.InlineKeyboardButton
	fullButtons = append(fullButtons, unisButtons...)
	fullButtons = append(fullButtons, paginator)
	if backPattern != "" {
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(makeBackButton(backPattern)))
	}
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(mainButton))

	unisFullMenu := tgbotapi.NewInlineKeyboardMarkup(
		fullButtons...
	)

	return unisFullMenu
}

func makeProfsMenu(profsNum int, profs []*Profile, pagesPattern string, backPattern string, curPage string) tgbotapi.InlineKeyboardMarkup {
	var profsButtons [][]tgbotapi.InlineKeyboardButton
	for _, prof := range profs {
		profsButtons = append(profsButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(makeProfOrSpecCode(prof.ProfileId) + " " + prof.Name, "specs&" + strconv.Itoa(prof.ProfileId) + pagesPattern + "#" + curPage + "#1")))
	}

	curPageNum, _ := strconv.Atoi(curPage)
	paginator := makePaginator(profsNum, 5, curPageNum, "profs" + pagesPattern)

	var fullButtons [][]tgbotapi.InlineKeyboardButton
	fullButtons = append(fullButtons, profsButtons...)
	fullButtons = append(fullButtons, paginator)
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(makeBackButton(backPattern)), tgbotapi.NewInlineKeyboardRow(mainButton))

	profsFullMenu := tgbotapi.NewInlineKeyboardMarkup(
		fullButtons...
	)

	return profsFullMenu
}

func makeSpecsMenu(specsNum int, specs []*Speciality, pagesPattern string, backPattern string, progsPattern string, curPage string) tgbotapi.InlineKeyboardMarkup {
	var specsButtons [][]tgbotapi.InlineKeyboardButton
	for _, spec := range specs {
		specsButtons = append(specsButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(makeProfOrSpecCode(spec.SpecialityId) + " " + spec.Name, "progs&" + strconv.Itoa(spec.SpecialityId) + progsPattern + "#" + curPage + "#1")))
	}

	curPageNum, _ := strconv.Atoi(curPage)
	paginator := makePaginator(specsNum, 5, curPageNum, "specs" + pagesPattern)

	var fullButtons [][]tgbotapi.InlineKeyboardButton
	fullButtons = append(fullButtons, specsButtons...)
	fullButtons = append(fullButtons, paginator)
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(makeBackButton(backPattern)), tgbotapi.NewInlineKeyboardRow(mainButton))

	specsFullMenu := tgbotapi.NewInlineKeyboardMarkup(
		fullButtons...
	)

	return specsFullMenu
}

func makeProgsMenu(progsNum int, progs []*Program, pagesPattern string, backPattern string, progPattern string, curPage string) tgbotapi.InlineKeyboardMarkup {
	var progsButtons [][]tgbotapi.InlineKeyboardButton
	for _, prog := range progs {
		progsButtons = append(progsButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(prog.Name, "getProg&" + prog.ProgramId.String() + progPattern)))
	}

	curPageNum, _ := strconv.Atoi(curPage)
	paginator := makePaginator(progsNum, 5, curPageNum, "progs" + pagesPattern)

	var fullButtons [][]tgbotapi.InlineKeyboardButton
	fullButtons = append(fullButtons, progsButtons...)
	fullButtons = append(fullButtons, paginator)
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(makeBackButton(backPattern)), tgbotapi.NewInlineKeyboardRow(mainButton))

	progsFullMenu := tgbotapi.NewInlineKeyboardMarkup(
		fullButtons...
	)

	return progsFullMenu
}

func makeProgMenu(backPattern string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(makeBackButton(backPattern)),
		tgbotapi.NewInlineKeyboardRow(mainButton),
	)
}

func makeUnisCompilationMenu(user *UserInfo) tgbotapi.InlineKeyboardMarkup {
	var fullButtons [][]tgbotapi.InlineKeyboardButton
	filter := false

	if len(user.Eges) == 0 {
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("ЕГЭ","ege#1")))
	} else {
		filter = true
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Изменить ЕГЭ", "ege#1")))
	}

	if !user.EntryTest {
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Готов ко вступительным","entry")))
	} else {
		filter = true
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Не готов ко вступительным", "entry")))
	}

	if user.City == 0 {
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Город","city#1")))
	} else {
		filter = true
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Изменить город","chOrCl&" + strconv.Itoa(CityState))))
	}

	if user.SpecialityId == 0 && user.ProfileId == 0 {
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Профиль", "pro#1")))
	} else {
		filter = true
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Изменить профиль/cпециальность", "chOrCl&" + strconv.Itoa(ProfileState))))
	}

	if !user.Dormatary {
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Нужно общежитие", "dorm")))
	} else {
		filter = true
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Не важно общежитие", "dorm")))
	}

	if !user.MilitaryDep {
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Нужна военная кафедра", "army")))
	} else {
		filter = true
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Не важна военная кафедра", "army")))
	}

	if user.Fee == math.MaxUint64 {
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Цена", "fee")))
	} else {
		filter = true
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Изменить цену", "chOrCl&" + strconv.Itoa(FeeState))))
	}

	if filter {
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сбросить всё", "clear&" + strconv.Itoa(UniState))))
	}

	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Поиск", "search#1")))

	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(mainButton))

	unisCompFullMenu := tgbotapi.NewInlineKeyboardMarkup(
		fullButtons...
	)

	return unisCompFullMenu
}

func makeChangeOrClearMenu(state int, user *UserInfo, curPage string) tgbotapi.InlineKeyboardMarkup {
	var fullButtons [][]tgbotapi.InlineKeyboardButton
	var backPattern string
	switch state {
	case FeeState:
		backPattern = "uni"
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Изменить максимальную цену", "fee")))
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сбросить максимальную цену", "clear&" + strconv.Itoa(state))))
	case CityState:
		backPattern = "uni"
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Изменить город", "city#1")))
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сбросить город", "clear&" + strconv.Itoa(state))))
	case ProfileState:
		backPattern = "uni"
		if user.SpecialityId != 0 {
			fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Изменить специальность", "spe&" + strconv.Itoa(user.ProfileId) + "#1")))
			fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сбросить специальность", "clear&" + strconv.Itoa(SpecialityState))))
		} else {
			fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Выбрать специальность", "spe&" + strconv.Itoa(user.ProfileId) + "#1")))
		}
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Изменить профиль", "pro#1")))
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сбросить профиль", "clear&" + strconv.Itoa(state))))
	case EgeState:
		backPattern = "ege#" + curPage
		subjs := getSubjsMapFromDb()
		for _, ege := range user.Eges {
			fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(subjs[ege.SubjId], "chOrCl&" + strconv.Itoa(SubjState) + "&" + strconv.Itoa(ege.SubjId) + "#" + curPage)))
		}
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сбросить всё", "clear&" + strconv.Itoa(EgeState))))
	case SubjState:
		backPattern = "chOrCl&" + strconv.Itoa(EgeState) + "#" + curPage
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Изменить баллы", "chPoints#" + curPage)))
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сбросить", "clear&" + strconv.Itoa(state))))
	}

	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(makeBackButton(backPattern)))
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(mainButton))

	changeFullMenu := tgbotapi.NewInlineKeyboardMarkup(
		fullButtons...
	)

	log.Println("ALRIGHT:", changeFullMenu.InlineKeyboard)
	return changeFullMenu
}

func makeCitiesMenu(citiesNum int, cities []*City, backPattern string, curPage string) tgbotapi.InlineKeyboardMarkup {
	var citiesButtons [][]tgbotapi.InlineKeyboardButton
	for _, city := range cities {
		citiesButtons = append(citiesButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(city.Name, "setCity&" + strconv.Itoa(city.CityId))))
	}

	curPageNum, _ := strconv.Atoi(curPage)
	paginator := makePaginator(citiesNum, 5, curPageNum, "city")

	var fullButtons [][]tgbotapi.InlineKeyboardButton
	fullButtons = append(fullButtons, citiesButtons...)
	fullButtons = append(fullButtons, paginator)
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(makeBackButton(backPattern)))
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(mainButton))

	citiesFullMenu := tgbotapi.NewInlineKeyboardMarkup(
		fullButtons...
	)

	return citiesFullMenu
}

func makeProfsPageMenu(profsNum int, profs []*Profile, backPattern string, curPage string) tgbotapi.InlineKeyboardMarkup {
	var profsButtons [][]tgbotapi.InlineKeyboardButton
	for _, prof := range profs {
		profsButtons = append(profsButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(makeProfOrSpecCode(prof.ProfileId) + " " + prof.Name, "proOrSpe&" + strconv.Itoa(prof.ProfileId) + "#" + curPage)))
	}

	curPageNum, _ := strconv.Atoi(curPage)
	paginator := makePaginator(profsNum, 5, curPageNum, "pro")

	var fullButtons [][]tgbotapi.InlineKeyboardButton
	fullButtons = append(fullButtons, profsButtons...)
	fullButtons = append(fullButtons, paginator)
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(makeBackButton(backPattern)))
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(mainButton))

	profsFullMenu := tgbotapi.NewInlineKeyboardMarkup(
		fullButtons...
	)

	return profsFullMenu
}

func makeSpecOrNotMenu(profsPage string, profId string) tgbotapi.InlineKeyboardMarkup {
	var fullButtons [][]tgbotapi.InlineKeyboardButton

	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Выбрать специальность", "spe&" + profId + "#" + profsPage + "#1")))
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Искать только по профилю", "setPro&" + profId)))
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(makeBackButton("pro#" + profsPage)))
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(mainButton))

	specOrNotFullMenu := tgbotapi.NewInlineKeyboardMarkup(
		fullButtons...
	)

	return specOrNotFullMenu
}

func makeSpecsPageMenu(specsNum int, specs []*Speciality, profId string, pagesPattern string, backPattern string, curPage string) tgbotapi.InlineKeyboardMarkup {
	var specsButtons [][]tgbotapi.InlineKeyboardButton
	for _, spec := range specs {
		specsButtons = append(specsButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(makeProfOrSpecCode(spec.SpecialityId) + " " + spec.Name, "setSpe&" + profId + "&" + strconv.Itoa(spec.SpecialityId))))
	}

	curPageNum, _ := strconv.Atoi(curPage)
	paginator := makePaginator(specsNum, 5, curPageNum, "spe&" + profId + pagesPattern)

	var fullButtons [][]tgbotapi.InlineKeyboardButton
	fullButtons = append(fullButtons, specsButtons...)
	fullButtons = append(fullButtons, paginator)
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(makeBackButton(backPattern)))
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(mainButton))

	specsFullMenu := tgbotapi.NewInlineKeyboardMarkup(
		fullButtons...
	)

	return specsFullMenu
}

func makeMainBackMenu(data string) tgbotapi.InlineKeyboardMarkup {
	var fullButtons [][]tgbotapi.InlineKeyboardButton
	if data != "" {
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(makeBackButton(data)))
	}
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(mainButton))

	mainBackFullMenu := tgbotapi.NewInlineKeyboardMarkup(
		fullButtons...
	)

	return mainBackFullMenu
}

func makeEgesMenu(subjsNum int, subjs []*Subject, isEges bool, curPage string) tgbotapi.InlineKeyboardMarkup {
	var egesButtons [][]tgbotapi.InlineKeyboardButton
	for _, subj := range subjs {
		egesButtons = append(egesButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(subj.Name, "subj&" + strconv.Itoa(subj.SubjectId) + "#" + curPage)))
	}

	curPageNum, _ := strconv.Atoi(curPage)
	paginator := makePaginator(subjsNum, 5, curPageNum, "ege")

	var fullButtons [][]tgbotapi.InlineKeyboardButton
	fullButtons = append(fullButtons, egesButtons...)
	if isEges {
		fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Изменить/сбросить ЕГЭ", "chOrCl&" + strconv.Itoa(EgeState) + "#" + curPage)))
	}
	fullButtons = append(fullButtons, paginator)
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Готово!", "uni")))
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(mainButton))

	egesFullMenu := tgbotapi.NewInlineKeyboardMarkup(
		fullButtons...
	)

	return egesFullMenu
}

func makePointsOrNotMenu(curPage string, subjId string) tgbotapi.InlineKeyboardMarkup {
	var fullButtons [][]tgbotapi.InlineKeyboardButton

	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Искать только по предмету", "setEge&" + subjId)))
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(makeBackButton("ege#" + curPage)))
	fullButtons = append(fullButtons, tgbotapi.NewInlineKeyboardRow(mainButton))

	pointsOrNotFullMenu := tgbotapi.NewInlineKeyboardMarkup(
		fullButtons...
	)

	return pointsOrNotFullMenu
}