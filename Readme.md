# json-validation-api

API written in go which allows you to upload/download JSON schemas and validate JSON documents against the schemas you uploaded.
Uses [gojsonschema](https://github.com/xeipuuv/gojsonschema) for document validation and [PostgreSQL](https://www.postgresql.org) for storage.

### Setup:

```
createdb -h localhost -p 5432 -U postgres -W jva // password=postgres
createdb -h localhost -p 5432 -U postgres -W jva_test // password=postgres
go get github.com/ndey96/jsonvalidation
cd $GOPATH/src/github.com/ndey96/jsonvalidation
go install && jsonvalidation // runs server on localhost:8080
```

### Endpoints

```
POST    /schema/SCHEMAID        - Upload JSON Schema with unique `SCHEMAID`
GET     /schema/SCHEMAID        - Download JSON Schema with unique `SCHEMAID`

POST    /validate/SCHEMAID      - Validate JSON document against the JSON Schema identified by `SCHEMAID`
```

### Sample Requests

- Upload Schema: `curl localhost:8080/schema/test-schema -X POST -d @json/test-schema.json`
- Download Schema: `curl localhost:8080/schema/test-schema`
- Validate Document `curl localhost:8080/validate/test-schema -X POST -d @json/test.json`

### Sample JSON
- [test-schema.json](https://github.com/ndey96/jsonvalidation/blob/master/json/test-schema.json)
- [test.json](https://github.com/ndey96/jsonvalidation/blob/master/json/test.json)
