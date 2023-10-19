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
var pagesToScrapeCount int = 10 // can be altered to be dynamic based on zoom later on
var userAgents []string = []string{
	"Mozilla/5.0 (Linux; Android 13; SM-S908B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 13; SM-S901U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 13; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 13; Pixel 6a) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 13; Pixel 6a) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36",
}


func InitialiseScraper(filterOptions url.Values) *[]propertyTypes.Property {
	if(DATA_ACCESS_URL == "") {
		panic("Data access URL not provided")
	}

	fmt.Println("Creating scrape urls...")
	scrapeUrl := createScrapeUrl(DATA_ACCESS_URL, filterOptions)
	pagesToScrape := createScraperQueue(scrapeUrl)
	
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

func createScraperQueue(url string) []string {
	pagesToScrape := []string{}

	for i := 0; i < pagesToScrapeCount; i++ {
		pagesToScrape = append(pagesToScrape, fmt.Sprintf(`%s&page=%v`, url, i + 1))
	}

	return pagesToScrape
}

func createScrapeUrl(dataAccessUrl string, filterOptions url.Values) string {
	if(filterOptions["type"] != nil) {
		dataAccessUrl += filterOptions["type"][0] + "/"
	} else {
		panic("No type provided - type (rent or buy) is required for search")
	}

	if(filterOptions["suburb"] != nil) {
		dataAccessUrl += filterOptions["suburb"][0] + "/"
	} else {
		panic("No suburb provided - suburb is required for search")
	}

	dataAccessUrl += "?excludedeposittaken=1"

	if(filterOptions["bedrooms"] != nil) {
		dataAccessUrl += "&bedrooms=" + filterOptions["bedrooms"][0]
	}

	if(filterOptions["bathrooms"] != nil) {
		dataAccessUrl += "&bathrooms=" + filterOptions["bathrooms"][0]
	}

	if(filterOptions["parking"] != nil) {
		dataAccessUrl += "&parking=" + filterOptions["parking"][0]
	}

	return dataAccessUrl
}