package pipes

import (
	"fmt"
	. "github.com/KonishchevDmitry/go-rss"
	. "github.com/KonishchevDmitry/rsspipes"
	"golang.org/x/net/html"
	"strings"
)

var tmUserId = "53f3e807c946dfd8936807a3c3c764c9"

func init() {
	Register("/habrahabr.rss", habrFeed)
}

func habrFeed() (feed *Feed, err error) {
	feed, err = getTmFeed("Хабр", "habr.com", "feed/posts/"+tmUserId)
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
	link := "https://" + domain + "/"
	rssLink := link + "rss/"

	type tmFeed struct {
		Url  string
		Name string
	}

	makeTmFeed := func(path string, name string) tmFeed {
		return tmFeed{rssLink + path + "/", name}
	}

	tmFeeds := []tmFeed{
		makeTmFeed(userFeedPath, "my"),
		makeTmFeed("best", ""),
		makeTmFeed("best/weekly", "weekly"),
		makeTmFeed("best/monthly", "monthly"),
	}

	futureFeeds := make([]FutureFeed, 0, len(tmFeeds))
	for feedId, _ := range tmFeeds {
		tmFeed := &tmFeeds[feedId]
		futureFeeds = append(futureFeeds, FutureFetch(func(url string) (feed *Feed, err error) {
			return fetchNamedFeed(url, tmFeed.Name)
		}, tmFeed.Url))
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

	for _, item := range feed.Items {
		if len(item.Category) != 0 {
			item.Description += fmt.Sprintf("<p>%s</p>", html.EscapeString(strings.Join(item.Category, " | ")))
		}
	}

	return
}

func fetchNamedFeed(url string, name string) (feed *Feed, err error) {
	feed, err = FetchUrl(url)
	if err != nil {
		return
	}

	if name != "" {
		titlePrefix := fmt.Sprintf("[%s] ", name)
		for _, item := range feed.Items {
			item.Title = titlePrefix + item.Title
		}
	}

	return
}
