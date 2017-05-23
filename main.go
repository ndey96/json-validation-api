package main

import (
    "log"
    "net/http"
)

var isTesting = false

func main() {
    router := NewRouter()
    log.Fatal(http.ListenAndServe(":8080", router))
}
