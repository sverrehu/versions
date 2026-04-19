package repos

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

// sample: https://gitlab.com/api/v4/projects/gitlab-org%2Fgitlab-runner/releases?page=1&per_page=100

type GitLabReleasesFetcher struct {
	ReleasesFetcher
	firstPage   int
	perPage     int
	credentials *config.Credentials
}

type fullGitLabReleasesResponse []struct {
	Name            string    `json:"name"`
	TagName         string    `json:"tag_name"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	ReleasedAt      time.Time `json:"released_at"`
	UpcomingRelease bool      `json:"upcoming_release"`
	Author          struct {
		ID          int    `json:"id"`
		Username    string `json:"username"`
		PublicEmail string `json:"public_email"`
		Name        string `json:"name"`
		State       string `json:"state"`
		Locked      bool   `json:"locked"`
		AvatarURL   string `json:"avatar_url"`
		WebURL      string `json:"web_url"`
	} `json:"author"`
	Commit struct {
		ID             string    `json:"id"`
		ShortID        string    `json:"short_id"`
		CreatedAt      time.Time `json:"created_at"`
		ParentIds      []string  `json:"parent_ids"`
		Title          string    `json:"title"`
		Message        string    `json:"message"`
		AuthorName     string    `json:"author_name"`
		AuthorEmail    string    `json:"author_email"`
		AuthoredDate   time.Time `json:"authored_date"`
		CommitterName  string    `json:"committer_name"`
		CommitterEmail string    `json:"committer_email"`
		CommittedDate  time.Time `json:"committed_date"`
		Trailers       struct {
		} `json:"trailers"`
		ExtendedTrailers struct {
		} `json:"extended_trailers"`
		WebURL string `json:"web_url"`
	} `json:"commit"`
	CommitPath string `json:"commit_path"`
	TagPath    string `json:"tag_path"`
	Assets     struct {
		Count   int `json:"count"`
		Sources []struct {
			Format string `json:"format"`
			URL    string `json:"url"`
		} `json:"sources"`
		Links []struct {
			ID             int    `json:"id"`
			Name           string `json:"name"`
			URL            string `json:"url"`
			DirectAssetURL string `json:"direct_asset_url"`
			LinkType       string `json:"link_type"`
		} `json:"links"`
	} `json:"assets"`
	Evidences []struct {
		Sha         string    `json:"sha"`
		Filepath    string    `json:"filepath"`
		CollectedAt time.Time `json:"collected_at"`
	} `json:"evidences"`
	Links struct {
		ClosedIssuesURL        string `json:"closed_issues_url"`
		ClosedMergeRequestsURL string `json:"closed_merge_requests_url"`
		MergedMergeRequestsURL string `json:"merged_merge_requests_url"`
		OpenedIssuesURL        string `json:"opened_issues_url"`
		OpenedMergeRequestsURL string `json:"opened_merge_requests_url"`
		Self                   string `json:"self"`
	} `json:"_links"`
}

func NewGitLabReleasesFetcher(credentials *config.Credentials) *GitLabReleasesFetcher {
	return &GitLabReleasesFetcher{
		firstPage:   1,
		perPage:     100,
		credentials: credentials,
	}
}

func (rf *GitLabReleasesFetcher) GetReleases(pkg string) (*internal.ReleasesResponse, error) {
	parts := regexp.MustCompile("[/]").Split(pkg, -1)
	if len(parts) != 2 {
		return nil, &ReleasesFetcherError{Err: fmt.Errorf("expected two parts, separated by '/' in GitLab releases package, got %s", pkg), IsParameterError: true}
	}
	return rf.getReleases(parts[0], parts[1])
}

func (rf *GitLabReleasesFetcher) getReleases(owner, repo string) (*internal.ReleasesResponse, error) {
	searchUrl := rf.getSearchUrl(owner, repo, rf.firstPage)
	body, err := webclient.Get(searchUrl, rf.credentials)
	if err != nil {
		return nil, err
	}
	if body == "" {
		return &internal.ReleasesResponse{}, nil
	}
	releases, err := rf.translateResponse(body, owner, repo)
	if err != nil {
		return nil, err
	}
	return releases, nil
}

func (rf *GitLabReleasesFetcher) getSearchUrl(owner, repo string, page int) string {
	return fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/releases?page=%d&per_page=%d",
		url.PathEscape(owner+"/"+repo), page, rf.perPage)
}

func (rf *GitLabReleasesFetcher) translateResponse(jsonResponse, owner, repo string) (*internal.ReleasesResponse, error) {
	var resp fullGitLabReleasesResponse
	err := json.Unmarshal([]byte(jsonResponse), &resp)
	if err != nil {
		return nil, err
	}
	sourceURL := "https://gitlab.com/" + url.PathEscape(owner) + "/" + url.PathEscape(repo)
	releases := internal.ReleasesResponse{
		Releases:  make([]internal.Release, 0, len(resp)),
		SourceURL: &sourceURL,
	}
	for _, result := range resp {
		release := internal.Release{
			Version:          result.TagName,
			ReleaseTimestamp: result.ReleasedAt,
			ChangelogURL:     &result.Links.Self,
		}
		releases.Releases = append(releases.Releases, release)
	}
	return &releases, nil
}
