package models

import (
	"time"
)

// WebhookPayload represents the top-level webhook event payload.
type PRWebhookPayload struct {
	Action            string       `json:"action"`
	CommitID          string       `json:"commit_id"`
	Number            int          `json:"number"`
	PullRequest       *PullRequest `json:"pull_request"`
	Repository        Repository   `json:"repository"`
	RequestedReviewer *User        `json:"requested_reviewer"`
	Review            interface{}  `json:"review"`
	Sender            User         `json:"sender"`
}

type RepoCreateWebhookPayload struct {
	Action       string     `json:"action"`
	Organization User       `json:"organization"`
	Repository   Repository `json:"repository"`
	Sender       User       `json:"sender"`
}

type RepoPushWebhookPayload struct {
	Ref          string     `json:"ref"`
	Before       string     `json:"before"`
	After        string     `json:"after"`
	CompareUrl   string     `json:"compare_url"`
	Commits      []Commit   `json:"commits"`
	TotalCommits int        `json:"total_commits"`
	HeadCommit   Commit     `json:"head_commit"`
	Repository   Repository `json:"repository"`
	Pusher       User       `json:"pusher"`
	Sender       User       `json:"sender"`
}

type CommitUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type Commit struct {
	ID           string     `json:"id"`
	Message      string     `json:"message"`
	URL          string     `json:"url"`
	Author       CommitUser `json:"author"`
	Committer    CommitUser `json:"committer"`
	Verification any        `json:"verification"`
	Timestamp    string     `json:"timestamp"`
	Added        []string   `json:"added"`
	Removed      []string   `json:"removed"`
	Modified     []string   `json:"modified"`
}

type CommitFile struct {
	FileName string `json:"filename"`
	Status   string `json:"status"`
}

type Stat struct {
	Additions int `json:"additions"`
	Deletions int `json:"deletions"`
	Total     int `json:"total"`
}

// PullRequest represents a pull request object.
type PullRequest struct {
	Additions               int64      `json:"additions"`
	AllowMaintainerEdit     bool       `json:"allow_maintainer_edit"`
	Assignee                *User      `json:"assignee"`
	Assignees               []User     `json:"assignees"`
	Base                    PRBranch   `json:"base"`
	Body                    string     `json:"body"`
	ChangedFiles            int        `json:"changed_files"`
	ClosedAt                *time.Time `json:"closed_at"`
	Comments                int        `json:"comments"`
	CreatedAt               time.Time  `json:"created_at"`
	Deletions               int64      `json:"deletions"`
	DiffURL                 string     `json:"diff_url"`
	Draft                   bool       `json:"draft"`
	DueDate                 *time.Time `json:"due_date"`
	Head                    PRBranch   `json:"head"`
	HTMLURL                 string     `json:"html_url"`
	ID                      int64      `json:"id"`
	IsLocked                bool       `json:"is_locked"`
	Labels                  []Label    `json:"labels"`
	MergeBase               string     `json:"merge_base"`
	MergeCommitSHA          *string    `json:"merge_commit_sha"`
	Mergeable               bool       `json:"mergeable"`
	Merged                  bool       `json:"merged"`
	MergedAt                *time.Time `json:"merged_at"`
	MergedBy                *User      `json:"merged_by"`
	Milestone               *Milestone `json:"milestone"`
	Number                  int        `json:"number"`
	PatchURL                string     `json:"patch_url"`
	PinOrder                int        `json:"pin_order"`
	RequestedReviewers      []User     `json:"requested_reviewers"`
	RequestedReviewersTeams []Team     `json:"requested_reviewers_teams"`
	State                   string     `json:"state"`
	Title                   string     `json:"title"`
	UpdatedAt               time.Time  `json:"updated_at"`
	URL                     string     `json:"url"`
	User                    User       `json:"user"`
}

// PRBranch represents the base or head branch of a pull request.
type PRBranch struct {
	Label  string     `json:"label"`
	Ref    string     `json:"ref"`
	Repo   Repository `json:"repo"`
	RepoID int64      `json:"repo_id"`
	SHA    string     `json:"sha"`
}

