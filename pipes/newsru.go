package pipes

import (
    "fmt"

    . "github.com/KonishchevDmitry/go-rss"
    . "github.com/KonishchevDmitry/rsspipes"
)

func init() {
    Register("/newsru.rss", newsruFeed)
}

func newsruFeed() (feed *Feed, err error) {
    feed, err = FetchUrlWithParams("http://rss.newsru.com/top/big/", GetParams{
        SkipContentTypeCheck: true,
    })
    if err != nil {
        return
    }

    Filter(feed, func(item *Item) bool {
        if len(item.Category) == 0 {
            return true
        }

        category := item.Category[0]

        switch category {
        case "Спорт":
            return false
        }

        item.Title = fmt.Sprintf("[%s] %s", category, item.Title)
        return true
    })

    return
}
