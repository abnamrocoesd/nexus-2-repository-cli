package backend

import (
	"io/ioutil"
	"log"
	"net/http"
)

func UrlToJsonStruct(url string, jsonStruct JsonStruct) (JsonStruct, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
		return nil, getErr
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
		return nil, readErr
	}
	jsonStruct, err = jsonStruct.Unmarshal(body)
	if err != nil {
		return nil, readErr
	}
	return jsonStruct, nil
}
