package requests

import (
	envHelper "back-end/mapSearchService/env"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var GOOGLE_API_KEY string = envHelper.GetEnvVar("GOOGLE_API_KEY")
var MAP_QUEST_API_KEY string = envHelper.GetEnvVar("MAP_QUEST_API_KEY")

func GetParsedData(url string) map[string]interface{} {
	request, err := http.NewRequest("GET", url, nil)

	if(err != nil) {
		fmt.Println(err)
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")

    client := &http.Client{}
    response, error := client.Do(request); if error != nil {
        fmt.Println(error)
    }

    responseBody, error := io.ReadAll(response.Body); if error != nil {
        fmt.Println(error)
    }

    formattedData := FormatJSON(responseBody)

    defer response.Body.Close()

    return formattedData
}

func FormatJSON(data []byte) map[string]interface{} {
    var out bytes.Buffer
    err := json.Indent(&out, data, "", " ")

    if err != nil {
        fmt.Println(err)
    }

    formattedData := string(out.Bytes())
	var result map[string]interface{}

	err = json.Unmarshal([]byte(formattedData), &result)

	if err != nil {
		panic(err)
	}

    return result
}

// Used by propertyScraperV1
func GetPropertyGeocode(address string) map[string]interface{} {
	if(GOOGLE_API_KEY == "") {
		panic("Google API key not provided")
	}

	addressSlice := strings.Split(address, " ")
	formattedAddress := strings.Join(addressSlice, "+")
	// url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s", formattedAddress, GOOGLE_API_KEY)
	url := fmt.Sprintf("https://www.mapquestapi.com/geocoding/v1/address?key=%s&location=%s", MAP_QUEST_API_KEY, formattedAddress)
	fmt.Printf("Geocoding: %s \n", url)
	geocodeData := GetParsedData(url)

	return geocodeData
}