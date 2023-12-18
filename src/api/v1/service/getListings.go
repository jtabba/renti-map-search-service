package propertiesServiceV1

import (
	"encoding/json"
	"fmt"
	"strings"

	middleware "back-end/mapSearchService/src/api/v1/middleware"
	dbClient "back-end/mapSearchService/src/database/v2-Psql"
	types "back-end/mapSearchService/src/types"
	utilities "back-end/mapSearchService/src/utilities"

	"github.com/lib/pq"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

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
			"-",
		))
	searchPostcode := searchSuburbSlice[len(searchSuburbSlice) - 1:][0]

	fmt.Println(searchSuburb,)

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

func SeparateListingsInDb(suburb string) (*[]float64, *map[string]interface{}) {
	searchSuburbSlice := strings.Split(suburb, "-")
	suburbsInSearchRadius := getSuburbsInSearchRadius(searchSuburbSlice) 
	scrapedSububurbsIds, suburbsToScrape := []float64{}, map[string]interface{}{}
	fmt.Println("IN RADIUS: ", suburbsInSearchRadius)

	for _, suburb := range *suburbsInSearchRadius {
		suburbName, suburbState, suburbPostcode := 
		suburb.([]interface{})[1].(string), 
		suburb.([]interface{})[2].(string), 
		suburb.([]interface{})[3].(string)
		cacheId := fmt.Sprintf("%s-%s-%s", suburbName, suburbState, suburbPostcode)
		cachedSuburb, found := middleware.ReadCache(cacheId)
		// fmt.Println("In radius: ", suburb, found, cachedSuburb)

		if found {
			// fmt.Println("BEFORE")
			// cachedSuburb := utilities.FormatJSON(cachedSuburb)
			scrapedSububurbsIds = append(scrapedSububurbsIds, float64(cachedSuburb.DatabaseId))
			// fmt.Println("AFTER", scrapedSububurbsIds)
		} else {
			suburbData := middleware.ScrapedSuburb{
					ID: cacheId,
					DatabaseId: suburb.([]interface{})[0].(int),
					SuburbName: suburbName,
					SuburbPostcode: suburbPostcode,
				}
			
			suburbsToScrape[cacheId] = suburbData
		}
	}

	return &scrapedSububurbsIds, &suburbsToScrape
}

func InsertScrapedListingsIntoDb(listings *[]types.Property) {
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

func GetListingsInDb(suburbsInDatabase *[]float64) *[]types.Property {
	db := dbClient.Connect()
	defer db.Close()

	selectSuburbsInRadiusQuery := `
		SELECT *
		FROM listings 
		WHERE suburb_id = ANY($1)
	`

	rows, err := db.Query(selectSuburbsInRadiusQuery, pq.Array(suburbsInDatabase));
	defer rows.Close()
	
	if err != nil {
		fmt.Println(err)
	}
	
	suburbsInSearchRadius := []types.Property{}

	for rows.Next() {
		var listing types.Property
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
			listing.Agency = utilities.FormatJSON(agencyBytes)
			listing.Geolocation = string(geolocationBytes)
			listing.Geocode = utilities.FormatJSON(geocodeBytes)

			suburbsInSearchRadius = append(suburbsInSearchRadius, listing)
		}
	}

	return &suburbsInSearchRadius
}
