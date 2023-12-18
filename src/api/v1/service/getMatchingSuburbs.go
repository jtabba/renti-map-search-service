package propertiesServiceV1

import (
	dbClient "back-end/mapSearchService/src/database/v2-Psql"
	types "back-end/mapSearchService/src/types"
	utilities "back-end/mapSearchService/src/utilities"
	"fmt"
)

func GetMatchingSuburbs(currentSearchSuburb string) *[]types.Suburb {
	db := dbClient.Connect()
	defer db.Close()
	
	suburbsThatMatchQuery := `
		SELECT 
			id, 
			suburb, 
			state, 
			postcode, 
			ST_AsText(geolocation), 
			geocode 
		FROM aus_suburbs
		WHERE (
			suburb ILIKE '%' || $1 || '%'
		) LIMIT 8
	`

	rows, err := db.Query(suburbsThatMatchQuery, currentSearchSuburb)
	defer rows.Close()

	if err != nil {
		fmt.Println(err)
	}

	suburbsThatMatch := []types.Suburb{}

	for rows.Next() {
		var suburb types.Suburb
		geocodeBytes := []byte{}
		geolocationBytes := []byte{}

		if err := rows.Scan(
			&suburb.ID, 
			&suburb.Suburb, 
			&suburb.State, 
			&suburb.Postcode, 
			&geolocationBytes, 
			&geocodeBytes,
		); err != nil {
			fmt.Println(err)
		} else {
 			suburb.Geolocation = string(geolocationBytes)
			suburb.Geocode = utilities.FormatJSON(geocodeBytes)

			suburbsThatMatch = append(suburbsThatMatch, suburb)
		}
	}

	fmt.Println("Query matches: ", len(suburbsThatMatch))

	return &suburbsThatMatch
}