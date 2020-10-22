package CollysFunc

import (
	"encoding/csv"
	"fmt"
	"strconv"

	"github.com/gocolly/colly"
)

type StoreInfo struct {
	Id               string
	Producto         string
	Precio           string
	Stock            string
	Garantia         string
	Ubicacion        string
	CantidadDeVentas string
}

//Set default values
var setHeader = 0
var idProduct = 0
var countPage = 1

func Create(c *colly.Collector, writer *csv.Writer) {

	// //Move to link publication
	c.OnHTML("div.ui-search-item__group.ui-search-item__group--title a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		//fmt.Printf("Link found: -> %s\n", link)
		c.Visit(e.Request.AbsoluteURL(link))

	})

	//InsertStoreName(c)

	//Found values for the store
	c.OnHTML("#root-app > div > div.layout-main.u-clearfix > div.layout-col.layout-col--right", func(e *colly.HTMLElement) {

		//Add product id
		idProduct = idProduct + 1

		//Set values
		producto := e.ChildText("#short-desc > div > header > h1")
		precio := e.ChildText("#productInfo > fieldset.item-price > span > span.price-tag-fraction")
		stock := e.ChildText("#dropdown-quantity > button > span.dropdown-quantity-available")
		garantia := e.ChildText("#root-app > div > div.layout-main.u-clearfix > div.layout-col.layout-col--right > div.layout-description-wrapper > section.ui-view-more.vip-section-warranty.vip-section-service > div:nth-child(3) > p.text-light.warranty__store")
		ubicacion := e.ChildText("div.card-section.seller-location > p.card-description.text-light")
		cantidadDeVentas := e.ChildText("#root-app div.layout-description-wrapper > section.ui-view-more.vip-section-seller-info.new-reputation > div.reputation-info.block > dl > dd:nth-child(1) > strong")

		// Set header if don't have
		if setHeader == 0 {
			// Write CSV header
			writer.Write([]string{"ID", "Tienda", "Producto", "Precio", "Stock", "Garantia", "Ubicacion", "Cantidad de Ventas"})
			setHeader = setHeader + 1
		}

		//Set Store struct
		infoStore := StoreInfo{
			Id:               strconv.Itoa(idProduct), //The array writer.Write only receives string,then convert id to string
			Producto:         producto,
			Precio:           precio,
			Stock:            stock,
			Garantia:         garantia,
			Ubicacion:        ubicacion,
			CantidadDeVentas: cantidadDeVentas,
		}

		//Write the file with obtains values
		writer.Write([]string{
			infoStore.Id,
			"",
			infoStore.Producto,
			infoStore.Precio,
			infoStore.Stock,
			infoStore.Garantia,
			infoStore.Ubicacion,
			infoStore.CantidadDeVentas,
		})
	})
}

func InsertStoreName(c *colly.Collector, writer *csv.Writer) {
	c.OnHTML("div.layout-description-wrapper > section.ui-view-more.vip-section-seller-info.new-reputation  a[href]", func(e *colly.HTMLElement) {

		link := e.Attr("href")

		//fmt.Printf("Link found: -> %s\n", link)

		//If link is a String and we know the link for all profiles is : https://perfil.mercadolibre.com.co/NAME-PROFAILE
		//This URL have 35 words until name profile
		writer.Write([]string{
			string(link[35:]),
		})
	})
}

func NextPage(c *colly.Collector, pageUntil int) {
	//NEXT PAGE
	c.OnHTML("#root-app > div > div > section > div.ui-search-pagination > ul > li.andes-pagination__button.andes-pagination__button--next > a[href]", func(e *colly.HTMLElement) {
		if countPage < pageUntil {
			link := e.Attr("href")
			// Visit link found on page
			countPage = countPage + 1
			fmt.Println("Siguiente pag: " + strconv.Itoa(countPage))
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})
}
