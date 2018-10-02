package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chrismytton/csvquery"
)

func readCsv(fileName string) ([][]string, error) {
	var records [][]string
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

type csvTableSpec map[string]string

func (c *csvTableSpec) String() string {
	return fmt.Sprint(*c)
}

func (c *csvTableSpec) Set(value string) error {
	parts := strings.Split(value, ":")
	table := parts[0]
	file := parts[1]
	(*c)[table] = file
	return nil
}

// Load CSV file into a sqlite database
func main() {
	var files = make(csvTableSpec)
	var query string
	flag.Var(&files, "table", "table to file mapping in the form tablename:file.csv")
	flag.StringVar(&query, "query", "", "the SQL query to run against the given tables")
	flag.Parse()

	tables := make(map[string][][]string)
	for tableName, fileName := range files {
		records, err := readCsv(fileName)
		if err != nil {
			log.Fatal(err)
		}
		tables[tableName] = records
	}
	q, err := csvquery.New()
	if err != nil {
		log.Fatal(err)
	}
	defer q.Close()

	for tableName, rows := range tables {
		err := q.Insert(tableName, rows)
		if err != nil {
			log.Fatal(err)
		}
	}

	result, err := q.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	writer := csv.NewWriter(os.Stdout)
	for _, row := range result {
		writer.Write(row)
	}
	writer.Flush()
}
