// Package asyncsorted fetches the first 30 most popular GitHub repositories
// about Go language using async http requests and outputs the sorted results.
//
// Instructions:
// The files struct.go and url.go are part of this package and thus are needed
// to be compiled together.
// Work in Progress!!!
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func asyncHTTPGets(user []*userObject) []*HTTPResponse {
	ch := make(chan *HTTPResponse)
	responses := make([]*HTTPResponse, len(user)) // Slice of Pointers to Struct
	for _, item := range user {
		index := item.index //Assigning variables from map obj.
		url := item.url
		login := item.login
		go func() { //Goroutine
			fmt.Printf("Fetching url: %s, ranking: %d \n", url, index+1)
			data, res := fetchData(url)
			ch <- &HTTPResponse{index, url, login, data, res} //Pointers to channel
		}()
	}

	for {
		select {
		case r, ok := <-ch:
			if ok {
				fmt.Printf("%s, ranking %d was fetched\n", r.url, r.index+1)
				fmt.Printf("%t\n", ok)
				responses[r.index] = r
			} else {
				return responses
			}
		case <-time.After(50 * time.Millisecond):
			fmt.Printf(".")
		}
	}
}

/* Function fetchData to make the http requests to Github API */
func fetchData(url string) (*json.Decoder, http.Response) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	check(err)
	req.Header = customHeader()
	resp, err := client.Do(req)
	check(err)
	data := resp.Body
	decoder := json.NewDecoder(data) // Parsing the JSON Object.
	return decoder, *resp
}

func main() {
	start := time.Now()
	u := customURL() // Replacing previous code below with this function.
	/* ====================================================================================
	  u, error := url.Parse("https://api.github.com/search/repositories") //Parsing the URL
		check(error)
		q := u.Query() //Getting url.Query() in order to specify the Query
		q.Add("q", "language:go") // Search for all Golang repositories
		q.Add("sort", "stars")    // Sort by amount of Stars
		q.Add("order", "desc")    // In descentant order
		u.RawQuery = q.Encode()   // Encoding the query to make it a encode string
	* ===================================================================================== */
	github := users{} // Struct of Github API // Before was declared as: `var github users`.
	decoder, res := fetchData(u.String())
	err := decoder.Decode(&github)
	check(err)
	defer res.Body.Close()
	defer fmt.Println("BOOOOOMMMMM ! ! !\n30 URLs fetched in ", time.Since(start))
	repos := len(github.Items)
	result := make([]*resObj, repos)     //Initializing slice of pointers to resObj struct
	names := make([]string, repos)       //Initializing empty arrays
	ghUser := make([]*userObject, repos) //Initializing empty arrays of pointers.(to be used as function parameter)
	for index, item := range github.Items {
		name := item.FullName
		login := item.Owner.Login
		u, _ := url.Parse("https://api.github.com")
		u.Path = "/users/" + login
		names[index] = name
		ghUser[index] = &userObject{
			index: index,
			url:   u.String(),
			login: login,
		}
	}
	results := asyncHTTPGets(ghUser)
	for index, item := range results {
		loc := &users{} // Before was declared: `var loc users`.
		decoder := item.data
		error := decoder.Decode(&loc)
		check(error)
		defer item.res.Body.Close()
		/* Object to be displayed in the output */
		result[index] = &resObj{
			FullName: names[item.index],
			Location: loc.Location,
			Ranking:  (item.index + 1),
		}
	}
	ar, err := json.MarshalIndent(result, "", "    ") /* Indenting the output(Json Prettifyied) */
	check(err)
	fmt.Println(string(ar))
}
