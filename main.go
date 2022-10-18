package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http" // HTTP Client and server implementations
	"net/url"
	"os" // Access OS functionality

	"github.com/joho/godotenv"
)

var tpl = template.Must(template.ParseFiles("index.html")) // template.ParseFiles parses index.html and validates it
// It's wrapped with template.Must so that the code panics if an error occurs when parsing

func indexHandler(w http.ResponseWriter, r *http.Request /* Register http requests handler*/) {
	// w.Write([]byte("<h1>Let's Go</h1>"))
	tpl.Execute(w, nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := u.Query()
	searchQuery := params.Get("q")
	page := params.Get("page")
	if page == "" {
		page = "1"
	}

	fmt.Println("Search Query is: ", searchQuery)
	fmt.Println("Page is: ", page)
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

	mux := http.NewServeMux() // Create a http request multiplexer assigned to the mux variable

	// Essentially, a request multiplexer matches the URL of incoming requests against a list of registered patterns, and calls the associated handler for the pattern whenever a match is found
	// Instantiate a file server object
	fs := http.FileServer(http.Dir("assets"))

	// tell router to use file server object for all paths beginning with the /assets/ prefix
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs)) // The http.StripPrefix() method modifies the request URL by stripping off the specified prefix before forwarding the handling of the request to the http.Handler in the second parameter.

	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux) // Start server
}
