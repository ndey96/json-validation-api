package main

import "net/http"

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
    Route{
        "Index",
        "GET",
        "/",
        Index,
    },
    Route{
        "SchemaIndex",
        "GET",
        "/schemas",
        SchemaIndex,
    },
    Route{
        "SchemaShow",
        "GET",
        "/schemas/{schemaId}",
        SchemaShow,
    },
        Route{
        "SchemaCreate",
        "POST",
        "/schemas",
        SchemaCreate,
    },
}
