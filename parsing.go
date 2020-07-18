package main

import (
	"fmt"
	//"fmt"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/charmap"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func getChildren(node *html.Node) []*html.Node {
	var children []*html.Node
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}
	return children
}

func isText(node *html.Node) bool {
	return node != nil && node.Type == html.TextNode
}

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func isElem(node *html.Node, tag string) bool {
	return node != nil && node.Type == html.ElementNode && node.Data == tag
}

func isDiv(node *html.Node, class string) bool {
	return isElem(node, "div") && getAttr(node, "class") == class
}

func findUnisNum() int {
	log.Println("finding number of universities")

	if response, err := http.Get(UniversitiesSite); err != nil {
		log.Println("request to " + UniversitiesSite + " failed", "error: ", err)
	} else {
		defer response.Body.Close()
		status := response.StatusCode
		log.Println("got response from " + UniversitiesSite, "status", status)
		if status == http.StatusOK {
			if doc, err := html.Parse(response.Body); err != nil {
				log.Println("invalid HTML from " + UniversitiesSite, "error", err)
			} else {
				log.Println("HTML from " + UniversitiesSite + " parsed successfully")
				return searchUnisNum(doc)
			}
		}
	}

	return -1
}

func searchUnisNum(node *html.Node) int {
	if isDiv(node, "optParent") {
		cs := getChildren(node)
		unisNumString := cs[3].FirstChild.FirstChild.Data
		unisNum, err := strconv.Atoi(unisNumString)
		if err != nil {
			log.Print("Unable to parse number of universities, got: " + unisNumString)
			return -1
		}

		return unisNum
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if unisNum := searchUnisNum(c); unisNum != -1 {
			return unisNum
		}
	}

	return -1
}

func parseUniversities() []*University {
	log.Println("parsing universities")

	unisNum := findUnisNum()
	unisPageNums := int(math.Ceil(float64(unisNum) / 15))
	pace := int(math.Ceil(float64(unisPageNums) / 10))

	var wg sync.WaitGroup

	var unis []*University

	pageString := "?page="

	for i := 1; i <= unisPageNums; i += pace + 1 {
		for j := i; j <= i + pace; j++ {
			wg.Add(1)
			go func(j int) { unis = append(unis, parsePage(&wg, UniversitiesSite + pageString + strconv.Itoa(j))...) }(j)
		}
		wg.Wait()
	}

	return unis
}

func parsePage(wg *sync.WaitGroup, pageUrl string) []*University {
	defer wg.Done()

	log.Println("sending request to " + pageUrl)
	if response, err := http.Get(pageUrl); err != nil {
		log.Println("request to " + pageUrl + " failed", "error: ", err)
	} else {
		defer response.Body.Close()
		status := response.StatusCode
		log.Println("got response from " + pageUrl, "status", status)
		if status == http.StatusOK {
			if doc, err := html.Parse(response.Body); err != nil {
				log.Println("invalid HTML from " + pageUrl, "error", err)
			} else {
				log.Println("HTML from " + pageUrl + " parsed successfully")
				return searchUniversities(doc)
			}
		}
	}

	return nil
}

func searchUniversities(node *html.Node) []*University {
	universitiesMainUrl := UniversitiesSite[:20]

	if isDiv(node, "sideContent") {
		var wg sync.WaitGroup
		var unis []*University
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if isDiv(c, "col-md-12 itemVuz") || isDiv(c, "col-md-12 itemVuzPremium") {
				if a := c.FirstChild.FirstChild.FirstChild; isElem(a, "a") {
					uniSite := universitiesMainUrl + getAttr(a, "href")
					wg.Add(1)
					go func(uniSite string) { uni := parseUniversity(&wg, uniSite); unis = append(unis, uni) }(uniSite)
				}
			}
		}

		wg.Wait()

		return unis
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if unis := searchUniversities(c); unis != nil {
			return unis
		}
	}

	return nil
}

