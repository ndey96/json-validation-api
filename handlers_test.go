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
  uploadValidSchema(t)
  downloadValidSchema(t)
  downloadInvalidSchema(t)
}

func TestUploadSchema(t *testing.T) {
  // setup
  db := OpenDBConn()
  defer db.Close()
  ResetTestSchemas(db, t)
  uploadValidSchema(t)
  uploadDuplicateValidSchema(t)
  uploadInvalidSchema(t)
}

func TestValidateDocument(t *testing.T) {
  db := OpenDBConn()
  defer db.Close()
  ResetTestSchemas(db, t)
  uploadValidSchema(t)
  validateValidDocument(t)
  validateInvalidDocument(t)
}

func downloadValidSchema(t *testing.T) {
  r := mux.NewRouter()
  r.HandleFunc("/schema/{schemaId}", http.HandlerFunc(DownloadSchema))
  req, err := http.NewRequest("GET", "/schema/test-schema", nil)
  FailIf(err, t)
  w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
  ExpectValue(http.StatusOK, w.Code, "status", t)
  rawSchema, err := ioutil.ReadFile("./json/test-schema.json")
  FailIf(err, t)
  expectedBody := string(rawSchema)
  ExpectValue(expectedBody, w.Body.String(), "body", t)
}

func downloadInvalidSchema(t *testing.T) {
  r := mux.NewRouter()
  r.HandleFunc("/schema/{schemaId}", http.HandlerFunc(DownloadSchema))
  req, err := http.NewRequest("GET", "/schema/does-not-exist", nil)
  FailIf(err, t)
  w := httptest.NewRecorder()
  r.ServeHTTP(w, req)
  ExpectValue(http.StatusNotFound, w.Code, "status", t)
  expectedBody := `{"action":"downloadSchema","id":"does-not-exist","status":"error","message":"No schema found"}`
  ExpectValue(expectedBody, strings.Trim(w.Body.String(), "\n"), "body", t)
}

func uploadValidSchema(t *testing.T) {
  r := mux.NewRouter()
  r.HandleFunc("/schema/{schemaId}", http.HandlerFunc(UploadSchema))
  rawSchema, err := ioutil.ReadFile("./json/test-schema.json")
  FailIf(err, t)
  req, err := http.NewRequest("POST", "/schema/test-schema", bytes.NewBuffer(rawSchema))
  FailIf(err, t)
  w := httptest.NewRecorder()
  r.ServeHTTP(w, req)
  ExpectValue(http.StatusCreated, w.Code, "status", t)
  expectedBody := `{"action":"uploadSchema","id":"test-schema","status":"success"}`
  ExpectValue(expectedBody, strings.Trim(w.Body.String(), "\n"), "body", t)
}

func uploadDuplicateValidSchema(t *testing.T) {
  r := mux.NewRouter()
  r.HandleFunc("/schema/{schemaId}", http.HandlerFunc(UploadSchema))
  rawSchema, err := ioutil.ReadFile("./json/test-schema.json")
  FailIf(err, t)
  req, err := http.NewRequest("POST", "/schema/test-schema", bytes.NewBuffer(rawSchema))
  FailIf(err, t)
  w := httptest.NewRecorder()
  r.ServeHTTP(w, req)
  ExpectValue(http.StatusBadRequest, w.Code, "status", t)
  expectedBody := `{"action":"uploadSchema","id":"test-schema","status":"error","message":"Schema ID already in use"}`
  ExpectValue(expectedBody, strings.Trim(w.Body.String(), "\n"), "body", t)
}

func uploadInvalidSchema(t *testing.T) {
  r := mux.NewRouter()
  r.HandleFunc("/schema/{schemaId}", http.HandlerFunc(UploadSchema))
  rawSchema, err := ioutil.ReadFile("./json/invalid-schema.json")
  FailIf(err, t)
  req, err := http.NewRequest("POST", "/schema/test-schema", bytes.NewBuffer(rawSchema))
  FailIf(err, t)
  w := httptest.NewRecorder()
  r.ServeHTTP(w, req)
  ExpectValue(http.StatusUnprocessableEntity, w.Code, "status", t)
  expectedBody := `{"action":"uploadSchema","id":"test-schema","status":"error","message":"Invalid JSON provided"}`
  ExpectValue(expectedBody, strings.Trim(w.Body.String(), "\n"), "body", t)
}

func validateValidDocument(t *testing.T) {
  r := mux.NewRouter()
  r.HandleFunc("/validate/{schemaId}", http.HandlerFunc(ValidateDocument))
  rawDocument, err := ioutil.ReadFile("./json/test.json")
  FailIf(err, t)
  req, err := http.NewRequest("POST", "/validate/test-schema", bytes.NewBuffer(rawDocument))
  FailIf(err, t)
  w := httptest.NewRecorder()
  r.ServeHTTP(w, req)
  ExpectValue(http.StatusOK, w.Code, "status", t)
  expectedBody := `{"action":"validateDocument","id":"test-schema","status":"success"}`
  ExpectValue(expectedBody, strings.Trim(w.Body.String(), "\n"), "body", t)
}

func validateInvalidDocument(t *testing.T) {
  r := mux.NewRouter()
  r.HandleFunc("/validate/{schemaId}", http.HandlerFunc(ValidateDocument))
  rawDocument, err := ioutil.ReadFile("./json/test2.json")
  FailIf(err, t)
  req, err := http.NewRequest("POST", "/validate/test-schema", bytes.NewBuffer(rawDocument))
  FailIf(err, t)
  w := httptest.NewRecorder()
  r.ServeHTTP(w, req)
  ExpectValue(http.StatusBadRequest, w.Code, "status", t)
  expectedBody := `{"action":"validateDocument","id":"test-schema","status":"error","message":"Document does not conform to schema"}`
  ExpectValue(expectedBody, strings.Trim(w.Body.String(), "\n"), "body", t)
}
