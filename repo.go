package main

import "fmt"

var currentId int

var schemas Schemas

// Give us some seed data
func init() {
  RepoCreateSchema(Schema{Id:"swag", Schema:[]byte(`{"Name":"Alice","Body":"Hello","Time":1294706395881547000}`)})
}

func RepoFindSchema(id string) Schema {
    for _, s := range schemas {
        if s.Id == id {
            return s
        }
    }
    // return empty Schema if not found
    return Schema{}
}

func RepoCreateSchema(s Schema) {
    schemas = append(schemas, s)
}

func RepoDestroySchema(id string) error {
    for i, s := range schemas {
        if s.Id == id {
            schemas = append(schemas[:i], schemas[i+1:]...)
            return nil
        }
    }
    return fmt.Errorf("Could not find Schema with id of %d to delete", id)
}
