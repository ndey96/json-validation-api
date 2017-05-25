package jsonvalidation

type Schema struct {
    Id      string    `json:"id"`
    Schema  string    `json:"schema"`
}

type Response struct {
  Action string `json:"action"`
  Id string `json:"id"`
  Status string `json:"status"`
}

type ResponseWithMessage struct {
  Action string `json:"action"`
  Id string `json:"id"`
  Status string `json:"status"`
  Message string `json:"message"`
}
