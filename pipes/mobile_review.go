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

	skipRubrics := map[string]bool{"Кухня сайта": true, "Обзоры новинок": true, "Штучки": true}

	Filter(feed, func(item *Item) bool {
		for _, sentence := range strings.Split(item.Title, ".") {
			if skipRubrics[strings.TrimSpace(sentence)] {
				return false
			}
		}

		return true
	})

	Limit(feed, 10)

	return
}
