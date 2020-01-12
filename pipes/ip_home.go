package pipes

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"

	. "github.com/KonishchevDmitry/go-rss"
	. "github.com/KonishchevDmitry/rsspipes"
)

var ipHomeUrlBuilder = UrlBuilder{"https://ip-home.net/"}

func init() {
	Register("/ip-home.rss", ipHomeFeed)
}

func ipHomeFeed() (feed *Feed, err error) {
	doc, err := FetchHtml(ipHomeUrlBuilder.BaseUrl)
	if err != nil {
		return
	}

	feed = &Feed{Title: "IP-Home", Link: ipHomeUrlBuilder.BaseUrl}

	newsBlock := doc.Find("div.news-content")
	if newsBlock.Size() != 1 {
		return nil, errors.New("Unable to find news block")
	}

	newsItems := newsBlock.Find("div.item")
	if newsItems.Size() == 0 {
		return nil, errors.New("Unable to find news item blocks")
	}

	newsItems.EachWithBreak(func(i int, item *goquery.Selection) bool {
		url, ok := item.Find("a").First().Attr("href")
		title := strings.TrimSpace(item.Find("a.show-news").Text())

		if !ok || title == "" {
			err = errors.New("Unable to find news URL/title")
			return false
		}

		url = ipHomeUrlBuilder.getUrl(url)
		description := getIPHomeNewsDescription(url)

		feed.Items = append(feed.Items, &Item{
			Title:       title,
			Link:        url,
			Description: description,
		})
		return true
	})

	return
}

func getIPHomeNewsDescription(url string) (description string) {
	var err error
	defer func() {
		if err != nil {
			description = fmt.Sprintf("Failed to fetch news description from %s: %s.", url, err)
		}
	}()

	doc, err := FetchHtml(url)
	if err != nil {
		return
	}

	contents := doc.Find("div.single-news-text")
	if contents.Size() != 1 {
		err = errors.New("Unable to find news block")
		return
	}

	description, err = getDescriptionFromSelection(contents, ipHomeUrlBuilder)
	return
}
