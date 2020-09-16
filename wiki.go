package main

import (
	"encoding/json"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func urlEncoded(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func wikiSearch(query string) string {
	queryUrl, err := urlEncoded(query)
	if err != nil {
		log.Fatal(err)
	}
	request := "https://en.wikipedia.org/w/api.php?action=opensearch&search=" + queryUrl + "&limit=1&origin=*&format=json"

	if response, err := http.Get(request); err != nil {
		log.Println("request to " + request + " failed", "error: ", err)
	} else {
		defer response.Body.Close()

		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		sr := []interface{}{}
		if err = json.Unmarshal(contents, &sr); err != nil {
			log.Fatal("Something going wrong, try to change your question")
		}

		if len(sr[1].([]interface{})) > 0 {
			uniWikiUrl := sr[3].([]interface{})[0].(string)

			uniSite := parseUniWikiUrl(uniWikiUrl)
			return uniSite
		}
	}

	return ""
}

func parseUniWikiUrl(uniWikiUrl string) string {
	if response, err := http.Get(uniWikiUrl); err != nil {
		log.Println("request to " + uniWikiUrl + " failed", "error: ", err)
	} else {
		defer response.Body.Close()
		status := response.StatusCode
		log.Println("got response from " + uniWikiUrl, "status", status)
		if status == http.StatusOK {
			if doc, err := html.Parse(response.Body); err != nil {
				log.Println("invalid HTML from " + uniWikiUrl, "error", err)
			} else {
				log.Println("HTML from " + uniWikiUrl + " parsed successfully")
				return uniWikiSiteSearch(doc)
			}
		}
	}

	return ""
}

func uniWikiSiteSearch(node *html.Node) string {
	if isElem(node, "th") && isText(node.FirstChild) && node.FirstChild.Data == "Website" {
		fc := node.NextSibling.FirstChild
		var uniUrl string
		if isElem(fc, "a") {
			uniUrl = getAttr(fc, "href")
		} else if isElem(fc, "span") {
			uniUrl = getAttr(fc.FirstChild, "href")
		}

		if strings.Contains(uniUrl, "phystech") {
			uniUrl = "https://mipt.ru"
		}

		splitted := strings.Split(uniUrl, ".")
		if (len(splitted) > 1) {
			uniUrlForSearch := splitted[len(splitted) - 2]
			sub := substrAfter(uniUrlForSearch, "//")
			if sub != "" {
				return sub
			}
			return uniUrlForSearch
		}

		return uniUrl
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if res := uniWikiSiteSearch(c); res != "" {
			return res
		}
	}

	return ""
}

func googleWikiSearch(query string) string {
	queryUrl, err := urlEncoded(query)
	if err != nil {
		log.Fatal(err)
	}

	request := "https://www.google.com/search?q=" + queryUrl + "&num=20"

	if response, err := http.Get(request); err != nil {
		log.Println("request to " + request + " failed", "error: ", err)
	} else {
		defer response.Body.Close()
		status := response.StatusCode
		log.Println("got response from " + request, "status", status)
		if status == http.StatusOK {
			if doc, err := html.Parse(response.Body); err != nil {
				log.Println("invalid HTML from " + request, "error", err)
			} else {
				log.Println("HTML from " + request + " parsed successfully")
				return googleUniWikiSiteSearch(doc)
			}
		}
	}

	return ""
}

func googleUniWikiSiteSearch(node *html.Node) string {
	if isElem(node, "a") && strings.Contains(getAttr(node, "href"), "en.wikipedia.org") {
		uniGoogleWikiUrl := getAttr(node, "href")
		uniWikiUrl := substrBetween(uniGoogleWikiUrl, "/url?q=", "&")
		encodedUrl := strings.ReplaceAll(uniWikiUrl, "%25", "%")
		decodedUrl, err := url.QueryUnescape(encodedUrl)
		if err != nil {
			log.Fatal(err)
		}

		uniSite := parseUniWikiUrl(decodedUrl)
		return uniSite
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if res := googleUniWikiSiteSearch(c); res != "" {
			return res
		}
	}

	return ""
}