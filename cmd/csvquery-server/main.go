package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/chrismytton/csvquery"
)

func errorResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf("error\n%s\n", message)))
}

func badRequest(w http.ResponseWriter, message string) {
	errorResponse(w, http.StatusBadRequest, message)
}

func internalServerError(w http.ResponseWriter, message string) {
	errorResponse(w, http.StatusInternalServerError, message)
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		return
	}
	log.Println(r.URL)
	origin := r.Header.Get("Origin")
	if origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	queryString := r.URL.Query()
	if len(queryString["table[]"]) == 0 {
		badRequest(w, "Please provide one or more 'table' parameters")
		return
	}
	if len(queryString["query"]) == 0 {
		badRequest(w, "Please provide a 'query' parameters")
		return
	}
	tables := queryString["table[]"]
	query := queryString["query"][0]
	q, err := csvquery.New()
	if err != nil {
		internalServerError(w, fmt.Sprintf("An unexpected error occurred: %s. Please try again later. If this error persists, please get in touch.", err))
		return
	}
	defer q.Close()

	for _, tableSpec := range tables {
		parts := strings.SplitN(tableSpec, ":", 2)
		tableName := parts[0]
		csvURL := parts[1]
		resp, err := http.Get(csvURL)
		if err != nil {
			badRequest(w, fmt.Sprintf("Couldn't GET URL %s: %s", csvURL, err))
			return
		}
		defer resp.Body.Close()
		r := csv.NewReader(resp.Body)
		records, err := r.ReadAll()
		if err != nil {
			badRequest(w, fmt.Sprintf("Couldn't parse CSV from URL %s: %s", csvURL, err))
			return
		}
		err = q.Insert(tableName, records)
		if err != nil {
			internalServerError(w, fmt.Sprintf("Error trying to load data: %s", err))
			return
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
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
