package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

var cfg = elasticsearch.Config{
	Addresses: []string{
		"http://localhost:9200",
	},
	Logger: &elastictransport.ColorLogger{Output: os.Stdout},
}
var es, err = elasticsearch.NewTypedClient(cfg)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /movies", getMovies)
	mux.HandleFunc("POST /movies", createMovie)
	mux.HandleFunc("PUT /movies/:title", updateMovie)
	mux.HandleFunc("DEL /movies/:title", deleteMovie)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalln("Server failed to start", err)
	}
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	genre := r.URL.Query().Get("genre")

	var esRes *search.Response
	var esErr error

	if title != "" {
		fmt.Printf("Searching by title - %s", title)
		esRes, esErr = es.Search().
			Index("movies").
			Request(&search.Request{
				Query: &types.Query{
					QueryString: &types.QueryStringQuery{
						Fields: []string{"title"},
						Query:  title,
					},
				},
			}).Do(context.Background())
	} else if genre != "" {
		fmt.Printf("Searching by genre - %s", genre)
		esRes, esErr = es.Search().
			Index("movies").
			Request(&search.Request{
				Query: &types.Query{
					QueryString: &types.QueryStringQuery{
						Fields: []string{"genre"},
						Query:  genre,
					},
				},
			}).Do(context.Background())
	} else {
		esRes, esErr = es.Search().
			Index("movies").
			Request(&search.Request{
				Query: &types.Query{
					MatchAll: &types.MatchAllQuery{},
				},
			}).Do(context.Background())
	}

	if esErr != nil {
		fmt.Printf("Elasticsearch Error: %s", esErr)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var hits []map[string]interface{}
	for _, hit := range esRes.Hits.Hits {
		var source map[string]interface{}
		err := json.Unmarshal(hit.Source_, &source)
		if err != nil {
			continue
		}
		hits = append(hits, source)
	}

	jsonRes, jsonErr := json.Marshal(map[string]interface{}{"movies": hits})

	if jsonErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	w.Write(jsonRes)
	return
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Printf("%s", b)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
}
