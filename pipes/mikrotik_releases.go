package pipes

import (
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"

	. "github.com/KonishchevDmitry/go-rss"
	. "github.com/KonishchevDmitry/rsspipes"
)

func init() {
	Register("/mikrotik-releases.rss", mikrotikReleasesFeed)
}

func mikrotikReleasesFeed() (feed *Feed, err error) {
	const url = "https://mikrotik.com/download/invalid-test/changelogs"

	doc, err := FetchHtml(url)
	if err != nil {
		return
	}

	feed = &Feed{Title: "MikroTik Releases", Link: url}

	header := doc.Find("h4:contains('Current release tree')")
	if header.Size() != 1 {
		return nil, errors.New("Unable to find current release tree")
	}

	releases := header.Closest("section").Find("ul.accordion li")
	if releases.Size() == 0 {
		return nil, errors.New("Unable to find release list")
	}

	releases.EachWithBreak(func(i int, release *goquery.Selection) bool {
		name := release.Find("a b:contains('Release')").First().Text()
		version := strings.TrimSpace(strings.Replace(name, "Release", "", -1))

		changelog := release.Find("p.chlog-p").First()
		description, renderErr := changelog.Html()
		description = strings.TrimSpace(description)

		if version == "" || description == "" || renderErr != nil {
			err = errors.New("Unable to find release description in release list")
			return false
		}

		feed.Items = append(feed.Items, &Item{
			Title:       name,
			Description: description,
			Guid:        Guid{Id: "release-" + version},
		})

		return true
	})

	return
}
