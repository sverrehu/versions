package internal

import (
	"time"
)

type Release struct {
	Version    string    `json:"version"`
	ReleasedAt time.Time `json:"released_at"`
	ReleaseURL *string   `json:"release_url,omitempty"`
	SourceURL  *string   `json:"source_url,omitempty"`
}
