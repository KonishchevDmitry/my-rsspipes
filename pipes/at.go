package pipes

import (
    "errors"
    "io/ioutil"
    "net/http"
    "os/user"
    "path"
    "runtime"
    "strings"

    . "github.com/KonishchevDmitry/go-rss"
    . "github.com/KonishchevDmitry/rsspipes"
)

func init() {
    if runtime.GOOS == "darwin" {
        Register("/at.rss", atFeed)
        Register("/bbs.rss", bbsFeed)
    }
}

func atFeed() (feed *Feed, err error) {
    urls := []string{
        "https://my.at.yandex-team.ru/rss/popular.xml",
        "https://konishchev.at.yandex-team.ru/rss/friends.xml",
    }

    cookies, err := getYandexTeamCookies()
    if err != nil {
        return
    }

    getParams := GetParams{
        Cookies: cookies,
        SkipContentTypeCheck: true,
    }

    fetchFunc := func(url string) (*Feed, error) {
        return FetchUrlWithParams(url, getParams)
    }

    futureFeeds := make([]FutureFeed, len(urls))
    for id, url := range(urls) {
        futureFeeds[id] = FutureFetch(fetchFunc, url)
    }

    feeds, err := GetFutures(futureFeeds...)
    if err != nil {
        return
    }

    feed = &Feed{
        Title: "Этушка",
        Link: "https://my.at.yandex-team.ru/",
        Image: feeds[0].Image,
    }

    Union(feed, feeds...)

    return
}

func bbsFeed() (feed *Feed, err error) {
    cookies, err := getYandexTeamCookies()
    if err != nil {
        return
    }

    feed, err = FetchUrlWithParams("https://clubs.at.yandex-team.ru/bbs/rss/posts.xml", GetParams{
        Cookies: cookies,
    })

    return
}

func getYandexTeamCookies() ([]*http.Cookie, error) {
    var data []byte

    user, err := user.Current()
    if err == nil {
        data, err = ioutil.ReadFile(path.Join(user.HomeDir, "Yandex", "session-id"))
    }

    if err != nil {
        return nil, errors.New("Unable to obtain session ID.")
    }

    return []*http.Cookie{&http.Cookie{
        Name: "Session_id",
        Value: strings.TrimSpace(string(data)),
    }}, nil
}
