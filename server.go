package main

import (
	"fmt"
	"net/http"
)

func startServerMode(port int, influxConfig ConnectionConfig) {

	fmt.Println("starting server")
	http.HandleFunc("/", handler)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "hellow %s!", request.URL.Path[1:])
}
