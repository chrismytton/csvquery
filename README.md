# CSV Query

Query CSV files with SQL.

## Install

    go get -u github.com/chrismytton/csvquery

## Usage

Make a GET request to the server, providing the following parameters.

- `table` - In the form `table_name:url`. This can be specified multiple times to use multiple tables.
- `query` - The SQL query to run against the specified tables.

## Known limitations

- _Everything_ is a string, which means you need to do things like `SELECT * FROM data WHERE id = "42"`, because the `id` field will be a string.
