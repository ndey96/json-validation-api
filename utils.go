package main

import "testing"

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
