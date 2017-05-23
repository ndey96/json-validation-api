package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gorilla/mux"
    "strings"
)

func init() {
  isTesting = true
}

func TestDownloadSchema(t *testing.T) {
  db := OpenDBConn()
  defer db.Close()
  _, err := db.Exec("DROP TABLE IF EXISTS schemas")
  FailIf(err, t)
  _, err = db.Exec("CREATE TABLE schemas(id TEXT UNIQUE, schema TEXT)")
  FailIf(err, t)
  _, err = db.Exec("INSERT INTO schemas (id, schema) VALUES('potato', 'tomato')")
  FailIf(err, t)

  req, err := http.NewRequest("GET", "/schema/potato", nil)
  FailIf(err, t)
  w := httptest.NewRecorder()
  handler := http.HandlerFunc(DownloadSchema)
  m := mux.NewRouter()
	m.HandleFunc("/schema/{schemaId}", handler)
	m.ServeHTTP(w, req)

  if status := w.Code; status != http.StatusOK {
      t.Errorf("handler returned wrong status code: got %v want %v",
          status, http.StatusOK)
  }
  expected := `"tomato"`
  if strings.Trim(w.Body.String(), "\n") != expected {
      t.Errorf("handler returned unexpected body: got %v want %v",
          w.Body.String(), expected)
  }
}
