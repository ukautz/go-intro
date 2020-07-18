// returns a `200 OK` response with the body `Hello World (<number>)` to any incoming request
package main

import (
	"fmt"
	"net/http"
)

const address = "127.0.0.1:12345"
var counter = 0

func main() {
	fmt.Println("Starting HTTP server at", address)
	http.HandleFunc("/", helloWorld)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}

func helloWorld(writer http.ResponseWriter, request *http.Request) {
	counter++
	response := fmt.Sprintf("Hello World (%d)", counter)
	writer.Write([]byte(response))
}
