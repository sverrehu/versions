package repos

import (
	"github.com/sverrehu/versions/internal"
	"github.com/sverrehu/versions/internal/config"
	"github.com/sverrehu/versions/internal/webclient"
)

type Fetcher interface {
	GetReleases(pkg string) (*internal.ReleasesResponse, error)
	getSearchUrl(owner, repo string, page int) string
	extractReleases(owner, repo, jsonResponse string) ([]internal.Release, error)
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

func (fb *FetcherBase) paginate(f Fetcher, releasesResponse *internal.ReleasesResponse, owner, repo string) error {
	page := fb.firstPage
	for {
		searchUrl := f.getSearchUrl(owner, repo, page)
		body, err := webclient.Get(searchUrl, fb.credentials)
		if err != nil {
			return err
		}
		if body == "" {
			break
		}
		releases, err := f.extractReleases(owner, repo, body)
		if err != nil {
			return err
		}
		if len(releases) == 0 {
			break
		}
		releasesResponse.Releases = append(releasesResponse.Releases, releases...)
		page++
		if fb.maxReleases > 0 && len(releasesResponse.Releases) > fb.maxReleases {
			releasesResponse.Releases = releasesResponse.Releases[:fb.maxReleases]
			break
		}
	}
	return nil
}
