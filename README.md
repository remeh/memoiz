# Setup

psql -U postgres
# drop database memoiz; drop role memoiz;
# \ir resources/schema.sql; \ir resources/bing.sql; \ir resources/kg.sql
go generate ./...
go build

# Tests

