package pipes

import (
    "crypto/md5"
    "database/sql"
    "encoding/hex"
    "fmt"
    "os/user"
    "path"
    "regexp"
    "strings"
    "time"

    "github.com/mattn/go-sqlite3"

    . "github.com/KonishchevDmitry/go-rss"
    . "github.com/KonishchevDmitry/rsspipes"
)

func init() {
    Register("/vk.rss", vkFeed)
}

func vkFeed() (feed *Feed, err error) {
    db, err := openQuotesDb()
    if err != nil {
        return
    }
    defer db.Close()

    feed, err = FetchUrl("http://konishchev.ru/social-rss/vk.rss?user_avatars=0")
    if err != nil {
        return
    }

    Filter(feed, func(item *Item) bool {
        // Filter out "New friend" items
        if item.HasCategory("type/friend") {
            return false
        }

        if item.HasCategory("source/group/club27121021") {
            // Filter "Лучшие мысли всех времен" group

            // All reposts are advertisement posts
            if item.HasCategory("type/repost") {
                return false
            }

            // Regular posts don't contain any links. All posts with links are advertisement posts.
            if strings.Contains(item.Description, "<a") {
                return false
            }

            return checkQuote(db, item.Description)
        } else if item.HasCategory("source/group/club55155418") {
            // Filter "Vert Dider" group

            if item.HasCategory("type/repost") && strings.Contains(item.Description, "#ЛекторийSetUp") {
                return false
            }

            for _, substring := range []string{
                "Расписание лектория",
                "Регистрация:",
                "Регистрация и билеты:",
                "Регистрация по ссылке:",
                "Зарегистрироваться на событие:",
                "Регистрация на мероприятие:",
            } {
                if strings.Contains(item.Description, substring) {
                    return false
                }
            }
        }

        return true
    })

    return
}

func openQuotesDb() (db *sql.DB, err error) {
    user, err := user.Current()
    if err != nil {
        return
    }

    db, err = sql.Open("sqlite3", path.Join(user.HomeDir, ".rsspipes.sqlite"))
    if err != nil {
        err = fmt.Errorf("Failed to open database: %s.", err)
        return
    }

    defer func() {
        if err != nil {
            db.Close()
        }
    }()

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS quotes (
            id TEXT NOT NULL PRIMARY KEY,
            time TIMESTAMP NOT NULL,
            text TEXT NOT NULL
        )
    `)
    if err != nil {
        err = fmt.Errorf("Failed to create a table: %s.", err)
        return
    }

    return
}

var fingerprintExtraCharsRe = regexp.MustCompile(`(\s+|<br>|[–().,:;!?'"-]+)`)
func getQuoteFingerprint(text string) string {
    fingerprint := fingerprintExtraCharsRe.ReplaceAllString(text, "")
    fingerprint = strings.Replace(fingerprint, "ё", "e", -1)
    fingerprint = strings.ToLower(fingerprint)

    hasher := md5.New()
    hasher.Write([]byte(fingerprint))
    return hex.EncodeToString(hasher.Sum(nil))
}

func checkQuote(db *sql.DB, text string) bool {
    curTime := time.Now()
    fingerprint := getQuoteFingerprint(text)

    rows, err := db.Query("SELECT time FROM quotes WHERE id == ?", fingerprint)
    if err != nil {
        log.Error("Failed to query a quote from database: %s.", err)
        return true
    }
    defer rows.Close()

    if rows.Next() {
        var quoteFirstSeenTime time.Time

        err = rows.Scan(&quoteFirstSeenTime)
        if err != nil {
            log.Error("Failed to fetch a quote info from database: %s.", err)
            return true
        }

        quoteAgeDays := curTime.Sub(quoteFirstSeenTime).Hours() / 24

        return quoteAgeDays < 3
    } else if err = rows.Err(); err != nil {
        log.Error("Failed to fetch a quote info from database: %s.", err)
        return true
    }

    _, err = db.Exec(`INSERT INTO quotes (id, time, text) VALUES (?, ?, ?)`, fingerprint, curTime, text)
    if err != nil {
        if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
            // It's OK - we've just got a race
        } else {
            log.Error("Failed to store a quote info to database: %s.", err)
        }
    }

    return true
}