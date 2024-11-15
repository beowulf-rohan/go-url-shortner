package model

import "time"

type Response struct {
	URL      string    `json:"url"`
	ShortURL string    `json:"short_url"`
	Expiry   time.Time `json:"expiry"`
}
