package model

import "time"

type Response struct {
	URL            string    `json:"url"`
	ShortenedURL   string    `json:"shortened_url"`
	Expiry         time.Duration       `json:"expiry"`
	RateRemaining  int       `json:"rate_remainsing"`
	RateLimitReset time.Duration `json:"rate_limit_reset"`
}
