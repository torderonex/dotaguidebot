package parser

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var count int

func Parser(hero string) ([]string, string, error) {
	currentUrl := heroFormat(hero)
	count = 1
	resp, err := http.Get(currentUrl)
	if resp.StatusCode != 200 {
		return []string{}, "", errors.New("This hero does not exist")
	}
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	roles := []string{}
	var parse func(*html.Node)
	parse = func(body *html.Node) {
		if body.Type == html.ElementNode && body.Data == "button" {
			for _, a := range body.Attr {
				if a.Key == "id" && strings.HasPrefix(a.Val, "tabs-") {
					if a.Val == "tabs-"+strconv.Itoa(count) {
						roles = append(roles, strip(body.LastChild.Data))
						count++
					} else {
						break
					}

				}
			}
		}
		for c := body.FirstChild; c != nil; c = c.NextSibling {
			parse(c)
		}

	}
	parse(body)
	return roles, currentUrl, nil

}

func strip(str string) string {
	var new string
	for _, chr := range str {
		if chr != ' ' && chr != 10 {
			new += string(chr)
		}
	}
	return new
}

func heroFormat(hero string) string {
	var url string
	hero = strings.Title(hero)
	path := "https://www.dota2protracker.com/hero/"
	split := strings.Split(hero, " ")
	if len(split) != 1 {
		url = path + split[0] + "%20" + split[1]
	} else {
		url = path + hero
	}
	return url
}
