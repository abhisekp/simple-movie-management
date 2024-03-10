package main

import (
	"fmt"
	. "github.com/abhisekp/simple_movie_management/constants"
	"github.com/abhisekp/simple_movie_management/src/controllers"
	. "github.com/abhisekp/simple_movie_management/src/services"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Movie Management App")

	Movies = append(Movies,
		Movie{
			ID:    "1",
			Title: "Zingala",
			ISBN:  "53566",
			Director: &Director{
				FirstName: "John",
				LastName:  "Smith",
			},
			RevSeq: 1,
		}, Movie{
			ID:    "2",
			Title: "Postuo",
			ISBN:  "44566",
			Director: &Director{
				FirstName: "Steven",
				LastName:  "Doe",
			},
			RevSeq: 1,
		},
	)
	SeqMovieId = 2
	MovieRevisions["1-1"] = Movies[0]
	MovieRevisions["2-1"] = Movies[1]
	SeqMovieRevisions["1"] = 1
	SeqMovieRevisions["2"] = 1

	for idx := range len(Movies) {
		movie := &Movies[idx]
		MovieIndexes[movie.ID] = movie
	}

	r := mux.NewRouter()
	r.HandleFunc("/movies", controllers.ListMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", controllers.GetMovieById).Methods("GET")
	r.HandleFunc("/movies/{id}/revisions", controllers.GetMovieRevisionsById).Methods("GET")
	r.HandleFunc("/movies/{id}", controllers.UpdateMovieById).Methods("PATCH").Headers(
		"Content-Type", "application/json",
		"Accept", "application/json",
	)
	r.HandleFunc("/movies/{id}", controllers.DeleteMovieById).Methods("DELETE")
	r.HandleFunc("/movies", controllers.CreateMovie).Methods("POST").Headers(
		"Content-Type", "application/json",
		"Accept", "application/json",
	)

	fmt.Println("Server started at http://localhost:7998")
	log.Fatalln(http.ListenAndServe(":7998", r))
}
