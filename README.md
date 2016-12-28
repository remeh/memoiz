# Setup

psql -U postgres
# drop database scratche; drop role scratche;
# \ir resources/schema.sql; \ir resources/bing.sql; \ir resources/kg.sql
go generate ./...
go build

# Tests

