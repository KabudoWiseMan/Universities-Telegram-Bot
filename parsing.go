package main
//
//import (
//	"context"
//	"github.com/chromedp/chromedp"
//	uuid "github.com/satori/go.uuid"
//	"golang.org/x/net/html"
//	"golang.org/x/text/encoding/charmap"
//	"log"
//	"math"
//	"net/http"
//	"strconv"
//	"strings"
//	"sync"
//	"time"
//)
//
//func getChildren(node *html.Node) []*html.Node {
//	var children []*html.Node
//	for c := node.FirstChild; c != nil; c = c.NextSibling {
//		children = append(children, c)
//	}
//	return children
//}
//
//func isText(node *html.Node) bool {
//	return node != nil && node.Type == html.TextNode
//}
//
//func getAttr(node *html.Node, key string) string {
//	for _, attr := range node.Attr {
//		if attr.Key == key {
//			return attr.Val
//		}
//	}
//	return ""
//}
//
//func isElem(node *html.Node, tag string) bool {
//	return node != nil && node.Type == html.ElementNode && node.Data == tag
//}
//
//func isDiv(node *html.Node, class string) bool {
//	return isElem(node, "div") && getAttr(node, "class") == class
//}
//
//func parseProfsNSpecs(specSite string) ([]*Profile, []*Speciality) {
//	log.Println("sending request to " + specSite)
//	if response, err := http.Get(specSite); err != nil {
//		log.Println("request to " + specSite + " failed", "error: ", err)
//	} else {
//		defer response.Body.Close()
//		status := response.StatusCode
//		log.Println("got response from " + specSite, "status", status)
//		if status == http.StatusOK {
//			if doc, err := html.Parse(response.Body); err != nil {
//				log.Println("invalid HTML from " + specSite, "error", err)
//			} else {
//				log.Println("HTML from " + specSite + " parsed successfully")
//				return searchProfsNSpecs(doc)
//			}
//		}
//	}
//
//	return nil, nil
//}
//
//func searchProfsNSpecs(node *html.Node) ([]*Profile, []*Speciality) {
//	if isElem(node, "table") {
//		tBody := node.LastChild
//		cs := getChildren(tBody)[1:]
//
//		var profs []*Profile
//		var specs []*Speciality
//
//		for _, elem := range cs {
//			if !isElem(elem, "tr") {
//				continue
//			}
//
//			code := ""
//			name := ""
//			isSpec := false
//
//			elemData := getChildren(elem)
//			for _, elemDataCs := range elemData {
//				if !isElem(elemDataCs, "td") {
//					continue
//				}
//
//				elemDataCss := getChildren(elemDataCs)
//				for _, elemDataCsss := range elemDataCss {
//					if !isElem(elemDataCsss, "p") {
//						continue
//					}
//
//					if getAttr(elemDataCsss, "class") == "s_1" {
//						lc := elemDataCsss.FirstChild
//						if isElem(lc, "a") {
//							code = strings.TrimSpace(lc.FirstChild.Data)
//						} else {
//							code = strings.TrimSpace(lc.Data)
//						}
//					}
//
//					if getAttr(elemDataCsss, "class") == "s_16" {
//						if name != "" {
//							isSpec = true
//						} else {
//							d := charmap.Windows1251.NewDecoder()
//							st, err := d.String(elemDataCsss.LastChild.Data)
//							if err != nil {
//								panic(err)
//							}
//							name = st
//						}
//					}
//				}
//			}
//
//			if code != "" {
//				profileNum, err := strconv.Atoi(code[0:2])
//				if err != nil {
//					log.Println("Couldn't convert profile ID, got: " + code[0:2])
//				}
//
//				isBachelorIdent, err := strconv.Atoi(code[3:5])
//				if err != nil {
//					log.Println("Couldn't convert isBachelor, got: " + code[6:])
//				}
//				isBachelor := isBachelorIdent == 3
//
//				specialityId, err := strconv.Atoi(strings.ReplaceAll(code, ".", ""))
//				if specialityId == 0 {
//					log.Println("SHIT! " + strings.ReplaceAll(code, ".", ""))
//				}
//				if err != nil {
//					log.Println("Couldn't convert speciality ID, got: " + code)
//					log.Println(err)
//				}
//
//				if isSpec {
//					spec := &Speciality{
//						SpecialityId: specialityId,
//						Name:         name,
//						Bachelor:     isBachelor,
//						ProfileId:    profileNum * 10000,
//					}
//					specs = append(specs, spec)
//				} else {
//					prof := &Profile{
//						ProfileId: profileNum * 10000,
//						Name:      name,
//					}
//					profs = append(profs, prof)
//				}
//			}
//		}
//
//		return profs, specs
//	}
//
//	for c := node.FirstChild; c != nil; c = c.NextSibling {
//		if profs, specs := searchProfsNSpecs(c); profs != nil {
//			return profs, specs
//		}
//	}
//
//	return nil, nil
//}
//
//func findUnisNum() int {
//	log.Println("finding number of universities")
//
//	if response, err := http.Get(UniversitiesSite); err != nil {
//		log.Println("request to " + UniversitiesSite + " failed", "error: ", err)
//	} else {
//		defer response.Body.Close()
//		status := response.StatusCode
//		log.Println("got response from " + UniversitiesSite, "status", status)
//		if status == http.StatusOK {
//			if doc, err := html.Parse(response.Body); err != nil {
//				log.Println("invalid HTML from " + UniversitiesSite, "error", err)
//			} else {
//				log.Println("HTML from " + UniversitiesSite + " parsed successfully")
//				return searchUnisNum(doc)
//			}
//		}
//	}
//
//	return -1
//}
//
//func searchUnisNum(node *html.Node) int {
//	if isDiv(node, "optParent") {
//		cs := getChildren(node)
//		unisNumString := cs[3].FirstChild.FirstChild.Data
//		unisNum, err := strconv.Atoi(unisNumString)
//		if err != nil {
//			log.Print("Unable to parse number of universities, got: " + unisNumString)
//			return -1
//		}
//
//		return unisNum
//	}
//
//	for c := node.FirstChild; c != nil; c = c.NextSibling {
//		if unisNum := searchUnisNum(c); unisNum != -1 {
//			return unisNum
//		}
//	}
//
//	return -1
//}
//
//func parseUniversities() []*University {
//	log.Println("parsing universities")
//
//	unisNum := findUnisNum()
//	unisPageNums := int(math.Ceil(float64(unisNum) / 15))
//	pace := int(math.Ceil(float64(unisPageNums) / 10))
//
//	var wg sync.WaitGroup
//
//	var unis []*University
//
//	pageString := "?page="
//
//	for i := 1; i <= unisPageNums; i += pace + 1 {
//		for j := i; j <= i + pace; j++ {
//			wg.Add(1)
//			go func(j int) { unis = append(unis, parsePage(&wg, UniversitiesSite + pageString + strconv.Itoa(j))...) }(j)
//		}
//		wg.Wait()
//	}
//
//	return unis
//}
//
//func parsePage(wg *sync.WaitGroup, pageUrl string) []*University {
//	defer wg.Done()
//
//	log.Println("sending request to " + pageUrl)
//	if response, err := http.Get(pageUrl); err != nil {
//		log.Println("request to " + pageUrl + " failed", "error: ", err)
//	} else {
//		defer response.Body.Close()
//		status := response.StatusCode
//		log.Println("got response from " + pageUrl, "status", status)
//		if status == http.StatusOK {
//			if doc, err := html.Parse(response.Body); err != nil {
//				log.Println("invalid HTML from " + pageUrl, "error", err)
//			} else {
//				log.Println("HTML from " + pageUrl + " parsed successfully")
//				return searchUniversities(doc)
//			}
//		}
//	}
//
//	return nil
//}
//
//func searchUniversities(node *html.Node) []*University {
//	universitiesMainUrl := UniversitiesSite[:20]
//
//	if isDiv(node, "sideContent") {
//		var wg sync.WaitGroup
//		var unis []*University
//		for c := node.FirstChild; c != nil; c = c.NextSibling {
//			if isDiv(c, "col-md-12 itemVuz") || isDiv(c, "col-md-12 itemVuzPremium") {
//				if a := c.FirstChild.FirstChild.FirstChild; isElem(a, "a") {
//					uniSite := universitiesMainUrl + getAttr(a, "href")
//					wg.Add(1)
//					go func(uniSite string) { uni := parseUniversity(&wg, uniSite); unis = append(unis, uni) }(uniSite)
//				}
//			}
//		}
//
//		wg.Wait()
//
//		return unis
//	}
//
//	for c := node.FirstChild; c != nil; c = c.NextSibling {
//		if unis := searchUniversities(c); unis != nil {
//			return unis
//		}
//	}
//
//	return nil
//}
//
//func parseUniversity(wg *sync.WaitGroup, uniSite string) *University {
//	defer wg.Done()
//
//	log.Println("sending request to " + uniSite)
//	if response, err := http.Get(uniSite); err != nil {
//		log.Println("request to " + uniSite + " failed", "error: ", err)
//	} else {
//		defer response.Body.Close()
//		status := response.StatusCode
//		log.Println("got response from " + uniSite, "status", status)
//		if status == http.StatusOK {
//			if doc, err := html.Parse(response.Body); err != nil {
//				log.Println("invalid HTML from " + uniSite, "error", err)
//			} else {
//				log.Println("HTML from " + uniSite + " parsed successfully")
//				return searchUniversity(doc, uniSite)
//			}
//		}
//	}
//
//	return nil
//}
//
//func searchUniversity(node *html.Node, uniSite string) *University {
//	if isDiv(node, "content clearfix") {
//		cs := getChildren(node)
//		universityId, err := strconv.Atoi(uniSite[25:])
//		if err != nil {
//			log.Println("couldn't get university id, got: " + uniSite[26:])
//		}
//		title, dormitary, militaryDep, description := searchUniInfo(cs[len(cs) - 2])
//		phone, adress, email, site, _ := searchUniInfo2(cs[len(cs) - 1])
//		uni := &University{
//			UniversityId: universityId,
//			Name: title,
//			Description: description,
//			Site: site,
//			Email: email,
//			Adress: adress,
//			Phone: phone,
//			MilitaryDep: militaryDep,
//			Dormitary: dormitary,
//		}
//
//		return uni
//	}
//
//	for c := node.FirstChild; c != nil; c = c.NextSibling {
//		if uni := searchUniversity(c, uniSite); uni != nil {
//			return uni
//		}
//	}
//
//	return nil
//}
//
//func searchUniInfo(node *html.Node) (string, bool, bool, string) {
//	if isDiv(node, "mainSlider-left") {
//		cs := getChildren(node)
//
//		title := strings.TrimSpace(cs[0].FirstChild.Data)
//		optionsDivCs := getChildren(cs[1])
//		dormitary := optionsDivCs[1].FirstChild.Data == CheckMark
//		militaryDep := optionsDivCs[5].FirstChild.Data == CheckMark
//
//		description := ""
//		for _, css := range cs {
//			if isDiv(css, "midVuztext") {
//				description = takeDescription(getChildren(css))
//			}
//		}
//
//		return title, dormitary, militaryDep, description
//	}
//
//	for c := node.FirstChild; c != nil; c = c.NextSibling {
//		if title, dorm, milDep, desc := searchUniInfo(c); title != "" {
//			return title, dorm, milDep, desc
//		}
//	}
//
//	return "", false, false, ""
//}
//
//func takeDescription(nodes []*html.Node) string {
//	if len(nodes) == 0 {
//		return ""
//	} else if len(nodes) == 1 {
//		node := nodes[0]
//		if isText(node) {
//			return node.Data + "\n"
//		} else if isElem(node, "li") {
//			return "— " + takeDescription(getChildren(node))
//		} else {
//			return takeDescription(getChildren(node))
//		}
//	} else {
//		return takeDescription(nodes[0:1]) + takeDescription(nodes[1:])
//	}
//}
//
////func takeUniDescription(node *html.Node) string {
////	description := ""
////
////	cs := getChildren(node)
////	for _, css := range cs {
////		if isText(css) {
////			description += strings.TrimSpace(css.Data) + "\n"
////		} else if isElem(css, "p") {
////			for c := css.FirstChild; c != nil; c = c.NextSibling {
////				if isText(c) {
////					description += strings.TrimSpace(c.Data) + "\n"
////				}
////			}
////		} else if isElem(css, "ul") {
////			listCs := getChildren(css)
////
////			for _, listCss := range listCs {
////				description += "— " + strings.TrimSpace(listCss.FirstChild.Data) + "\n"
////			}
////		}
////	}
////
////	if len(description) > 0 {
////		return description[:len(description) - 1]
////	}
////
////	return description
////}
//
//func searchUniInfo2(node *html.Node) (string, string, string, string, bool) {
//	if isDiv(node, "col-lg-6 col-md-6 col-xs-12 col-sm-6") {
//		cs := getChildren(node)
//
//		phone, adress, email, site := "", "", "", ""
//
//		phoneContent := cs[1].LastChild.FirstChild
//		if isText(phoneContent) {
//			phone = phoneContent.Data
//		}
//		adressContent := cs[2].LastChild.FirstChild
//		if isText(adressContent) {
//			adress = adressContent.Data
//		}
//		emailContent := cs[3].LastChild.FirstChild
//		if isText(emailContent) {
//			email = emailContent.Data
//		}
//		siteContent := cs[4].LastChild.FirstChild
//		if isText(siteContent) {
//			if siteContent.Data == "http://susu.ac.ru" {
//				site = "http://www.susu.ru"
//			} else if siteContent.Data == "https://www.ba.hse.ru/" {
//				site = "https://www.hse.ru/"
//			} else {
//				site = siteContent.Data
//			}
//		}
//
//		return phone, adress, email, site, true
//	}
//
//	for c := node.FirstChild; c != nil; c = c.NextSibling {
//		if phone, addr, email, site, ret := searchUniInfo2(c); ret {
//			return phone, addr, email, site, ret
//		}
//	}
//
//	return "", "", "", "", false
//}
//
//func parseFaculties(unis []*University) []*Faculty {
//	var wg sync.WaitGroup
//
//	var facs []*Faculty
//
//	facsString := "podrazdeleniya"
//
//	if len(unis) == 0 {
//		unis = getUnisIdsNamesFromDb(false)
//	}
//
//	unisNum := len(unis)
//	pace := 15
//
//	for i := 0; i < unisNum; i += pace + 1 {
//		for j := i; j <= i + pace; j++ {
//			if j >= unisNum {
//				break
//			}
//			uniId := unis[j].UniversityId
//			uniIdString := strconv.Itoa(uniId)
//			wg.Add(1)
//			go func() { facs = append(facs, parseFacultyPage(&wg, UniversitiesSite + uniIdString + "/" + facsString, uniId)...) }()
//		}
//		wg.Wait()
//	}
//
//	return facs
//}
//
//func parseFacultyPage(wg *sync.WaitGroup, facPageSite string, uniId int) []*Faculty {
//	defer wg.Done()
//
//	log.Println("sending request to " + facPageSite)
//	if response, err := http.Get(facPageSite); err != nil {
//		log.Println("request to " + facPageSite + " failed", "error: ", err)
//	} else {
//		defer response.Body.Close()
//		status := response.StatusCode
//		log.Println("got response from " + facPageSite, "status", status)
//		if status == http.StatusOK {
//			if doc, err := html.Parse(response.Body); err != nil {
//				log.Println("invalid HTML from " + facPageSite, "error", err)
//			} else {
//				log.Println("HTML from " + facPageSite + " parsed successfully")
//				return searchFaculties(doc, uniId)
//			}
//		}
//	}
//
//	return nil
//}
//
//func searchFaculties(node *html.Node, uniId int) []*Faculty {
//	universitiesMainUrl := UniversitiesSite[:20]
//
//	if isDiv(node, "tab-pane active") && getAttr(node, "id") == "fak" {
//		var wg sync.WaitGroup
//		var facs []*Faculty
//		for c := node.FirstChild; c != nil; c = c.NextSibling {
//			if isDiv(c, "col-md-12 itemVuz") {
//				if a := c.FirstChild.FirstChild.FirstChild; isElem(a, "a") {
//					facSite := universitiesMainUrl + getAttr(a, "href")
//					wg.Add(1)
//					go func(facSite string) { fac := parseFaculty(&wg, facSite, uniId); facs = append(facs, fac) }(facSite)
//				}
//			}
//		}
//
//		wg.Wait()
//
//		return facs
//	}
//
//	for c := node.FirstChild; c != nil; c = c.NextSibling {
//		if facs := searchFaculties(c, uniId); facs != nil {
//			return facs
//		}
//	}
//
//	return nil
//}
//
//func parseFaculty(wg *sync.WaitGroup, facSite string, uniId int) *Faculty {
//	defer wg.Done()
//
//	log.Println("sending request to " + facSite)
//	if response, err := http.Get(facSite); err != nil {
//		log.Println("request to " + facSite + " failed", "error: ", err)
//	} else {
//		defer response.Body.Close()
//		status := response.StatusCode
//		log.Println("got response from " + facSite, "status", status)
//		if status == http.StatusOK {
//			if doc, err := html.Parse(response.Body); err != nil {
//				log.Println("invalid HTML from " + facSite, "error", err)
//			} else {
//				log.Println("HTML from " + facSite + " parsed successfully")
//				return searchFaculty(doc, facSite, uniId)
//			}
//		}
//	}
//
//	return nil
//}
//
//func searchFaculty(node *html.Node, facSite string, uniId int) *Faculty {
//	if isDiv(node, "content clearfix") {
//		cs := getChildren(node)
//		facultyId, err := strconv.Atoi(facSite[25:])
//		if err != nil {
//			log.Println("couldn't get university id, got: " + facSite[26:])
//		}
//		title, _, _, description := searchUniInfo(cs[len(cs) - 2])
//		phone, adress, email, site, _ := searchUniInfo2(cs[len(cs) - 1])
//		fac := &Faculty{
//			FacultyId: facultyId,
//			Name: title,
//			Description: description,
//			Site: site,
//			Email: email,
//			Adress: adress,
//			Phone: phone,
//			UniversityId: uniId,
//		}
//
//		return fac
//	}
//
//	for c := node.FirstChild; c != nil; c = c.NextSibling {
//		if fac := searchFaculty(c, facSite, uniId); fac != nil {
//			return fac
//		}
//	}
//
//	return nil
//}
//
////func parseStudyForms() map[string]int {
////	universitiesMainUrl := UniversitiesSite[:20]
////
////	log.Println("sending request to " + universitiesMainUrl)
////	if response, err := http.Get(universitiesMainUrl); err != nil {
////		log.Println("request to " + universitiesMainUrl + " failed", "error: ", err)
////	} else {
////		defer response.Body.Close()
////		status := response.StatusCode
////		log.Println("got response from " + universitiesMainUrl, "status", status)
////		if status == http.StatusOK {
////			if doc, err := html.Parse(response.Body); err != nil {
////				log.Println("invalid HTML from " + universitiesMainUrl, "error", err)
////			} else {
////				log.Println("HTML from " + universitiesMainUrl + " parsed successfully")
////				return searchStudyForms(doc)
////			}
////		}
////	}
////
////	return nil
////}
////
////func searchStudyForms(node *html.Node) map[string]int {
////	if isDiv(node, "panel-group") && getAttr(node, "id") == "acc" {
////		studyForms := make(map[string]int)
////
////		for c := node.FirstChild; c != nil; c = c.NextSibling {
////			if isDiv(c, "panel panel-default") {
////				first := c.FirstChild
////				if getAttr(first, "id") == "uslll" {
////					uslDiv := first.FirstChild
////					for uc := uslDiv.FirstChild; uc != nil; uc = uc.NextSibling {
////						if isElem(uc, "br") {
////							break
////						}
////
////						if isElem(uc, "label") {
////							studyForm := strings.TrimSpace(uc.LastChild.Data)
////							studyForms[studyForm] = len(studyForms) + 1
////						}
////					}
////				}
////			}
////		}
////
////		return studyForms
////	}
////
////	for c := node.FirstChild; c != nil; c = c.NextSibling {
////		if studyForms := searchStudyForms(c); studyForms != nil {
////			return studyForms
////		}
////	}
////
////	return nil
////}
//
//func parseSubjs() map[string]int {
//	subjsUrl := UniversitiesSite[:20] + "/kakie-ege-nuzhno-sdavat"
//
//	log.Println("sending request to " + subjsUrl)
//	if response, err := http.Get(subjsUrl); err != nil {
//		log.Println("request to " + subjsUrl + " failed", "error: ", err)
//	} else {
//		defer response.Body.Close()
//		status := response.StatusCode
//		log.Println("got response from " + subjsUrl, "status", status)
//		if status == http.StatusOK {
//			if doc, err := html.Parse(response.Body); err != nil {
//				log.Println("invalid HTML from " + subjsUrl, "error", err)
//			} else {
//				log.Println("HTML from " + subjsUrl + " parsed successfully")
//				return searchSubjs(doc)
//			}
//		}
//	}
//
//	return nil
//}
//
//func searchSubjs(node *html.Node) map[string]int {
//	if isDiv(node, "col-md-12 teloSpecFilter") {
//		subjs := make(map[string]int)
//
//		for c := node.FirstChild; c != nil; c = c.NextSibling {
//			if isElem(c, "label") {
//				name := getAttr(c, "title")
//				if name == "Вступительные" {
//					continue
//				}
//				_, ok := subjs[name]
//				if !ok {
//					subjs[name] = len(subjs) + 1
//				}
//			}
//		}
//
//		return subjs
//	}
//
//	for c := node.FirstChild; c != nil; c = c.NextSibling {
//		if subjs := searchSubjs(c); subjs != nil {
//			return subjs
//		}
//	}
//
//	return nil
//}
//
//func parsePrograms(facs []*Faculty, specs []*Speciality, subjs map[string]int) ([]*Program, []*MinEgePoints, []*EntranceTest) {
//	var wg sync.WaitGroup
//
//	var progs []*Program
//	var minPoints []*MinEgePoints
//	var entrTests []*EntranceTest
//
//	if len(facs) == 0 {
//		facs = getFacsIdsFromDb()
//	}
//	if len(specs) == 0 {
//		specs = getSpecsIdsFromDb()
//	}
//	if len(subjs) == 0 {
//		subjs = getRevSubjsMapFromDb()
//	}
//
//	specsIds := make(map[int]bool)
//	for _, spec := range specs {
//		specsIds[spec.SpecialityId] = true
//	}
//
//	facsNum := len(facs)
//	pace := 10
//
//	for i := 2700; i < facsNum; i += pace + 1 {
//		for j := i; j <= i + pace; j++ {
//			if j >= facsNum {
//				break
//			}
//			facId := facs[j].FacultyId
//			facIdString := strconv.Itoa(facId)
//			wg.Add(1)
//			go func() {
//				prog, minPoint, entrTest := parseProgramPages(&wg, UniversitiesSite + facIdString, facId, specsIds, subjs)
//				progs = append(progs, prog...)
//				minPoints = append(minPoints, minPoint...)
//				entrTests = append(entrTests, entrTest...)
//			}()
//		}
//		wg.Wait()
//	}
//
//	return progs, minPoints, entrTests
//}
//
//func findProgsNum(uniProgsSite string) int {
//	log.Println("finding number of programs")
//
//	if response, err := http.Get(uniProgsSite); err != nil {
//		log.Println("request to " + uniProgsSite + " failed", "error: ", err)
//	} else {
//		defer response.Body.Close()
//		status := response.StatusCode
//		log.Println("got response from " + uniProgsSite, "status", status)
//		if status == http.StatusOK {
//			if doc, err := html.Parse(response.Body); err != nil {
//				log.Println("invalid HTML from " + uniProgsSite, "error", err)
//			} else {
//				log.Println("HTML from " + uniProgsSite + " parsed successfully")
//				return searchProgsNum(doc)
//			}
//		}
//	}
//
//	return -1
//}
//
//func searchProgsNum(node *html.Node) int {
//	if isDiv(node, "dropdown") {
//		a := node.FirstChild
//		if isElem(a, "a") && getAttr(a, "id") == "dropdownMenuLink" {
//			progsNumString := a.FirstChild.LastChild.FirstChild.Data
//			progsNum, err := strconv.Atoi(progsNumString)
//			if err != nil {
//				log.Print("Unable to parse number of universities, got: " + progsNumString)
//				return -1
//			}
//
//			return progsNum
//		}
//	}
//
//	for c := node.FirstChild; c != nil; c = c.NextSibling {
//		if progsNum := searchProgsNum(c); progsNum != -1 {
//			return progsNum
//		}
//	}
//
//	return -1
//}
//
//func parseProgramPages(wg *sync.WaitGroup, uniSite string, facId int, specsIds map[int]bool, subjs map[string]int) ([]*Program, []*MinEgePoints, []*EntranceTest) {
//	defer wg.Done()
//
//	progsNum := findProgsNum(uniSite)
//
//	unisPageNums := int(math.Ceil(float64(progsNum) / 10))
//
//	var wg2 sync.WaitGroup
//
//	var progs []*Program
//	var minPoints []*MinEgePoints
//	var entrTests []*EntranceTest
//
//	pageString := "programs/bakispec?page="
//
//	for i := 1; i <= unisPageNums; i ++ {
//		wg2.Add(1)
//		go func(i int) {
//			prog, minPoint, entrTest := parseProgramPage(&wg2, uniSite + "/" + pageString + strconv.Itoa(i), facId, specsIds, subjs)
//			progs = append(progs, prog...)
//			minPoints = append(minPoints, minPoint...)
//			entrTests = append(entrTests, entrTest...)
//		}(i)
//	}
//	wg2.Wait()
//
//	return progs, minPoints, entrTests
//}
//
//func parseProgramPage(wg *sync.WaitGroup, progPageSite string, facId int, specsIds map[int]bool, subjs map[string]int) ([]*Program, []*MinEgePoints, []*EntranceTest) {
//	defer wg.Done()
//
//	log.Println("sending request to " + progPageSite)
//	if response, err := http.Get(progPageSite); err != nil {
//		log.Println("request to " + progPageSite + " failed", "error: ", err)
//	} else {
//		defer response.Body.Close()
//		status := response.StatusCode
//		log.Println("got response from " + progPageSite, "status", status)
//		if status == http.StatusOK {
//			if doc, err := html.Parse(response.Body); err != nil {
//				log.Println("invalid HTML from " + progPageSite, "error", err)
//			} else {
//				log.Println("HTML from " + progPageSite + " parsed successfully")
//				return searchPrograms(doc, facId, specsIds, subjs)
//			}
//		}
//	}
//
//	return nil, nil, nil
//}
//
//func searchPrograms(node *html.Node, facId int, specsIds map[int]bool, subjs map[string]int) ([]*Program, []*MinEgePoints, []*EntranceTest) {
//	universitiesMainUrl := UniversitiesSite[:20]
//
//	if isElem(node, "div") && getAttr(node, "id") == "refrdiv" {
//		var wg sync.WaitGroup
//		var progs []*Program
//		var minPoints []*MinEgePoints
//		var entrTests []*EntranceTest
//		for c := node.FirstChild; c != nil; c = c.NextSibling {
//			if isDiv(c, "col-md-12 shadowForItem") {
//				if getAttr(c, "style") != "" {
//					break
//				}
//				if a := getChildren(c.FirstChild.FirstChild)[1].FirstChild.FirstChild; isElem(a, "a") {
//					progSite := universitiesMainUrl + getAttr(a, "href")
//					progName := a.FirstChild.Data
//					wg.Add(1)
//					go func(progSite string) {
//						prog, minPoint, entrTest := parseProgram(&wg, progSite, facId, progName, subjs)
//						if _, ok := specsIds[prog.SpecialityId]; ok {
//							progs = append(progs, prog)
//							minPoints = append(minPoints, minPoint...)
//							entrTests = append(entrTests, entrTest...)
//						}
//					}(progSite)
//				}
//			}
//		}
//
//		wg.Wait()
//
//		return progs, minPoints, entrTests
//	}
//
//	for c := node.FirstChild; c != nil; c = c.NextSibling {
//		if progs, minPoints, entrTests := searchPrograms(c, facId, specsIds, subjs); progs != nil {
//			return progs, minPoints, entrTests
//		}
//	}
//
//	return nil, nil, nil
//}
//
//func parseProgram(wg *sync.WaitGroup, progSite string, facId int, progName string, subjs map[string]int) (*Program, []*MinEgePoints, []*EntranceTest) {
//	defer wg.Done()
//
//	log.Println("sending request to " + progSite)
//	if response, err := http.Get(progSite); err != nil {
//		log.Println("request to " + progSite + " failed", "error: ", err)
//	} else {
//		defer response.Body.Close()
//		status := response.StatusCode
//		log.Println("got response from " + progSite, "status", status)
//		if status == http.StatusOK {
//			if doc, err := html.Parse(response.Body); err != nil {
//				log.Println("invalid HTML from " + progSite, "error", err)
//			} else {
//				log.Println("HTML from " + progSite + " parsed successfully")
//				return searchProgram(doc, progSite, facId, progName, subjs)
//			}
//		}
//	}
//
//	return nil, nil, nil
//}
//
//func searchProgram(node *html.Node, progSite string, facId int, progName string, subjs map[string]int) (*Program, []*MinEgePoints, []*EntranceTest) {
//	if isDiv(node, "content clearfix") {
//		splitted := strings.Split(progSite, "/")
//		progNum, err := strconv.Atoi(splitted[len(splitted) - 1])
//		if err != nil {
//			log.Println("couldn't get university id, got: " + splitted[len(splitted) - 1])
//		}
//
//		var progInfo, progInfo2 *html.Node
//		for c := node.FirstChild; c != nil; c = c.NextSibling {
//			if isDiv(c, "mainBlock") {
//				progInfo = c
//			} else if isDiv(c, "wrap") {
//				progInfo2 = c
//			}
//		}
//
//		programId, err := uuid.NewV4()
//		if err != nil {
//			log.Println("Something went wrong with UUID:", err)
//		}
//		specialityId, freePassPoints, freePlaces, paidPlaces, fee := searchProgInfo(progInfo)
//		paidPassPoints, studyForm, studyLanguage, studyBase, studyYears, description, minPoints, entrTests := searchProgInfo2(progInfo2, programId, subjs)
//		prog := &Program{
//			ProgramId: programId,
//			ProgramNum: progNum,
//			Name: progName,
//			Description: description,
//			FreePlaces: freePlaces,
//			PaidPlaces: paidPlaces,
//			Fee: fee,
//			FreePassPoints: freePassPoints,
//			PaidPassPoints: paidPassPoints,
//			StudyForm: studyForm,
//			StudyLanguage: studyLanguage,
//			StudyBase: studyBase,
//			StudyYears: studyYears,
//			FacultyId: facId,
//			SpecialityId: specialityId,
//		}
//
//		return prog, minPoints, entrTests
//	}
//
//	for c := node.FirstChild; c != nil; c = c.NextSibling {
//		if prog, minPoints, entrTests := searchProgram(c, progSite, facId, progName, subjs); prog != nil {
//			return prog, minPoints, entrTests
//		}
//	}
//
//	return nil, nil, nil
//}
//
//func searchProgInfo(node *html.Node) (int, int, int, int, float64) {
//	if isDiv(node, "mainSlider-left") {
//		cs := getChildren(node)
//
//		specDivCs := getChildren(cs[1])
//		speciality := specDivCs[1].FirstChild.FirstChild.Data
//		specialityCode := speciality[len(speciality) - 9 : len(speciality) - 1]
//		specialityId, err := strconv.Atoi(strings.ReplaceAll(specialityCode, ".", ""))
//		if err != nil {
//			log.Println("Couldn't convert speciality ID, got: " + specialityCode)
//		}
//
//		freePassPoints := -1
//		freePlaces := -1
//		paidPlaces := -1
//		fee := float64(-1)
//		optDivCs := getChildren(cs[2])[2]
//		if isDiv(optDivCs, "optParent") {
//			for c := optDivCs.FirstChild; c != nil; c = c.NextSibling {
//				if freePassPoints == -1 {
//					freePassPointsData := getChildren(c)[1].FirstChild.Data
//					if freePassPointsData == "нет" {
//						freePassPoints = 0
//					} else {
//						freePassPoints, err = strconv.Atoi(freePassPointsData)
//						if err != nil {
//							log.Println("couldn't get Free pass points, got: " + freePassPointsData)
//						}
//					}
//				} else if freePlaces == -1 {
//					freePlacesData := c.FirstChild.FirstChild.Data
//					if freePlacesData == "нет" {
//						freePlaces = 0
//					} else {
//						freePlaces, err = strconv.Atoi(freePlacesData)
//						if err != nil {
//							log.Println("couldn't get Free places, got: " + freePlacesData)
//						}
//					}
//				} else if paidPlaces == -1 {
//					paidPlacesData := c.FirstChild.FirstChild.Data
//					if paidPlacesData == "нет" {
//						paidPlaces = 0
//					} else {
//						paidPlaces, err = strconv.Atoi(paidPlacesData)
//						if err != nil {
//							log.Println("couldn't get Paid places, got: " + paidPlacesData)
//						}
//					}
//				} else if fee == -1 {
//					feeData := c.FirstChild.FirstChild.Data
//					if feeData == "—" {
//						fee = 0
//					} else {
//						fee, err = strconv.ParseFloat(feeData, 64)
//						if err != nil {
//							log.Println("couldn't get Fee, got: " + feeData)
//						}
//					}
//				}
//			}
//		}
//
//		return specialityId, freePassPoints, freePlaces, paidPlaces, fee
//	}
//
//	for c := node.FirstChild; c != nil; c = c.NextSibling {
//		if specialityId, freePassPoints, freePlaces, paidPlaces, fee := searchProgInfo(c); freePassPoints != -1 {
//			return specialityId, freePassPoints, freePlaces, paidPlaces, fee
//		}
//	}
//
//	return -1, -1, -1, -1, -1
//}
//
//func searchProgInfo2(node *html.Node, programId uuid.UUID, subjs map[string]int) (int, string, string, string, string, string, []*MinEgePoints, []*EntranceTest) {
//	if isDiv(node, "sideContent progpagege") {
//		paidPassPoints := -1
//		var studyForm, studyLanguage, studyBase, studyYears, description string
//		var minPoints []*MinEgePoints
//		var entrTests []*EntranceTest
//
//		for c := node.FirstChild; c != nil; c = c.NextSibling {
//			if isDiv(c, "tab-content") {
//				for cs := c.FirstChild; cs != nil; cs = cs.NextSibling {
//					if getAttr(cs, "id") == "filial" {
//						paidInfoDiv := cs.FirstChild.LastChild.FirstChild
//						for css := paidInfoDiv.FirstChild; css != nil; css = css.NextSibling {
//							if isElem(css, "strong") {
//								splitted := strings.Split(css.FirstChild.Data, " ")
//								paidPassPoints, _ = strconv.Atoi(splitted[len(splitted) - 1])
//							}
//						}
//					} else if getAttr(cs, "id") == "fak" {
//						minPointsDiv := cs.LastChild
//						if isDiv(minPointsDiv, "col-md-3 col-sm-6 varEgeProg") {
//							isEntrance := false
//							testNames := make(map[string]bool)
//							for css := minPointsDiv.FirstChild; css != nil; css = css.NextSibling {
//								if isDiv(css, "cpPara") {
//									minPointsInfo := css.FirstChild
//									if isText(minPointsInfo) {
//										splitted := strings.Split(minPointsInfo.Data, " - ")
//										subjMinPoints, err := strconv.Atoi(splitted[len(splitted) - 1])
//										if err != nil {
//											log.Println("couldn't convers Min ege points, got: " + splitted[1])
//										}
//										if isEntrance {
//											testName := strings.TrimSpace(strings.Join(splitted[:len(splitted) - 1], " "))
//											if _, ok := testNames[testName]; ok {
//												continue
//											}
//											testNames[testName] = true
//											entrTest := &EntranceTest{
//												ProgramId: programId,
//												TestName: testName,
//												MinPoints: subjMinPoints,
//											}
//											entrTests = append(entrTests, entrTest)
//										} else {
//											subj := strings.Split(splitted[0], " (")[0]
//											if subj == "Английский" {
//												subj = "Иностранный язык"
//											}
//											subjectId, ok := subjs[subj]
//											if !ok {
//												log.Println("couldn't find subject key, got: " + subj)
//											}
//											progMinPoints := &MinEgePoints{
//												ProgramId: programId,
//												SubjectId: subjectId,
//												MinPoints: subjMinPoints,
//											}
//											minPoints = append(minPoints, progMinPoints)
//										}
//									}
//								} else if isText(css.FirstChild) && css.FirstChild.Data == "Вступительные испытания" {
//									isEntrance = true
//								}
//							}
//						}
//					}
//				}
//			} else if isDiv(c, "podrInfo") {
//				cs := getChildren(c)
//
//				studyForms := strings.TrimSpace(cs[3].LastChild.Data)
//				if len(studyForms) > 0 {
//					studyForm = studyForms[ : len(studyForms) - 1]
//				}
//
//				studyLangs := strings.TrimSpace(cs[4].LastChild.Data)
//				if len(studyLangs) > 0 {
//					studyLanguage = studyLangs[ : len(studyLangs) - 1]
//				}
//
//				studyBaseInfo := strings.TrimSpace(cs[5].LastChild.Data)
//				if len(studyBaseInfo) > 0 {
//					studyBase = studyBaseInfo[ : len(studyBaseInfo) - 1]
//				}
//
//				studyYearsInfo := strings.TrimSpace(cs[6].LastChild.Data)
//				if len(studyYearsInfo) > 0 {
//					studyYears = studyYearsInfo[ : len(studyYearsInfo) - 1]
//				}
//			} else if isDiv(c, "mainTitleBlTelo") && getAttr(c, "id") == "chemy" {
//				var nodes []*html.Node
//				for cs := c.NextSibling; cs != nil; cs = cs.NextSibling {
//					nodes = append(nodes, cs)
//				}
//				description = takeDescription(nodes)
//			}
//		}
//
//		return paidPassPoints, studyForm, studyLanguage, studyBase, studyYears, description, minPoints, entrTests
//	}
//
//	for c := node.FirstChild; c != nil; c = c.NextSibling {
//		if paidPassPoints, studyForm, studyLanguage, studyBase, studyYears, description, minPoints, entrTests := searchProgInfo2(c, programId, subjs); paidPassPoints != -1 {
//			return paidPassPoints, studyForm, studyLanguage, studyBase, studyYears, description, minPoints, entrTests
//		}
//	}
//
//	return -1, "", "", "", "", "", nil, nil
//}
//
//func parseRatingQS() []*RatingQS {
//	log.Println("sending request to " + RatingQsSite)
//	if response, err := http.Get(RatingQsSite); err != nil {
//		log.Println("request to " + RatingQsSite + " failed", "error: ", err)
//	} else {
//		defer response.Body.Close()
//		status := response.StatusCode
//		log.Println("got response from " + RatingQsSite, "status", status)
//		if status == http.StatusOK {
//			if doc, err := html.Parse(response.Body); err != nil {
//				log.Println("invalid HTML from " + RatingQsSite, "error", err)
//			} else {
//				log.Println("HTML from " + RatingQsSite + " parsed successfully")
//				return searchRatingQS(doc)
//			}
//		}
//	}
//
//	return nil
//}
//
//func searchRatingQS(node *html.Node) []*RatingQS {
//	if isDiv(node, "uni-list") {
//		mainQSSite := RatingQsSite[:31]
//		unisListString := ""
//		for c := node.FirstChild; c != nil; c = c.NextSibling {
//			if isElem(c, "a") {
//				unisListString = getAttr(c, "href")
//			}
//		}
//		if unisListString == "" {
//			return nil
//		}
//		unisHtml, err := parseRatingQSListWithChrome(mainQSSite + unisListString)
//		if err != nil {
//			return nil
//		}
//		return parseRatingQSList(unisHtml)
//	}
//
//	for c := node.FirstChild; c != nil; c = c.NextSibling {
//		if unisNames := searchRatingQS(c); unisNames != nil {
//			return unisNames
//		}
//	}
//
//	return nil
//}
//
//func parseRatingQSListWithChrome(unisListUrl string) (string, error) {
//	log.Println("sending request through Chrome to " + unisListUrl)
//	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
//	defer cancel()
//
//	ctx, cancel = context.WithTimeout(ctx, 15 * time.Second)
//	defer cancel()
//
//	var f []byte
//	var unisHtml string
//	err := chromedp.Run(ctx,
//		chromedp.Navigate(unisListUrl),
//		chromedp.SetValue(`select[name="qs-rankings_length"]`, "-1"),
//		chromedp.EvaluateAsDevTools(`$("select[name='qs-rankings_length']").change()`, &f),
//		chromedp.Sleep(time.Second * 2),
//		chromedp.OuterHTML(`table[id="qs-rankings"]`, &unisHtml),
//	)
//	if err != nil {
//		log.Println("invalid HTML through Chrome from " + unisListUrl, "error", err)
//		return "", err
//	}
//
//	log.Println("HTML from " + unisListUrl + " parsed successfully through Chrome")
//
//	return unisHtml, nil
//}
//
//func parseRatingQSList(unisHtml string) []*RatingQS {
//	if doc, err := html.Parse(strings.NewReader(unisHtml)); err != nil {
//		log.Println("invalid HTML, error", err)
//	} else {
//		log.Println("HTML parsed successfully")
//		return searchRatingQSList(doc)
//	}
//
//	return nil
//}
//
//func searchRatingQSList(node *html.Node) []*RatingQS {
//	if isElem(node, "tbody") {
//		var ratingQS []*RatingQS
//		var wg sync.WaitGroup
//		for c := node.FirstChild; c != nil; c = c.NextSibling {
//			wg.Add(1)
//			go func(c *html.Node) {
//				defer wg.Done()
//				cs := getChildren(c)
//				if cs[2].FirstChild.FirstChild.Data == "Russia" {
//					uniName := getChildren(cs[1].FirstChild)[1].FirstChild.Data
//					uniMarkElem := cs[0].FirstChild.FirstChild.FirstChild
//					var uniLowMark, uniHighMark int
//					splitted := strings.Split(uniMarkElem.LastChild.Data, "-")
//					if len(splitted) < 2 {
//						mark, err := strconv.Atoi(uniMarkElem.LastChild.Data)
//						if err != nil {
//							log.Println("couldn't convert uni qs small rating, got: " + uniMarkElem.LastChild.Data)
//						}
//						uniLowMark = mark
//						uniHighMark = mark
//					} else {
//						splitted := strings.Split(uniMarkElem.LastChild.Data, "-")
//						if len(splitted) < 2 {
//							log.Println("something's wrong with split, got: " + uniMarkElem.LastChild.Data)
//						}
//						highMark, err := strconv.Atoi(splitted[0])
//						if err != nil {
//							log.Println("couldn't convert uni qs high rating, got: " + uniMarkElem.LastChild.Data)
//						}
//						lowMark, err := strconv.Atoi(splitted[1])
//						if err != nil {
//							log.Println("couldn't convert uni qs low rating, got: " + uniMarkElem.LastChild.Data)
//						}
//						uniHighMark = highMark
//						uniLowMark = lowMark
//					}
//
//					uniSite := wikiSearch(uniName)
//					if uniSite == "" {
//						uniSite = googleWikiSearch(uniName + " wiki")
//					}
//
//					uniId := getUniIdFromDb(uniSite)
//
//					uniRatingQS := &RatingQS{
//						UniversityId: uniId,
//						HighMark: uniHighMark,
//						LowMark: uniLowMark,
//					}
//
//					ratingQS = append(ratingQS, uniRatingQS)
//				}
//			}(c)
//
//		}
//		wg.Wait()
//
//		return ratingQS
//	}
//
//	for c := node.FirstChild; c != nil; c = c.NextSibling {
//		if ratingQS := searchRatingQSList(c); ratingQS != nil {
//			return ratingQS
//		}
//	}
//
//	return nil
//}
//
//func parseCities() map[int]string {
//	universitiesMainUrl := UniversitiesSite[:20]
//	log.Println("sending request to " + universitiesMainUrl)
//	if response, err := http.Get(universitiesMainUrl); err != nil {
//		log.Println("request to " + universitiesMainUrl + " failed", "error: ", err)
//	} else {
//		defer response.Body.Close()
//		status := response.StatusCode
//		log.Println("got response from " + universitiesMainUrl, "status", status)
//		if status == http.StatusOK {
//			if doc, err := html.Parse(response.Body); err != nil {
//				log.Println("invalid HTML from " + universitiesMainUrl, "error", err)
//			} else {
//				log.Println("HTML from " + universitiesMainUrl + " parsed successfully")
//				return searchCities(doc)
//			}
//		}
//	}
//
//	return nil
//}
//
//func searchCities(node *html.Node) map[int]string {
//	if isElem(node, "select") && getAttr(node, "name") == "city[]" {
//		cities := make(map[int]string)
//		for c := node.FirstChild; c != nil; c = c.NextSibling {
//			if isElem(c, "option") {
//				key, err := strconv.Atoi(getAttr(c, "value"))
//				if err != nil {
//					log.Println("couldn't convert city key, got: " + getAttr(c, "value"))
//				}
//				var city string
//				splitted := strings.Split(c.FirstChild.Data, "(")
//				if len(splitted) > 1 {
//					cityWithBr := splitted[1]
//					city = cityWithBr[:len(cityWithBr) - 1]
//				} else {
//					city = splitted[0]
//				}
//
//				cities[key] = city
//			}
//		}
//
//		return cities
//	}
//
//	for c := node.FirstChild; c != nil; c = c.NextSibling {
//		if cities := searchCities(c); cities != nil {
//			return cities
//		}
//	}
//
//	return nil
//}
//
////func main() {
////	// Pfors and Specs
////	//profsBach, specs := parseProfsNSpecs(BachelorSpecialitiesSite)
////	//profsSpec, specsSpec := parseProfsNSpecs(SpecialistSpecialitiesSite)
////	//
////	//profsMap := make(map[Profile]bool)
////	//for _, p := range profsBach {
////	//	profsMap[*p] = true
////	//}
////	//for _, p := range profsSpec {
////	//	profsMap[*p] = true
////	//}
////	//profsBach = nil
////	//profsSpec = nil
////	//
////	//var profs []*Profile
////	//for profM, _ := range profsMap {
////	//	prof := &Profile{
////	//		ProfileId: profM.ProfileId,
////	//		Name: profM.Name,
////	//	}
////	//	profs = append(profs, prof)
////	//}
////	//profsMap = nil
////	//
////	//specs = append(specs, specsSpec...)
////	//specsSpec = nil
////	//
////	//insertProfsNSpecs(profs, specs)
////	//fmt.Println("Profiles:\n")
////	//for _, prof := range profs {
////	//	fmt.Println(*prof)
////	//}
////	//fmt.Println("\n\nSpecialities:\n")
////	//for _, spec := range specs {
////	//	fmt.Println(*spec)
////	//}
////
////	// Universities
////	//log.Println("Downloader started")
////	//unis := parseUniversities()
////	//fmt.Println(len(unis))
////	//if len(unis) != 0 {
////	//	insertUnis(unis)
////	//}
////	//for _, uni := range unis {
////	//	fmt.Println(uni.UniversityId)
////	//}
////
////	// Faculties
////	//log.Println("Downloader started")
////	//var unis []*University // needs to be deleted
////	//facs := parseFaculties(unis)
////	//if len(facs) != 0 {
////	//	insertFacs(facs)
////	//}
////	//for _, fac := range facs {
////	//	fmt.Println(*fac)
////	//}
////
////	// Subjects
////	//log.Println("Downloader started")
////	//subjs := parseSubjs()
////	//if len(subjs) != 0 {
////	//	insertSubjs(subjs)
////	//}
////	//fmt.Println("Subjects:\n")
////	//for subj, k := range subjs {
////	//	fmt.Println(k, subj)
////	//}
////
////	// Programs
////	//f, err := os.OpenFile("log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
////	//if err != nil {
////	//	log.Fatalf("error opening file: %v", err)
////	//}
////	//defer f.Close()
////	//log.SetOutput(f)
////
////	//log.Println("Downloader started")
////	//var facs []*Faculty // needs to be deleted
////	//var specs []*Speciality // needs to be deleted
////	//var subjs map[string]int // needs to be deleted
////	//progs, minPoints, entrTests := parsePrograms(facs, specs, subjs)
////	//if len(progs) != 0 {
////	//	insertProgsNInfo(progs, minPoints, entrTests)
////	//}
////	//fmt.Println("Programs:\n")
////	//fmt.Println(len(progs))
////	//for _, prog := range progs {
////	//	fmt.Println(*prog)
////	//}
////	//fmt.Println("\n\nMin points:\n")
////	//fmt.Println(len(minPoints))
////	//for _, minPoint := range minPoints {
////	//	fmt.Println(*minPoint)
////	//}
////	//fmt.Println("\n\nEntrance test:\n")
////	//fmt.Println(len(entrTests))
////	//for _, entrTest := range entrTests {
////	//	fmt.Println(*entrTest)
////	//}
////
////	// Rating QS
////	log.Println("Downloader started")
////	ratingQS := parseRatingQS()
////	if len(ratingQS) != 0 {
////		insertRatingQS(ratingQS)
////	}
////	for _, uniRatingQS := range ratingQS {
////		fmt.Println(*uniRatingQS)
////	}
////
////	//// Cities
////	//log.Println("Downloader started")
////	//cities := parseCities()
////	//if len(cities) != 0 {
////	//	insertCities(cities)
////	//}
////	//fmt.Println(len(cities))
////	//for k, c := range cities {
////	//	fmt.Println(k, c)
////	//}
////}
//
//
