package models

type AiConfiuguration struct {
	Model   string
	Prompt  string
	System  string
	Stream  bool
	Options map[string]interface{}
}

type GetModelReviewMutationOptions struct {
	Files []FileDiff
}

// owner, repo, pr_number, comments, head_sha
type PostReviewToGitMutationOptions struct {
	Owner             string
	Repository        string
	PullRequestNumber int
	Comments          []AIReviewComment
	HeadSha           string
}

type AiModelProperties struct {
	Name string
	// maximum number of characters that gets sent to model (also leaving some space for system prompt)
	ModelMaximumCharacters int
}

type GetModelReviewMutationResponse struct {
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	Response           string `json:"response"`
	Done               bool   `json:"done"`
	Context            []int  `json:"context"`
	TotalDuration      int    `json:"total_duration"`
	LoadDuration       int    `json:"load_duration"`
	PromptEvalCount    int    `json:"prompt_eval_count"`
	PromptEvalDuration int    `json:"prompt_eval_duration"`
	EvalCount          int    `json:"eval_count"`
	EvalDuration       int    `json:"eval_duration"`
}

type GitReviewPayload struct {
	CommitId string             `json:"commit_id"`
	Body     string             `json:"body"`
	Event    string             `json:"event"`
	Comments []GitReviewComment `json:"comments"`
}

type AIReviewComment struct {
	Comment  string `json:"comment"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Severity string `json:"severity"`
}

type GitReviewComment struct {
	Path string `json:"path"`
	Line int    `json:"line"`
	Side string `json:"side"`
	Body string `json:"body"`
}
