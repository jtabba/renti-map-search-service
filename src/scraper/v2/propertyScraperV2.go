package propertyScraperv2

import (
	"fmt"
	"log"
	"strings"

	"back-end/mapSearchService/src/api/v1/middleware"
	propertyTypes "back-end/mapSearchService/src/types"
	utilities "back-end/mapSearchService/src/utilities"

	"github.com/gocolly/colly/v2"
)


func ScrapeInParallel(pagesToScrape []interface{}, collector *colly.Collector) *[]propertyTypes.Property {
	properties := []propertyTypes.Property{}
	var suburbCacheData middleware.ScrapedSuburb

	collector.OnError(func(_ *colly.Response, err error) {
		log.Println("Failed to scrape: ", err)
	})

	collector.OnHTML("#__NEXT_DATA__", func(element *colly.HTMLElement) {
		rawListingData := utilities.FormatJSON([]byte(element.Text))
		formattedListingData := rawListingData["props"].(map[string]interface{})["pageProps"].(map[string]interface{})["componentProps"].(map[string]interface{})["listingsMap"]

		for _, listing := range formattedListingData.(map[string]interface{}) {
			inspectionOpenTime := listing.(map[string]interface{})["listingModel"].(map[string]interface{})["inspection"].(map[string]interface{})["openTime"]

			if(inspectionOpenTime == nil) {
				continue
			}

			listingData := listing.(map[string]interface{})["listingModel"].(map[string]interface{})

			beds := listingData["features"].(map[string]interface{})["beds"]; if beds == nil {
				beds = 0.0
			}

			parking := listingData["features"].(map[string]interface{})["parking"]; if parking == nil {
				parking = 0.0
			}

			geocode := map[string]interface{}{
				"lat": listingData["address"].(map[string]interface{})["lat"].(float64),
				"lng": listingData["address"].(map[string]interface{})["lng"].(float64),
			}

			imagesSlice := make([]string, len(listingData["images"].([]interface {})))
			for i, v := range listingData["images"].([]interface {}) {
				imagesSlice[i] = v.(string)
			}

			suburb := strings.Replace(listingData["address"].(map[string]interface{})["suburb"].(string), " ", "-", -1)
			postcode := listingData["address"].(map[string]interface{})["postcode"].(string)
			state := listingData["address"].(map[string]interface{})["state"].(string)

			property := propertyTypes.Property{
				Suburb: suburb,
				Postcode: postcode,
				State: state,
				Geocode: geocode,
				Images: imagesSlice,
				Beds: beds.(float64),
				Parking: parking.(float64),
				Country: "Australia",
				WeeklyPrice: listingData["price"].(string),
				InspectionOpenTime: inspectionOpenTime.(string),
				Agency: listingData["branding"].(map[string]interface{}),
				Geolocation: fmt.Sprintf("POINT(%v %v)", geocode["lng"], geocode["lat"]),
				Address: listingData["address"].(map[string]interface{})["street"].(string),
				Baths: listingData["features"].(map[string]interface{})["baths"].(float64),
				InspectionCloseTime: listingData["inspection"].(map[string]interface{})["closeTime"].(string),
				PropertyType: listingData["features"].(map[string]interface{})["propertyTypeFormatted"].(string),
			}

			properties = append(properties, property)
			cacheId := fmt.Sprintf("%s-%s-%s", suburb, state, postcode)
			middleware.SetCache(cacheId, suburbCacheData)
		}
	})
		
	for _, pageToScrape := range pagesToScrape {
		suburbCacheData = pageToScrape.([]interface{})[1].(middleware.ScrapedSuburb) 
		collector.Visit(pageToScrape.([]interface{})[0].(string))
	}
	
	collector.Wait()

	return &properties
}
