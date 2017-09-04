package main

import (
	"encoding/csv"
	"github.com/chrismytton/csvsql"
	"log"
	"os"
)

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
	q, err := csvsql.New(tables)
	if err != nil {
		log.Fatal(err)
	}
	defer q.Close()

	result, err := q.Query("select id, name, age from test join ages on ages.person_id = test.id")
	if err != nil {
		log.Fatal(err)
	}

	writer := csv.NewWriter(os.Stdout)
	for _, row := range result {
		writer.Write(row)
	}
	writer.Flush()
}
