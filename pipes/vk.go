package pipes

import (
    "strings"

    . "github.com/KonishchevDmitry/go-rss"
    . "github.com/KonishchevDmitry/rsspipes"
)

func init() {
    Register("/vk.rss", vkFeed)
}

func vkFeed() (feed *Feed, err error) {
    feed, err = FetchUrl("http://konishchev.ru/social-rss/vk.rss")
    if err != nil {
        return
    }

    Filter(feed, func(item *Item) bool {
        // Filter out "New friend" items
        if item.HasCategory("type/friend") {
            return false
        }

        // Filter "Лучшие мысли всех времен" group
        if item.HasCategory("source/group/club27121021") {
            spamMarkers := []string{
                "Статьи, расширяющие понимание мира:",
                "Для тех, кто хочет изменить свою жизнь - ",
                "Для тех, кто хочет системно работать над собой - ",
                "Для тех, кто хочет начать системно работать над собой – ",
            }

            for _, marker := range(spamMarkers) {
                if strings.Contains(item.Description, marker) {
                    return false
                }
            }
        }

        return true
    })

    return
}