package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strings"
)

// Load CSV file into a sqlite database
func main() {
	// Read CSV file from disk and parse
	file, err := os.Open("test.csv")
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Extract various info from CSV
	headers := records[0]
	headerString := strings.Join(headers, ", ")
	csvRows := records[1:]

	// Create the correct number of placeholder question marks for using in the
	// prepared statement.
	questionMarks := make([]string, len(csvRows))
	for i := 0; i < len(csvRows); i++ {
		questionMarks[i] = "?"
	}
	rowQuestionMarks := strings.Join(questionMarks, ", ")

	// Open an in-memory sqlite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a table for the CSV to live in
	// TODO: Make table name configurable
	sqlStatement := fmt.Sprintf("create table test (%s)", headerString)
	_, err = db.Exec(sqlStatement)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStatement)
		return
	}

	// Insert CSV data into sqlite db
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	sqlStatement = fmt.Sprintf("insert into test(%s) values(%s)", headerString, rowQuestionMarks)
	stmt, err := tx.Prepare(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for _, row := range csvRows {
		rowCopy := make([]interface{}, len(row))
		for i, d := range row {
			rowCopy[i] = d
		}
		_, err = stmt.Exec(rowCopy...)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()

	// Get the data back out of the CSV
	sqlRows, err := db.Query("select id, name from test")
	if err != nil {
		log.Fatal(err)
	}
	defer sqlRows.Close()

	// Dump the data back out to CSV
	colNames, err := sqlRows.Columns()
	if err != nil {
		log.Fatal(err)
	}
	writer := csv.NewWriter(os.Stdout)
	writer.Write(colNames)
	readCols := make([]interface{}, len(colNames))
	writeCols := make([]string, len(colNames))
	for i, _ := range writeCols {
		readCols[i] = &writeCols[i]
	}
	for sqlRows.Next() {
		err := sqlRows.Scan(readCols...)
		if err != nil {
			log.Fatal(err)
		}
		writer.Write(writeCols)
	}
	if err = sqlRows.Err(); err != nil {
		log.Fatal(err)
	}
	writer.Flush()
}
