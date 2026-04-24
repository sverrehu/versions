package repos

// Sample releases: https://api.github.com/repos/prometheus/prometheus/releases?page=1&per_page=100
// Sample tags: https://api.github.com/repos/confluentinc/libserdes/tags?page=1&per_page=100

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/sverrehu/versions/internal"
	"github.com/sverrehu/versions/internal/config"
	"github.com/sverrehu/versions/internal/state"
	"github.com/sverrehu/versions/internal/webclient"
)

type GitHubReleasesFetcher struct {
	FetcherBase
}

type GitHubTagsFetcher struct {
	FetcherBase
}

type fullGitHubReleasesResponse []struct {
	URL       string `json:"url"`
	AssetsURL string `json:"assets_url"`
	UploadURL string `json:"upload_url"`
	HTMLURL   string `json:"html_url"`
	ID        int    `json:"id"`
	Author    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		UserViewType      string `json:"user_view_type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Immutable       bool      `json:"immutable"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	PublishedAt     time.Time `json:"published_at"`
	Assets          []struct {
		URL      string `json:"url"`
		ID       int    `json:"id"`
		NodeID   string `json:"node_id"`
		Name     string `json:"name"`
		Label    string `json:"label"`
		Uploader struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			UserViewType      string `json:"user_view_type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
		ContentType        string    `json:"content_type"`
		State              string    `json:"state"`
		Size               int       `json:"size"`
		Digest             string    `json:"digest"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadURL string    `json:"browser_download_url"`
	} `json:"assets"`
	TarballURL string `json:"tarball_url"`
	ZipballURL string `json:"zipball_url"`
	Body       string `json:"body"`
	Reactions  struct {
		URL        string `json:"url"`
		TotalCount int    `json:"total_count"`
		Num1       int    `json:"+1"`
		Num10      int    `json:"-1"`
		Laugh      int    `json:"laugh"`
		Hooray     int    `json:"hooray"`
		Confused   int    `json:"confused"`
		Heart      int    `json:"heart"`
		Rocket     int    `json:"rocket"`
		Eyes       int    `json:"eyes"`
	} `json:"reactions,omitempty"`
	MentionsCount int    `json:"mentions_count,omitempty"`
	DiscussionURL string `json:"discussion_url,omitempty"`
}

type fullGitHubTagsResponse []struct {
	Name       string `json:"name"`
	ZipballURL string `json:"zipball_url"`
	TarballURL string `json:"tarball_url"`
	Commit     struct {
		Sha string `json:"sha"`
		URL string `json:"url"`
	} `json:"commit"`
	NodeID string `json:"node_id"`
}

type fullGitHubCommitResponse struct {
	Sha    string `json:"sha"`
	NodeID string `json:"node_id"`
	Commit struct {
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Committer struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"committer"`
		Message string `json:"message"`
		Tree    struct {
			Sha string `json:"sha"`
			URL string `json:"url"`
		} `json:"tree"`
		URL          string `json:"url"`
		CommentCount int    `json:"comment_count"`
		Verification struct {
			Verified   bool      `json:"verified"`
			Reason     string    `json:"reason"`
			Signature  string    `json:"signature"`
			Payload    string    `json:"payload"`
			VerifiedAt time.Time `json:"verified_at"`
		} `json:"verification"`
	} `json:"commit"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	CommentsURL string `json:"comments_url"`
	Author      struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		UserViewType      string `json:"user_view_type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Committer struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		UserViewType      string `json:"user_view_type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"committer"`
	Parents []struct {
		Sha     string `json:"sha"`
		URL     string `json:"url"`
		HTMLURL string `json:"html_url"`
	} `json:"parents"`
	Stats struct {
		Total     int `json:"total"`
		Additions int `json:"additions"`
		Deletions int `json:"deletions"`
	} `json:"stats"`
	Files []struct {
		Sha         string `json:"sha"`
		Filename    string `json:"filename"`
		Status      string `json:"status"`
		Additions   int    `json:"additions"`
		Deletions   int    `json:"deletions"`
		Changes     int    `json:"changes"`
		BlobURL     string `json:"blob_url"`
		RawURL      string `json:"raw_url"`
		ContentsURL string `json:"contents_url"`
		Patch       string `json:"patch"`
	} `json:"files"`
}

func NewGitHubReleasesFetcher(datasource *config.Datasource) *GitHubReleasesFetcher {
	return &GitHubReleasesFetcher{
		FetcherBase: *NewFetcherBase(1, 100, datasource.MaxReleases, datasource.Credentials),
	}
}

func NewGitHubTagsFetcher(datasource *config.Datasource) *GitHubTagsFetcher {
	return &GitHubTagsFetcher{
		FetcherBase: *NewFetcherBase(1, 100, datasource.MaxReleases, datasource.Credentials),
	}
}

func (rf *GitHubReleasesFetcher) GetReleases(pkg string) (*internal.ReleasesResponse, error) {
	parts := regexp.MustCompile("[/]").Split(pkg, -1)
	if len(parts) != 2 {
		return nil, &FetcherError{Err: fmt.Errorf("expected two parts, separated by '/' in GitHub releases package, got %s", pkg), IsParameterError: true}
	}
	return rf.getReleases(parts[0], parts[1])
}

func (rf *GitHubReleasesFetcher) getSearchUrl(owner, repo string, page int) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/releases?page=%d&per_page=%d",
		url.PathEscape(owner), url.PathEscape(repo), page, rf.perPage)
}

