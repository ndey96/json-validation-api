package main

import (
  "testing"
)

func init() {
  isTesting = true
}

func TestStorageRetrieveSchema(t *testing.T) {
  db := OpenDBConn()
  defer db.Close()
  ResetTestSchemas(db, t)
  _, err := db.Exec("INSERT INTO schemas (id, schema) VALUES('potato', 'tomato')")
  FailIf(err, t)
  result := StorageRetrieveSchema("potato")
  if (result.Id != "potato" || result.Schema != "tomato") {
    t.Fatal()
  }
}

func TestStorageWriteSchema(t *testing.T) {
  db := OpenDBConn()
  defer db.Close()
  ResetTestSchemas(db, t)
  testSchema := Schema{"potato", "tomato"}
  err := StorageWriteSchema(testSchema)
  FailIf(err, t)
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
