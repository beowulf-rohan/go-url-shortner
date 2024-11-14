package model

import "time"

type Response struct {
	URL          string    `json:"url"`
	ShortenedURL string    `json:"shortened_url"`
	Expiry       time.Time `json:"expiry"`
}
