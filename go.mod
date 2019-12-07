module github.com/KonishchevDmitry/my-rsspipes

go 1.13

replace github.com/KonishchevDmitry/go-rss => ./go-rss

replace github.com/KonishchevDmitry/rsspipes => ./rsspipes

require (
	github.com/KonishchevDmitry/go-rss v0.0.0-00010101000000-000000000000
	github.com/KonishchevDmitry/rsspipes v0.0.0-00010101000000-000000000000
	github.com/PuerkitoBio/goquery v1.1.1-0.20170324135448-ed7d758e9a34
	github.com/andybalholm/cascadia v0.0.0-20161224141413-349dd0209470 // indirect
	github.com/mattn/go-sqlite3 v1.2.1-0.20170407154627-cf7286f069c3
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7 // indirect
	golang.org/x/net v0.0.0-20170503120255-feeb485667d1
	golang.org/x/text v0.0.0-20170427093521-470f45bf29f4 // indirect
)
