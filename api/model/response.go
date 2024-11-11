package model

import "time"

type Response struct {
	URL             string    `json:"url"`
	ShortenedURL    string    `json:"shortened_url"`
	Expiry          time.Time `json:"expiry"`
	XRateRemainsing int       `json:"rate_remainsing"`
	XRateLimitReset time.Time `json:"rate_limit_reset"`
}
