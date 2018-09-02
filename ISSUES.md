# Should be able to use URLs from CLI tool

At the moment it just works with local files. Should work with URLs to CSV files as well.

# Hangs when run without arguments and starts to use up lots of CPU

Should print a friendly error message and exit when run without arguments.

# Everything is a string

- _Everything_ is a string, which means you need to do things like `SELECT * FROM data WHERE id = "42"`, because the `id` field will be a string.

Or `SELECT * FROM data WHERE CAST(age as integer) > 22`.

# Need a license
