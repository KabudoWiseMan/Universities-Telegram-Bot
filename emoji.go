package main

import (
	"html"
	"strconv"
)

const (
	CheckEmoji = 9989
	CrossEmoji = 10060
	CryingEmoji = 128546
	WinkEmoji = 128521
)

func makeEmoji(i int) string {
	return html.UnescapeString("&#" + strconv.Itoa(i) + ";")
}
