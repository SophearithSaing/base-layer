package movie

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type movie struct {
	ID        bson.ObjectID `bson:"_id" json:"id"`
	CreatorID bson.ObjectID `bson:"creatorId" json:"creatorId"`
	Date      time.Time     `bson:"date" json:"date"`
	Metadata  metadata      `bson:"metadata" json:"metadata"`
	Rating    int           `bson:"rating" json:"rating"`
	CreatedAt time.Time     `bson:"createAt" json:"createAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type metadata struct {
	ID               int      `json:"id"`
	Overview         string   `json:"overview"`
	OriginalLanguage string   `json:"original_language"`
	OriginalTitle    string   `json:"original_title"`
	Video            bool     `json:"video"`
	Title            string   `json:"title"`
	GenreIDs         []string `json:"genre_ids"`
	PosterPath       []string `json:"poster_path"`
	BackdropPath     []string `json:"backdrop_path"`
	ReleaseDate      string   `json:"release_date"`
	VoteAverage      float64  `json:"vote_average"`
	Popularity       float64  `json:"popularity"`
	Adult            bool     `json:"adult"`
	VoteCount        int      `json:"vote_count"`
}

type createMoviePayload struct {
	Date     time.Time `bson:"date" json:"date"`
	Metadata metadata  `bson:"metadata" json:"metadata"`
	Rating   int       `bson:"rating" json:"rating"`
}

type updatePayload[T any] struct {
	Payload   T         `bson:",inline"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

type updateMoviePayload struct {
	Metadata *metadata `bson:"metadata,omitempty" json:"metadata,omitempty"`
	Rating   int       `bson:"rating,omitempty" json:"rating,omitempty"`
}