func parseUniversity(wg *sync.WaitGroup, uniSite string) *University {
	defer wg.Done()

	log.Println("sending request to " + uniSite)
	if response, err := http.Get(uniSite); err != nil {
		log.Println("request to " + uniSite + " failed", "error: ", err)
	} else {
		defer response.Body.Close()
		status := response.StatusCode
		log.Println("got response from " + uniSite, "status", status)
		if status == http.StatusOK {
			if doc, err := html.Parse(response.Body); err != nil {
				log.Println("invalid HTML from " + uniSite, "error", err)
			} else {
				log.Println("HTML from " + uniSite + " parsed successfully")
				return searchUniversity(doc, uniSite)
			}
		}
	}

	return nil
}

func searchUniversity(node *html.Node, uniSite string) *University {
	if isDiv(node, "content clearfix") {
		cs := getChildren(node)
		universityId, err := strconv.Atoi(uniSite[25:])
		if err != nil {
			log.Println("couldn't get university id, got: " + uniSite[26:])
		}
		title, dormitary, militaryDep, description := searchUniInfo(cs[len(cs) - 2])
		phone, adress, email, site := searchUniInfo2(cs[len(cs) - 1])
		uni := &University{
			UniversityId: universityId,
			Name: title,
			Description: description,
			Site: site,
			Email: email,
			Adress: adress,
			Phone: phone,
			MilitaryDep: militaryDep,
			Dormitary: dormitary,
		}

		return uni
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if uni := searchUniversity(c, uniSite); uni != nil {
			return uni
		}
	}

	return nil
}

func searchUniInfo(node *html.Node) (string, bool, bool, string) {
	if isDiv(node, "mainSlider-left") {
		cs := getChildren(node)

		title := strings.TrimSpace(cs[0].FirstChild.Data)
		optionsDivCs := getChildren(cs[1])
		dormitary := optionsDivCs[1].FirstChild.Data == CheckMark
		militaryDep := optionsDivCs[5].FirstChild.Data == CheckMark

		description := ""
		for _, css := range cs {
			if isDiv(css, "midVuztext") {
				description = takeUniDescription(css)
				//descriptionContent := css.FirstChild
				//if isText(descriptionContent) {
				//	description = strings.TrimSpace(descriptionContent.Data)
				//} else if isElem(descriptionContent, "p") {
				//	description = strings.TrimSpace(descriptionContent.FirstChild.Data)
				//}
			}
		}

		return title, dormitary, militaryDep, description
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if title, dorm, milDep, desc := searchUniInfo(c); title != "" {
			return title, dorm, milDep, desc
		}
	}

	return "", false, false, ""
}

func takeUniDescription(node *html.Node) string {
	description := ""

	cs := getChildren(node)
	for _, css := range cs {
		if isText(css) {
			description += strings.TrimSpace(css.Data) + "\n"
		} else if isElem(css, "p") {
			for c := css.FirstChild; c != nil; c = c.NextSibling {
				if isText(c) {
					description += strings.TrimSpace(c.Data) + "\n"
				}
			}
		} else if isElem(css, "ul") {
			listCs := getChildren(css)

			for _, listCss := range listCs {
				description += "— " + strings.TrimSpace(listCss.FirstChild.Data) + "\n"
			}
		}
	}

	if len(description) > 0 {
		return description[:len(description) - 1]
	}

	return description
}

func searchUniInfo2(node *html.Node) (string, string, string, string) {
	if isDiv(node, "col-lg-6 col-md-6 col-xs-12 col-sm-6") {
		cs := getChildren(node)

		phone, adress, email, site := "", "", "", ""

		phoneContent := cs[1].LastChild.FirstChild
		if isText(phoneContent) {
			phone = phoneContent.Data
		}
		adressContent := cs[2].LastChild.FirstChild
		if isText(adressContent) {
			adress = adressContent.Data
		}
		emailContent := cs[3].LastChild.FirstChild
		if isText(emailContent) {
			email = emailContent.Data
		}
		siteContent := cs[4].LastChild.FirstChild
		if isText(siteContent) {
			site = siteContent.Data
		}

		return phone, adress, email, site
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if phone, addr, email, site := searchUniInfo2(c); phone != "" {
			return phone, addr, email, site
		}
	}

	return "", "", "", ""
}

