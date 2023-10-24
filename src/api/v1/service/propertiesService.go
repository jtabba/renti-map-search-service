package propertiesServiceV1

import (
	"back-end/mapSearchService/src/api/v1/middleware"
	dbClient "back-end/mapSearchService/src/database/v2-Psql"
	"back-end/mapSearchService/src/requests"
	propertyScraperBase "back-end/mapSearchService/src/scraper"
	propertyTypes "back-end/mapSearchService/src/types"
	"encoding/json"
	"fmt"

	"strings"

	"net/url"

	"github.com/lib/pq"
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
	scrapedSububurbsIds, suburbsToScrape := separateListingsInDb(filterOptions["suburb"][0])
	properties := []propertyTypes.Property{}
	
	if(len(*suburbsToScrape) > 0) {
		fmt.Println("Scraping: ", suburbsToScrape)
		suburbListings := propertyScraperBase.InitialiseScraper(filterOptions, *suburbsToScrape)
		
		properties = append(properties, *suburbListings...)

		insertScrapedListingsIntoDb(&properties)
	}

	if(len(*scrapedSububurbsIds) > 0) {
		listingsInDatabase := getListingsInDb(scrapedSububurbsIds)
		properties = append(properties, *listingsInDatabase...) 
	}

	return &properties
}

func getSuburbsInSearchRadius(searchSuburbSlice []string) *[]interface{} {
	db := dbClient.Connect()
	defer db.Close()

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
		SELECT id, suburb, state, postcode 
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
		var suburb string
		var state string
		var postcode string
		
		if err := rows.Scan(&id, &suburb, &state, &postcode); err != nil {
			fmt.Println(err)
		} else {
			suburbsInSearchRadius = append(suburbsInSearchRadius, []interface{}{id, suburb, state, postcode})
		}
	}

	return &suburbsInSearchRadius
}

func separateListingsInDb(suburb string) (*[]float64, *[]string) {
	searchSuburbSlice := strings.Split(suburb, "-")
	suburbsInSearchRadius := getSuburbsInSearchRadius(searchSuburbSlice)
	scrapedSububurbsIds := []float64{}
	suburbsToScrape := []string{}

	for _, suburb := range *suburbsInSearchRadius {
		suburbName := suburb.([]interface{})[1].(string)
		suburbState := suburb.([]interface{})[2].(string)
		suburbPostcode := suburb.([]interface{})[3].(string)
		cacheId := fmt.Sprintf("%s-%s", suburbName, suburbPostcode)

		res, found := middleware.CheckCache(cacheId)

		if(found) {
			cachedData := requests.FormatJSON(res)

			// fmt.Println("Found in cache: ", cachedData)

			scrapedSububurbsIds = append(scrapedSububurbsIds, cachedData["databaseId"].(float64))


			// query db where suburb_id = suburb ID
			// select * from listings* where suburb_id IN (val1, val2 ...)	
			// this is a loop - build query here and then select elsewhere all at once
		} else {
			suburbQueryString := fmt.Sprintf("%s-%s-%s", suburbName, suburbState, suburbPostcode)
			suburbsToScrape = append(suburbsToScrape, suburbQueryString)
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

 	fmt.Println("burbs: ",suburbsToScrape)
	return &scrapedSububurbsIds, &suburbsToScrape
}

func insertScrapedListingsIntoDb(listings *[]propertyTypes.Property) {
	db := dbClient.Connect()
	defer db.Close()

	insertListingQuery := `
		WITH suburb_details AS (
			SELECT *
			FROM aus_suburbs
			WHERE suburb = $8 AND state = $9 AND postcode = $11
		)
		INSERT INTO listings (
			weekly_price,
			inspection_open_time,
			inspection_close_time,
			geocode,
			geolocation,
			images,
			address,
			suburb,
			state,
			country,
			postcode,
			beds,
			baths,
			parking,
			property_type,
			agency,
			suburb_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, 
			(SELECT id FROM suburb_details)
		)
	`

	for _, listing := range(*listings) {
		geocodeJSON, err := json.Marshal(listing.Geocode)
		agencyJSON, err := json.Marshal(listing.Agency)

		if(err != nil) {
			panic(err)
		}

		_, err = db.Exec(
			insertListingQuery,
			listing.WeeklyPrice,
			listing.InspectionOpenTime,
			listing.InspectionCloseTime,
			geocodeJSON,
			listing.Geolocation,
			pq.Array(listing.Images),
			listing.Address,
			listing.Suburb,
			listing.State,
			listing.Country,
			listing.Postcode,
			listing.Beds,
			listing.Baths,
			listing.Parking,
			listing.PropertyType,
			agencyJSON,
		)
		
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Printf("Added %v records \n", len(*listings))
}

func getListingsInDb(suburbsInDatabase *[]float64) *[]propertyTypes.Property {
	db := dbClient.Connect()
	defer db.Close()

	selectSuburbsInRadiusQuery := `
		SELECT * 
		FROM listings 
		WHERE suburb_id = ANY($1)
	`

	rows, e := db.Query(selectSuburbsInRadiusQuery, pq.Array(suburbsInDatabase));
	defer rows.Close()
	
	if e != nil {
		fmt.Println(e)
	}
	
	suburbsInSearchRadius := []propertyTypes.Property{}

	for rows.Next() {
		var listing propertyTypes.Property
		var imagesBytes pq.StringArray
		agencyBytes := []byte{}
		geocodeBytes := []byte{}
		geolocationBytes := []byte{}
		
		if err := rows.Scan(
			&listing.ID,
			&listing.Address,
			&listing.SuburbId,
			&geolocationBytes,
			&geocodeBytes,
			&imagesBytes,
			&listing.Baths,
			&listing.Beds,
			&listing.Parking,
			&agencyBytes,
			&listing.WeeklyPrice,
			&listing.Suburb,
			&listing.State,
			&listing.Postcode,
			&listing.Country,
			&listing.InspectionOpenTime,
			&listing.InspectionCloseTime,
			&listing.PropertyType,
		); err != nil {
			fmt.Println(err)
		} else {
			listing.Images = []string(imagesBytes)
			listing.Agency = requests.FormatJSON(agencyBytes)
			listing.Geolocation = string(geolocationBytes)
			listing.Geocode = requests.FormatJSON(geocodeBytes)


			suburbsInSearchRadius = append(suburbsInSearchRadius, listing)
		}
	}

	return &suburbsInSearchRadius
}
