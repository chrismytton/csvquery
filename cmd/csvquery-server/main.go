package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/chrismytton/csvquery"
)

func requestHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		return
	}
	log.Println(r.URL)
	queryString := r.URL.Query()
	if len(queryString["table"]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error\nPlease provide one or more 'table' parameters"))
		return
	}
	if len(queryString["query"]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error\nPlease provide a 'query' parameter"))
		return
	}
	tables := queryString["table"]
	query := queryString["query"][0]
	q, err := csvquery.New()
	if err != nil {
		log.Fatal(err)
	}
	defer q.Close()

	for _, tableSpec := range tables {
		parts := strings.SplitN(tableSpec, ":", 2)
		tableName := parts[0]
		csvURL := parts[1]
		resp, err := http.Get(csvURL)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		r := csv.NewReader(resp.Body)
		records, err := r.ReadAll()
		if err != nil {
			log.Fatal(err)
		}
		err = q.Insert(tableName, records)
		if err != nil {
			log.Fatal(err)
		}
	}

	result, err := q.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error\n%s\n", err)))
		return
	}

	writer := csv.NewWriter(w)
	for _, row := range result {
		writer.Write(row)
	}
	writer.Flush()
}

func main() {
	http.HandleFunc("/", requestHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
