## Build

Build the binary
```sh
$ go build -v -x main.go
```

## SQL Queries

Uses `sqlc` to generate go from sql files.

`sqlc` config is in sqlc.yaml.

Generate the `internal/queries` package:
`sqlc generate`

Queries are defined in query.sql.
Queries are automatically converted to functions.

Tables are defined by the `db/migrations` using dbmate. Run `sqlc generate` to re-generate the go code after any migration alterations.
