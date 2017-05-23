package main

import (
  "testing"
)

func init() {
  isTesting = true
}

func setupTestStorageRetrieveSchema(t *testing.T) func(t *testing.T) {
  db := OpenDBConn()
  _, err := db.Exec("DROP TABLE IF EXISTS schemas")
  FailIf(err, t)
  _, err = db.Exec("CREATE TABLE schemas(id TEXT UNIQUE, schema TEXT)")
  FailIf(err, t)
  _, err = db.Exec("INSERT INTO schemas (id, schema) VALUES('potato', 'tomato')")
  FailIf(err, t)
	return func(t *testing.T) {
    db.Close()
	}
}

func TestStorageRetrieveSchema(t *testing.T) {
  teardown := setupTestStorageRetrieveSchema(t)
  defer teardown(t)
  result := StorageRetrieveSchema("potato")
  if (result.Id != "potato" || result.Schema != "tomato") {
    t.Fatal()
  }
}

func setupTestStorageWriteSchema(t *testing.T) {
  db := OpenDBConn()
  defer db.Close()
  _, err := db.Exec("DROP TABLE IF EXISTS schemas")
  FailIf(err, t)
  _, err = db.Exec("CREATE TABLE schemas(id TEXT UNIQUE, schema TEXT)")
  FailIf(err, t)
}

func TestStorageWriteSchema(t *testing.T) {
  setupTestStorageWriteSchema(t)
  testSchema := Schema{"potato", "tomato"}
  err := StorageWriteSchema(testSchema)
  FailIf(err, t)
  db := OpenDBConn()
  defer db.Close()
  var id, schema string
  db.QueryRow("SELECT id, schema FROM schemas WHERE id='potato'").Scan(&id, &schema)
  if (id != "potato" || schema != "tomato") {
    t.Fatal()
  }
  err = StorageWriteSchema(testSchema)
  if err.Error() != "Schema ID already in use" {
    t.Fatal(err)
  }
}
