package main

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

const (
	ITEM_LIST_URL        = "listado.mercadolibre.com.co"
	ITEM_DESCRIPTION_URL = "www.mercadolibre.com.co"
)

const (
	ITEM_A_CLASS = "a.ui-search-item__group__element.shops__items-group-details.ui-search-link"
	ITEM_TITLE   = "h1.ui-pdp-title"
)

var links = []string{
	"https://listado.mercadolibre.com.co/portatil",
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains(
			ITEM_LIST_URL,
			ITEM_DESCRIPTION_URL,
		),
	)

	// parse item list first
	c.OnResponse(func(r *colly.Response) {
		switch r.Request.URL.Hostname() {
		case ITEM_LIST_URL:
			visitItems(c, r.Body)
		case ITEM_DESCRIPTION_URL:
			scrapItem(c, r.Body)
		}
	})

	// startpoint
	c.Visit("https://listado.mercadolibre.com.co/portatil")
}

func visitItems(c *colly.Collector, body []byte) {
	queryBody, err := newDocumentFromBody(body)

	if err != nil {
		panic(err)
	}

	queryBody.Find(ITEM_A_CLASS).Each(func(i int, s *goquery.Selection) {
		href, found := s.Attr("href")

		if found {
			//log.Printf("Found item: %s", href)
			c.Visit(href)
		}
	})
}

func scrapItem(c *colly.Collector, body []byte) {
	queryBody, err := newDocumentFromBody(body)

	if err != nil {
		panic(err)
	}

	title := queryBody.Find("h1.ui-pdp-title").Text()

	images := 0
	queryBody.Find("div.ui-pdp-gallery__column").Find("img").Each(func(i int, s *goquery.Selection) {
		images += 1
	})

	log.Println(title, images)
}

func newDocumentFromBody(body []byte) (*goquery.Document, error) {
	return goquery.NewDocumentFromReader(strings.NewReader(string(body)))
}
