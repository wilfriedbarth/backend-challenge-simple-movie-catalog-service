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
	mux.HandleFunc("GET /movies/{id}", getMovie)
	mux.HandleFunc("POST /movies", createMovie)
	mux.HandleFunc("PUT /movies/{id}", updateMovie)
	mux.HandleFunc("DELETE /movies/{id}", deleteMovie)

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
		fmt.Printf("Get Movies By Title Request - %s", title)
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
		fmt.Printf("Get Movies By Genre Request - %s", genre)
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
		fmt.Println("Get All Movies Request")
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

	var movies []map[string]interface{}
	for _, hit := range esRes.Hits.Hits {
		var movie map[string]interface{}
		err := json.Unmarshal(hit.Source_, &movie)
		movie["id"] = hit.Id_
		if err != nil {
			continue
		}
		movies = append(movies, movie)
	}

	jsonRes, jsonErr := json.Marshal(map[string]interface{}{"movies": movies})

	if jsonErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
	return
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Printf("Get Movie Request - %s", id)
	esRes, esErr := es.Get("movies", id).Do(context.Background())

	if esErr != nil {
		fmt.Printf("Elasticsearch Error: %s", esErr)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var movie map[string]interface{}
	err := json.Unmarshal(esRes.Source_, &movie)
	if err != nil {
		fmt.Printf("JSON Unmarshal Error: %s", esErr)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	movie["id"] = esRes.Id_

	jsonRes, jsonErr := json.Marshal(movie)
	if jsonErr != nil {
		fmt.Printf("JSON Marshal Error: %s", esErr)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
	return
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var movie map[string]interface{}
	json.Unmarshal(b, &movie)

	fmt.Printf("Create Movie Request - %s", b)

	esRes, esErr := es.Index("movies").Request(movie).Do(context.Background())

	if esErr != nil {
		fmt.Printf("Elasticsearch Error: %s", esErr)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Request result - %s", esRes.Result.String())

	if esRes.Result.String() != "created" {
		fmt.Printf("Elasticsearch failed to create document")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("201 - Status Created"))
	return
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var movie map[string]interface{}
	json.Unmarshal(b, &movie)

	fmt.Printf("Update Movie Request - id %s, body %s", id, b)

	esRes, esErr := es.Index("movies").Id(id).Request(movie).Do(context.Background())

	if esErr != nil {
		fmt.Printf("Elasticsearch Error: %s", esErr)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Request result - %s", esRes.Result.String())

	if esRes.Result.String() != "updated" {
		fmt.Printf("Elasticsearch failed to update document")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	return
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	_, esErr := es.Delete("movies", id).Do(context.Background())

	if esErr != nil {
		fmt.Printf("Elasticsearch Error: %s", esErr)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK"))
	return
}
