goRESTfulAPIExamples
====================

go RESTful API with examples

#goRESTfulAPIExamples

# Install Go
# Install MySQL

# Installation

        git clone https://github.com/vijaysgit/goRESTfulAPIExamples
        go get github.com/go-sql-driver/mysql


Note: A valid MySQL user account is required. Create a database "currencydb" and a table "currencyTable" with the dump file "gocurrencydb.sql" using the following command.

	      mysql -u user -p password < gocurrencydb.sql

        cd goRESTfulAPIExamples
        

Run the programs in the following order:

#1 To read all currencies in JSON format from the "testdata.txt" file and store it into the database:

	      go run goReadnSaveJSONtoDB.go


#2 To fetch all currencies from the database with a latest timestamp and send it as a JSON response:

	      go run goSendRESTfulJSONfromDB.go

	      Open "http://127.0.0.1:8090" in the browser to view the JSON response.


#3 To receive the JSON response and display it into a simple HTML table:

	      go run goReadnDisplayRESTfulJSON.go

	      Open "http://127.0.0.1:9090" in your browser to view all currencies in a table.
