package main

import (
	"io"
	"net/http"

	honeycomb "github.com/honeycombio/honeycomb-go-magic"
)

func main() {
	honeycomb.Init(honeycomb.Config{
		WriteKey: "abcabc123123",
		Dataset:  "http-mux",
	})
	globalmux := http.NewServeMux()
	globalmux.HandleFunc("/hello", hello)
	http.ListenAndServe(":8080", honeycomb.InstrumentMuxHandler(globalmux))
}

func hello(w http.ResponseWriter, r *http.Request) {
	honeycomb.AddField(r.Context(), "custom", "Wheee")
	io.WriteString(w, "Hello world!")
}

// Example events created:
// {
//   "data": {
//     "durationMs": 0.045241,
//     "handlerName": "main.hello",
//     "handlerPattern": "/hello",
//     "handlerType": "http.HandlerFunc",
//     "request.content_length": 0,
//     "request.host": "",
//     "request.method": "GET",
//     "request.path": "/hello",
//     "request.proto": "HTTP/1.1",
//     "request.remote_addr": "[::1]:57594",
//     "request.user_agent": "curl/7.54.0",
//     "response.status_code": 200
//   },
//   "time": "2018-03-08T15:28:18.351099-08:00"
// }
// {
//   "data": {
//     "durationMs": 0.077218,
//     "handlerName": "net/http.NotFound",
//     "handlerPattern": "",
//     "handlerType": "http.HandlerFunc",
//     "request.content_length": 0,
//     "request.host": "",
//     "request.method": "GET",
//     "request.path": "/bar",
//     "request.proto": "HTTP/1.1",
//     "request.remote_addr": "[::1]:57593",
//     "request.user_agent": "curl/7.54.0",
//     "response.status_code": 404
//   },
//   "time": "2018-03-08T15:28:21.458141-08:00"
// }
