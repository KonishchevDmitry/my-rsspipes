package pipes

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os/user"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"

	. "github.com/KonishchevDmitry/go-rss"
	. "github.com/KonishchevDmitry/rsspipes"
)

func init() {
	Register("/vk.rss", vkNewsFeed)
	Register("/vert-dider.rss", vertDiderFeed)
}

func vkNewsFeed() (feed *Feed, err error) {
	feed, err = vkFeed()
	if err != nil {
		return
	}

	db, err := openSeenDb()
	if err != nil {
		return
	}
	defer db.Close()

	Filter(feed, func(item *Item) bool {
		// Filter out "Vert Dider" group (it has it's own RSS)
		if item.HasCategory("source/group/club55155418") {
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

			allow, err := checkAndRememberQuote(db, item.Description)
			if err != nil {
				log.Errorf("Failed to check a quote against seen quotes database: %s.", err)
			}

			return allow
		}

		return true
	})

	return
}

func vertDiderFeed() (feed *Feed, err error) {
	feed, err = vkFeed()
	if err != nil {
		return
	}

	db, err := openSeenDb()
	if err != nil {
		return
	}
	defer db.Close()

	feed.Title = "Vert Dider"

	Filter(feed, func(item *Item) bool {
		// Filter "Vert Dider" group
		if !item.HasCategory("source/group/club55155418") {
			return false
		}

		var allowPost bool
		const videoAttachmentPrefix = "attachment/video/"

		for _, category := range item.Category {
			if !strings.HasPrefix(category, videoAttachmentPrefix) {
				continue
			}

			videoId, err := strconv.ParseInt(category[len(videoAttachmentPrefix):], 10, 64)
			if err != nil {
				log.Errorf("Got an invalid category: %q.", category)
				allowPost = true
				continue
			}

			allowVideo, err := checkAndRememberVertDiderVideo(db, videoId)
			if err != nil {
				log.Errorf("Failed to check a video against seen videos database: %s.", err)
			}

			allowPost = allowPost || allowVideo
		}

		return allowPost
	})

	return
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

		return true
	})

	return
}

func openSeenDb() (db *sql.DB, err error) {
	user, err := user.Current()
	if err != nil {
		return
	}

	db, err = sql.Open("sqlite3", path.Join(user.HomeDir, ".rsspipes.sqlite"))
	if err != nil {
		err = fmt.Errorf("Failed to open seen items database: %s.", err)
		return
	}

	defer func() {
		if err != nil {
			db.Close()
		}
	}()

	err = createTable(db, `
        CREATE TABLE IF NOT EXISTS quotes (
            id TEXT NOT NULL PRIMARY KEY,
            time TIMESTAMP NOT NULL,
            text TEXT NOT NULL
        )
    `)
	if err != nil {
		return
	}

	err = createTable(db, `
        CREATE TABLE IF NOT EXISTS vert_dider_videos (
            id INTEGER NOT NULL PRIMARY KEY,
            time TIMESTAMP NOT NULL
        )
    `)
	if err != nil {
		return
	}

	return
}

func createTable(db *sql.DB, query string) error {
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("Failed to create a table: %s.", err)
	} else {
		return nil
	}
}

var quoteFingerprintExtraCharsRe = regexp.MustCompile(`(\s+|<br>|[–().,:;!?'"-]+)`)

func getQuoteFingerprint(text string) string {
	fingerprint := quoteFingerprintExtraCharsRe.ReplaceAllString(text, "")
	fingerprint = strings.Replace(fingerprint, "ё", "e", -1)
	fingerprint = strings.ToLower(fingerprint)

	hasher := md5.New()
	hasher.Write([]byte(fingerprint))
	return hex.EncodeToString(hasher.Sum(nil))
}

func checkAndRememberQuote(db *sql.DB, text string) (bool, error) {
	fingerprint := getQuoteFingerprint(text)

	rows, err := db.Query("SELECT time FROM quotes WHERE id == ?", fingerprint)
	if err != nil {
		return true, err
	}
	defer rows.Close()

	if rows.Next() {
		return checkSeenItem(rows)
	} else if err := rows.Err(); err != nil {
		return true, err
	}

	if _, err := db.Exec(`INSERT INTO quotes (id, time, text) VALUES (?, ?, ?)`, fingerprint, time.Now(), text); err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
			// It's OK - we've just got a race
		} else {
			return true, err
		}
	}

	return true, nil
}

func checkAndRememberVertDiderVideo(db *sql.DB, videoId int64) (bool, error) {
	rows, err := db.Query("SELECT time FROM vert_dider_videos WHERE id == ?", videoId)
	if err != nil {
		return true, err
	}
	defer rows.Close()

	if rows.Next() {
		return checkSeenItem(rows)
	} else if err := rows.Err(); err != nil {
		return true, err
	}

	if _, err := db.Exec(`INSERT INTO vert_dider_videos (id, time) VALUES (?, ?)`, videoId, time.Now()); err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
			// It's OK - we've just got a race
		} else {
			return true, err
		}
	}

	return true, nil
}

func checkSeenItem(rows *sql.Rows) (bool, error) {
	var firstSeenTime time.Time
	if err := rows.Scan(&firstSeenTime); err != nil {
		return true, err
	}

	itemAgeDays := time.Now().Sub(firstSeenTime).Hours() / 24

	return itemAgeDays < 3, nil
}
