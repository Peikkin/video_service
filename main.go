package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Movie struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var movies []Movie

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	router := mux.NewRouter()

	router.HandleFunc("/movies", createMovie).Methods("POST")
	router.HandleFunc("/movies", getMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	router.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	router.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	log.Info().Msg("Запуск сервера на порту :8080")
	log.Fatal().Err(http.ListenAndServe(":8080", router)).Msg("Ошибка запуска сервера")

}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "aplication/json")
	var movie Movie
	json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = uuid.New()
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "aplication/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID.String() == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			json.NewDecoder(r.Body).Decode(&movie)
			movie.ID, _ = uuid.Parse(params["id"])
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			break
		}
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "aplication/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID.String() == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "aplication/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID.String() == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "aplication/json")
	json.NewEncoder(w).Encode(movies)
}
