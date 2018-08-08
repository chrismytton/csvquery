package csvquery

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

// CSVDatabase represents a query-able CSV backed database.
type CSVDatabase struct {
	sqliteDb *sql.DB
}

// New returns a *CSVDatabase.
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

// Insert populates the given table with the CSV records provided.
func (c *CSVDatabase) Insert(tableName string, records [][]string) (err error) {
	table := &csvTable{tableName, records}
	_, err = c.sqliteDb.Exec(table.createStatement())
	if err != nil {
		return
	}

	// Insert CSV data into sqlite db
	tx, err := c.sqliteDb.Begin()
	if err != nil {
		return
	}
	stmt, err := tx.Prepare(table.insertStatement())
	if err != nil {
		return
	}
	defer stmt.Close()
	for _, row := range table.rows() {
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

// Close the underlying SQLite database.
func (c *CSVDatabase) Close() {
	c.sqliteDb.Close()
}

// Query takes a string of SQL and runs it against the current database, returning the result.
func (c *CSVDatabase) Query(query string) (result [][]string, err error) {
	sqlRows, err := c.sqliteDb.Query(query)
	if err != nil {
		return
	}
	defer sqlRows.Close()

	colNames, err := sqlRows.Columns()
	if err != nil {
		return
	}
	result = append(result, colNames)

	readCols := make([]interface{}, len(colNames))
	writeCols := make([]string, len(colNames))
	for i := range writeCols {
		readCols[i] = &writeCols[i]
	}
	for sqlRows.Next() {
		err = sqlRows.Scan(readCols...)
		if err != nil {
			return
		}
		cols := make([]string, len(writeCols))
		copy(cols, writeCols)
		result = append(result, cols)
	}
	if err = sqlRows.Err(); err != nil {
		return
	}
	return
}

type csvTable struct {
	name    string
	records [][]string
}

func (c *csvTable) createStatement() string {
	return fmt.Sprintf("CREATE TABLE '%s' (%s)", c.name, c.headerString())
}

func (c *csvTable) insertStatement() string {
	return fmt.Sprintf("INSERT INTO '%s' (%s) VALUES (%s)", c.name, c.headerString(), c.rowQuestionMarks())
}

func (c *csvTable) headerString() string {
	return strings.Join(c.headers(), ", ")
}

func (c *csvTable) headers() []string {
	return c.records[0]
}

func (c *csvTable) rows() [][]string {
	return c.records[1:]
}

func (c *csvTable) rowQuestionMarks() string {
	// Create the correct number of placeholder question marks for using in the
	// prepared statement.
	questionMarks := make([]string, len(c.headers()))
	for i := 0; i < len(questionMarks); i++ {
		questionMarks[i] = "?"
	}
	return strings.Join(questionMarks, ", ")
}
