package constants

type RevId string
type MovieId string

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
