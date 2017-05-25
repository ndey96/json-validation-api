package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "io/ioutil"
    "io"
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
    s := StorageRetrieveSchema(vars["schemaId"])
    if len(s.Id) == 0  || len(s.Schema) == 0 {
      w.WriteHeader(http.StatusNotFound)
      res := ResponseWithMessage{Action: action, Status: "error", Id: id, Message: "No schema found"}
      err := json.NewEncoder(w).Encode(res)
      PanicIf(err)
      return
    }
    w.WriteHeader(http.StatusOK)
    err := json.NewEncoder(w).Encode(s.Schema)
    PanicIf(err)
}

func UploadSchema(w http.ResponseWriter, r *http.Request) {
  action := "uploadSchema"
  vars := mux.Vars(r)
  id := vars["schemaId"]
  schema := Schema{Id: id}
  body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
  PanicIf(err)
  err = r.Body.Close()
  PanicIf(err)
  var js map[string]interface{}
  if err := json.Unmarshal(body, &js); err != nil {
    w.WriteHeader(http.StatusUnprocessableEntity)
    res := ResponseWithMessage{Action: action, Status: "error", Id: id, Message: "Invalid JSON provided"}
    err = json.NewEncoder(w).Encode(res)
    PanicIf(err)
    return
  }
  schema.Schema = string(body[:])
  err = StorageWriteSchema(schema)
  if (err != nil) {
    w.WriteHeader(http.StatusBadRequest)
    res := ResponseWithMessage{Action: action, Status: "error", Id: id, Message: err.Error()}
    err = json.NewEncoder(w).Encode(res)
    PanicIf(err)
    return
  }
  w.WriteHeader(http.StatusCreated)
  res := Response{Action: action, Status: "success", Id: id}
  err = json.NewEncoder(w).Encode(res)
  PanicIf(err)
}

func ValidateDocument(w http.ResponseWriter, r *http.Request) {
  action := "validateDocument"
  vars := mux.Vars(r)
  schemaId := vars["schemaId"]

  document, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
  PanicIf(err)
  err = r.Body.Close()
  PanicIf(err)

  var js map[string]interface{}
  if err := json.Unmarshal(document, &js); err != nil {
    w.WriteHeader(http.StatusUnprocessableEntity)
    res := ResponseWithMessage{Action: action, Status: "error", Id: schemaId, Message: "Document is not valid JSON"}
    err = json.NewEncoder(w).Encode(res)
    PanicIf(err)
    return
  }
  schema := StorageRetrieveSchema(vars["schemaId"])
  documentLoader := gojsonschema.NewStringLoader(string(document[:]))
  schemaLoader := gojsonschema.NewStringLoader(schema.Schema)
  result, err := gojsonschema.Validate(schemaLoader, documentLoader)
  PanicIf(err)

  if result.Valid() {
    w.WriteHeader(http.StatusOK)
    res := Response{Action: action, Status: "success", Id: schemaId}
    err = json.NewEncoder(w).Encode(res)
    PanicIf(err)
  } else {
    w.WriteHeader(http.StatusBadRequest)
    res := ResponseWithMessage{Action: action, Status: "error", Id: schemaId, Message: "Document does not conform to schema"}
    err = json.NewEncoder(w).Encode(res)
    PanicIf(err)
  }
}
