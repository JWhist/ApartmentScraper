package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

var (
	searchString  string = "https://www.apartments.com/new-bern-nc/min-%s-bedrooms-over-%s/"
	defaultAmount int    = 10
	defaultOutput string = "homes.json"
)

type Home struct {
	Title string
	Price string
	Range string
	Link  string
}

func (q *Home) String() string {
	return fmt.Sprintf("Title: %s, Price: %s, Range: %s, Link: %s", q.Title, q.Price, q.Range, q.Link)
}

func usage() {
	fmt.Println("Usage: ./scraper #-bedrooms price [amount] [output]")
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 && len(os.Args) != 3 && len(os.Args) != 4 && len(os.Args) != 5 {
		usage()
	}

	bedrooms := os.Args[1]
	price := os.Args[2]
	page := 2
	var amount int = defaultAmount
	var output string = defaultOutput
	var err error

	if len(os.Args) >= 4 {
		amount, err = strconv.Atoi(os.Args[3])
		if err != nil {
			usage()
		}
	}

	if len(os.Args) >= 5 {
		output = os.Args[4]
	}

	var homes []Home

	c := colly.NewCollector(
		colly.UserAgent("Chrome/79"),
	)

	c.OnHTML("li.mortar-wrapper", func(e *colly.HTMLElement) {
		fmt.Println("Found record!")
		title := e.ChildText("span.js-placardTitle")
		price := e.ChildText("p.property-pricing")
		pRange := e.ChildText("div.price-range")
		link := e.ChildAttr("a.property-link", "href")

		if title == "" {
			return
		}

		if price == "" && pRange == "" {
			return
		}

		if price == "" {
			price = "N/A"
		} else if pRange == "" {
			pRange = "N/A"
		}

		homes = append(homes, Home{
			Title: strings.Trim(title, "\n"),
			Price: strings.Trim(price, "\n"),
			Range: strings.Trim(pRange, "\n"),
			Link:  strings.Trim(link, "\n"),
		})

		fmt.Println(title)
		fmt.Println(price)
		fmt.Println(link)
		fmt.Println(pRange + "\n")
	})

	// click next only if we don't have enough quotes
	pageTag := fmt.Sprintf(`[data-page=%d]`, page)
	c.OnHTML(pageTag, func(e *colly.HTMLElement) {
		if len(homes) < amount {
			e.Request.Visit(e.Attr("href"))
		}
		page++
	})

	fmt.Println("Launching Scraper !")
	link := fmt.Sprintf(searchString, bedrooms, price)
	fmt.Println("Visiting: " + link)
	c.Visit(link)

	fmt.Printf("Scraped %d homes.\n\n", len(homes))

	var homesString []string

	for _, home := range homes {
		homesString = append(homesString, home.String())
	}

	toWrite, err := json.MarshalIndent(homesString, "", "  ")
	if err != nil {
		fmt.Println("Can't marshall: " + err.Error())
		os.Exit(1)
	}

	err = os.WriteFile(output, toWrite, 0644)
	if err != nil {
		fmt.Println("Can't write to file: " + err.Error())
		os.Exit(1)
	}
}
