package pipes

import (
    "fmt"
    "regexp"
    "strings"

    . "github.com/KonishchevDmitry/go-rss"
    . "github.com/KonishchevDmitry/rsspipes"
)

func init() {
    Register("/tv-shows.rss", tvShowsFeed)
    Register("/kate-tv-shows.rss", kateTvShowsFeed)
}

func tvShowsFeed() (feed *Feed, err error) {
    return getTvShowsFeed([]string{"Continuum", "The Walking Dead"})
}

func kateTvShowsFeed() (feed *Feed, err error) {
    tvShows, err := getKateTvShows()
    if err != nil {
        err = fmt.Errorf("Failed to fetch Kate's TV shows list: %s", err)
        return
    }

    return getTvShowsFeed(tvShows)
}

func getTvShowsFeed(tvShows []string) (feed *Feed, err error) {
    feed, err = FetchUrl("http://lostfilm.tv/rssdd.xml")
    if err != nil {
        return
    }

    tvShowsMap := make(map[string]bool)
    for _, tvShow := range tvShows {
        tvShowsMap[unifyTvShowName(tvShow)] = true
    }

    // Possible title formats:
    // Вечность (Forever). Первый сезон полностью (The Complete First Season). (S01)
    // Вечность (Forever). Первый сезон полностью (The Complete First Season) [1080p]. (S01)
    titleRe, err := regexp.Compile(`^` +
        `([^(]+)` +     // Localized name
        `\(([^)]+)\)`)  // Name

    if err != nil {
        return
    }

    Filter(feed, func(item *Item) bool {
        if !strings.Contains(item.Title, "[MP4]") {
            return false
        }

        match := titleRe.FindStringSubmatch(item.Title)
        if match == nil {
            return true
        }

        name, localizedName := unifyTvShowName(match[2]), unifyTvShowName(match[1])

        return tvShowsMap[name] || tvShowsMap[localizedName]
    })

    return
}

func getKateTvShows() (tvShows []string, err error) {
    _, data, err := FetchData("https://www.dropbox.com/s/6danram2fmm3c9u/tv-shows.txt?dl=1", []string{"text/plain"})
    if err != nil {
        return
    }

    for _, line := range strings.Split(data, "\n") {
        tvShow := stripSpaces(line)
        if tvShow != "" {
            tvShows = append(tvShows, tvShow)
        }
    }

    return
}

var spaceRe = regexp.MustCompile(`\s+`)
func stripSpaces(value string) string {
    return strings.TrimSpace(spaceRe.ReplaceAllString(value, " "))
}

func unifyTvShowName(name string) string {
    return strings.ToLower(stripSpaces(name))
}