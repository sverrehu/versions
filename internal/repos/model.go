package repos

import (
	"github.com/sverrehu/gotest/versions/internal"
	"github.com/sverrehu/gotest/versions/internal/config"
)

type Fetcher interface {
	GetReleases(pkg string) (*internal.ReleasesResponse, error)
}

type FetcherBase struct {
	Fetcher
	firstPage   int
	perPage     int
	maxReleases int
	credentials *config.Credentials
}

type FetcherError struct {
	Err              error
	IsParameterError bool
}

func (e *FetcherError) Error() string {
	return e.Err.Error()
}

func NewFetcherBase(firstPage, perPage, maxReleases int, credentials *config.Credentials) *FetcherBase {
	fb := &FetcherBase{
		firstPage:   firstPage,
		perPage:     perPage,
		maxReleases: maxReleases,
		credentials: credentials,
	}
	if maxReleases > 0 && maxReleases < perPage {
		fb.perPage = maxReleases
	}
	return fb
}
