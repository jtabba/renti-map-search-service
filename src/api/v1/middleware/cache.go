package middleware

import "fmt"

// "encoding/json"
// "fmt"
// "log"
// "time"

// "github.com/patrickmn/go-cache"

// type ScrapedSuburbsCacheType struct {
// 	activeScrapedSuburbs *cache.Cache
// }

type ScrapedSuburb struct {
	ID 				string 		`json:"id"`
	DatabaseId 		int 		`json:"databaseId"`
	SuburbName 		string 		`json:"suburbName"`
	SuburbPostcode 	string 		`json:"suburbPostcode"`
}

var scrapedSuburbsCache map[string]ScrapedSuburb = make(map[string]ScrapedSuburb)

func UpdateCache(id string, suburb ScrapedSuburb) {
	cacheRes := checkCache(id)

	if(!cacheRes["found"].(bool)) {
		fmt.Println("Not found in cache. Adding: ", id, suburb)
		scrapedSuburbsCache[id] = suburb

		return
	}

	// fmt.Println("Found in cache. No operation needed. ", scrapedSuburbsCache)

	return
}

func checkCache(id string) map[string]interface{} {
	data, found := scrapedSuburbsCache[id]; 
	
	if(found) {
		// fmt.Println("Found in cache: ", id, scrapedSuburbsCache)
	}
	
	return map[string]interface{}{
		"found": found,
		"data": data,
	}
}

func ReadCache(id string) (ScrapedSuburb, bool) {
	cacheRes := checkCache(id)

	fmt.Println("Read cache: ", cacheRes)

	return cacheRes["data"].(ScrapedSuburb), cacheRes["found"].(bool)
}

// const (
// 	defaultExpiration = 22 * time.Hour
// 	purgeTime = 24 * time.Hour
// )

// func InitialiseSuburbsCache() *ScrapedSuburbsCacheType {
// 	Cache := cache.New(defaultExpiration, purgeTime)

// 	return &ScrapedSuburbsCacheType{
// 		activeScrapedSuburbs: Cache,
// 	}
// }

// func (cache *ScrapedSuburbsCacheType) read(id string) (item []byte, found bool) {
// 	scrapedSuburb, found := cache.activeScrapedSuburbs.Get(id); 
	
// 	if found {
// 		fmt.Println("Found in cache: ", id)

// 		res, err := json.Marshal(scrapedSuburb.(ScrapedSuburb)); if err != nil {
// 			log.Fatal("Error unmarshaling suburb from cache. ID: ", id)
// 		}
		
// 		return res, true
// 	}

// 	return nil, false
// }

// func (suburbsCache *ScrapedSuburbsCacheType) update(id string, suburb ScrapedSuburb) {
// 	suburbsCache.activeScrapedSuburbs.Set(id, suburb, cache.DefaultExpiration)
// }

// func CheckCache(id string) ([]byte, bool) {
// 	return scrapedSuburbsCache.read(id)
// }

// func SetCache(id string, suburb ScrapedSuburb) {
// 	scrapedSuburbsCache.update(id, suburb)
// }

// var scrapedSuburbsCache = InitialiseSuburbsCache()