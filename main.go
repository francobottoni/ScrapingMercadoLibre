package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	controller "scrappingMercadoLibre/controller"

	"github.com/gocolly/colly/v2"
)

func main() {
	//Here set the name OutputCSV Scraping
	fName := "mercadolibreTelefonosCO.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)

	//Set Header
	writer.Write([]string{"ID", "Tienda", "Producto", "Precio", "Stock", "Garantia", "Ubicacion", "Cantidad de Ventas"})

	//Here set the name OutputCSV best stores
	fNameStores := "mercadolibreBestStores.csv"
	fileStores, err := os.Create(fNameStores)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fNameStores, err)
		return
	}
	defer fileStores.Close()
	writerStore := csv.NewWriter(fileStores)

	//Set Header
	writerStore.Write([]string{"Tienda", "Ventas"})

	//Set here how many pages do you want Scrape
	var pageUntil = 4

	//Create a Collector
	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	controller.Create(c, writer, writerStore) //Call the controller Create - > found values, save CSV
	controller.NextPage(c, pageUntil)         //Call next page until page we send

	//Place here the site you want to Scrape
	//Start scraping on https://celulares.mercadolibre.com.ar/telefonos
	c.Visit("https://celulares.mercadolibre.com.co/telefonos")

	log.Printf("Scraping finished, check file %q for results\n", fName)
}
