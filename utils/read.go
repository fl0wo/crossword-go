// File: utils/read.go
package utils

import (
	"encoding/json"
	"io/ioutil"
)

// Data represents the structure of each object in the JSON array
type Data struct {
	Nome string   `json:"nome"`
	Desc []string `json:"desc"`
}

func ReadWords() []Data {
	// Read the JSON file
	content, err := ioutil.ReadFile("./assets/data.json")
	if err != nil {
		return nil
	}

	// Create a slice to store multiple Data structs
	var payload []Data

	// Unmarshal the JSON data into the slice
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return nil
	}

	return payload
}
