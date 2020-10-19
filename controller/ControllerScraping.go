package CollysFunc

import (
	"encoding/csv"
	"os"

	"github.com/gocolly/colly"
)

type StoreInfo struct {
	Producto         string
	Precio           string
	Stock            string
	Garantia         string
	Ubicacion        string
	CantidadDeVentas string
}

type StoreInfoCompleted struct {
	Store     string
	StoreInfo StoreInfo
}

//Instance value for set the header in lane 72
var setHeader = 0
var page = 1

func Create(c *colly.Collector, file *os.File) {
	// //Move to link publication
	c.OnHTML("div.ui-search-item__group.ui-search-item__group--title a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		//fmt.Printf("Link found: -> %s\n", link)

		c.Visit(e.Request.AbsoluteURL(link))
	})

	//Found values for the store
	c.OnHTML("#root-app > div > div.layout-main.u-clearfix > div.layout-col.layout-col--right", func(e *colly.HTMLElement) {

		//Set values
		producto := e.ChildText("#short-desc > div > header > h1")
		precio := e.ChildText("#productInfo > fieldset.item-price > span > span.price-tag-fraction")
		stock := e.ChildText("#dropdown-quantity > button > span.dropdown-quantity-available")
		garantia := e.ChildText("#root-app > div > div.layout-main.u-clearfix > div.layout-col.layout-col--right > div.layout-description-wrapper > section.ui-view-more.vip-section-warranty.vip-section-service > div:nth-child(3) > p.text-light.warranty__store")
		ubicacion := e.ChildText("div.card-section.seller-location > p.card-description.text-light")
		cantidadDeVentas := e.ChildText("#root-app div.layout-description-wrapper > section.ui-view-more.vip-section-seller-info.new-reputation > div.reputation-info.block > dl > dd:nth-child(1) > strong")

		//Set Store struct
		infoStore := StoreInfo{
			Producto:         producto,
			Precio:           precio,
			Stock:            stock,
			Garantia:         garantia,
			Ubicacion:        ubicacion,
			CantidadDeVentas: cantidadDeVentas,
		}

		//Write the file with the obtains values
		writer := csv.NewWriter(file)
		defer writer.Flush()

		// Set header if don't have
		if setHeader == 0 {
			// Write CSV header
			writer.Write([]string{"Producto", "Precio", "Stock", "Garantia", "Ubicacion", "Cantidad de Ventas"})
			setHeader = setHeader + 1
		}

		// Write CSV header
		writer.Write([]string{
			infoStore.Producto,
			infoStore.Precio,
			infoStore.Stock,
			infoStore.Garantia,
			infoStore.Ubicacion,
			infoStore.CantidadDeVentas,
		})

	})

	//NEXT PAGE
	c.OnHTML("#root-app > div > div > section > div.ui-search-pagination > ul > li.andes-pagination__button.andes-pagination__button--next > a[href]", func(e *colly.HTMLElement) {

		if page < 7 {
			link := e.Attr("href")
			// Visit link found on page
			page = page + 1
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

}

// func MoveToLinkReputation(c *colly.Collector, s *StoreInfo) {
// 	//Move to link reputacion in the publiation
// 	c.OnHTML("#root-app > div > div.layout-main.u-clearfix > div.layout-col.layout-col--right > div.layout-description-wrapper > section.ui-view-more.vip-section-seller-info.new-reputation a[href]", func(e *colly.HTMLElement) {
// 		link := e.Attr("href")

// 		//fmt.Printf("Link found: -> %s\n", link)
// 		c.Visit(e.Request.AbsoluteURL(link))
// 	})

// 	//Found name Store
// 	c.OnHTML("div.store-info-title", func(e *colly.HTMLElement) {

// 		fmt.Println("Store: " + e.ChildText("#store-info__name"))

// 	})
// }
