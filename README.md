# Notes

## Sqlc
### 1. All of sql files are stored under */db/queries*
### 2. After writing sql use `sqlc generate` cli to compile sql to go

## Migrations
### 1. When wanting to introduce new things(constraints, tables, ...) that involved the db use the `migrate create -ext sql -dir db/migrations -seq <Name of the migration>` to create 2 files under */db/migrations*
### 2. After creating migrations file add the changes that are needed in the up file
### 3. If wanting to test on a db local use `migrate -path db/migrations -database <connection_string> -verbose up` to migrate first
### 4. When wanting to drop everything use `migrate -path db/migrations -database <connection_string> -verbose drop`
