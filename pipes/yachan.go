package pipes

import (
    "bytes"
    "runtime"
    "strings"

    "github.com/PuerkitoBio/goquery"
    "golang.org/x/net/html"

    . "github.com/KonishchevDmitry/go-rss"
    . "github.com/KonishchevDmitry/rsspipes"
)

func init() {
    if runtime.GOOS == "darwin" {
        Register("/yachan.rss", yachanFeed)
    }
}

func yachanFeed() (feed *Feed, err error) {
    baseUrl := "https://yachan.dev.yandex.net"

    feed, err = FetchUrlWithParams(baseUrl + "/bbs/~/feed/auth/9881/Anonymous.rss", GetParams{
        SkipCertificateCheck: true,
    })

    if err != nil {
        return
    }

    feed.Title = "YaChan"

    for _, item := range feed.Items {
        if strings.HasPrefix(item.Link, "/") {
            item.Link = baseUrl + item.Link
        }

        if doc, err := html.Parse(bytes.NewReader([]byte(item.Description))); err == nil {
            doc := goquery.NewDocumentFromNode(doc)

            title := doc.Find("span.replytitle").Text()
            if title != "" {
                item.Title = title
            }

            description := doc.Find("blockquote.postbody")
            if description.Size() != 0 {
                if description, err := description.Html(); err == nil {
                    item.Description = description
                }
            }
        }
    }

    return
}
