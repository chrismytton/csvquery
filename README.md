# CSVSQL

Query CSV files with SQL.

## Install

    go get -u github.com/chrismytton/csvsql

## Known limitations

- _Everything_ is a string, which means you need to do things like `SELECT * FROM data WHERE id = "42"`, because the `id` field will be a string.
