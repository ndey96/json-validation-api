package main

import (
    "database/sql"
    "errors"
    _ "github.com/lib/pq"
)

func init() {
  createTableIfNotExist()
}

func createTableIfNotExist() {
  db := OpenDBConn()
  defer db.Close()
  _, err := db.Exec(CREATE_SCHEMAS_TABLE_IF_NOT_EXISTS_QUERY)
  PanicIf(err)
}

func OpenDBConn() *sql.DB {
  dbConnStr := APP_DB_CONN_STR
  if isTesting {
    dbConnStr = TEST_DB_CONN_STR
  }
  db, err := sql.Open("postgres", dbConnStr)
  PanicIf(err)
  err = db.Ping()
  PanicIf(err)
  return db
}

func StorageRetrieveSchema(schemaId string) (Schema, error) {
  db := OpenDBConn()
  defer db.Close()
  var id, schema string
  err := db.QueryRow("SELECT id, schema FROM schemas WHERE id=$1", schemaId).Scan(&id,&schema)
  if err == sql.ErrNoRows {
    err = errors.New("Schema not found")
    return Schema{}, err
  } else if err != nil {
    return Schema{}, err
  }
  return Schema{Id:id, Schema:schema}, nil
}

func StorageWriteSchema(s Schema) error {
  db := OpenDBConn()
  defer db.Close()
  var id string
  err := db.QueryRow("SELECT id FROM schemas WHERE id=$1", s.Id).Scan(&id)
  if err == sql.ErrNoRows {
    _, err = db.Exec("INSERT INTO schemas(id,schema) VALUES($1, $2)", s.Id, s.Schema)
    return err
  } else if err != nil && err != sql.ErrNoRows {
    panic(err)
  } else {
    return errors.New("Schema ID already in use")
  }
}
