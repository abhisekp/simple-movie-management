package utils

import (
	"fmt"
	. "github.com/abhisekp/simple_movie_management/constants"
)

func GetRevId(movieId MovieId, revSeq int) RevId {
	return RevId(fmt.Sprintf("%s-%d", movieId, revSeq))
}
