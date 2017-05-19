package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "io"
    "io/ioutil"
    "github.com/gorilla/mux"
    "github.com/xeipuuv/gojsonschema"
)

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Server is alive!")
}

func DownloadSchema(w http.ResponseWriter, r *http.Request) {
    action := "downloadSchema"
    vars := mux.Vars(r)
    id := vars["schemaId"]
    s := RepoFindSchema(vars["schemaId"])
    if len(s.Id) == 0 {
      w.Header().Set("Content-Type", "application/json; charset=UTF-8")
      w.WriteHeader(http.StatusNotFound)
      res := ResponseWithMessage{Action: action, Status: "error", Id: id, Message: "No schema found"}
      if err := json.NewEncoder(w).Encode(res); err != nil {
        panic(err)
      }
      return
    }
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(s); err != nil {
        panic(err)
    }
}

func UploadSchema(w http.ResponseWriter, r *http.Request) {
  action := "uploadSchema"
  vars := mux.Vars(r)
  id := vars["schemaId"]
  schema := Schema{Id: id}
  body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
  if err != nil {
    panic(err)
  }
  if err := r.Body.Close(); err != nil {
    panic(err)
  }
  var js map[string]interface{}
  if err := json.Unmarshal(body, &js); err != nil {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusUnprocessableEntity)
    res := ResponseWithMessage{Action: action, Status: "error", Id: id, Message: "Invalid JSON provided"}
    if err := json.NewEncoder(w).Encode(res); err != nil {
      panic(err)
    }
    return
  }
  schema.Schema = body
  RepoCreateSchema(schema)
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusCreated)
  res := Response{Action: action, Status: "success", Id: id}
  if err := json.NewEncoder(w).Encode(res); err != nil {
    panic(err)
  }
}

func ValidateDocument(w http.ResponseWriter, r *http.Request) {
  action := "validateDocument"
  vars := mux.Vars(r)
  schemaId := vars["schemaId"]

  document, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
  if err != nil {
    panic(err)
  }
  if err := r.Body.Close(); err != nil {
    panic(err)
  }
  var js map[string]interface{}
  if err := json.Unmarshal(document, &js); err != nil {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusUnprocessableEntity)
    res := ResponseWithMessage{Action: action, Status: "error", Id: schemaId, Message: "Document is not valid JSON"}
    if err := json.NewEncoder(w).Encode(res); err != nil {
      panic(err)
    }
    return
  }
  schema := RepoFindSchema(vars["schemaId"])
  schemaLoader := gojsonschema.NewStringLoader(string(document[:]))
  documentLoader := gojsonschema.NewStringLoader(string(schema.Schema[:]))
  result, err := gojsonschema.Validate(schemaLoader, documentLoader)
  if err != nil {
    panic(err)
  }

  if result.Valid() {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    res := Response{Action: action, Status: "success", Id: schemaId}
    if err := json.NewEncoder(w).Encode(res); err != nil {
      panic(err)
    }
    return
  } else {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusUnprocessableEntity)
    res := ResponseWithMessage{Action: action, Status: "error", Id: schemaId, Message: "Document does not conform to schema"}
    if err := json.NewEncoder(w).Encode(res); err != nil {
      panic(err)
    }
    return
  }
}
