all: scraper

scraper: main.go
	go build -o scraper main.go