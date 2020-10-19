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

	//Create a Collector
	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	//Ready to call the controller - > found values, next_pages , outpoutCSV
	controller.Create(c, file)

	//Place here the site you want to scrape
	//Start scraping on https://celulares.mercadolibre.com.ar/telefonos
	c.Visit("https://celulares.mercadolibre.com.co/telefonos")

	log.Printf("Scraping finished, check file %q for results\n", fName)
}
