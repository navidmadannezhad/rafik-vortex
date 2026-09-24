package controller

import (
	"airun/code-reviewer/models"
	"airun/code-reviewer/services"
	"airun/code-reviewer/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequestPRReview(requestContext *gin.Context) {
	body := models.WebhookPayload{}
	headers := map[string]interface{}{}

	err := requestContext.BindJSON(&body)
	if err != nil {
		requestContext.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	err2 := requestContext.BindHeader(&headers)
	if err2 != nil {
		requestContext.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error": err2.Error(),
			},
		)
		return
	}

	isPullRequestEvent := body.PullRequest != nil
	if !isPullRequestEvent {
		requestContext.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ResolveErrorMessage("Not a pull request event"),
		)
		return
	}

	// pr = payload["pull_request"]
	// repo = payload["repository"]
	// owner = repo["owner"]["login"]
	// repo_name = repo["name"]
	// pr_number = pr["number"]
	// head_sha = pr["head"]["sha"]

	pullRequest := body.PullRequest
	repository := body.Repository
	owner := repository.Owner
	repositoryName := repository.Name
	pullRequestNumber := pullRequest.Number
	// headSha := pullRequest.Head.SHA
	fmt.Printf("Received PR! Number: %d, Repository: %s\n", pullRequestNumber, repositoryName)

	resp, err := services.GetPullRequestDiffQuery(
		models.GitDiffApiUrlOptions{
			PullRequestNumber: pullRequestNumber,
			RepositoryName:    repository.Name,
			OwnerName:         owner.Login,
		},
	)
	if err != nil {
		requestContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ResolveErrorMessage(err.Error()),
		)
		return
	}

	rawResponseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		requestContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ResolveErrorMessage(err.Error()),
		)
		return
	}

	files := utils.ParsePullRequestDiff(string(rawResponseBody))

	var aiResponseBody = models.GetModelReviewMutationResponse{}
	payload := models.GetModelReviewMutationOptions{
		Files: files,
	}
	aiResponse, err := services.GetModelReviewMutation(payload)
	if err != nil {
		requestContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ResolveErrorMessage(err.Error()),
		)
		return
	}
	defer aiResponse.Body.Close()

	err = json.NewDecoder(aiResponse.Body).Decode(&aiResponseBody)
	if err != nil {
		requestContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ResolveErrorMessage(err.Error()),
		)
		return
	}
	comments := utils.ParseAIReviewComments(aiResponseBody.Response)

	postReviewResponse, err := services.PostReviewToGitMutation(
		models.PostReviewToGitMutationOptions{
			Owner:             owner.Login,
			Repository:        repositoryName,
			PullRequestNumber: pullRequestNumber,
			HeadSha:           pullRequest.Head.SHA,
			Comments:          comments,
		},
	)
	if err != nil {
		requestContext.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ResolveErrorMessage(err.Error()),
		)
		return
	}
	defer postReviewResponse.Body.Close()

	requestContext.JSON(
		http.StatusOK,
		utils.ResolveSuccessMessage("Review completed"),
	)
}
