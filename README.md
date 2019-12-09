# colly-postgres-storage

A PostgreSQL storage back end for the Colly web crawling/scraping framework https://go-colly.org

Example Usage:

```go
package main

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/zolamk/colly-postgres-storage/colly/postgres"
)

func main() {

	c := colly.NewCollector()

	storage := &postgres.Storage{
        URI:      "postgres://username:password@localhost:5432/database",
        VisitedTable: "colly_visited",
        CookiesTable: "colly_cookies",
	}

	if err := c.SetStorage(storage); err != nil {
		panic(err)
	}

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("http://go-colly.org/")
}

```
