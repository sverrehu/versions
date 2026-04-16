package internal

import (
	"time"
)

type OldRelease struct {
	Version    string    `json:"version"`
	ReleasedAt time.Time `json:"released_at"`
	ReleaseURL *string   `json:"release_url,omitempty"`
	SourceURL  *string   `json:"source_url,omitempty"`
}

type Release struct {
	Version          string    `json:"version"`
	IsDeprecated     *bool     `json:"isDeprecated,omitempty"`
	ReleaseTimestamp time.Time `json:"releaseTimestamp"`
	ChangelogURL     *string   `json:"changelogUrl,omitempty"`
	SourceURL        *string   `json:"sourceUrl,omitempty"`
	SourceDirectory  *string   `json:"sourceDirectory,omitempty"`
	Digest           *string   `json:"digest,omitempty"`
	IsStable         *bool     `json:"isStable,omitempty"`
}

type ReleasesResponse struct {
	Releases        []Release `json:"releases"`
	SourceURL       *string   `json:"sourceUrl,omitempty"`
	SourceDirectory *string   `json:"sourceDirectory,omitempty"`
	ChangelogURL    *string   `json:"changelogUrl,omitempty"`
	Homepage        *string   `json:"homepage,omitempty"`
}
