/*
	This go program read JSON data stored in a text file and save it to a MySQL database.
*/

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
	"time"
	"flag"
)

var dbhost 		= flag.String("dbhost", "localhost", "The mysql hostname/ip address.")
var dbport *int = flag.Int("dbport", 3306, "The mysql port number.")
var dbuser 		= flag.String("dbuser", "root", "The mysql username to access the database.")
var dbpass 		= flag.String("dbpass", "root", "The mysql password to access the database.")
var dbname 		= flag.String("dbname", "currencydb", "The mysql database name.")
var debug *bool = flag.Bool("debug", false, "Print extra debugging info.")


func insertToDB(allCurrencies map[string]interface{}, curTimestamp int64) {

		// Convert the timestamp
		timeStamp := time.Unix(curTimestamp, 0)


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
	
		
    	// Construct the SQL insert query using a string and interface 
		sqlString := "INSERT INTO currencyTable(timeStamp, curCode, curValue) VALUES "
		sqlValues := []interface{}{}

		for currencyName, _ := range allCurrencies {
			sqlString += "(?, ?, ?),"
			sqlValues = append(sqlValues, timeStamp, currencyName, allCurrencies[currencyName].(float64))
		}
		sqlString = strings.TrimSuffix(sqlString, ",")

		// Create the SQL insert statement
		stmt, _ := db.Prepare(sqlString)

		// Execute the SQL insert statement
		_ , err = stmt.Exec(sqlValues...)
		if err != nil {
			fmt.Println(err)
			return
		}
	
}


func checkError(err error) {
		if err != nil {
			panic(err)
		}
}

func getCurrencies() (currString string) {

		// Open the currency text file stored locally - "testdata.txt"
		dataFile := "testdata.txt"
		dfile, err := os.Open(dataFile)
		if err != nil {
			fmt.Println(dataFile, err)
			return
		}

    	// Read data from the text file 
		defer dfile.Close()
		body, err := ioutil.ReadAll(dfile)
		if err != nil {
			fmt.Println("file read error")
			return
		}

		// Convert the data to a string format and remove extra spaces
		currString = string(body)
		currString= strings.TrimSpace(currString)
		return currString

}


func main() {

		// Function call to get JSON data stored in a text file 
		currString := getCurrencies()

		// Decode the JSON data 
		data := make(map[string]interface{})
		err := json.Unmarshal([]byte(currString), &data)
		if err != nil {
			panic(err)
		}

		// Extract the currencies and timestamp from JSON data
		allCurrencies := data["rates"].(map[string]interface{})
		curTimestamp := int64((data["timestamp"]).(float64))

		// Convert the timestamp and calculate the time difference
		timeStamp := time.Unix(curTimestamp, 0)
		duration := time.Since(timeStamp)
	
		// Check the time difference to store the currencies in to the database only once in a day
		if duration.Hours() > 24 {
			// Funciton call to insert the currencies into the database
			insertToDB(allCurrencies, curTimestamp)
			fmt.Println("Currencies stored in to the database")
		} else {
			fmt.Println("Currencies already stored in to the database today. Try tomorrow...")
		}

}
