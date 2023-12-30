package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"` // Struct embedding (pointer only)
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Contain sample movies
var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request)  {

}

func getMovie(w http.ResponseWriter, r *http.Request)  {

}

func createMovie(w http.ResponseWriter, r *http.Request)  {
  
}

func updateMovie(w http.ResponseWriter, r *http.Request)  {
  
}

func deleteMovie(w http.ResponseWriter, r *http.Request)  {

}

func main() {
  r := mux.NewRouter()

  // Adding movies to the movie list
  movies = append(movies, Movie{
    ID       : "1",
    Isbn     : "438227",
    Title    : "Movie One",
    Director : &Director{
      Firstname : "John",
      Lastname  : "Doe",
    },
  })
  movies = append(movies, Movie{
    ID       : "2",
    Isbn     : "438228",
    Title    : "Movie Two",
    Director : &Director{
      Firstname : "Ada",
      Lastname  : "Lovelace",
    },
  })

  r.HandleFunc("/movies", getMovies).Methods("GET")
  r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
  r.HandleFunc("/movies", createMovie).Methods("POST")
  r.HandleFunc("/movies/{id}", updateMovie).Methods("POST")
  r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

  fmt.Println("Starting server at port 8000")
  log.Fatal(http.ListenAndServe(":8000", r))
}
