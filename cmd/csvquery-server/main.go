package main

import (
	"encoding/csv"
	"github.com/chrismytton/csvquery"
	"log"
	"net/http"
	"strings"
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
		csvUrl := parts[1]
		resp, err := http.Get(csvUrl)
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
		log.Fatal(err)
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
