package content

import (
	"strings"
	"github.com/PuerkitoBio/goquery"
	"log"
)

type RedditContent struct {
	Url string
}

func (reddit RedditContent) callAPI() ([]Content, error) {
	resp, err := getWebserviceResponse(reddit.Url)
	if( err != nil ) {
		log.Println("Error while calling url: "+reddit.Url)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Println("Error while calling API")
		return nil, err
	}

	rv := make([]Content, 0)

	doc.Find("div.search-result,div.entry").Each(func(i int, selec *goquery.Selection) {
		// ignore sticky posts
		if selec.HasClass("stickied") {
			return
		}

		if len(rv) > 20 {
			return
		}

		t := selec.Find("a.search-title,a.title")
		title := t.First().Text()

		// Limit size of content
		if( len(title) + urlLength > 140 ) {
			title = title[0:139-urlLength] + "…"
		}

		l := selec.Find("a.search-link,a.link")
		externalLink, _ := l.First().Attr("href")

		// self posts
		if strings.HasPrefix(externalLink, "/r/") {
 			externalLink = "https://reddit.com" + externalLink
 		}

		rv = append(rv, Content{
			Text: title,
			Url:  externalLink,
		})
	})

	return rv, nil
}
