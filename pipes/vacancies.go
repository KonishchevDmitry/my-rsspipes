package pipes

import (
    "regexp"

    . "github.com/KonishchevDmitry/go-rss"
    . "github.com/KonishchevDmitry/rsspipes"
)

func init() {
    Register("/vacancies.rss", vacanciesFeed)
}

func vacanciesFeed() (feed *Feed, err error) {
    urls := []string{
        "http://brainstorage.me/rss/backend.xml",
        "http://itmozg.ru/search/vacancy?VacancySearchParams%5Bkeyword%5D=&VacancySearchParams%5Bregion%5D=Москва&VacancySearchParams%5Bsalary%5D=&rss=true",
        "http://hh.ru/search/vacancy/rss?items_on_page=100&specialization=1.221&area=1&enable_snippets=true&no_magic=true&clusters=true&employment=full&search_period=30",
    }

    futureFeeds := make([]FutureFeed, len(urls))
    for id, url := range(urls) {
        futureFeeds[id] = FutureFetch(FetchUrl, url)
    }

    feeds, err := GetFutures(futureFeeds...)
    if err != nil {
        return
    }

    langRe, err := regexp.Compile(`\b(?i:python|go|golang)\b`)
    if err != nil {
        return
    }

    feed = &Feed{Title: "Вакансии"}
    Union(feed, feeds...)
    Filter(feed, func(item *Item) bool {
        return langRe.Match([]byte(item.Title))
    })

    return
}