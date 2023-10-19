package dbClient

// import (
// 		dbClient "back-end/mapSearchService/src/database/v2-Psql"
// )

func test() {
	db := Connect()
	defer db.Close()

	
}