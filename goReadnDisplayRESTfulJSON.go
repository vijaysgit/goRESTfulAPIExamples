/*
	This go program will read RESTFul response in JSON format and display all currencies in a table on client
*/

package main

import (
	
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"	
)


func getJSONData() (currString  string) {

		// Get all currencies from a RESTful JSON response
		resp, err := http.Get("http://127.0.0.1:8090")
		if err != nil {
			fmt.Println("http get error", err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("http read error")
			return
		}

		currString  = string(body)
		currString  = strings.TrimSpace(currString )

		fmt.Println("here is india")

		return currString

}


func parseJSONData() map[string]interface{} {
	
		// Get JSON data
		currString := getJSONData()

		// Decode the JSON data
		parsedJSONData := make(map[string]interface{})
		err := json.Unmarshal([]byte(currString), &parsedJSONData)
		if err != nil {
			panic(err)
		}

		return parsedJSONData

}


func dispCurrencyTable(response http.ResponseWriter, request *http.Request) {


		// Get all currencies from JSON response
		resData := make(map[string]interface{})
		resData = parseJSONData()

		// Get all currency names/codes 
		currencies := resData["Rate"].(map[string]interface{})
		var currCode []string
		for currIndex := range currencies {
			currCode = append(currCode, currIndex)
		}

		// Sort currency names/codes in alphabetical order
		sort.Strings(currCode)

		// Construct a HTML table using a string to display all currencies
		tableString := "<head> <title>Currency Table</title> </head> <body> <table border= \"1\" width=\"100\"> <tr> <th>Currency Name</th> <th>Value</th> </tr>"

		for _ , currencyCode := range currCode {
			tableString = tableString + "<tr> <td>" + currencyCode + "</td> <td>" + 
										strconv.FormatFloat(currencies[currencyCode].(float64), 'f', 2, 64) + "</td> </tr>"
		}
		tableString = tableString + "</table> </body>"

		// Display the currency table on client/browser 
		fmt.Fprintf(response, tableString)

}


func checkError(err error) {

		if err != nil {
			panic(err)
		}
}


func main() {

		// Display currency data from the JSON response in a table
		http.HandleFunc("/", dispCurrencyTable)
	
		// Set listening port
		http.ListenAndServe("127.0.0.1:9090", nil)
		
	
}
