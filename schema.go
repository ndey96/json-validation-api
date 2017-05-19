package main

type Schema struct {
    Id      string    `json:"id"`
    Schema  []byte    `json:"schema"`
}

type Schemas []Schema
