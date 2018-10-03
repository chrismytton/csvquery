# CSV Query

Query CSV files with SQL.

## Install

    go get -u github.com/chrismytton/csvquery/cmd/csvquery

## Usage

Given two CSV files:

**people.csv**

| id | name  |
| -- | ----- |
| 1  | Alice |
| 2  | Bob   |

**ages.csv**

| person_id | age |
| --------- | --- |
| 1         | 20  |
| 2         | 30  |

You can then query these two files like so:

    csvquery -table people:people.csv -table ages:ages.csv \
      -query 'select id,name,age from people join ages on ages.person_id = people.id'

- `-table` - In the form `table_name:file.csv`. This can be specified multiple times to use multiple tables.
- `-query` - The SQL query to run against the specified tables.

## Server version

There is also a very experimental HTTP server, which you can install with the following:

    go get -u github.com/chrismytton/csvquery/cmd/csvquery-server

This accepts one or more `?table` parameters, but instead of pointing to files they point to URLs, and then a `?query` parameter including the query. These parameters will need to be URL-encoded.
