package pipes

import (
    . "github.com/KonishchevDmitry/go-rss"
    . "github.com/KonishchevDmitry/rsspipes"
)

var tmUserId = "53f3e807c946dfd8936807a3c3c764c9"

func init() {
    Register("/geektimes.rss", geektimesFeed)
    Register("/habrahabr.rss", habrahabrFeed)
}

func geektimesFeed() (*Feed, error) {
    return getTmFeed("Geektimes", "geektimes.ru", "/feed/" + tmUserId)
}

func habrahabrFeed() (*Feed, error) {
    return getTmFeed("Хабрахабр", "habrahabr.ru", "/feed/posts/" + tmUserId)
}

func getTmFeed(name string, domain string, userFeedPath string) (feed *Feed, err error) {
    link := "http://" + domain + "/"
    rssLink := link + "rss/"
    feedPaths := []string{userFeedPath, "best", "best/weekly", "best/monthly"}

    futureFeeds := make([]FutureFeed, len(feedPaths))
    for id, feedPath := range(feedPaths) {
        futureFeeds[id] = FutureFetch(FetchUrl, rssLink + feedPath + "/")
    }

    feeds, err := GetFutures(futureFeeds...)
    if err != nil {
        return
    }

    feed = &Feed{
        Title: name,
        Link: link,
        Image: feeds[0].Image,
    }

    Union(feed, feeds...)
    return
}