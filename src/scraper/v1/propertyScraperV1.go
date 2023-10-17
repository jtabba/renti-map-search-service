package propertyScraperV1

import (
	requests "back-end/mapSearchService/src/requests"
	"fmt"
	"log"
	"strconv"
	"strings"

	propertyTypes "back-end/mapSearchService/src/types"

	"github.com/gocolly/colly/v2"
	"github.com/google/uuid"
)

func ScrapeInParallel(url string, pagesToScrape []string, collector *colly.Collector) []propertyTypes.Property {
	properties := []propertyTypes.Property{}

	collector.OnError(func(_ *colly.Response, err error) {
		log.Println("Error: ", err)
	})

	collector.OnHTML(".css-8tedj6 li", func(element *colly.HTMLElement) {
		if(element.ChildText(".css-hwihpw") == "" || element.ChildText(".css-iqrvhs") == "") {
			return
		}

		address := element.ChildText(".css-iqrvhs")
		postcode := strings.Split(address, " ")[len(strings.Split(address, " ")) - 1]
		propertyGeocodeRes := requests.GetPropertyGeocode(address)

		if res, ok := propertyGeocodeRes["results"].([]interface{})[0].(map[string]interface{})["locations"].([]interface{})[0].(map[string]interface{}); ok {
			geocode := propertyTypes.Geocode{
				Lat: res["latLng"].(map[string]interface{})["lat"].(float64),
				Lng: res["latLng"].(map[string]interface{})["lng"].(float64),
			}
			// images := []string{}
			// images := element.ChildAttr("div", ".css-6yavch").ChildAttrs(".css-1ie6g1l", "style") 
			
			utilitiesArr := []string{}
			utilities := propertyTypes.Utilities{}
			
			element.ForEach(".css-6yavch", func(_ int, image *colly.HTMLElement) {
				fmt.Println("image: ", image)
				// images = image.ChildAttrs("listing-card-lazy-image img", "src")
			})

			element.ForEach(".css-1ie6g1l", func(_ int, div *colly.HTMLElement) {
				utilitiesArr = append(utilitiesArr, div.ChildText(".css-lvv8is"))
			})

			for i, util := range utilitiesArr {
				utilArr := strings.Split(util, " ")
				
				if val, err := strconv.ParseUint(utilArr[0], 10, 32); err == nil {
					if(i == 0) {
						utilities.Beds = uint(val)
					} else if (i == 1) {
						utilities.Baths = uint(val)
					} else if (i == 2) {
						utilities.Parking = uint(val)
					}
				} else {
					if(i == 2 ) {
						utilities.Parking = 0
					}
					// fmt.Println(util, err)
				}
 			}

			property := propertyTypes.Property{
				// FormattedAddress: address,
				// WeeklyPrice: element.ChildText(".css-mgq8yx"),
				// InspectionTime: element.ChildText(".css-hwihpw"),
				ID: uuid.New(),
				Geocode: geocode,
				Location: propertyTypes.GeoJSON{
					Type: "Point",
					Coordinates: []float64{
						geocode.Lng, geocode.Lat, 
					},
				},
				// Images: images,
				// City: res["adminArea5"].(string),
				State: res["adminArea3"].(string),
				Country: res["adminArea1"].(string),
				Postcode: postcode,
				// Beds: utilities.Beds,
				// Baths: utilities.Baths,
				// Parking: utilities.Parking,
				PropertyType: element.ChildText(".css-693528"),
				// Utilities: utilities,
			}

			properties = append(properties, property)
		}
	})
		
	for _, pageToScrape := range pagesToScrape {
		collector.Visit(pageToScrape)
	}
	
	collector.Wait()

	return properties
}