func parseProfsNSpecs(specSite string) ([]*Profile, []*Speciality) {
	log.Println("sending request to " + specSite)
	if response, err := http.Get(specSite); err != nil {
		log.Println("request to " + specSite + " failed", "error: ", err)
	} else {
		defer response.Body.Close()
		status := response.StatusCode
		log.Println("got response from " + specSite, "status", status)
		if status == http.StatusOK {
			if doc, err := html.Parse(response.Body); err != nil {
				log.Println("invalid HTML from " + specSite, "error", err)
			} else {
				log.Println("HTML from " + specSite + " parsed successfully")
				return searchProfsNSpecs(doc)
			}
		}
	}

	return nil, nil
}

func searchProfsNSpecs(node *html.Node) ([]*Profile, []*Speciality) {
	if isElem(node, "table") {
		tBody := node.LastChild
		cs := getChildren(tBody)[1:]

		var profs []*Profile
		var specs []*Speciality

		for _, elem := range cs {
			if !isElem(elem, "tr") {
				continue
			}

			code := ""
			name := ""
			isSpec := false

			elemData := getChildren(elem)
			for _, elemDataCs := range elemData {
				if !isElem(elemDataCs, "td") {
					continue
				}

				elemDataCss := getChildren(elemDataCs)
				for _, elemDataCsss := range elemDataCss {
					if !isElem(elemDataCsss, "p") {
						continue
					}

					if getAttr(elemDataCsss, "class") == "s_1" {
						lc := elemDataCsss.FirstChild
						if isElem(lc, "a") {
							code = lc.FirstChild.Data
						} else {
							code = lc.Data
						}
					}

					if getAttr(elemDataCsss, "class") == "s_16" {
						if name != "" {
							isSpec = true
						} else {
							d := charmap.Windows1251.NewDecoder()
							st, err := d.String(elemDataCsss.LastChild.Data)
							if err != nil {
								panic(err)
							}
							name = st
						}
					}
				}
			}

			if code != "" {
				profileNum, err := strconv.Atoi(code[0:2])
				if err != nil {
					log.Println("Couldn't convert profile ID, got: " + code[0:2])
				}

				isBachelorIdent, err := strconv.Atoi(code[3:5])
				if err != nil {
					log.Println("Couldn't convert isBachelor, got: " + code[6:])
				}
				isBachelor := isBachelorIdent == 3

				specialityId, err := strconv.Atoi(strings.ReplaceAll(code, ".", ""))
				if err != nil {
					log.Println("Couldn't convert speciality ID, got: " + code[6:])
				}

				if isSpec {
					spec := &Speciality{
						SpecialityId: specialityId,
						Name:         name,
						Bachelor:     isBachelor,
						ProfileId:    profileNum * 10000,
					}
					specs = append(specs, spec)
				} else {
					prof := &Profile{
						ProfileId: profileNum * 10000,
						Name:      name,
					}
					profs = append(profs, prof)
				}
			}
		}

		return profs, specs
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if profs, specs := searchProfsNSpecs(c); profs != nil {
			return profs, specs
		}
	}

	return nil, nil
}

func main() {
	log.Println("Downloader started")
	unis := parseUniversities()
	fmt.Println(len(unis))
	if len(unis) == 739 {
		insertUnis(unis)
	}
	for _, uni := range unis {
		fmt.Println(uni.UniversityId)
	}

	// READ WHOLE PAGE
	//url := BachelorSpecialitiesSite
	//fmt.Printf("HTML code of %s ...\n", url)
	//resp, err := http.Get(url)
	//// handle the error if there is one
	//if err != nil {
	//	panic(err)
	//}
	//// do this now so it won't be forgotten
	//defer resp.Body.Close()
	//// reads html as a slice of bytes
	//html, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	panic(err)
	//}
	//
	//f, err := os.Create("/Users/vsevolodmolchanov/Downloads/specs.html")
	//if err != nil {
	//	panic(err)
	//}
	//defer f.Close()
	//
	//_, err = f.Write(html)
	//if err != nil {
	//	panic(err)
	//}

	profsBach, specsBach := parseProfsNSpecs(BachelorSpecialitiesSite)
	profsSpec, specsSpec := parseProfsNSpecs(SpecialistSpecialitiesSite)

	profs := make(map[Profile]bool)
	for _, p := range profsBach {
		profs[*p] = true
	}
	for _, p := range profsSpec {
		profs[*p] = true
	}

	insertProfsNSpecs(profs, specsBach, specsSpec)
}
