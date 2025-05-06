package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	// Create a new router
	r := mux.NewRouter()

	// Define a route for the profile URL
	r.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello Joaquín!, you´ve requested your profile data")
	})

	// Define a route for the profile URL with a variable
	r.HandleFunc("/profile/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		fmt.Fprintf(w, "Hello %s!, you´ve requested your profile data", name)
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	// Start the server on port 80 and use the router
	http.ListenAndServe(":80", r)
}
