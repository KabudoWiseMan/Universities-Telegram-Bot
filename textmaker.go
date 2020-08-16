package main

import (
	"strconv"
	"strings"
)

func makeTextUnisQS(unisQS []*UniversityQS) string {
	var res string
	for _, uniQS := range unisQS {
		res += "*" + uniQS.Mark + "* " + uniQS.Name + "\n\n"
	}

	return res[:len(res) - 2]
}

func makeTextUnis(unis []*University) string {
	var res string
	for _, uni := range unis {
		res += uni.Name + "\n\n"
	}

	return res[:len(res) - 2]
}

func makeTextUni(uni *University, ratingQS string) string {
	res := "*" + uni.Name + "*"
	if uni.Description != "" {
		descriptionForMarkdown := strings.ReplaceAll(uni.Description, "*", "")
		descriptionForMarkdown = strings.ReplaceAll(descriptionForMarkdown, "[", "\\[")
		descriptionForMarkdown = strings.ReplaceAll(descriptionForMarkdown, "`", "")
		res += "\n\n" + descriptionForMarkdown
	}

	if ratingQS != "" {
		res += "\n\n*Рейтинг QS:* " + ratingQS
	}

	if strings.Contains(uni.Site, " ") {
		res += "\n\n*Сайты:* " + strings.ReplaceAll(uni.Site, "_", "\\_")
	}

	if uni.Phone != "" {
		res += "\n\n*Телефон:* " + uni.Phone
	}
	if uni.Email != "" {
		res += "\n\n*E-mail:* " + strings.ReplaceAll(uni.Email, "_", "\\_")
	}
	if uni.Adress != "" {
		res += "\n\n*Адрес:* " + uni.Adress
	}

	res += "\n\n*Военная кафедра:* "
	if uni.MilitaryDep {
		res += makeEmoji(CheckEmoji)
	} else {
		res += makeEmoji(CrossEmoji)
	}

	res += "\n\n*Общежитие:* "
	if uni.Dormitary {
		res += makeEmoji(CheckEmoji)
	} else {
		res += makeEmoji(CrossEmoji)
	}

	return res
}

func makeTextFacs(facs []*Faculty) string {
	var res string
	for _, fac := range facs {
		res += fac.Name + "\n\n"
	}

	return res[:len(res) - 2]
}

func makeTextFac(fac *Faculty) string {
	res := "*" + fac.Name + "*"
	if fac.Description != "" {
		descriptionForMarkdown := strings.ReplaceAll(fac.Description, "*", "")
		descriptionForMarkdown = strings.ReplaceAll(descriptionForMarkdown, "[", "\\[")
		descriptionForMarkdown = strings.ReplaceAll(descriptionForMarkdown, "`", "")
		res += "\n\n" + descriptionForMarkdown
	}

	if strings.Contains(fac.Site, " ") {
		res += "\n\n*Сайты:* " + strings.ReplaceAll(fac.Site, "_", "\\_")
	}

	if fac.Phone != "" {
		res += "\n\n*Телефон:* " + fac.Phone
	}
	if fac.Email != "" {
		res += "\n\n*E-mail:* " + strings.ReplaceAll(fac.Email, "_", "\\_")
	}
	if fac.Adress != "" {
		res += "\n\n*Адрес:* " + fac.Adress
	}

	return res
}

func makeProfOrSpecCode(profOrSpecId int) string {
	strId := strconv.Itoa(profOrSpecId)
	if len(strId) == 5 {
		strId = "0" + strId
	}

	return strId[:2] + "." + strId[2:4] + "." + strId[4:6]
}

func makeTextProfs(profs []*Profile) string {
	var res string
	for _, prof := range profs {
		res += "*" + makeProfOrSpecCode(prof.ProfileId) + "* " + prof.Name + "\n\n"
	}

	return res[:len(res) - 2]
}

func makeTextSpecs(specs []*Speciality) string {
	var res string
	for _, spec := range specs {
		var bachelorStr string
		if spec.Bachelor {
			bachelorStr = "Бакалавриат"
		} else {
			bachelorStr = "Специалитет"
		}
		res += "*" + makeProfOrSpecCode(spec.SpecialityId) + "* " + spec.Name + " *" + bachelorStr + "*\n\n"
	}

	return res[:len(res) - 2]
}

func makeTextProgs(progs []*Program) string {
	var res string
	for _, prog := range progs {
		res += prog.Name + " *" + makeProfOrSpecCode(prog.SpecialityId) + "*\n\n"
	}

	return res[:len(res) - 2]
}

func makeTextProg(prog *ProgramInfo) string {
	res := "*" + prog.Name + "*"

	res += "\n\n*Специальность:* " + prog.SpecialityName + " (" + makeProfOrSpecCode(prog.SpecialityId) + ")"
	var bachelorStr string
	if prog.Bachelor {
		bachelorStr = "Бакалавриат"
	} else {
		bachelorStr = "Специалитет"
	}
	res += "\n*Квалификация:* " + bachelorStr

	res += "\n\n*Бюджетных мест:* "
	if prog.FreePlaces != 0 {
		res += strconv.Itoa(prog.FreePlaces) + "\n*Проходной балл на бюджет:* "
		if prog.FreePassPoints != 0 {
			res += strconv.Itoa(prog.FreePassPoints)
		} else {
			res += makeEmoji(QuestionEmoji)
		}
	} else {
		res += makeEmoji(CrossEmoji)
	}

	res += "\n\n*Платных мест:* "
	if prog.PaidPlaces != 0 {
		res += strconv.Itoa(prog.PaidPlaces) + "\n*Проходной балл на платное:* "
		if prog.PaidPassPoints != 0 {
			res += strconv.Itoa(prog.PaidPassPoints)
		} else {
			res += makeEmoji(QuestionEmoji)
		}
		res += "\n*Цена за год:* "
		if prog.Fee != 0 {
			res += strconv.Itoa(prog.Fee) + "₽"
		} else {
			res += makeEmoji(QuestionEmoji)
		}
	} else {
		res += makeEmoji(CrossEmoji)
	}

	res += "\n\n*Минимальные баллы ЕГЭ:*"
	if prog.EGEs != "" {
		res += "\n" + prog.EGEs
	} else {
		res += " " + makeEmoji(QuestionEmoji)
	}

	if prog.EntranceTests != "" {
		res += "\n\n*Вступительные испытания:*\n" + prog.EntranceTests
	}

	if prog.StudyForm != "" {
		res += "\n\n*Форма обучения:* " + prog.StudyForm
	}
	if prog.StudyYears != "" {
		res += "\n*Срок обучения:* " + prog.StudyYears
	}

	if prog.StudyBase != "" {
		res += "\n\n*На базе:* " + prog.StudyBase
	}

	if prog.StudyLanguage != "" {
		res += "\n\n*Язык обучения:* " + prog.StudyLanguage
	}

	if prog.Description != "" {
		descriptionForMarkdown := strings.ReplaceAll(prog.Description, "*", "")
		descriptionForMarkdown = strings.ReplaceAll(descriptionForMarkdown, "[", "\\[")
		descriptionForMarkdown = strings.ReplaceAll(descriptionForMarkdown, "`", "")
		res += "\n\n" + descriptionForMarkdown
	}

	return res
}

func makeTextEges(eges []Ege, subjs map[int]string, indent string) string {
	var res string
	for _, ege := range eges {
		res += indent + subjs[ege.SubjId] + " " + strconv.Itoa(int(ege.MinPoints)) + "\n"
	}

	return res[:len(res) - 1]
}