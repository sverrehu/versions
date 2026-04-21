package repos

import (
	"github.com/sverrehu/gotest/versions/internal"
	"github.com/sverrehu/gotest/versions/internal/config"
)

type ReleasesFetcher interface {
	GetReleases(pkg string) (*internal.ReleasesResponse, error)
}

type FetcherBase struct {
	ReleasesFetcher
	firstPage   int
	perPage     int
	credentials *config.Credentials
}

type ReleasesFetcherError struct {
	Err              error
	IsParameterError bool
}

func (e *ReleasesFetcherError) Error() string {
	return e.Err.Error()
}
