package CollysFunc

import (
	"encoding/csv"
	"fmt"
	"strconv"

	"github.com/gocolly/colly/v2"
)

type StoreInfo struct {
	Id               string
	Store            string
	Producto         string
	Precio           string
	Stock            string
	Garantia         string
	Ubicacion        string
	CantidadDeVentas string
}

type BestSeller struct {
	Store            string
	CantidadDeVentas string
}

//Set default values
var setHeader = 0
var idProduct = 0
var countPage = 1

var bestStores [10]BestSeller
var isHaveStore bool = true
var isEmpty int = -1
var isMenorPos int = 23         //Set random number
var cantVentasAux int = 1000000 //Set default number

func Create(c *colly.Collector, writer *csv.Writer, writerStore *csv.Writer) {

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
			setHeader = setHeader + 1
		}

		//Call the function that returns Store name
		store := FoundStoreName(e)

		//Set Store struct
		infoStore := StoreInfo{
			Id:               strconv.Itoa(idProduct), //The array writer.Write only receives string,then convert id to string
			Store:            store,
			Producto:         producto,
			Precio:           precio,
			Stock:            stock,
			Garantia:         garantia,
			Ubicacion:        ubicacion,
			CantidadDeVentas: cantidadDeVentas,
		}

		//Call function BestSeller
		GenerateBestSellers(infoStore, writerStore)

		//Write the file with obtains values
		writer.Write([]string{
			infoStore.Id,
			infoStore.Store,
			infoStore.Producto,
			infoStore.Precio,
			infoStore.Stock,
			infoStore.Garantia,
			infoStore.Ubicacion,
			infoStore.CantidadDeVentas,
		})
		writer.Flush()
	})
}

func FoundStoreName(e *colly.HTMLElement) string {

	//Find the seller's profile link and save it to a string
	link := e.ChildAttr("div.layout-description-wrapper > section.ui-view-more.vip-section-seller-info.new-reputation > a", "href")

	//If link is a String and we know the link for all profiles is : https://perfil.mercadolibre.com.co/NAME-PROFAILE
	//Then we return a string that is assembled from character 35 onwards
	return string(link[35:])

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

func GenerateBestSellers(values StoreInfo, writerStore *csv.Writer) {

	//If the store has already exist in the Array Aux set isHaveStore false
	for i := 0; i < len(bestStores); i++ {
		if bestStores[i].CantidadDeVentas == values.CantidadDeVentas && bestStores[i].Store == values.Store {
			isHaveStore = false
		}
	}

	if isHaveStore == true {

		for i := 0; i < len(bestStores); i++ {

			cantVentas, _ := strconv.Atoi(bestStores[i].CantidadDeVentas)

			//If array bestStores is empty, set this pos in isEmpty
			if bestStores[i].Store == "" {
				isEmpty = i
			} else {
				for j := 1; j < len(bestStores); j++ {
					if bestStores[j].CantidadDeVentas != "" {
						cantVentas1, _ := strconv.Atoi(bestStores[j].CantidadDeVentas)
						fmt.Println(cantVentas)
						fmt.Println(cantVentas1)
						if cantVentas < cantVentas1 {
							if cantVentas < cantVentasAux {
								//Set the lowest value and then replace it
								isMenorPos = i
								cantVentasAux = cantVentas
							}
						}
					}
				}
			}
		}

		if isMenorPos != 23 {
			s, _ := strconv.Atoi(bestStores[isMenorPos].CantidadDeVentas)
			s1, _ := strconv.Atoi(values.CantidadDeVentas)
			//Replace the lower value with a new higher value
			if s < s1 {
				bestStores[isMenorPos] = BestSeller{
					Store:            values.Store,
					CantidadDeVentas: values.CantidadDeVentas,
				}
			}
		}

		//If array is Empty, load a value
		if isEmpty != -1 {
			bestStores[isEmpty] = BestSeller{
				Store:            values.Store,
				CantidadDeVentas: values.CantidadDeVentas,
			}
			isEmpty = -1
		}

	}

	fmt.Println(bestStores)
	//we set isHaveStore back to true
	isHaveStore = true
}
