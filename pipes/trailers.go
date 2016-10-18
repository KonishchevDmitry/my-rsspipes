package pipes

import (
	"net/url"
	"regexp"
	"strings"

	. "github.com/KonishchevDmitry/go-rss"
	. "github.com/KonishchevDmitry/rsspipes"
)

func init() {
	Register("/trailers.rss", trailersFeed)
}

func trailersFeed() (feed *Feed, err error) {
	feedUrl := "https://www.youtube.com/feeds/videos.xml?channel_id=UC7JOTODmlIFJc8u_3BTrXZw"

	// Use this service to convert Atom feed to RSS
	feedUrl = "http://www.devtacular.com/utilities/atomtorss/?url=" + url.QueryEscape(feedUrl)

	// The service returns an invalid content type, so fetch the feed manually
	_, data, err := FetchData(feedUrl, []string{"text/html"})
	if err != nil {
		return
	}

	// Fix an invalid XML header
	data = strings.Replace(data, `<?xml version="1.0" encoding="utf-16"?>`, `<?xml version="1.0" encoding="utf-8"?>`, 1)

	feed, err = Parse([]byte(data))
	if err != nil {
		return
	}

	feed.Title = "Трейлеры"
	feed.Description = feed.Title

	nthTrailerRe := regexp.MustCompile(`трейлер \d`)

	Filter(feed, func(item *Item) bool {
		title := strings.ToLower(item.Title)
		ok :=
			// We're interested only in trailers
			strings.Index(title, "трейлер") != -1 &&

				// Skip teasers:
				// Пит и его дракон – Русский Тизер-Трейлер (2016)
				strings.Index(title, "тизер") == -1 &&

				// Only interested in first trailer:
				// Первый мститель: Противостояние - Русский Трейлер 2 (финальный, 2016)
				!nthTrailerRe.MatchString(title)

		return ok
	})

	return
}
