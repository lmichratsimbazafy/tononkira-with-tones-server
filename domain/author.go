package domain

import "time"

type Author struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Songs     []Lyrics  `json:"songs"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
