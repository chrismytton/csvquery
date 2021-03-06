**WORK IN PROGRESS - there will be things that are broken or don't work**

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

Which should give you the following result:

| id | name  | age |
| -- | ----- | --- |
| 1  | Alice | 20  |
| 2  | Bob   | 30  |

See `csvquery -help` for help.

## Server version

There is also a very experimental HTTP server, which you can install with the following:

    go get -u github.com/chrismytton/csvquery/cmd/csvquery-server

This accepts one or more `?table` parameters, but instead of pointing to files they point to URLs, and then a `?query` parameter including the query. These parameters will need to be URL-encoded.

## Compiling for Windows on macOS

I had to do the following to get things cross compiling to Windows from macOS.

    brew install mingw-w64

    GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc go build -v github.com/chrismytton/csvquery/cmd/csvquery
