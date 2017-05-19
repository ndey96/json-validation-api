package main

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
