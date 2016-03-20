package pipes

import (
    "fmt"
    "strings"

    . "github.com/KonishchevDmitry/go-rss"
    . "github.com/KonishchevDmitry/rsspipes"
)

func init() {
    Register("/meduza.rss", meduzaFeed)
}

func meduzaFeed() (feed *Feed, err error) {
    feed, err = FetchUrl("https://meduza.io/rss/all")
    if err != nil {
        return
    }

    url_prefix := "https://meduza.io/"

    Filter(feed, func(item *Item) bool {
        category := "unknown"

        if strings.HasPrefix(item.Link, url_prefix) {
            path := item.Link[len(url_prefix):]
            pos := strings.Index(path, "/")

            if pos > 0 {
                category = path[:pos]
            }
        }

        if category != "news" {
            item.Title = fmt.Sprintf("[%s] %s", category, item.Title)
        }

        return category != "shapito"
    })

    return
}
