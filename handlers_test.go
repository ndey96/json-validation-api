package main

import (
    "net/http"
    "net/http/httptest"
    "io/ioutil"
    "testing"
    "github.com/gorilla/mux"
    "strings"
    "bytes"
)

func init() {
  isTesting = true
}

func TestDownloadSchema(t *testing.T) {
  // setup
  db := OpenDBConn()
  defer db.Close()
  ResetTestSchemas(db, t)
  _, err := db.Exec("INSERT INTO schemas (id, schema) VALUES('potato', 'tomato')")
  FailIf(err, t)
  req, err := http.NewRequest("GET", "/schema/potato", nil)
  FailIf(err, t)
  w := httptest.NewRecorder()
  handler := http.HandlerFunc(DownloadSchema)
  m := mux.NewRouter()
	m.HandleFunc("/schema/{schemaId}", handler)
	m.ServeHTTP(w, req)

  ExpectValue(http.StatusOK, w.Code, "status", t)
  ExpectValue(`"tomato"`, strings.Trim(w.Body.String(), "\n"), "body", t)
}

func TestUploadSchema(t *testing.T) {
  // setup
  db := OpenDBConn()
  defer db.Close()
  ResetTestSchemas(db, t)
  handler := http.HandlerFunc(UploadSchema)
  m := mux.NewRouter()
  m.HandleFunc("/schema/{schemaId}", handler)
  rawSchema, err := ioutil.ReadFile("./json/test-schema.json")
  FailIf(err, t)

  // upload a schema
  req, err := http.NewRequest("POST", "/schema/test-schema", bytes.NewBuffer(rawSchema))
  FailIf(err, t)
  w := httptest.NewRecorder()
  m.ServeHTTP(w, req)
  ExpectValue(http.StatusCreated, w.Code, "status", t)
  expectedBody := `{"action":"uploadSchema","id":"test-schema","status":"success"}`
  ExpectValue(expectedBody, strings.Trim(w.Body.String(), "\n"), "body", t)

  // upload a duplicate schema
  req, err = http.NewRequest("POST", "/schema/test-schema", bytes.NewBuffer(rawSchema))
  FailIf(err, t)
  w = httptest.NewRecorder()
  m.ServeHTTP(w, req)
  ExpectValue(http.StatusBadRequest, w.Code, "status", t)
  expectedBody = `{"action":"uploadSchema","id":"test-schema","status":"error","message":"Schema ID already in use"}`
  ExpectValue(expectedBody, strings.Trim(w.Body.String(), "\n"), "body", t)

  // upload an invalid schema
  rawSchema, err = ioutil.ReadFile("./json/invalid-schema.json")
  FailIf(err, t)
  req, err = http.NewRequest("POST", "/schema/test-schema", bytes.NewBuffer(rawSchema))
  FailIf(err, t)
  w = httptest.NewRecorder()
  m.ServeHTTP(w, req)
  ExpectValue(http.StatusUnprocessableEntity, w.Code, "status", t)
  expectedBody = `{"action":"uploadSchema","id":"test-schema","status":"error","message":"Invalid JSON provided"}`
  ExpectValue(expectedBody, strings.Trim(w.Body.String(), "\n"), "body", t)
}

func TestValidateDocument(t *testing.T) {
  // setup
  db := OpenDBConn()
  defer db.Close()
  ResetTestSchemas(db, t)
  m := mux.NewRouter()

  // upload a schema
  handler := http.HandlerFunc(UploadSchema)
  m.HandleFunc("/schema/{schemaId}", handler)
  rawSchema, err := ioutil.ReadFile("./json/test-schema.json")
  FailIf(err, t)
  req, err := http.NewRequest("POST", "/schema/test-schema", bytes.NewBuffer(rawSchema))
  FailIf(err, t)
  w := httptest.NewRecorder()
  m.ServeHTTP(w, req)
  ExpectValue(http.StatusCreated, w.Code, "status", t)
  expectedBody := `{"action":"uploadSchema","id":"test-schema","status":"success"}`
  ExpectValue(expectedBody, strings.Trim(w.Body.String(), "\n"), "body", t)


  handler = http.HandlerFunc(ValidateDocument)
  m.HandleFunc("/validate/{schemaId}", handler)

  // validate valid document
  rawDocument, err := ioutil.ReadFile("./json/test.json")
  FailIf(err, t)
  req, err = http.NewRequest("POST", "/validate/test-schema", bytes.NewBuffer(rawDocument))
  FailIf(err, t)
  w = httptest.NewRecorder()
  m.ServeHTTP(w, req)
  ExpectValue(http.StatusOK, w.Code, "status", t)
  expectedBody = `{"action":"validateDocument","id":"test-schema","status":"success"}`
  ExpectValue(expectedBody, strings.Trim(w.Body.String(), "\n"), "body", t)

  // invalidate invalid document
  rawDocument, err = ioutil.ReadFile("./json/test2.json")
  FailIf(err, t)
  req, err = http.NewRequest("POST", "/validate/test-schema", bytes.NewBuffer(rawDocument))
  FailIf(err, t)
  w = httptest.NewRecorder()
  m.ServeHTTP(w, req)
  ExpectValue(http.StatusBadRequest, w.Code, "status", t)
  expectedBody = `{"action":"validateDocument","id":"test-schema","status":"error","message":"Document does not conform to schema"}`
  ExpectValue(expectedBody, strings.Trim(w.Body.String(), "\n"), "body", t)
}
