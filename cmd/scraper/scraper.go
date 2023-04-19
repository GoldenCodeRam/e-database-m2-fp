package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ITEM_LIST_URL        = "listado.mercadolibre.com.co"
	ITEM_DESCRIPTION_URL = "www.mercadolibre.com.co"
)

var contentToScrap = []Content{
	{
		Type:  "Laptops",
		Url:   "https://listado.mercadolibre.com.co/portátiles",
		Items: &[]Item{},
	},
	{
		Type:  "Cameras",
		Url:   "https://listado.mercadolibre.com.co/cámaras",
		Items: &[]Item{},
	},
	{
		Type:  "Phones",
		Url:   "https://listado.mercadolibre.com.co/teléfonos",
		Items: &[]Item{},
	},
	{
		Type:  "Televisions",
		Url:   "https://listado.mercadolibre.com.co/televisores",
		Items: &[]Item{},
	},
	{
		Type:  "Earphones",
		Url:   "https://listado.mercadolibre.com.co/audífonos",
		Items: &[]Item{},
	},
}

type Content struct {
	Type  string
	Url   string
	Items *[]Item
}

type Item struct {
	Name        string   `json:"name"`
	Price       int      `json:"price"`
	Images      []string `json:"images"`
	Reviews     int      `json:"reviews"`
    ReviewScore float32  `json:"reviewScore"`
	Description string   `json:"description"`
	Url         string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URL")).SetServerAPIOptions(serverApi)

	client, err := mongo.Connect(context.TODO(), opts)

    coll := client.Database("products").Collection("products")

	log.Println("Starting scraper...")
	for _, content := range contentToScrap {
		itemListQuery := fmt.Sprintf(`
        LET doc = DOCUMENT("%s")
        LET nextPage = "li.andes-pagination__button.andes-pagination__button--next.shops__pagination-button a"
        LET itemClass = "a.ui-search-item__group__element.shops__items-group-details.ui-search-link"
        LET pages = 2

        LET result = (
            FOR page IN 1..pages
            LET clicked = page == 1 ? false : CLICK(doc, nextPage)
            LET wait = clicked ? WAIT_NAVIGATION(doc, 10000) : false

            LET items = (
                FOR el IN ELEMENTS(doc, itemClass)
                    RETURN TRIM(el.attributes.href)
            )

            RETURN items
        )

        RETURN FLATTEN(result)
        `, content.Url)

		comp := compiler.New()
		program, _ := comp.Compile(itemListQuery)

		ctx := context.Background()

		ctx = drivers.WithContext(ctx, cdp.NewDriver(), drivers.AsDefault())

		out, err := program.Run(ctx)
		if err != nil {
			panic(err)
		}

		itemUrls := []string{}

		json.Unmarshal(out, &itemUrls)

		log.Printf("Items found: %d for category: %s", len(itemUrls), content.Type)
		log.Println("Starting item scraping...")
		for i, url := range itemUrls {

			log.Printf("Scraping %d from %d...", i+1, len(itemUrls))

			itemQuery := fmt.Sprintf(`
            LET doc = DOCUMENT("%s")

            LET mainSelector = "div.ui-pdp-price__second-line span meta"
            LET nameSelector = "h1.ui-pdp-title"
            LET priceSelector = "div.ui-pdp-price__second-line span meta"
            LET imagesSelector = "figure.ui-pdp-gallery__figure img"
            LET reviewsSelector = "span.ui-pdp-review__amount"
            LET descriptionSelector = "p.ui-pdp-description__content"
            LET reviewScoreSelector = "p.ui-review-capability__rating__average"

            WAIT_ELEMENT(doc, mainSelector)

            LET images = (
                FOR image IN ELEMENTS(doc, imagesSelector)
                RETURN image.attributes.src
            )

            RETURN {
                name: ELEMENT(doc, nameSelector).innerText,
                price: TO_INT(ELEMENT(doc, priceSelector).attributes.content),
                images: images,
                reviews: TO_INT(REGEX_MATCH(ELEMENT(doc, reviewsSelector).innerText, "[0-9]+")),
                reviewScore: TO_FLOAT(ELEMENT(doc, reviewScoreSelector).innerText),
                description: ELEMENT(doc, descriptionSelector).innerText,
            }
            `, url)

			prg, _ := comp.Compile(itemQuery)

			o, _ := prg.Run(ctx)

			item := Item{
				Url: url,
			}
			json.Unmarshal(o, &item)

            coll.InsertOne(context.Background(), item)
		}
	}
}
