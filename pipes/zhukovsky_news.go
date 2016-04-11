package pipes

import (
    "errors"
    "fmt"
    "regexp"
    "strings"

    . "github.com/KonishchevDmitry/go-rss"
    . "github.com/KonishchevDmitry/rsspipes"

    "github.com/PuerkitoBio/goquery"
)

func init() {
    Register("/zhukovsky-news.rss", zhukovskyNewsFeed)
}

func zhukovskyNewsFeed() (feed *Feed, err error) {
    doc, err := FetchHtml(getZhukovskyNewsUrl("/articles/95/"))
    if err != nil {
        return
    }

    items, err := getZhukovskyNewsArticles(doc)
    if err != nil {
        return
    }

    feed = &Feed{Title: "Жуковские ВЕСТИ", Link: getZhukovskyNewsUrl("/"), Items: items}

    return
}

func getZhukovskyNewsUrl(url string) string {
    if strings.HasPrefix(url, "/") {
        url = "http://zhukvesti.info" + url
    }

    return url
}

func getZhukovskyNewsArticles(doc *goquery.Document) (items []*Item, err error) {
    doc.Find("div.news-list div.article").EachWithBreak(func(i int, article *goquery.Selection) bool {
        var item *Item

        item, err = getZhukovskyNewsArticle(article)
        if err != nil {
            return false
        }

        items = append(items, item)

        return true
    })

    return
}

func getZhukovskyNewsArticle(article *goquery.Selection) (item *Item, err error) {
    title := article.Find("div.newscontent h2").First()

    url, _ := title.Find("a").First().Attr("href")
    if url == "" {
        err = fmt.Errorf("Can't find URL of the following article:\n%s", getSelectionHtml(article))
        return
    }
    url = getZhukovskyNewsUrl(url)

    description, err := getZhukovskyNewsArticleDescription(url)
    if err != nil {
        return
    }

    item = &Item{
        Title: title.Text(),
        Link: url,
        Description: description,
    }

    return
}

func getZhukovskyNewsArticleDescription(url string) (description string, err error) {
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

    article.Find("a").Each(func(i int, link *goquery.Selection) {
        url, exists := link.Attr("href")
        if exists {
            link.SetAttr("href", getZhukovskyNewsUrl(url))
        }
    })

    article.Find("img").Each(func(i int, image *goquery.Selection) {
        url, exists := image.Attr("src")
        if exists {
            image.SetAttr("src", getZhukovskyNewsUrl(url))
        }
    })

    description, err = article.Html()
    if err != nil {
        return
    }

    scriptRe, err := regexp.Compile(`(?is:<script(?:\s[^>]*)?>.*?</script\s*>)`)
    if err != nil {
        return
    }

    description = scriptRe.ReplaceAllString(description, "")

    return
}

func getSelectionHtml(selection *goquery.Selection) string {
    html, err := selection.Html()
    if err != nil {
        html = fmt.Sprintf("[Failed to render the HTML: %s]", err)
    }

    return html
}