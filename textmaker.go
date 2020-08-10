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

func makeTextUni(uni University) string {
	res := "*" + uni.Name + "*"
	if uni.Description != "" {
		res += "\n\n" + uni.Description
	}

	ratingQS := getUniQSRateFromDb(uni.UniversityId)
	if ratingQS != "" {
		res += "\n\n*Рейтинг QS:* " + ratingQS
	}

	if strings.Contains(uni.Site, " ") {
		res += "\n\n*Сайты:* " + uni.Site
	}

	if uni.Phone != "" {
		res += "\n\n*Телефон:* " + uni.Phone
	}
	if uni.Email != "" {
		res += "\n\n*E-mail:* " + uni.Email
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

func makeTextFac(fac Faculty) string {
	res := "*" + fac.Name + "*"
	if fac.Description != "" {
		res += "\n\n" + fac.Description
	}

	if strings.Contains(fac.Site, " ") {
		res += "\n\n*Сайты:* " + fac.Site
	}

	if fac.Phone != "" {
		res += "\n\n*Телефон:* " + fac.Phone
	}
	if fac.Email != "" {
		res += "\n\n*E-mail:* " + fac.Email
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