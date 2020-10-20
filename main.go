package main

import (
	"fmt"
	"log"
	"os"
	controller "scrappingMercadoLibre/controller"

	"github.com/gocolly/colly"
)

func main() {
	//Here set the name Output
	fName := "mercadoliberTelefonosCO.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()

	//Set here how many pages do you want Scrape
	var pageUntil = 5

	//Create a Collector
	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	//Ready to call the controller - > found values, next_pages , outpoutCSV
	controller.Create(c, file, pageUntil)

	//Place here the site you want to Scrape
	//Start scraping on https://celulares.mercadolibre.com.ar/telefonos
	c.Visit("https://celulares.mercadolibre.com.co/telefonos")

	log.Printf("Scraping finished, check file %q for results\n", fName)
}
