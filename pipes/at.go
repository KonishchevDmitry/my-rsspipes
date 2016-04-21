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
    commonErr := errors.New("Unable to obtain intranet cookies.")

    user, err := user.Current()
    if err != nil {
        return nil, commonErr
    }

    data, err = ioutil.ReadFile(path.Join(user.HomeDir, "Yandex", "intranet-cookies"))
    if err != nil {
        return nil, commonErr
    }

    cookies := []*http.Cookie{}

    for _, line := range strings.Split(string(data), "\n") {
        line = strings.TrimSpace(line)
        if len(line) == 0 {
            continue
        }

        cookie := strings.SplitN(line, "=", 2)
        if len(cookie) != 2 {
            return nil, commonErr
        }

        cookies = append(cookies, &http.Cookie{
            Name: cookie[0],
            Value: cookie[1],
        })
    }

    return cookies, nil
}
