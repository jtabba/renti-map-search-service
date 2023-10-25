package propertyScraperBase

import (
	"fmt"
	"math/rand"
	"net/url"

	envHelper "back-end/mapSearchService/env"
	propertyScraperv2 "back-end/mapSearchService/src/scraper/v2"
	propertyTypes "back-end/mapSearchService/src/types"

	"github.com/gocolly/colly/v2"
)

var DATA_ACCESS_URL string = envHelper.GetEnvVar("DATA_ACCESS_URL")
var userAgents []string = []string{
	"Mozilla/5.0 (Linux; Android 13; SM-S908B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 13; SM-S901U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 13; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 13; Pixel 6a) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 13; Pixel 6a) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36",
}

func InitialiseScraper(filterOptions url.Values, suburbsToScrape []interface{}) *[]propertyTypes.Property {
	if(DATA_ACCESS_URL == "") {
		panic("Data access URL not provided")
	}

	fmt.Println("Creating scrape urls...")
	pagesToScrape := createScrapeUrls(DATA_ACCESS_URL, filterOptions, suburbsToScrape)
	
	fmt.Println("Initialising scraper...")
	collector := colly.NewCollector(colly.Async(true))
	collector.Limit(&colly.LimitRule{
		Parallelism: 4,
	})
	randomInt := rand.Intn(len(userAgents))
	collector.UserAgent = userAgents[randomInt]

	fmt.Println("Scraping...")
	properties := propertyScraperv2.ScrapeInParallel(pagesToScrape, collector)
	fmt.Printf("\nComplete! Scraped %v properties \n\n", len(*properties))

	return properties
}

func createScrapeUrls(dataAccessUrl string, filterOptions url.Values, suburbsToScrape []interface{}) []interface{} {
	for i, suburb := range(suburbsToScrape) {
		processedUrl := dataAccessUrl
		suburbDetails := suburb.([]interface{})[0].(string)

		if(filterOptions["type"] != nil) {
			processedUrl += filterOptions["type"][0] + "/"
		} else {
			panic("No type provided - type (rent or buy) is required for search")
		}

		if(filterOptions["suburb"] != nil) {
			processedUrl += suburbDetails + "/"
		} else {
			panic("No suburb provided - suburb is required for search")
		}

		// Don't include rented/sold properties AND only search for properties which match the exact suburb
		processedUrl += "?excludedeposittaken=1&ssubs=0"

		if(filterOptions["bedrooms"] != nil) {
			processedUrl += "&bedrooms=" + filterOptions["bedrooms"][0]
		}

		if(filterOptions["bathrooms"] != nil) {
			processedUrl += "&bathrooms=" + filterOptions["bathrooms"][0]
		}

		if(filterOptions["parking"] != nil) {
			processedUrl += "&parking=" + filterOptions["parking"][0]
		}

		suburbsToScrape[i].([]interface{})[0] = processedUrl
	}

	return suburbsToScrape
}