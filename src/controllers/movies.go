package controllers

import (
	"encoding/json"
	"fmt"
	. "github.com/abhisekp/simple_movie_management/constants"
	. "github.com/abhisekp/simple_movie_management/src/services"
	"github.com/abhisekp/simple_movie_management/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func GetMovieRevisionsById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := MovieId(mux.Vars(r)["id"])
	fmt.Println("id: ", id)

	if _, ok := MovieIndexes[id]; !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	movieRevisionsById := make([]Movie, 0, SeqMovieRevisions[id])

	for seq := range SeqMovieRevisions[id] {
		if movieRevision, ok := MovieRevisions[utils.GetRevId(id, seq+1)]; ok {
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
		_newMovieData, ok := MovieRevisions[utils.GetRevId(idToUpdate, newMovieData.RevSeq)]
		if !ok {
			http.Error(w, "Revision not found", http.StatusNotFound)
			return
		}
		newMovieData = _newMovieData
	}

	oldMovie, ok := MovieIndexes[idToUpdate]
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
		SeqMovieRevisions[oldMovie.ID] += 1
		revId := utils.GetRevId(oldMovie.ID, SeqMovieRevisions[oldMovie.ID])
		oldMovie.RevSeq = SeqMovieRevisions[oldMovie.ID]
		MovieRevisions[revId] = *oldMovie
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

	_, existing := MovieIndexes[newMovie.ID]

	if newMovie.ID == "" || existing {
		for ; MovieIndexes[MovieId(strconv.Itoa(SeqMovieId))] != nil; SeqMovieId += 1 {
		}
		newMovie.ID = MovieId(strconv.Itoa(SeqMovieId))
	}
	Movies = append(Movies, newMovie)
	MovieIndexes[newMovie.ID] = &newMovie

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

	for idx, movie := range Movies {
		if movie.ID == MovieId(idToDelete) {
			deletedMovie = movie
			Movies = append(Movies[:idx], Movies[idx+1:]...)
			delete(MovieIndexes, MovieId(idToDelete))

			// Delete Revisions
			if seqMovieRevisionMax, ok := SeqMovieRevisions[MovieId(idToDelete)]; ok {
				for revSeq := range seqMovieRevisionMax {
					delete(MovieRevisions, utils.GetRevId(MovieId(idToDelete), revSeq))
				}
				delete(SeqMovieRevisions, MovieId(idToDelete))
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

	foundMovie, foundMovieOk := MovieIndexes[MovieId(idToFind)]
	if !foundMovieOk {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	revSeq, _ := strconv.Atoi(r.URL.Query().Get("revSeq"))
	if revSeq > 0 {
		fmt.Println("revSeq: ", revSeq)
		revId := utils.GetRevId(MovieId(idToFind), revSeq)
		if movieRevision, ok := MovieRevisions[revId]; ok {
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

func ListMovies(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(Movies)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
