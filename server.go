package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!, you´ve requested the path: %s", r.URL.Path)
	})
	http.ListenAndServe(":80", nil)
}
