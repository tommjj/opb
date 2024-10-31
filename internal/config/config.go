package config

import (
	"errors"
)

var (
	ErrInvalid = errors.New("arguments is invalid")
)

type Config struct {
	// Browser app path example:"/c/br.exe"
	Browser string `json:"browser"`

	// SearchDefault name of SearchDefault
	// example:"google"
	// Search list search links
	// example:"google: https://www.google.com/search?q=$"
	Search map[string]string `json:"search,omitempty"`

	// Links is a map data
	// map[link name]url
	// example:"yt: https://www.youtube.com"
	Links map[string]string `json:"links,omitempty"`
}
