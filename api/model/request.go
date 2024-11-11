package model

import "time"

type Request struct {
	URL          string    `json:"url"`
	ShortenedURL string    `json:"shortened_url"`
	Expiry       time.Time `json:"expiry"`
}
