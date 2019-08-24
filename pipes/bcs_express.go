package pipes

import (
	"errors"
	"fmt"
	"strings"

	. "github.com/KonishchevDmitry/go-rss"
	. "github.com/KonishchevDmitry/rsspipes"

	"github.com/PuerkitoBio/goquery"
)

var bcsExpressUrlBuilder = UrlBuilder{"https://bcs-express.ru"}
var bcsExpressCategoryBlacklist = map[string]bool{
	"Дивиденды":        true,
	"Российский рынок": true,
	"Рынок нефти":      true,
	"Теханализ":        true,
}

func init() {
	Register("/bcs-express.rss", bcsExpressFeed)
}

func bcsExpressFeed() (feed *Feed, err error) {
	doc, err := FetchHtml(bcsExpressUrlBuilder.getUrl("/category"))
	if err != nil {
		return
	}

	items, err := getBcsExpressArticles(doc)
	if err != nil {
		return
	}

	feed = &Feed{
		Title:       "БКС Экспресс",
		Link:        bcsExpressUrlBuilder.getUrl("/"),
		Description: "Биржевые новости и аналитика от БКС Экспресс",
		Image: &Image{
			Title: "БКС Экспресс",
			Link:  bcsExpressUrlBuilder.getUrl("/"),
			Url:   bcsExpressUrlBuilder.getUrl("/favicon-16x16.png"),
		},
		Items: items,
	}

	return
}

func getBcsExpressArticles(doc *goquery.Document) (items []*Item, err error) {
	foundArticles := false

	doc.Find("div.feed div.feed-item").EachWithBreak(func(i int, article *goquery.Selection) bool {
		var item *Item

		item, err = getBcsExpressArticle(article)
		if err != nil {
			return false
		}
		foundArticles = true

		if item != nil {
			items = append(items, item)
		}

		return true
	})

	if err == nil && !foundArticles {
		err = errors.New("Unable to find the articles")
		return
	}

	return
}

func getBcsExpressArticle(article *goquery.Selection) (*Item, error) {
	url, _ := article.Find("a.feed-item__head").First().Attr("href")
	title := article.Find("div.feed-item__title").First().Text()
	summary := article.Find("div.feed-item__summary").First().Text()
	category := strings.TrimSpace(article.Find("div.rubric").First().Text())

	if url == "" || title == "" || category == "" {
		return nil, fmt.Errorf("Can't parse the following article:\n%s", getSelectionHtml(article))
	}
	url = bcsExpressUrlBuilder.getUrl(url)

	if bcsExpressCategoryBlacklist[category] {
		return nil, nil
	}

	if category == "Инвестидеи" && !strings.Contains(title, "Яндекс") && !strings.Contains(title, "YNDX") {
		return nil, nil
	}

	description, err := getBcsExpressArticleDescription(url)
	if err != nil {
		log.Errorf("Failed to get article description from %s: %s.", url, err)
		description = summary
	}

	return &Item{
		Title:       fmt.Sprintf("[%s] %s", category, title),
		Link:        url,
		Description: description,
	}, nil
}

func getBcsExpressArticleDescription(url string) (description string, err error) {
	doc, err := FetchHtml(url)
	if err != nil {
		return
	}

	return getDescriptionFromSelection(doc.Find("div.article__text"), bcsExpressUrlBuilder)
}
