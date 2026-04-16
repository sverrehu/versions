package repos

import (
	"github.com/sverrehu/gotest/versions/internal"
	"github.com/sverrehu/gotest/versions/internal/config"
)

type ReleasesFetcher interface {
	GetReleases(pkg string, credentials *config.Credentials) (*internal.ReleasesResponse, error)
}

type ReleasesFetcherError struct {
	Err              error
	IsParameterError bool
}

func (e *ReleasesFetcherError) Error() string {
	return e.Err.Error()
}
