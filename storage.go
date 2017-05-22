package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

const (
    DB_USER     = "postgres"
    DB_PASSWORD = "postgres"
    DB_NAME     = "jva"
)

func StorageRetrieveSchema(schemaId string) Schema {
  dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
  db, err := sql.Open("postgres", dbinfo)
  PanicIf(err)
  err = db.Ping()
  PanicIf(err)
  defer db.Close()
  var id string
  var schema string
  err = db.QueryRow("SELECT id, schema FROM schemas WHERE id=$1", schemaId).Scan(&id,&schema)
  if err == sql.ErrNoRows {
    return Schema{}
  }
  PanicIf(err)
  return Schema{Id:id, Schema:schema}
}

func StorageWriteSchema(s Schema) {
  dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
  db, err := sql.Open("postgres", dbinfo)
  defer db.Close()
  PanicIf(err)
  _, err = db.Exec("INSERT INTO schemas(id,schema) VALUES($1, $2)", s.Id, s.Schema)
  PanicIf(err)
}

func PanicIf(err error) {
    if err != nil {
        panic(err)
    }
}
