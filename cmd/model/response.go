package model

import "time"

type Response struct {
	URL       string    `json:"url"`
	ShortURL  string    `json:"short_url"`
	CreatedAt time.Time `json:"created_at"`
	Expiry    time.Time `json:"expiry"`
}
