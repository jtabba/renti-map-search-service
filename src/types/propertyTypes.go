package propertyTypes

import (
	"github.com/google/uuid"
)

// type GeocodeData struct {
// 	FormattedAddress string `json:"formattedAddress"`
// 	Latitude		 string `json:"latitude"`
// 	Longitude		 string `json:"longitude"`
// 	Country			 string `json:"country"`
// 	City 			 string `json:"city"`
// 	StateCode		 string `json:"stateCode"`
// 	ZipCode			 string `json:"zipCode"`
// 	StreetName		 string `json:"streetName"`
// 	StreetNumber	 string `json:"streetNumber"`
// 	CountryCode		 string `json:"countryCode"`
// 	Provider		 string `json:"provider"`
// }

type Property struct {
	Geocode 			Geocode		`json:"geocode"`
	Location 			GeoJSON 	`json:"location"`
	Images  			[]interface {}    `json:"images"`
	Baths 				float64       	`json:"baths"`
	Beds  				float64       	`json:"beds"`
	Parking 			float64       	`json:"parking"`
	Agency 				map[string]interface{} 		`json:"agency"`
	SuburbId 			int 		`json:"suburb_id"`
	// Size	  			uint       	`json:"size"`
	// Utilities 			[]string  	`json:"utilities"`
	WeeklyPrice 		string     	`json:"weekly_price"`
	ID					uuid.UUID    	`json:"id"`
	Address 			string 		`json:"address"`
	Suburb 				string 		`json:"suburb"`
	// City  				string    	`json:"city"`
	State  				string    	`json:"state"`
	Country  			string    	`json:"country"`
	Postcode  			string    	`json:"postcode"`
	InspectionOpenTime 		string 		`json:"inspection_open_time"`
	InspectionCloseTime 	string 		`json:"inspection_close_time"`
	PropertyType  		string 		`json:"property_type"`
	// AgentNumber 		string 		`json:"number"`
}

type Utilities struct {
	Baths 			uint       	`json:"baths"`
	Beds 			uint       	`json:"beds"`
	Parking 		uint       	`json:"parking"`
} 

// type Property struct {
// 	Address 			string 		`json:"address"`;
// 	InspectionTime 		string 		`json:"inspection_time"`;
// 	WeeklyPrice 		string    	`json:"weeklyPrice"`
// }

type GeoJSON struct {
	Type 			string 			`json:"type"`
	Coordinates 	[]float64 		`json:"coordinates"`
}

type Geocode struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}