func (rf *GitHubReleasesFetcher) getReleases(owner, repo string) (*internal.ReleasesResponse, error) {
	releasesResponse := internal.ReleasesResponse{
		Releases:  make([]internal.Release, 0),
		SourceURL: new("https://github.com/" + url.PathEscape(owner) + "/" + url.PathEscape(repo)),
	}
	err := rf.paginate(rf, &releasesResponse, owner, repo)
	if err != nil {
		return nil, err
	}
	return &releasesResponse, nil
}

func (rf *GitHubReleasesFetcher) extractReleases(_, _, jsonResponse string) ([]internal.Release, error) {
	var resp fullGitHubReleasesResponse
	err := json.Unmarshal([]byte(jsonResponse), &resp)
	if err != nil {
		return nil, err
	}
	releases := make([]internal.Release, 0, len(resp))
	for _, result := range resp {
		release := internal.Release{
			Version:          result.TagName,
			ReleaseTimestamp: result.PublishedAt,
			ChangelogURL:     &result.HTMLURL,
		}
		releases = append(releases, release)
	}
	return releases, nil
}

func (rf *GitHubTagsFetcher) GetReleases(pkg string) (*internal.ReleasesResponse, error) {
	parts := regexp.MustCompile("[/]").Split(pkg, -1)
	if len(parts) != 2 {
		return nil, &FetcherError{Err: fmt.Errorf("expected two parts, separated by '/' in GitHub releases package, got %s", pkg), IsParameterError: true}
	}
	return rf.getReleases(parts[0], parts[1])
}

func (rf *GitHubTagsFetcher) getSearchUrl(owner, repo string, page int) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/tags?page=%d&per_page=%d",
		url.PathEscape(owner), url.PathEscape(repo), page, rf.perPage)
}

func (rf *GitHubTagsFetcher) getReleases(owner, repo string) (*internal.ReleasesResponse, error) {
	releasesResponse := internal.ReleasesResponse{
		Releases:  make([]internal.Release, 0),
		SourceURL: new("https://github.com/" + url.PathEscape(owner) + "/" + url.PathEscape(repo)),
	}
	err := rf.paginate(rf, &releasesResponse, owner, repo)
	if err != nil {
		return nil, err
	}
	return &releasesResponse, nil
}

func (rf *GitHubTagsFetcher) extractReleases(owner, repo, jsonResponse string) ([]internal.Release, error) {
	var resp fullGitHubTagsResponse
	err := json.Unmarshal([]byte(jsonResponse), &resp)
	if err != nil {
		return nil, err
	}
	releases := make([]internal.Release, 0, len(resp))
	for _, result := range resp {
		var timestamp time.Time
		timestamp, err = rf.fetchTimestamp(owner, repo, result.Commit.Sha)
		if err != nil {
			return nil, err
		}
		release := internal.Release{
			Version:          result.Name,
			ReleaseTimestamp: timestamp,
		}
		releases = append(releases, release)
	}
	return releases, nil
}

func (rf *GitHubTagsFetcher) fetchTimestamp(owner, repo, commitSha string) (time.Time, error) {
	const datasource = "github"
	timestamp := state.GetCommitTimestamp(datasource, commitSha)
	if timestamp != nil {
		return *timestamp, nil
	}
	commitUrl := rf.getCommitUrl(owner, repo, commitSha)
	jsonResponse, err := webclient.Get(commitUrl, rf.credentials)
	if err != nil {
		return time.Time{}, err
	}
	var resp fullGitHubCommitResponse
	err = json.Unmarshal([]byte(jsonResponse), &resp)
	if err != nil {
		return time.Time{}, err
	}
	timestamp = &resp.Commit.Committer.Date
	state.PutCommitTimestamp(datasource, commitSha, *timestamp)
	return *timestamp, nil
}

func (rf *GitHubTagsFetcher) getCommitUrl(owner, repo, commitSha string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s",
		url.PathEscape(owner), url.PathEscape(repo), url.PathEscape(commitSha))
}
