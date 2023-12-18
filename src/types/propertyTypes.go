package types

type Property struct {
	Geocode 				map[string]interface{} 		`json:"geocode"`
	Geolocation 			string 						`json:"geolocation"`
	Images  				[]string    				`json:"images"`
	Baths 					float64       				`json:"baths"`
	Beds  					float64       				`json:"beds"`
	Parking 				float64       				`json:"parking"`
	Agency 					map[string]interface{} 		`json:"agency"`
	SuburbId 				int 						`json:"suburb_id"`
	WeeklyPrice 			string      					`json:"weekly_price"`
	ID						string    					`json:"id"`
	Address 				string 						`json:"address"`
	Suburb 					string 						`json:"suburb"`
	State  					string    					`json:"state"`
	Country  				string    					`json:"country"`
	Postcode  				string    					`json:"postcode"`
	InspectionOpenTime 		string 						`json:"inspection_open_time"`
	InspectionCloseTime 	string 						`json:"inspection_close_time"`
	PropertyType  			string 						`json:"property_type"`
}

type Suburb struct {
	ID 						int 						`json:"id"`
	Suburb 					string 						`json:"suburb"`
	State  					string    					`json:"state"`
	Postcode  				string    					`json:"postcode"`
	Geolocation 			string 						`json:"geolocation"`
	Geocode 				map[string]interface{} 		`json:"geocode"`
}