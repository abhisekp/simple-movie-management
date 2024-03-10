package services

import . "github.com/abhisekp/simple_movie_management/constants"

var (
	Movies       []Movie
	SeqMovieId   = 0
	MovieIndexes = map[MovieId]*Movie{}
)

var (
	MovieRevisions    = map[RevId]Movie{}
	SeqMovieRevisions = map[MovieId]int{}
)
