package repos

// Sample: https://hub.docker.com/v2/repositories/library/ubuntu/tags?page=1&page_size=100

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/sverrehu/versions/internal"
	"github.com/sverrehu/versions/internal/config"
)

type QuayReleasesFetcher struct {
	FetcherBase
}

type fullQuayResponse struct {
	Tags []struct {
		Name                          string `json:"name"`
		Reversion                     bool   `json:"reversion"`
		StartTs                       int64  `json:"start_ts"`
		ManifestDigest                string `json:"manifest_digest"`
		IsManifestList                bool   `json:"is_manifest_list"`
		Size                          int    `json:"size"`
		LastModified                  string `json:"last_modified"`
		CosignSignatureManifestDigest string `json:"cosign_signature_manifest_digest"`
		IsSparse                      bool   `json:"is_sparse,omitempty"`
		EndTs                         int    `json:"end_ts,omitempty"`
		Expiration                    string `json:"expiration,omitempty"`
	} `json:"tags"`
	Page          int  `json:"page"`
	HasAdditional bool `json:"has_additional"`
}

func NewQuayReleasesFetcher(datasource *config.Datasource) *QuayReleasesFetcher {
	return &QuayReleasesFetcher{
		FetcherBase: *NewFetcherBase(1, 100, datasource.MaxReleases, datasource.Credentials),
	}
}

func (rf *QuayReleasesFetcher) GetReleases(pkg string) (*internal.ReleasesResponse, error) {
	parts := strings.Split(pkg, "/")
	if len(parts) != 2 {
		return nil, &FetcherError{Err: fmt.Errorf("expected two parts, separated by '/' in Quay package, got %s", pkg), IsParameterError: true}
	}
	return rf.getReleases(parts[0], parts[1])
}

func (rf *QuayReleasesFetcher) getSearchUrl(repo, image string, page int) string {
	return fmt.Sprintf("https://quay.io/api/v1/repository/%s/%s/tag/?page=%d&limit=%d&onlyActiveTags=true",
		url.PathEscape(repo), url.PathEscape(image), page, rf.perPage)
}

func (rf *QuayReleasesFetcher) getReleases(owner, repo string) (*internal.ReleasesResponse, error) {
	releasesResponse := internal.ReleasesResponse{
		Releases: make([]internal.Release, 0),
	}
	err := rf.paginate(rf, &releasesResponse, owner, repo)
	if err != nil {
		return nil, err
	}
	return &releasesResponse, nil
}

func (rf *QuayReleasesFetcher) extractReleases(_, _, jsonResponse string) ([]internal.Release, bool, error) {
	var resp fullQuayResponse
	err := json.Unmarshal([]byte(jsonResponse), &resp)
	if err != nil {
		return nil, true, err
	}
	releases := make([]internal.Release, 0, len(resp.Tags))
	for _, result := range resp.Tags {
		if !result.IsManifestList {
			continue
		}
		release := internal.Release{
			Version:          result.Name,
			ReleaseTimestamp: time.Unix(result.StartTs, 0),
		}
		releases = append(releases, release)
	}
	return releases, len(resp.Tags) == 0, nil
}
