package propertiesServiceV1

import (
	dbClient "back-end/mapSearchService/src/database/v2-Psql"
	propertyScraperBase "back-end/mapSearchService/src/scraper"
	propertyTypes "back-end/mapSearchService/src/types"
	"fmt"

	// envHelper "back-end/mapSearchService/env" OS.GETENV????
	// "fmt"
	"net/url"
	// _ "github.com/lib/pq"
)

// check database stale time
// if stale scrape all new data
// else check
// IDEA - set geocode on the suburb entry then query nearby areas by proximity

func GetFilteredProperties(filterOptions url.Values) *[]propertyTypes.Property {
	db := dbClient.Connect()
	
	insertTest := `INSERT INTO SuburbsTest ("name") VALUES ($1)`
	// geo := propertyTypes.GeoJSON{ Type: "Point", Coordinates: []float64{-32.443345, 31.323455}}
	// geoCode := propertyTypes.Geocode{ Lat: -32.443345, Lng: 31.323455 }

	_, e := db.Exec(insertTest, "BELMORE"); if e != nil {
		fmt.Println(e)
	}

	fmt.Println("Added 1 record")

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



