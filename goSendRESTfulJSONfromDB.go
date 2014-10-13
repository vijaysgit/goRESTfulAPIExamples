/*
	This go program fetch data from MySQL database and send a RESTful JSON response.
*/

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"flag"
	_ "github.com/go-sql-driver/mysql"
)

var dbhost 		= flag.String("dbhost", "localhost", "The mysql hostname/ip address.")
var dbport *int = flag.Int("dbport", 3306, "The mysql port number.")
var dbuser 		= flag.String("dbuser", "root", "The mysql username to access the database.")
var dbpass 		= flag.String("dbpass", "root", "The mysql password to access the database.")
var dbname 		= flag.String("dbname", "currencydb", "The mysql database name.")
var debug *bool = flag.Bool("debug", false, "Print extra debugging info.")


type Data struct {
	Rate Rates
}

// Currency with code and rate
type Rates map[string]float64

func getCurrencies(response http.ResponseWriter, request *http.Request) {

		// Set MIME type to JSON in response header
		response.Header().Set("Content-type", "application/json")

		// Call a function to generate a JSON data with a list of currencies 
		respJSON, err := sendJSONResponse()
		if err != nil {
			panic(err)
		}

		// Display the JSON data on browswer/client 
		fmt.Fprintf(response, string(respJSON))

}


// Select the list of currencies from the database based on the latest timestamp.
func sendJSONResponse() ([]byte, error) {
		
		curRates := make(map[string]float64)

		// Construct the database connection using a string
		connectString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", *dbuser, *dbpass, *dbhost, *dbport, *dbname)
		if *debug {
			fmt.Printf("connectString:%s\n", connectString)
		}

		// Open the database
		db, err := sql.Open("mysql", connectString)
		defer db.Close()
		if err != nil {
			fmt.Println("Database connection failed", err)
		}

		// Get the latest timestamp from the database
		rows, err := db.Query("SELECT timeStamp FROM currencyTable ORDER BY timeStamp DESC LIMIT 1")
		if err != nil {
			fmt.Println("timeStamp query failed", err)
		}

		var latestTimeStamp string
		for rows.Next() {
    		err = rows.Scan(&latestTimeStamp)
		}

		// Get all currencies from the database with the latest timestamp
		rows, err = db.Query("SELECT curCode,curValue FROM currencyTable where timeStamp=?", latestTimeStamp)
		if err != nil {
			fmt.Println("Currency query failed", err)
		}

		defer rows.Close()

		for rows.Next() {
			var currCode string
			var currValue float64
			if err := rows.Scan(&currCode, &currValue); err != nil {
				fmt.Println("Rows scan failed", err)
			}
			curRates[currCode] = currValue
		}
		
		// Return all currencies in JSON data format 
		allCurrencies := Data{curRates}
		return json.MarshalIndent(allCurrencies, "", "  ")

}

func main() {

		// Call function to get currencies in JSON format
		http.HandleFunc("/", getCurrencies)

		// Set listening port
		http.ListenAndServe("127.0.0.1:8090", nil)

}