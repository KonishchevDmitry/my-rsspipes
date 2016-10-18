package pipes

import (
	"strings"

	. "github.com/KonishchevDmitry/go-rss"
	. "github.com/KonishchevDmitry/rsspipes"
)

var tmUserId = "53f3e807c946dfd8936807a3c3c764c9"

func init() {
	Register("/geektimes.rss", geektimesFeed)
	Register("/habrahabr.rss", habrahabrFeed)
}

func geektimesFeed() (*Feed, error) {
	return getTmFeed("Geektimes", "geektimes.ru", "/feed/"+tmUserId)
}

func habrahabrFeed() (feed *Feed, err error) {
	feed, err = getTmFeed("Хабрахабр", "habrahabr.ru", "/feed/posts/"+tmUserId)
	if err != nil {
		return nil, err
	}

	blogBlacklist := []string{
		"Блог компании PVS-Studio",
		"Блог компании Vivaldi Technologies AS",
	}

	Filter(feed, func(item *Item) bool {
		for _, blogName := range blogBlacklist {
			if item.HasCategory(blogName) {
				return false
			}
		}

		return true
	})

	return
}

func getTmFeed(name string, domain string, userFeedPath string) (feed *Feed, err error) {
	link := "http://" + domain + "/"
	rssLink := link + "rss/"
	feedPaths := []string{userFeedPath, "best", "best/weekly", "best/monthly"}

	futureFeeds := make([]FutureFeed, len(feedPaths))
	for id, feedPath := range feedPaths {
		futureFeeds[id] = FutureFetch(FetchUrl, rssLink+feedPath+"/")
	}

	subFeeds, err := GetFutures(futureFeeds...)
	if err != nil {
		return
	}

	// URL may contain query parameters which are different in each subfeed:
	// https://habrahabr.ru/post/279703/?utm_source=habrahabr&utm_medium=rss&utm_campaign=interesting
	//
	// Strip them in GUID to make union work right.
	for _, subFeed := range subFeeds {
		for _, item := range subFeed.Items {
			if item.Guid.IsPermaLink == nil || !*item.Guid.IsPermaLink {
				continue
			}

			query_params_index := strings.Index(item.Guid.Id, "?")
			if query_params_index != -1 {
				item.Guid.Id = item.Guid.Id[:query_params_index]
			}
		}
	}

	feed = &Feed{
		Title: name,
		Link:  link,
		Image: subFeeds[0].Image,
	}

	Union(feed, subFeeds...)
	return
}
