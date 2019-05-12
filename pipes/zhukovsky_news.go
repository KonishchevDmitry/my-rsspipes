package pipes

import (
	"errors"
	"fmt"
	"strings"

	. "github.com/KonishchevDmitry/go-rss"
	. "github.com/KonishchevDmitry/rsspipes"

	"github.com/PuerkitoBio/goquery"
)

var zhukovskyNewsUrlBuilder = UrlBuilder{"http://zhukvesti.info"}

func init() {
	Register("/zhukovsky-news.rss", zhukovskyNewsFeed)
}

func zhukovskyNewsFeed() (feed *Feed, err error) {
	doc, err := FetchHtml(zhukovskyNewsUrlBuilder.getUrl("/articles/95/"))
	if err != nil {
		return
	}

	items, err := getZhukovskyNewsArticles(doc)
	if err != nil {
		return
	}

	feed = &Feed{Title: "Жуковские ВЕСТИ", Link: zhukovskyNewsUrlBuilder.getUrl("/"), Items: items}

	return
}

func getZhukovskyNewsArticles(doc *goquery.Document) (items []*Item, err error) {
	doc.Find("div.news-list div.article").EachWithBreak(func(i int, article *goquery.Selection) bool {
		var item *Item
		var spam bool

		item, spam, err = getZhukovskyNewsArticle(article)
		if err != nil {
			return false
		}

		if !spam {
			items = append(items, item)
		}

		return true
	})

	if err == nil && len(items) == 0 {
		err = errors.New("Unable to find the arcticles")
	}

	return
}

func getZhukovskyNewsArticle(article *goquery.Selection) (item *Item, spam bool, err error) {
	title := article.Find("div.newscontent h2").First()

	url, _ := title.Find("a").First().Attr("href")
	if url == "" {
		err = fmt.Errorf("Can't find URL of the following article:\n%s", getSelectionHtml(article))
		return
	}
	url = zhukovskyNewsUrlBuilder.getUrl(url)

	description, spam, err := getZhukovskyNewsArticleDescription(url)
	if err != nil {
		return
	}

	item = &Item{
		Title:       title.Text(),
		Link:        url,
		Description: description,
	}

	return
}

func getZhukovskyNewsArticleDescription(url string) (description string, spam bool, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Failed to fetch article from %s: %s", url, err)
		}
	}()

	doc, err := FetchHtml(url)
	if err != nil {
		return
	}

	article := doc.Find("div#newscontainer")
	if article.Size() != 1 {
		err = errors.New("Unable to find the article container tag.")
		return
	}

	article.Find("h1").First().Remove()
	article.Find("p.article-info").Remove()

	article.Find("div#id_ya_direct").Remove()
	article.Find("div#ilikeit").Remove()
	article.Find("div#comments-block").Remove()

	disabledCommentsStubSeparator := article.Find("hr").First()
	disabledCommentsStub :=
		disabledCommentsStubSeparator.PrevAllFiltered("br").
			AddSelection(disabledCommentsStubSeparator).
			AddSelection(disabledCommentsStubSeparator.NextAll())
	disabledCommentsPrefix := "В связи с увеличившимся количеством комментариев, " +
		"подпадающих под антиэкстремисткое законодательство"

	if strings.HasPrefix(disabledCommentsStub.Text(), disabledCommentsPrefix) {
		disabledCommentsStub.Remove()
	}

	spam = article.Find("div, p").FilterFunction(func(i int, block *goquery.Selection) bool {
		return strings.TrimSpace(block.Text()) == "На правах рекламы"
	}).Size() != 0

	description, err = getDescriptionFromSelection(article, zhukovskyNewsUrlBuilder)

	return
}
