package main

type Schema struct {
    Id      string    `json:"id"`
    Schema  string    `json:"schema"`
}

type Schemas []Schema
