package propertiesControllerV1

import (
	"fmt"
	"net/http"

	propertiesServiceV1 "back-end/mapSearchService/src/api/v1/service"
	propertyScraperBase "back-end/mapSearchService/src/scraper"
	propertyTypes "back-end/mapSearchService/src/types"

	"github.com/gin-gonic/gin"
) 


func GetScheduledListings(context *gin.Context) {
	reqParams := context.Request.URL.Query()
	scrapedSububurbsIds, suburbsToScrape := propertiesServiceV1.SeparateListingsInDb(reqParams["suburb"][0])
	properties := []propertyTypes.Property{}
	
	if(len(*suburbsToScrape) > 0) {
		fmt.Println("Suburbs to scrape: ", *suburbsToScrape)
		suburbListings := propertyScraperBase.InitialiseScraper(reqParams, *suburbsToScrape)
		properties = append(properties, *suburbListings...)

		propertiesServiceV1.InsertScrapedListingsIntoDb(&properties)
	}

	if(len(*scrapedSububurbsIds) > 0) {
		fmt.Println("Suburbs in database: ", scrapedSububurbsIds)
		listingsInDatabase := propertiesServiceV1.GetListingsInDb(scrapedSububurbsIds)
		properties = append(properties, *listingsInDatabase...) 
	}

	context.JSON(http.StatusOK, properties)
	properties = nil
}