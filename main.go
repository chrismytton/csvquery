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

func csvQuery(tables map[string][][]string, query string) (result [][]string, err error) {
	// Open an in-memory sqlite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return result, err
	}
	defer db.Close()

	for tableName, records := range tables {
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

		// Create a table for the CSV to live in
		// TODO: Make table name configurable
		sqlStatement := fmt.Sprintf("CREATE TABLE %s (%s)", tableName, headerString)
		_, err = db.Exec(sqlStatement)
		if err != nil {
			return result, err
		}

		// Insert CSV data into sqlite db
		tx, err := db.Begin()
		if err != nil {
			return result, err
		}
		sqlStatement = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, headerString, rowQuestionMarks)
		stmt, err := tx.Prepare(sqlStatement)
		if err != nil {
			return result, err
		}
		defer stmt.Close()
		for _, row := range csvRows {
			rowCopy := make([]interface{}, len(row))
			for i, d := range row {
				rowCopy[i] = d
			}
			_, err = stmt.Exec(rowCopy...)
			if err != nil {
				return result, err
			}
		}
		tx.Commit()
	}

	// Get the data back out of the CSV
	sqlRows, err := db.Query(query)
	if err != nil {
		return result, err
	}
	defer sqlRows.Close()

	// Dump the data back out to CSV
	colNames, err := sqlRows.Columns()
	if err != nil {
		return result, err
	}
	result = append(result, colNames)

	readCols := make([]interface{}, len(colNames))
	writeCols := make([]string, len(colNames))
	for i, _ := range writeCols {
		readCols[i] = &writeCols[i]
	}
	for sqlRows.Next() {
		err := sqlRows.Scan(readCols...)
		if err != nil {
			return result, err
		}
		cols := make([]string, len(writeCols))
		copy(cols, writeCols)
		result = append(result, cols)
	}
	if err = sqlRows.Err(); err != nil {
		return result, err
	}
	return result, nil
}

func readCsv(fileName string) ([][]string, error) {
	records := [][]string{}
	// Read CSV file from disk and parse
	file, err := os.Open(fileName)
	if err != nil {
		return records, err
	}
	r := csv.NewReader(file)
	records, err = r.ReadAll()
	if err != nil {
		return records, err
	}
	return records, nil
}

// Load CSV file into a sqlite database
func main() {
	files := map[string]string{
		"test": "test.csv",
		"ages": "ages.csv",
	}
	tables := make(map[string][][]string)
	for tableName, fileName := range files {
		records, err := readCsv(fileName)
		if err != nil {
			log.Fatal(err)
		}
		tables[tableName] = records
	}
	result, err := csvQuery(tables, "select id, name, age from test join ages on ages.person_id = test.id")
	if err != nil {
		log.Fatal(err)
	}

	writer := csv.NewWriter(os.Stdout)
	for _, row := range result {
		writer.Write(row)
	}
	writer.Flush()
}
