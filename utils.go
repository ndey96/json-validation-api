package main

func PanicIf(err error) {
    if err != nil {
        panic(err)
    }
}