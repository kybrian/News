package main

import (
	"html/template"
	"log"
	"net/http" // HTTP Client and server implementations
	"os"       // Access OS functionality

	"github.com/joho/godotenv"
)

var tpl = template.Must(template.ParseFiles("index.html")) // template.ParseFiles parses index.html and validates it
// It's wrapped with template.Must so that the code panics if an error occurs when parsing

func indexHandler(w http.ResponseWriter, r *http.Request /* Register http requests handler*/) {
	// w.Write([]byte("<h1>Let's Go</h1>"))
	tpl.Execute(w, nil)
}

func main() {
	// Read variables from environment
	err := godotenv.Load() // read the .env file and loads the set environment variables into the environment making them accessible through os.Getenv() method
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Writing our web server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Instantiate a file server object
	fs := http.FileServer(http.Dir("assets"))

	mux := http.NewServeMux() // Create a http request multiplexer assigned to the mux variable

	// Essentially, a request multiplexer matches the URL of incoming requests against a list of registered patterns, and calls the associated handler for the pattern whenever a match is found

	mux.HandleFunc("/", indexHandler)
	// tell router to use file server object for all paths beginning with the /assets/ prefix
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs)) // The http.StripPrefix() method modifies the request URL by stripping off the specified prefix before forwarding the handling of the request to the http.Handler in the second parameter.
	http.ListenAndServe(":"+port, mux)                       // Start server
}
