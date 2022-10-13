package ghsearch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

type SearchRepoResult struct {
	TotalCount        int                     `json:"total_count"`
	IncompleteResults bool                    `json:"incomplete_results"`
	Items             []SearchRepoResultItems `json:"items"`
}

type SearchRepoResultItems struct {
	ID       int    `json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"owner"`
	Private          bool      `json:"private"`
	HTMLURL          string    `json:"html_url"`
	Description      string    `json:"description"`
	Fork             bool      `json:"fork"`
	URL              string    `json:"url"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	PushedAt         time.Time `json:"pushed_at"`
	Homepage         string    `json:"homepage"`
	Size             int       `json:"size"`
	StargazersCount  int       `json:"stargazers_count"`
	WatchersCount    int       `json:"watchers_count"`
	Language         string    `json:"language"`
	ForksCount       int       `json:"forks_count"`
	OpenIssuesCount  int       `json:"open_issues_count"`
	MasterBranch     string    `json:"master_branch"`
	DefaultBranch    string    `json:"default_branch"`
	Score            float32   `json:"score"`
	ArchiveURL       string    `json:"archive_url"`
	AssigneesURL     string    `json:"assignees_url"`
	BlobsURL         string    `json:"blobs_url"`
	BranchesURL      string    `json:"branches_url"`
	CollaboratorsURL string    `json:"collaborators_url"`
	CommentsURL      string    `json:"comments_url"`
	CommitsURL       string    `json:"commits_url"`
	CompareURL       string    `json:"compare_url"`
	ContentsURL      string    `json:"contents_url"`
	ContributorsURL  string    `json:"contributors_url"`
	DeploymentsURL   string    `json:"deployments_url"`
	DownloadsURL     string    `json:"downloads_url"`
	EventsURL        string    `json:"events_url"`
	ForksURL         string    `json:"forks_url"`
	GitCommitsURL    string    `json:"git_commits_url"`
	GitRefsURL       string    `json:"git_refs_url"`
	GitTagsURL       string    `json:"git_tags_url"`
	GitURL           string    `json:"git_url"`
	IssueCommentURL  string    `json:"issue_comment_url"`
	IssueEventsURL   string    `json:"issue_events_url"`
	IssuesURL        string    `json:"issues_url"`
	KeysURL          string    `json:"keys_url"`
	LabelsURL        string    `json:"labels_url"`
	LanguagesURL     string    `json:"languages_url"`
	MergesURL        string    `json:"merges_url"`
	MilestonesURL    string    `json:"milestones_url"`
	NotificationsURL string    `json:"notifications_url"`
	PullsURL         string    `json:"pulls_url"`
	ReleasesURL      string    `json:"releases_url"`
	SSHURL           string    `json:"ssh_url"`
	StargazersURL    string    `json:"stargazers_url"`
	StatusesURL      string    `json:"statuses_url"`
	SubscribersURL   string    `json:"subscribers_url"`
	SubscriptionURL  string    `json:"subscription_url"`
	TagsURL          string    `json:"tags_url"`
	TeamsURL         string    `json:"teams_url"`
	TreesURL         string    `json:"trees_url"`
	CloneURL         string    `json:"clone_url"`
	MirrorURL        string    `json:"mirror_url"`
	HooksURL         string    `json:"hooks_url"`
	SvnURL           string    `json:"svn_url"`
	Forks            int       `json:"forks"`
	OpenIssues       int       `json:"open_issues"`
	Watchers         int       `json:"watchers"`
	HasIssues        bool      `json:"has_issues"`
	HasProjects      bool      `json:"has_projects"`
	HasPages         bool      `json:"has_pages"`
	HasWiki          bool      `json:"has_wiki"`
	HasDownloads     bool      `json:"has_downloads"`
	Archived         bool      `json:"archived"`
	Disabled         bool      `json:"disabled"`
	Visibility       string    `json:"visibility"`
	License          struct {
		Key     string `json:"key"`
		Name    string `json:"name"`
		URL     string `json:"url"`
		SpdxID  string `json:"spdx_id"`
		NodeID  string `json:"node_id"`
		HTMLURL string `json:"html_url"`
	} `json:"license"`
}

// SearchRepo search repositories
// https://docs.github.com/cn/rest/search#search-repositories
/*
	curl \
	-H "Accept: application/vnd.github+json" \
	-H "Authorization: Bearer xxx" \
	https://api.github.com/search/repositories

*/
func SearchRepo(token string, page int, lang string, keywords ...string) ([]SearchRepoResultItems, error) {
	client := resty.New()

	q := "language:" + lang
	for _, v := range keywords {
		q += "+"
		q += v
	}
	resp, err := client.R().
		SetHeader("Accept", "application/vnd.github+json").
		SetHeader("Authorization", "Bearer "+token).
		Get("https://api.github.com/search/repositories?" + "q=" + q + "&sort=star&order=desc&page=" + strconv.Itoa(page))

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("error:%s", string(resp.Body()))
	}

	//fmt.Println(string(resp.Body()))
	sr := SearchRepoResult{}
	if err := json.Unmarshal(resp.Body(), &sr); err != nil {
		return nil, err
	}
	sort.Slice(sr.Items, func(i, j int) bool {
		return sr.Items[i].StargazersCount > sr.Items[j].StargazersCount
	})
	return sr.Items, nil
}

type SearchCodeResult struct {
	TotalCount        int                     `json:"total_count"`
	IncompleteResults bool                    `json:"incomplete_results"`
	Items             []SearchCodeResultItems `json:"items"`
}
type SearchCodeResultItems struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	Sha        string `json:"sha"`
	URL        string `json:"url"`
	GitURL     string `json:"git_url"`
	HTMLURL    string `json:"html_url"`
	Repository struct {
		ID       int    `json:"id"`
		NodeID   string `json:"node_id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Owner    struct {
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
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"owner"`
		Private          bool   `json:"private"`
		HTMLURL          string `json:"html_url"`
		Description      string `json:"description"`
		Fork             bool   `json:"fork"`
		URL              string `json:"url"`
		ForksURL         string `json:"forks_url"`
		KeysURL          string `json:"keys_url"`
		CollaboratorsURL string `json:"collaborators_url"`
		TeamsURL         string `json:"teams_url"`
		HooksURL         string `json:"hooks_url"`
		IssueEventsURL   string `json:"issue_events_url"`
		EventsURL        string `json:"events_url"`
		AssigneesURL     string `json:"assignees_url"`
		BranchesURL      string `json:"branches_url"`
		TagsURL          string `json:"tags_url"`
		BlobsURL         string `json:"blobs_url"`
		GitTagsURL       string `json:"git_tags_url"`
		GitRefsURL       string `json:"git_refs_url"`
		TreesURL         string `json:"trees_url"`
		StatusesURL      string `json:"statuses_url"`
		LanguagesURL     string `json:"languages_url"`
		StargazersURL    string `json:"stargazers_url"`
		ContributorsURL  string `json:"contributors_url"`
		SubscribersURL   string `json:"subscribers_url"`
		SubscriptionURL  string `json:"subscription_url"`
		CommitsURL       string `json:"commits_url"`
		GitCommitsURL    string `json:"git_commits_url"`
		CommentsURL      string `json:"comments_url"`
		IssueCommentURL  string `json:"issue_comment_url"`
		ContentsURL      string `json:"contents_url"`
		CompareURL       string `json:"compare_url"`
		MergesURL        string `json:"merges_url"`
		ArchiveURL       string `json:"archive_url"`
		DownloadsURL     string `json:"downloads_url"`
		IssuesURL        string `json:"issues_url"`
		PullsURL         string `json:"pulls_url"`
		MilestonesURL    string `json:"milestones_url"`
		NotificationsURL string `json:"notifications_url"`
		LabelsURL        string `json:"labels_url"`
		DeploymentsURL   string `json:"deployments_url"`
		ReleasesURL      string `json:"releases_url"`
	} `json:"repository"`
	Score float32 `json:"score"`
}

// SearchCode search code
// https://docs.github.com/cn/rest/search#search-code
/*
	curl \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer <YOUR-TOKEN>" \
  https://api.github.com/search/code

*/
func SearchCode(token string, page int, lang string, keywords ...string) ([]SearchCodeResultItems, error) {
	client := resty.New()

	q := "language:" + lang
	for _, v := range keywords {
		q += "+"
		q += v
	}
	resp, err := client.R().
		SetHeader("Accept", "application/vnd.github+json").
		SetHeader("Authorization", "Bearer "+token).
		Get("https://api.github.com/search/code?" + "q=" + q + "&sort=star&order=desc&page=" + strconv.Itoa(page))

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("error:%s", string(resp.Body()))
	}

	//fmt.Println(string(resp.Body()))
	sr := SearchCodeResult{}
	if err := json.Unmarshal(resp.Body(), &sr); err != nil {
		return nil, err
	}
	sort.Slice(sr.Items, func(i, j int) bool {
		return sr.Items[i].Score > sr.Items[j].Score
	})
	return sr.Items, nil
}
