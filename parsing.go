package main

import (
	"fmt"
	"golang.org/x/net/html"
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

	var wg sync.WaitGroup

	var unis []*University

	pageString := "?page="

	for i := 1; i <= unisPageNums; i++ {
		wg.Add(1)
		go func(i int) { unis = append(unis, parsePage(&wg, UniversitiesSite + pageString + strconv.Itoa(i))...) }(i)
	}

	wg.Wait()

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
		var unis []*University
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if isDiv(c, "col-md-12 itemVuz") || isDiv(c, "col-md-12 itemVuzPremium") {
				if a := c.FirstChild.FirstChild.FirstChild; isElem(a, "a") {
					uniSite := universitiesMainUrl + getAttr(a, "href")
					uni := parseUniversity(uniSite)
					unis = append(unis, uni)
				}
			}
		}

		return unis
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if unis := searchUniversities(c); unis != nil {
			return unis
		}
	}

	return nil
}

func parseUniversity(uniSite string) *University {
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
				return searchUniversity(doc)
			}
		}
	}

	return nil
}

func searchUniversity(node *html.Node) *University {
	if isDiv(node, "content clearfix") {
		cs := getChildren(node)
		title, dormitary, militaryDep, description := searchUniInfo(cs[len(cs) - 2])
		phone, adress, email, site := searchUniInfo2(cs[len(cs) - 1])
		uni := &University{
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
		if uni := searchUniversity(c); uni != nil {
			return uni
		}
	}

	return nil
}

func searchUniInfo(node *html.Node) (string, bool, bool, string) {
	if isDiv(node, "mainSlider-left") {
		cs := getChildren(node)
		title := strings.TrimSpace(cs[0].FirstChild.Data)
		description := ""
		descriptionContent := cs[6].FirstChild
		if isText(descriptionContent) {
			description = strings.TrimSpace(descriptionContent.Data)
		}
		optionsDivCs := getChildren(cs[1])
		dormitary := optionsDivCs[1].FirstChild.Data == CheckMark
		militaryDep := optionsDivCs[5].FirstChild.Data == CheckMark

		return title, dormitary, militaryDep, description
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if title, dorm, milDep, desc := searchUniInfo(c); title != "" {
			return title, dorm, milDep, desc
		}
	}

	return "", false, false, ""
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

func main() {
	log.Println("Downloader started")
	unis := parseUniversities()
	fmt.Println(len(unis))
	for _, uni := range unis {
		fmt.Println(*uni)
	}
}
