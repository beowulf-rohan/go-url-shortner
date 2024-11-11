package model

import "time"

type Response struct {
	URL            string    `json:"url"`
	ShortenedURL   string    `json:"shortened_url"`
	Expiry         time.Time `json:"expiry"`
	RateRemainsing int       `json:"rate_remainsing"`
	RateLimitReset time.Time `json:"rate_limit_reset"`
}
