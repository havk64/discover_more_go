package main

import (
	"encoding/json"
	"net/http"
)

// Structure to fetch and parse the JSON data.
type users struct {
	Location string `json:"location"`
	Items    []struct {
		FullName string `json:"full_name"`
		Owner    struct {
			Login string `json:"login"`
		}
	}
}

// The HTTPResponse struct is used to get the make the async requests."
type HTTPResponse struct {
	index int
	url   string
	login string
	data  *json.Decoder
	res   http.Response
	//err   error
}

// userObject is the struct used as parameter to call the asyncHTTPGets
// func and has the index to be used to sort the result in a proper way.
type userObject struct {
	index int
	url   string
	login string
}

// resObj is the struct used to save the results to be displayed in the output.
type resObj struct {
	FullName string `json:"full_name"`
	Location string `json:"location"`
	Ranking  int    `json:"ranking"`
}
