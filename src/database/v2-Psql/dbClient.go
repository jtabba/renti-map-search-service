package dbClient

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host    	= "localhost"
	port   	 	= 5436
	user   		= "postgres"
	password 	= "docker"
	dbname 		= "map_search_listings"
)

func Connect() *sql.DB {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", connectionString); if err != nil {
		panic(err)
	}

	// defer db.Close()

	if err = db.Ping(); err != nil {
		panic(err)
	}

	fmt.Printf("Successfully connected to PostgreSQL database %s at port %v \n", dbname, port)

	return db
}