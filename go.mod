module github.com/KonishchevDmitry/my-rsspipes

go 1.13

replace github.com/KonishchevDmitry/go-rss => ./go-rss

replace github.com/KonishchevDmitry/rsspipes => ./rsspipes

require (
	github.com/KonishchevDmitry/go-rss v0.0.0-20191207114205-40a828964875
	github.com/KonishchevDmitry/rsspipes v0.0.0-20180804154338-fea73d718c32
	github.com/PuerkitoBio/goquery v1.5.0
	github.com/andybalholm/cascadia v1.1.0 // indirect
	github.com/mattn/go-sqlite3 v2.0.1+incompatible
	golang.org/x/net v0.0.0-20191207000613-e7e4b65ae663
	golang.org/x/text v0.3.2 // indirect
)
