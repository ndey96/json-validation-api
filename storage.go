package main

import (
    "database/sql"
    "fmt"
    "errors"
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
  var id, schema string
  err = db.QueryRow("SELECT id, schema FROM schemas WHERE id=$1", schemaId).Scan(&id,&schema)
  if err == sql.ErrNoRows {
    return Schema{}
  }
  PanicIf(err)
  return Schema{Id:id, Schema:schema}
}

func StorageWriteSchema(s Schema) error {
  dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
  db, err := sql.Open("postgres", dbinfo)
  PanicIf(err)
  err = db.Ping()
  PanicIf(err)
  defer db.Close()
  var id string
  err = db.QueryRow("SELECT id FROM schemas WHERE id=$1", s.Id).Scan(&id)
  if err == sql.ErrNoRows {
    _, err = db.Exec("INSERT INTO schemas(id,schema) VALUES($1, $2)", s.Id, s.Schema)
    return err
  } else if err != nil && err != sql.ErrNoRows {
    panic(err)
  } else {
    return errors.New("Schema ID already in use")
  }
}
