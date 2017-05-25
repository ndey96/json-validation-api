package jsonvalidation

import (
  "database/sql"
  "testing"
)

func PanicIf(err error) {
    if err != nil {
        panic(err)
    }
}

func FailIf(err error, t *testing.T) {
  if err != nil {
    t.Fatal(err)
  }
}

func ExpectValue(expected interface{}, actual interface{}, valueType string, t *testing.T) {
  if expected != actual {
    t.Errorf("Got unexpected %v:\nactual: %v\nexpected: %v", valueType, actual, expected)
  }
}

func ResetTestSchemas(db *sql.DB, t *testing.T) {
  _, err := db.Exec(DROP_SCHEMAS_TABLE_IF_EXISTS_QUERY)
  FailIf(err, t)
  _, err = db.Exec(CREATE_SCHEMAS_TABLE_QUERY)
  FailIf(err, t)
}
