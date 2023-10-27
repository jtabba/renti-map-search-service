package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/patrickmn/go-cache"
)

type ScrapedSuburbsCacheType struct {
	activeScrapedSuburbs *cache.Cache
}

type ScrapedSuburb struct {
	ID 				string 		`json:"id"`
	DatabaseId 		int 		`json:"databaseId"`
	SuburbName 		string 		`json:"suburbName"`
	SuburbPostcode 	string 		`json:"suburbPostcode"`
}

const (
	defaultExpiration = 22 * time.Hour
	purgeTime = 24 * time.Hour
)

func InitialiseSuburbsCache() *ScrapedSuburbsCacheType {
	Cache := cache.New(defaultExpiration, purgeTime)

	return &ScrapedSuburbsCacheType{
		activeScrapedSuburbs: Cache,
	}
}

func (cache *ScrapedSuburbsCacheType) read(id string) (item []byte, found bool) {
	scrapedSuburb, found := cache.activeScrapedSuburbs.Get(id); 
	
	if found {
		fmt.Println("Found in cache: ", id)

		res, err := json.Marshal(scrapedSuburb.(ScrapedSuburb)); if err != nil {
			log.Fatal("Error unmarshaling suburb from cache. ID: ", id)
		}
		
		return res, true
	}

	return nil, false
}

func (suburbsCache *ScrapedSuburbsCacheType) update(id string, suburb ScrapedSuburb) {
	suburbsCache.activeScrapedSuburbs.Set(id, suburb, cache.DefaultExpiration)
}

func CheckCache(id string) ([]byte, bool) {
	res, found := scrapedSuburbsCache.read(id)

	if found {
		return res, true
	}

	return nil, false
}

func SetCache(id string, suburb ScrapedSuburb) {
	scrapedSuburbsCache.update(id, suburb)
}

var scrapedSuburbsCache = InitialiseSuburbsCache()