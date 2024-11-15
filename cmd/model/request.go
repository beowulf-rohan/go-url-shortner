package model

import "time"

type Request struct {
	URL      string        `json:"url"`
	ShortURL string        `json:"short_url"`
	Expiry   time.Duration `json:"expiry"`
}
