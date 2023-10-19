package propertyScraperv2

import (
	// "fmt"
	"log"

	requests "back-end/mapSearchService/src/requests"
	propertyTypes "back-end/mapSearchService/src/types"

	"github.com/gocolly/colly/v2"
	"github.com/google/uuid"
)


func ScrapeInParallel(pagesToScrape []string, collector *colly.Collector) *[]propertyTypes.Property {
	properties := []propertyTypes.Property{}

	collector.OnRequest(func(request *colly.Request) {
		// fmt.Println("Scraping: ", request.URL.String())
	})

	collector.OnError(func(_ *colly.Response, err error) {
		log.Println("Failed to scrape: ", err)
	})

	collector.OnHTML("#__NEXT_DATA__", func(element *colly.HTMLElement) {
		rawListingData := requests.FormatJSON([]byte(element.Text))
		formattedListingData := rawListingData["props"].(map[string]interface{})["pageProps"].(map[string]interface{})["componentProps"].(map[string]interface{})["listingsMap"]

		// fmt.Println(formattedListingData)
		for _, listing := range formattedListingData.(map[string]interface{}) {
			inspectionOpenTime := listing.(map[string]interface{})["listingModel"].(map[string]interface{})["inspection"].(map[string]interface{})["openTime"]

			if(inspectionOpenTime == nil) {
				continue
			}

			beds := listing.(map[string]interface{})["listingModel"].(map[string]interface{})["features"].(map[string]interface{})["beds"]
			if (beds == nil) {
				beds = 0.0
			}

			parking := listing.(map[string]interface{})["listingModel"].(map[string]interface{})["features"].(map[string]interface{})["parking"]
			if(parking == nil) {
				parking = 0.0
			}

			geocode := propertyTypes.Geocode{
				Lat: listing.(map[string]interface{})["listingModel"].(map[string]interface{})["address"].(map[string]interface{})["lat"].(float64),
				Lng: listing.(map[string]interface{})["listingModel"].(map[string]interface{})["address"].(map[string]interface{})["lng"].(float64),
			}

			property := propertyTypes.Property{
				WeeklyPrice: listing.(map[string]interface{})["listingModel"].(map[string]interface{})["price"].(string),
				InspectionOpenTime: inspectionOpenTime.(string),
				InspectionCloseTime: listing.(map[string]interface{})["listingModel"].(map[string]interface{})["inspection"].(map[string]interface{})["closeTime"].(string),
				ID: uuid.New(),
				Geocode: geocode,
				Location: propertyTypes.GeoJSON{
					Type: "Point",
					Coordinates: []float64{
						geocode.Lng, geocode.Lat, 
					},
				},
				Images: listing.(map[string]interface{})["listingModel"].(map[string]interface{})["images"].([]interface {}),
				Address: listing.(map[string]interface{})["listingModel"].(map[string]interface{})["address"].(map[string]interface{})["street"].(string),
				Suburb: listing.(map[string]interface{})["listingModel"].(map[string]interface{})["address"].(map[string]interface{})["suburb"].(string),
				State: listing.(map[string]interface{})["listingModel"].(map[string]interface{})["address"].(map[string]interface{})["state"].(string),
				Country: "Australia",
				Postcode: listing.(map[string]interface{})["listingModel"].(map[string]interface{})["address"].(map[string]interface{})["postcode"].(string),
				Beds: beds.(float64),
				Baths: listing.(map[string]interface{})["listingModel"].(map[string]interface{})["features"].(map[string]interface{})["baths"].(float64),
				Parking: parking.(float64),
				PropertyType: listing.(map[string]interface{})["listingModel"].(map[string]interface{})["features"].(map[string]interface{})["propertyTypeFormatted"].(string),
				Agency: listing.(map[string]interface{})["listingModel"].(map[string]interface{})["branding"].(map[string]interface{}),
			}

			properties = append(properties, property)

 			// fmt.Println("Scraped successfully")
		}
	})
		
	for _, pageToScrape := range pagesToScrape {
		collector.Visit(pageToScrape)
	}
	
	collector.Wait()

	return &properties
}
