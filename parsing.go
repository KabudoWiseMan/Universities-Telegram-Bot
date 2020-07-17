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
	if isDiv(node, "sideContent") {
		var unis []*University
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if isDiv(c, "col-md-12 itemVuz") {
				title := strings.TrimSpace(c.FirstChild.FirstChild.FirstChild.LastChild.FirstChild.Data)
				uni := &University{
					Name: title,
				}
				unis = append(unis, uni)
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

func main() {
	log.Println("Downloader started")
	unis := parseUniversities()
	fmt.Println(len(unis))
	for _, uni := range unis {
		fmt.Println(*uni)
	}
}
