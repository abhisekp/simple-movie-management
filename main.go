package main

import (
	"encoding/json"
	"fmt"
	. "github.com/abhisekp/simple_movie_management/constants"
	. "github.com/abhisekp/simple_movie_management/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Movie struct {
	ID       MovieId   `json:"id"`
	Title    string    `json:"title,omitempty"`
	ISBN     string    `json:"isbn,omitempty"`
	Director *Director `json:"director,omitempty"`
	RevSeq   int       `json:"revSeq,omitempty"`
}

type Director struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var (
	movies       []Movie
	seqMovieId   = 0
	movieIndexes = map[MovieId]*Movie{}
)

var (
	movieRevisions    = map[RevId]Movie{}
	seqMovieRevisions = map[MovieId]int{}
)

func main() {
	fmt.Println("Movie Management App")

	movies = append(movies,
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
	seqMovieId = 2
	movieRevisions["1-1"] = movies[0]
	movieRevisions["2-1"] = movies[1]
	seqMovieRevisions["1"] = 1
	seqMovieRevisions["2"] = 1

	for idx := range len(movies) {
		movie := &movies[idx]
		movieIndexes[movie.ID] = movie
	}

	r := mux.NewRouter()
	r.HandleFunc("/movies", ListMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", GetMovieById).Methods("GET")
	r.HandleFunc("/movies/{id}/revisions", GetMovieRevisionsById).Methods("GET")
	r.HandleFunc("/movies/{id}", UpdateMovieById).Methods("PATCH").Headers(
		"Content-Type", "application/json",
		"Accept", "application/json",
	)
	r.HandleFunc("/movies/{id}", DeleteMovieById).Methods("DELETE")
	r.HandleFunc("/movies", CreateMovie).Methods("POST").Headers(
		"Content-Type", "application/json",
		"Accept", "application/json",
	)

	fmt.Println("Server started at http://localhost:7998")
	log.Fatalln(http.ListenAndServe(":7998", r))
}

func GetMovieRevisionsById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := MovieId(mux.Vars(r)["id"])
	fmt.Println("id: ", id)

	if _, ok := movieIndexes[id]; !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	movieRevisionsById := make([]Movie, 0, seqMovieRevisions[id])

	for seq := range seqMovieRevisions[id] {
		if movieRevision, ok := movieRevisions[GetRevId(id, seq+1)]; ok {
			movieRevisionsById = append(movieRevisionsById, movieRevision)
		}
	}

	err := json.NewEncoder(w).Encode(movieRevisionsById)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func UpdateMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idToUpdate := MovieId(mux.Vars(r)["id"])
	fmt.Println("id: ", idToUpdate)

	var newMovieData Movie
	err := json.NewDecoder(r.Body).Decode(&newMovieData)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if newMovieData.RevSeq != 0 {
		_newMovieData, ok := movieRevisions[GetRevId(idToUpdate, newMovieData.RevSeq)]
		if !ok {
			http.Error(w, "Revision not found", http.StatusNotFound)
			return
		}
		newMovieData = _newMovieData
	}

	oldMovie, ok := movieIndexes[idToUpdate]
	if ok {

		if newMovieData.Title != "" {
			oldMovie.Title = newMovieData.Title
		}

		if newMovieData.ISBN != "" {
			oldMovie.ISBN = newMovieData.ISBN
		}

		if newMovieData.Director != nil && (newMovieData.Director.FirstName != "" || newMovieData.Director.LastName != "") {
			oldMovie.Director = newMovieData.Director
		}

		// Store revisions
		seqMovieRevisions[oldMovie.ID] += 1
		revId := GetRevId(oldMovie.ID, seqMovieRevisions[oldMovie.ID])
		oldMovie.RevSeq = seqMovieRevisions[oldMovie.ID]
		movieRevisions[revId] = *oldMovie
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(oldMovie)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newMovie Movie
	err := json.NewDecoder(r.Body).Decode(&newMovie)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	_, existing := movieIndexes[newMovie.ID]

	if newMovie.ID == "" || existing {
		for ; movieIndexes[MovieId(strconv.Itoa(seqMovieId))] != nil; seqMovieId += 1 {
		}
		newMovie.ID = MovieId(strconv.Itoa(seqMovieId))
	}
	movies = append(movies, newMovie)
	movieIndexes[newMovie.ID] = &newMovie

	err = json.NewEncoder(w).Encode(newMovie)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func DeleteMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idToDelete := mux.Vars(r)["id"]
	fmt.Println("id: ", idToDelete)

	var deletedMovie Movie

	for idx, movie := range movies {
		if movie.ID == MovieId(idToDelete) {
			deletedMovie = movie
			movies = append(movies[:idx], movies[idx+1:]...)
			delete(movieIndexes, MovieId(idToDelete))

			// Delete Revisions
			if seqMovieRevisionMax, ok := seqMovieRevisions[MovieId(idToDelete)]; ok {
				for revSeq := range seqMovieRevisionMax {
					delete(movieRevisions, GetRevId(MovieId(idToDelete), revSeq))
				}
				delete(seqMovieRevisions, MovieId(idToDelete))
			}
			break
		}
	}

	if deletedMovie == (Movie{}) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	err := json.NewEncoder(w).Encode(deletedMovie)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func GetMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	idToFind := params["id"]
	fmt.Println("id: ", idToFind)

	foundMovie, foundMovieOk := movieIndexes[MovieId(idToFind)]
	if !foundMovieOk {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	revSeq, _ := strconv.Atoi(r.URL.Query().Get("revSeq"))
	if revSeq > 0 {
		fmt.Println("revSeq: ", revSeq)
		revId := GetRevId(MovieId(idToFind), revSeq)
		if movieRevision, ok := movieRevisions[revId]; ok {
			foundMovie = &movieRevision
		} else {
			http.Error(w, "Revision Not found", http.StatusNotFound)
			return
		}
	}

	err := json.NewEncoder(w).Encode(foundMovie)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func ListMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(movies)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
