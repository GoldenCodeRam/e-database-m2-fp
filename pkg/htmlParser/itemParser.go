package htmlparser

import (
	"log"

	"github.com/PuerkitoBio/goquery"
)

const (
    STAR_ITEM_CLASS_FULL = "ui-pdp-icon ui-pdp-icon--star-full"
    STAR_ITEM_CLASS_HALF = "ui-pdp-icon ui-pdp-icon--star-half"
)

func ParseItemStars(s *goquery.Selection) float64 {
    stars := 0.0

    log.Println(s.Html())

    s.Find("svg").Each(func(i int, s *goquery.Selection) {
        starClass, _ := s.Attr("class")
        switch starClass {
        case STAR_ITEM_CLASS_FULL:
            stars += 1
        case STAR_ITEM_CLASS_HALF:
            stars += 0.5
        }
    })
    return stars
}
