package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http" // HTTP Client and server implementations
	"net/url"
	"os" // Access OS functionality
	"strconv"
	"time"

	"example.com/m/news"
	"github.com/joho/godotenv"
)

var tpl = template.Must(template.ParseFiles("index.html")) // template.ParseFiles parses index.html and validates it
// It's wrapped with template.Must so that the code panics if an error occurs when parsing

func indexHandler(w http.ResponseWriter, r *http.Request /* Register http requests handler*/) {
	// w.Write([]byte("<h1>Let's Go</h1>"))
	tpl.Execute(w, nil)
}

func searchHandler(newsapi *news.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		results, err := newsapi.FetchEverything(searchQuery, page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		nextPage, err := strconv.Atoi(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		search := &Search{
			Query:      searchQuery,
			NextPage:   nextPage,
			TotalPages: int(math.Ceil(float64(results.TotalResults) / float64(newsapi.PageSize))),
			Results:    results,
		}

		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf.WriteTo(w)

		fmt.Printf("%+v", results)

		fmt.Println("Search Query is: ", searchQuery)
		fmt.Println("Page is: ", page)
	}
}

type Search struct {
	Query      string
	NextPage   int
	TotalPages int
	Results    *news.Results
}

// var newsapi *news.Client

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

	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		log.Fatal("Env: api key must be set")
	}

	myClient := &http.Client{Timeout: 10 * time.Second}
	newsapi := news.NewClient(myClient, apiKey, 20)

	mux := http.NewServeMux() // Create a http request multiplexer assigned to the mux variable

	// Essentially, a request multiplexer matches the URL of incoming requests against a list of registered patterns, and calls the associated handler for the pattern whenever a match is found
	// Instantiate a file server object
	fs := http.FileServer(http.Dir("assets"))

	// tell router to use file server object for all paths beginning with the /assets/ prefix
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs)) // The http.StripPrefix() method modifies the request URL by stripping off the specified prefix before forwarding the handling of the request to the http.Handler in the second parameter.

	mux.HandleFunc("/search", searchHandler(newsapi))
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux) // Start server
}
