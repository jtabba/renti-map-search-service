package dbClient

import (
	"database/sql"
	"encoding/json"
	"fmt"

	propertyTypes "back-end/mapSearchService/src/types"

	_ "github.com/lib/pq"
)

const (
	host    	= "localhost"
	port   	 	= 5438
	user   		= "postgres"
	password 	= "docker"
	dbname 		= "renti_db"
)

func Connect() *sql.DB {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", connectionString); if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	fmt.Printf("Successfully connected to PostgreSQL database %s at port %v \n", dbname, port)

	return db
}

func inputSuburbData() {
	db := Connect()
	defer db.Close()

	insertTest := `INSERT INTO NSW_Suburbs ("name", "state", "postcode", "geolocation", "geocode") VALUES ($1, $2, $3, $4, $5)`
	geoCode := propertyTypes.Geocode{ Lat: -33.8754959, Lng: 151.2047 }
	geocodeJson, _ := json.Marshal(&geoCode)

	_, e := db.Exec(
		insertTest, 
		"North Ryde", "NSW", 2113, "POINT(151.1314 -33.8002167)", 
		geocodeJson,
	); 
	
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println("Successfully inserted data")
	}
}