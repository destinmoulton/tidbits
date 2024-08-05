## SQL Queries

Uses `sqlc` to generate go from sql files.

`sqlc` config is in sqlc.yaml.

Tables are defined in schema.sql
Queries are defined in query.sql

Generate the `internal/queries`:
`sqlc generate`