# CSV Query

Query CSV files with SQL.

## Install

    go get -u github.com/chrismytton/csvquery

## Usage

Make a GET request to the server, providing the following parameters.

- `table` - In the form `table_name:url`. This can be specified multiple times to use multiple tables.
- `query` - The SQL query to run against the specified tables.