// Repository represents a git repository.
type Repository struct {
	AllowFastForwardOnlyMerge     bool            `json:"allow_fast_forward_only_merge"`
	AllowManualMerge              bool            `json:"allow_manual_merge"`
	AllowMergeCommits             bool            `json:"allow_merge_commits"`
	AllowRebase                   bool            `json:"allow_rebase"`
	AllowRebaseExplicit           bool            `json:"allow_rebase_explicit"`
	AllowRebaseUpdate             bool            `json:"allow_rebase_update"`
	AllowSquashMerge              bool            `json:"allow_squash_merge"`
	Archived                      bool            `json:"archived"`
	ArchivedAt                    time.Time       `json:"archived_at"`
	AutodetectManualMerge         bool            `json:"autodetect_manual_merge"`
	AvatarURL                     string          `json:"avatar_url"`
	CloneURL                      string          `json:"clone_url"`
	CreatedAt                     time.Time       `json:"created_at"`
	DefaultAllowMaintainerEdit    bool            `json:"default_allow_maintainer_edit"`
	DefaultBranch                 string          `json:"default_branch"`
	DefaultDeleteBranchAfterMerge bool            `json:"default_delete_branch_after_merge"`
	DefaultMergeStyle             string          `json:"default_merge_style"`
	Description                   string          `json:"description"`
	Empty                         bool            `json:"empty"`
	Fork                          bool            `json:"fork"`
	ForksCount                    int             `json:"forks_count"`
	FullName                      string          `json:"full_name"`
	HasActions                    bool            `json:"has_actions"`
	HasCode                       bool            `json:"has_code"`
	HasIssues                     bool            `json:"has_issues"`
	HasPackages                   bool            `json:"has_packages"`
	HasProjects                   bool            `json:"has_projects"`
	HasPullRequests               bool            `json:"has_pull_requests"`
	HasReleases                   bool            `json:"has_releases"`
	HasWiki                       bool            `json:"has_wiki"`
	HTMLURL                       string          `json:"html_url"`
	ID                            int64           `json:"id"`
	IgnoreWhitespaceConflicts     bool            `json:"ignore_whitespace_conflicts"`
	Internal                      bool            `json:"internal"`
	InternalTracker               InternalTracker `json:"internal_tracker"`
	Language                      string          `json:"language"`
	LanguagesURL                  string          `json:"languages_url"`
	Licenses                      []string        `json:"licenses"`
	Link                          string          `json:"link"`
	Mirror                        bool            `json:"mirror"`
	MirrorInterval                string          `json:"mirror_interval"`
	MirrorUpdated                 time.Time       `json:"mirror_updated"`
	Name                          string          `json:"name"`
	ObjectFormatName              string          `json:"object_format_name"`
	OpenIssuesCount               int             `json:"open_issues_count"`
	OpenPRCounter                 int             `json:"open_pr_counter"`
	OriginalURL                   string          `json:"original_url"`
	Owner                         User            `json:"owner"`
	Permissions                   Permissions     `json:"permissions"`
	Private                       bool            `json:"private"`
	ProjectsMode                  string          `json:"projects_mode"`
	ReleaseCounter                int             `json:"release_counter"`
	Size                          int64           `json:"size"`
	SSHURL                        string          `json:"ssh_url"`
	StarsCount                    int             `json:"stars_count"`
	Template                      bool            `json:"template"`
	Topics                        []string        `json:"topics"`
	UpdatedAt                     time.Time       `json:"updated_at"`
	URL                           string          `json:"url"`
	WatchersCount                 int             `json:"watchers_count"`
	Website                       string          `json:"website"`
}

// InternalTracker holds the internal issue tracker settings for a repository.
type InternalTracker struct {
	AllowOnlyContributorsToTrackTime bool `json:"allow_only_contributors_to_track_time"`
	EnableIssueDependencies          bool `json:"enable_issue_dependencies"`
	EnableTimeTracker                bool `json:"enable_time_tracker"`
}

// Permissions represents the permissions the authenticated user has on a repository.
type Permissions struct {
	Admin bool `json:"admin"`
	Pull  bool `json:"pull"`
	Push  bool `json:"push"`
}

// User represents a git user (owner, sender, assignee, reviewer, etc.).
type User struct {
	Active            bool      `json:"active"`
	AvatarURL         string    `json:"avatar_url"`
	Created           time.Time `json:"created"`
	Description       string    `json:"description"`
	Email             string    `json:"email"`
	FollowersCount    int       `json:"followers_count"`
	FollowingCount    int       `json:"following_count"`
	FullName          string    `json:"full_name"`
	HTMLURL           string    `json:"html_url"`
	ID                int64     `json:"id"`
	IsAdmin           bool      `json:"is_admin"`
	Language          string    `json:"language"`
	LastLogin         time.Time `json:"last_login"`
	Location          string    `json:"location"`
	Login             string    `json:"login"`
	LoginName         string    `json:"login_name"`
	ProhibitLogin     bool      `json:"prohibit_login"`
	Restricted        bool      `json:"restricted"`
	SourceID          int64     `json:"source_id"`
	StarredReposCount int       `json:"starred_repos_count"`
	Username          string    `json:"username"`
	Visibility        string    `json:"visibility"`
	Website           string    `json:"website"`
}

// Label represents a pull request label.
type Label struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
	Exclusive   bool   `json:"exclusive"`
}

// Milestone represents a pull request milestone.
type Milestone struct {
	ID           int64      `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	State        string     `json:"state"`
	OpenIssues   int        `json:"open_issues"`
	ClosedIssues int        `json:"closed_issues"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	ClosedAt     *time.Time `json:"closed_at"`
	DueOn        *time.Time `json:"due_on"`
}

