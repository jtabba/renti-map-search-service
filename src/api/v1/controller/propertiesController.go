package propertiesControllerV1

import (
	"net/http"

	propertiesServiceV1 "back-end/mapSearchService/src/api/v1/service"
	propertyTypes "back-end/mapSearchService/src/types"

	"github.com/gin-gonic/gin"
)


func GetMultipleProperties(context *gin.Context) {
	// thoughts:
		// check DB for suburb DATABASE (?)
		// if there have a cool down period where data is not re-scraped
		// if not or cool down period has passed, scrape data and store/update in DB
			// updated ata can be optimised by only updating the properties that have changed
		// ONLY query geocode API if property does not already exist in DB
	reqParams := context.Request.URL.Query()
	var properties *[]propertyTypes.Property = propertiesServiceV1.GetFilteredProperties(reqParams)

	context.JSON(http.StatusOK, properties)
	properties = nil
}