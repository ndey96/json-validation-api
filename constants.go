package main

const (
    APP_DB_CONN_STR = "user=postgres password=postgres dbname=jva sslmode=disable"
    TEST_DB_CONN_STR = "user=postgres password=postgres dbname=jva_test sslmode=disable"
    CREATE_SCHEMAS_TABLE_QUERY = "CREATE TABLE schemas(id TEXT UNIQUE, schema TEXT)"
    CREATE_SCHEMAS_TABLE_IF_NOT_EXISTS_QUERY = "CREATE TABLE IF NOT EXISTS schemas(id TEXT UNIQUE, schema TEXT)"
    DROP_SCHEMAS_TABLE_IF_EXISTS_QUERY = "DROP TABLE IF EXISTS schemas"
)