// Team represents a team requested for review.
type Team struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Organization string `json:"organization"`
}

type GitDiffApiUrlOptions struct {
	OwnerName         string
	RepositoryName    string
	PullRequestNumber int
}

type PostReviewApiUrlOptions struct {
	OwnerName         string
	RepositoryName    string
	PullRequestNumber int
}

type GetGitFileContentApiUrlOptions struct {
	FilePath       string
	OwnerName      string
	RepositoryName string
}

type FileDiff struct {
	Path  string
	Hunks string
}

type RepoTechnology string

const (
	React   RepoTechnology = "react"
	Vue     RepoTechnology = "vue"
	Django  RepoTechnology = "django"
	FastAPI RepoTechnology = "fastapi"
	Nest    RepoTechnology = "nest"
	Gin     RepoTechnology = "gin"
	Php     RepoTechnology = "php"
	Unknown RepoTechnology = "unknown"
)

type GetGitFileContentQueryOptions struct {
	FilePath       string `json:"file_path"`
	OwnerName      string `json:"owner"`
	RepositoryName string `json:"repo"`
}

type GetProjectDirectoryListQueryOptions struct {
	OwnerName      string
	RepositoryName string
	RootPath       *string
}

type Author struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Committer Author

type Dates struct {
	Author    string
	Committer string
}

type CreateFileInRepoMutationBody struct {
	Author    Author    `json:"author"`
	Branch    string    `json:"branch"`
	Committer Committer `json:"committer"`
	Content   string    `json:"content"`
	Dates     Dates     `json:"dates"`
	Message   string    `json:"message"`
	NewBranch string    `json:"new_branch"`
	SignOff   bool      `json:"signoff"`
}

type CreateFileInRepoMutationOptions struct {
	OwnerName      string                       `json:"owner"`
	RepositoryName string                       `json:"repo"`
	FilePath       string                       `json:"filepath"`
	Body           CreateFileInRepoMutationBody `json:"body"`
}

type CreateFileInRepoApiUrlMutationOptions struct {
	OwnerName      string
	RepositoryName string
	FilePath       string
}

type Operation string

const (
	Create Operation = "create"
	Update Operation = "update"
	Delete Operation = "delete"
	Rename Operation = "rename"
	Upload Operation = "upload"
)

type ChangeFileOperation struct {
	Content   string    `json:"content"`
	FromPath  string    `json:"from_path,omitempty"`
	Operation Operation `json:"operation"`
	Path      string    `json:"path"`
	Sha       string    `json:"sha,omitempty"`
}

type ModifyFilesInRepoApiUrlOptions struct {
	OwnerName      string `json:"owner"`
	RepositoryName string `json:"repo"`
}

type ModifyFilesInRepoMutationOptions struct {
	OwnerName      string                `json:"owner"`
	RepositoryName string                `json:"repo"`
	Body           ModifyFilesInRepoBody `json:"body"`
}

type ModifyFilesInRepoBody struct {
	Author    *Author               `json:"author,omitempty"`
	Branch    string                `json:"branch"`
	Committer *Committer            `json:"committer,omitempty"`
	Files     []ChangeFileOperation `json:"files"`
	Dates     *Dates                `json:"dates,omitempty"`
	Message   string                `json:"message,omitempty"`
	NewBranch string                `json:"new_branch,omitempty"`
	SignOff   bool                  `json:"signoff,omitempty"`
}

type GitDirectory struct {
	Name              string `json:"name"`
	Path              string `json:"path"`
	Sha               string `json:"sha"`
	LastCommitSha     string `json:"last_commit_sha"`
	LastCommitterDate string `json:"last_committer_date"`
	LastAuthorDate    string `json:"last_author_date"`
	Type              string `json:"type"`
	Size              int64  `json:"size"`
	Encoding          string `json:"encoding"`
	Content           string `json:"content"`
	Target            string `json:"target"`
	Url               string `json:"url"`
	HtmlUrl           string `json:"html_url"`
	GitUrl            string `json:"git_url"`
	DownloadUrl       string `json:"download_url"`
	SubmoduleGitUrl   string `json:"submodule_git_url"`
	Links             Links  `json:"_links"`
}

type Links struct {
	Self string `json:"self"`
	Git  string `json:"git"`
	Html string `json:"html"`
}

type GetProjectDirectoryListQueryResponse struct {
	DirContents  []GitDirectory `json:"dir_contents"`
	FileContents GitFileMeta    `json:"file_contents"`
}

type GitFileMeta struct {
	// this struct has more field in real life, for now we just need these
	Content string `json:"content"`
	Sha     string `json:"sha"`
	HtmlUrl string `json:"html_url"`
}
