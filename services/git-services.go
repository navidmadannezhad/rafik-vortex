package services

import (
	"airun/code-reviewer/models"
	"airun/code-reviewer/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

func GetPullRequestDiffQuery(options models.GitDiffApiUrlOptions) (*http.Response, error) {
	url := utils.GetDiffApiUrl(options)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	gitToken := utils.GetEnv("GIT_TOKEN", "")
	req.Header.Add("Authorization", fmt.Sprintf("token %v", gitToken))
	req.Header.Add("Accept", "application/vnd.github+json")
	response, err := HttpClient.Do(req)
	if err != nil {
		return &http.Response{}, err
	}
	return response, nil
}

func PostReviewToGitMutation(options models.PostReviewToGitMutationOptions) (*http.Response, error) {

	comments := options.Comments

	hasCritical := utils.CommentListHasCritical(comments)
	var event string
	if hasCritical {
		event = "REQUEST_CHANGES"
	} else {
		event = "COMMENT"
	}

	var reviewComments []models.GitReviewComment
	for _, comment := range comments {
		line := comment.Line
		lineNumber := math.Max(float64(line), 1)

		gitReviewComment := models.GitReviewComment{
			Path: comment.File,
			Line: int(lineNumber),
			Side: "RIGHT",
			Body: comment.Comment,
		}
		reviewComments = append(reviewComments, gitReviewComment)
	}

	var bodyText string
	if len(reviewComments) == 0 {
		bodyText = "No Issues Found"
	} else {
		bodyText = fmt.Sprintf(`
		Automated AI Code Review - Number of issues found: %v
		`, len(reviewComments))
	}

	gitReviewPayload := models.GitReviewPayload{
		CommitId: options.HeadSha,
		Body:     bodyText,
		Event:    event,
		Comments: reviewComments,
	}
	jsonPayload, err := json.Marshal(gitReviewPayload)
	if err != nil {
		return nil, err
	}

	postReviewApiUrlOptions := models.PostReviewApiUrlOptions{
		OwnerName:         options.Owner,
		RepositoryName:    options.Repository,
		PullRequestNumber: options.PullRequestNumber,
	}
	url := utils.PostReviewApiUrl(postReviewApiUrlOptions)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, err
	}
	gitToken := utils.GetEnv("GIT_TOKEN", "")
	req.Header.Add("Authorization", fmt.Sprintf("token %v", gitToken))
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")
	response, err := HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return response, nil
}
