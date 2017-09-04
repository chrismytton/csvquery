package csvsql

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

type CSVTable struct {
	name string
	rows [][]string
}

func (c *CSVTable) CreateStatement() string {
	return fmt.Sprintf("CREATE TABLE %s (%s)", c.name, c.headerString())
}

func (c *CSVTable) headerString() string {
	return strings.Join(c.headers(), ", ")
}

func (c *CSVTable) headers() []string {
	return c.rows[0]
}

type CSVDatabase struct {
	sqliteDb *sql.DB
}

func New() (c *CSVDatabase, err error) {
	c = &CSVDatabase{}
	// Open an in-memory sqlite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return
	}
	c.sqliteDb = db
	return
}

func (c *CSVDatabase) Insert(tableName string, records [][]string) (err error) {
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
	sqlStatement := fmt.Sprintf("CREATE TABLE %s (%s)", tableName, headerString)
	_, err = c.sqliteDb.Exec(sqlStatement)
	if err != nil {
		return
	}

	// Insert CSV data into sqlite db
	tx, err := c.sqliteDb.Begin()
	if err != nil {
		return
	}
	sqlStatement = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, headerString, rowQuestionMarks)
	stmt, err := tx.Prepare(sqlStatement)
	if err != nil {
		return
	}
	defer stmt.Close()
	for _, row := range csvRows {
		// Turn row into a slice of interface{}
		rowCopy := make([]interface{}, len(row))
		for i, d := range row {
			rowCopy[i] = d
		}
		_, err = stmt.Exec(rowCopy...)
		if err != nil {
			return
		}
	}
	tx.Commit()
	return
}

func (c *CSVDatabase) Close() {
	c.sqliteDb.Close()
}

func (c *CSVDatabase) Query(query string) (result [][]string, err error) {
	// Get the data back out of the CSV
	sqlRows, err := c.sqliteDb.Query(query)
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
