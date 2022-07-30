package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	Id       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var movies []Movie

func GetMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func GetMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.Id == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "404 Not Found", http.StatusNotFound)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	json.NewDecoder(r.Body).Decode(&movie)
	movie.Id = strconv.Itoa(rand.Intn(100000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.Id == params["id"] {
			// remove movie
			movies = append(movies[:index], movies[index+1:]...)

			// update movie
			var movie Movie
			json.NewDecoder(r.Body).Decode(&movie)
			movie.Id = params["id"]
			movies = append(movies, movie)

			// return updated movie
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
	http.Error(w, "404 Not Found", http.StatusNotFound)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.Id == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			json.NewEncoder(w).Encode(true)
			return
		}
	}
	http.Error(w, "404 Not Found", http.StatusNotFound)
}

func main() {
	r := mux.NewRouter()

	// initial movies
	movies = append(movies, Movie{Id: "1", Isbn: "525637", Title: "Iron Man", Director: &Director{FirstName: "Tony", LastName: "Stark"}})
	movies = append(movies, Movie{Id: "2", Isbn: "846284", Title: "Superman", Director: &Director{FirstName: "Clarke", LastName: "Kent"}})
	// movies = append(movies, Movie{Id: "2", Isbn: "235481", Title: "Spiderman", Director: &Director{FirstName: "Peter", LastName: "Parker"}})

	r.HandleFunc("/movies", GetMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", GetMovie).Methods("GET")
	r.HandleFunc("/movies", CreateMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", UpdateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Starting server at http://localhost:8000")
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", r))
}
