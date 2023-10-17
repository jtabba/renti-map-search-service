package propertiesServiceV1

import (
	propertyScraperBase "back-end/mapSearchService/src/scraper"
	propertyTypes "back-end/mapSearchService/src/types"

	// envHelper "back-end/mapSearchService/env"
	// "fmt"
	"net/url"
)



func GetFilteredProperties(filterOptions url.Values) *[]propertyTypes.Property {
	// db := dbClient.Connect()
	// dbName := envHelper.GetEnvVar("PROP_DB")
	// collectionName := envHelper.GetEnvVar("PROP_COLLECTION")
	// collection := db.Database(dbName).Collection(collectionName)
	// latitude, err := strconv.ParseFloat(filterOptions["lat"][0], 8)
	// longitude, err := strconv.ParseFloat(filterOptions["lng"][0], 5)
	// // bedrooms, err := strconv.Atoi(filterOptions["bedrooms"][0])

	// if(err != nil) {
	// 	log.Println(err)
	// }

	// testCoord := bson.D{{Key: "type", Value: "Point"}, {Key: "coordinates", Value: []float64{longitude, latitude}}}
	// testFilter := bson.D{
	// 	// {Key: "bedrooms", Value: bedrooms},
	// 	{Key: "location", Value: bson.D{
	// 		{Key: "$near", Value: bson.D{
	// 			{Key: "$geometry", Value: testCoord},
	// 			{Key: "$minDistance", Value: 0},
	// 			{Key: "$maxDistance", Value: 2500},
	// 		}},
	// 	}},
	// }

	// cursor, err := collection.Find(dbClient.Ctx, testFilter)

	// if err != nil {
	// 	log.Println(err)
	// }

	// var properties []bson.M

	// err = cursor.All(dbClient.Ctx, &properties); 
	
	// if (err != nil) {
	// 	log.Println(err)
	// }

	// defer dbClient.Disconnect(db)

	properties := propertyScraperBase.InitialiseScraper(filterOptions)

	return properties
}



