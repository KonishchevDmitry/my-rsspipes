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
    feed, err = FetchUrl("http://konishchev.ru/social-rss/vk.rss?user_avatars=0")
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
            // All reposts are advertisement posts
            if item.HasCategory("type/repost") {
                return false
            }

            // Regular posts don't contain any links. All posts with links are advertisement posts.
            if strings.Contains(item.Description, "<a") {
                return false
            }
        }

        return true
    })

    return
}