package domain

import (
	"time"
)

type Lyrics struct {
	ID        string    `json:"id"`
	Lyrics    []string  `json:"lyrics"`
	Authors   []Author  `json:"authors"`
	Tone      string    `json:"tone"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type LyricsService interface {
	List(params PaginationOptions) (PaginateResult[Lyrics], error)
	Create(lyrics *Lyrics) error
}
