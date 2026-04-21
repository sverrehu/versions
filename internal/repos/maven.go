package repos

// Sample: https://central.sonatype.com/solrsearch/select?wt=json&q=g:com.fasterxml.jackson.core+AND+a:jackson-core&sort=v+desc

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/sverrehu/gotest/versions/internal"
	"github.com/sverrehu/gotest/versions/internal/config"
	"github.com/sverrehu/gotest/versions/internal/webclient"
)

type MavenReleasesFetcher struct {
	FetcherBase
}

type fullSonatypeResponse struct {
	Response struct {
		Docs []struct {
			V         string `json:"v"`
			Timestamp int64  `json:"timestamp"`
		} `json:"docs"`
	} `json:"response"`
}

func NewMavenReleasesFetcher(datasource *config.Datasource) *MavenReleasesFetcher {
	return &MavenReleasesFetcher{
		FetcherBase: *NewFetcherBase(1, 100, datasource.MaxReleases, datasource.Credentials),
	}
}

func (rf *MavenReleasesFetcher) GetReleases(pkg string) (*internal.ReleasesResponse, error) {
	parts := regexp.MustCompile("[/]").Split(pkg, -1)
	if len(parts) != 2 {
		return nil, &FetcherError{Err: fmt.Errorf("expected two parts, separated by ':' or '/' in maven package, got %s", pkg), IsParameterError: true}
	}
	return rf.getReleases(parts[0], parts[1])
}

func (rf *MavenReleasesFetcher) getReleases(groupId, artifactId string) (*internal.ReleasesResponse, error) {
	searchUrl := rf.getSearchUrl(groupId, artifactId)
	body, err := webclient.Get(searchUrl, rf.credentials)
	if err != nil {
		return nil, err
	}
	if body == "" {
		return &internal.ReleasesResponse{}, nil
	}
	releasesResponse, err := rf.translateResponse(body)
	if err != nil {
		return nil, err
	}
	if rf.maxReleases > 0 && len(releasesResponse.Releases) > rf.maxReleases {
		releasesResponse.Releases = releasesResponse.Releases[:rf.maxReleases]
	}
	return releasesResponse, nil
}

func (rf *MavenReleasesFetcher) getSearchUrl(groupId, artifactId string) string {
	return fmt.Sprintf("https://central.sonatype.com/solrsearch/select?wt=json&q=g:%s+AND+a:%s&sort=v+desc",
		url.QueryEscape(groupId), url.QueryEscape(artifactId))
}

func (rf *MavenReleasesFetcher) translateResponse(jsonResponse string) (*internal.ReleasesResponse, error) {
	var resp fullSonatypeResponse
	err := json.Unmarshal([]byte(jsonResponse), &resp)
	if err != nil {
		return nil, err
	}
	releases := internal.ReleasesResponse{
		Releases: make([]internal.Release, 0, len(resp.Response.Docs)),
	}
	for _, doc := range resp.Response.Docs {
		release := internal.Release{
			Version:          doc.V,
			ReleaseTimestamp: time.UnixMilli(doc.Timestamp),
		}
		releases.Releases = append(releases.Releases, release)
	}
	return &releases, nil
}
