package pipes

import (
    "strings"

    . "github.com/KonishchevDmitry/go-rss"
    . "github.com/KonishchevDmitry/rsspipes"
)

func init() {
    Register("/mobile-review.rss", mobileReviewFeed)
}

func mobileReviewFeed() (feed *Feed, err error) {
    feed, err = FetchUrl("http://www.mobile-review.com/podcasts/rss.xml")
    if err != nil {
        return
    }

    Filter(feed, func(item *Item) bool {
        skipRubrics := []string{"Кухня сайта", "Обзоры новинок", "Штучки"}

        for _, rubric := range(skipRubrics) {
            if strings.HasPrefix(item.Title, rubric + ". ") {
                return false
            }
        }

        return true
    })

    Limit(feed, 10)

    return
}