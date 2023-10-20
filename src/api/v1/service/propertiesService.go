package propertiesServiceV1

import (
	"back-end/mapSearchService/src/api/v1/middleware"
	dbClient "back-end/mapSearchService/src/database/v2-Psql"
	propertyScraperBase "back-end/mapSearchService/src/scraper"
	propertyTypes "back-end/mapSearchService/src/types"
	"encoding/json"
	"fmt"

	"strings"

	"net/url"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// check database stale time
// if stale scrape all new data
// else check
// IDEA - set geocode on the suburb entry then query nearby areas by proximity
type Suburb struct {
	Id 			int 		`json:"id"`
	Name 		string 		`json:"name"`
	State 		string 		`json:"state"`
	Postcode 	int 		`json:"postcode"`
	Geolocation string 		`json:"geolocation"`
	Geocode 	string 		`json:"geocode"`
}

func GetFilteredProperties(filterOptions url.Values) *[]propertyTypes.Property {
	searchSuburbSlice := strings.Split(filterOptions["suburb"][0], "-")
	suburbsInSearchRadius := getSuburbsInSearchRadius(searchSuburbSlice)
	// suburbsInDatabase := []interface{}{}
	suburbs := []middleware.ScrapedSuburb{}

	for _, suburb := range *suburbsInSearchRadius {
		suburbName := suburb.([]interface{})[1].(string)
		suburbPostcode := suburb.([]interface{})[2].(string)
		cacheId := fmt.Sprintf("%s-%s", suburbName, suburbPostcode)

		res, found := middleware.CheckCache(cacheId)

		if(found) {
			err := json.Unmarshal([]byte(res), &suburbs); if (err != nil) {
				fmt.Println(err)
			}


			// query db where suburb_id = suburb ID
			// select * from listings* where suburb_id IN (val1, val2 ...)	
			// this is a loop - build query here and then select elsewhere all at once
		} else {
			// CREATE QUERY STRING FOR SEARCH (SUBURB-STATE-POSTCODE, ...)
			// Scrape all as one in single request, partitioned and scraped as 1
			// will need to do outside of loop

			// this will have to move until data is actually scraped and sent to db vv
			scrapedSuburb := middleware.ScrapedSuburb{
				ID: cacheId,
				DatabaseId: suburb.([]interface{})[0].(int),
				SuburbName: suburbName,
				SuburbPostcode: suburbPostcode,
			}
			
			middleware.SetCache(cacheId, scrapedSuburb)
		}
	}

	fmt.Println("burbs: ",suburbs)

	properties := propertyScraperBase.InitialiseScraper(filterOptions)

	return properties
}

func getSuburbsInSearchRadius(searchSuburbSlice []string) *[]interface{} {
	db := dbClient.Connect()

	searchState := cases.
		Upper(language.Und, cases.NoLower).
		String(searchSuburbSlice[len(searchSuburbSlice) - 2:len(searchSuburbSlice) - 1][0])
	searchSuburb := cases.
		Upper(language.Und, cases.NoLower).
		String(strings.Join(
			searchSuburbSlice[:len(searchSuburbSlice) - 2], 
			" ",
		))
	searchPostcode := searchSuburbSlice[len(searchSuburbSlice) - 1:][0]

	fmt.Println(searchSuburb, searchState, searchPostcode)
	selectSuburbsInRadiusQuery := `
		SELECT id, suburb, postcode 
		FROM aus_suburbs 
		WHERE ST_DWithin(
			geolocation, (
				SELECT geolocation 
				FROM aus_suburbs 
				WHERE suburb = $1 and state = $2 and postcode = $3
			), 2500
		)
	`
	rows, e := db.Query(selectSuburbsInRadiusQuery, searchSuburb, searchState, searchPostcode);
	defer rows.Close()
	
	if e != nil {
		fmt.Println(e)
	}
	
	suburbsInSearchRadius := []interface{}{}

	for rows.Next() {
		var id int
		var name string
		var postcode string
		
		if err := rows.Scan(&id, &name, &postcode); err != nil {
			fmt.Println(err)
		} else {
			suburbsInSearchRadius = append(suburbsInSearchRadius, []interface{}{id, name, postcode})
		}
	}

	return &suburbsInSearchRadius
}

	// if err := rows.Err(); err != nil {
	// 	fmt.Println(err)
	// }

	// suburb := []Suburb{}

	// resJson := json.Unmarshal([]byte(res), v)


	// fmt.Println("Added 1 record")

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