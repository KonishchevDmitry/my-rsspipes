module github.com/KonishchevDmitry/my-rsspipes

go 1.13

replace github.com/KonishchevDmitry/go-rss => ./go-rss

replace github.com/KonishchevDmitry/rsspipes => ./rsspipes

require (
	github.com/KonishchevDmitry/go-rss v0.0.0-20170617091005-cf1fe9a72c9e
	github.com/KonishchevDmitry/rsspipes v0.0.0-00010101000000-000000000000
	github.com/PuerkitoBio/goquery v1.5.0
	github.com/mattn/go-sqlite3 v1.2.1-0.20170407154627-cf7286f069c3
	golang.org/x/net v0.0.0-20191207000613-e7e4b65ae663
)
