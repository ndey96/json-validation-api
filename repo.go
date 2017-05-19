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

func RepoCreateSchema(s Schema) (sch Schema, err error) {
    if len(s.Id) == 0 && len(s.Schema) == 0 {
      sch = Schema{Id: "BAD"}
      err = fmt.Errorf("Error creating schema")
      return
    }
    schemas = append(schemas, s)
    sch = s
    return
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
