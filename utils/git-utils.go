package utils

import (
	"airun/code-reviewer/configs"
	"airun/code-reviewer/models"
	"fmt"
	"log"
	"strings"
)

var unallowedParseFormats = []string{".lock", "-lock.json", ".png", ".jpg", ".svg", ".min.js", ".map"}

func ParsePullRequestDiff(rawDiff string) []models.FileDiff {
	var files []models.FileDiff
	var currentFile *models.FileDiff

	for _, line := range strings.Split(rawDiff, "\n") {
		if strings.HasPrefix(line, "diff --git") {
			currentFile = nil

			parts := strings.SplitN(line, " b/", 2)
			if len(parts) < 2 {
				continue
			}

			filepath := strings.TrimSpace(parts[1])
			if filepath == "" || hasAnySuffix(filepath, unallowedParseFormats) {
				continue
			}

			files = append(files, models.FileDiff{Path: filepath})
			currentFile = &files[len(files)-1]
		} else if currentFile != nil {
			currentFile.Hunks += line + "\n"
		}
	}

	return files
}

func hasAnySuffix(s string, suffixes []string) bool {
	for _, suf := range suffixes {
		if strings.HasSuffix(s, suf) {
			return true
		}
	}
	return false
}

func GetDiffApiUrl(options models.GitDiffApiUrlOptions) string {
	gitUrl := GetEnv("GIT_URL", "https://github.com")
	url := fmt.Sprintf("%v/api/v1/repos/%v/%v/pulls/%v.diff", gitUrl, options.OwnerName, options.RepositoryName, options.PullRequestNumber)
	return url
}

func PostReviewApiUrl(options models.PostReviewApiUrlOptions) string {
	gitUrl := GetEnv("GIT_URL", "https://github.com")
	url := fmt.Sprintf("%v/api/v1/repos/%v/%v/pulls/%v/reviews", gitUrl, options.OwnerName, options.RepositoryName, options.PullRequestNumber)
	return url
}

func GetGitFileContentApiUrl(options models.GetGitFileContentApiUrlOptions) string {
	gitUrl := GetEnv("GIT_URL", "https://github.com")
	log.Printf("GetGitFileContentApiUrl: %v/%v/%v/raw/branch/main/%v", gitUrl, options.OwnerName, options.RepositoryName, options.FilePath)
	url := fmt.Sprintf("%v/%v/%v/raw/branch/main/%v", gitUrl, options.OwnerName, options.RepositoryName, options.FilePath)
	return url
}

func GetProjectDirectoryListApiUrl(options models.GetProjectDirectoryListQueryOptions) string {
	rootPath := ""
	if options.RootPath != nil {
		rootPath = *options.RootPath
	}
	gitUrl := GetEnv("GIT_URL", "https://github.com")
	url := fmt.Sprintf("%v/api/v1/repos/%v/%v/contents-ext/%v", gitUrl, options.OwnerName, options.RepositoryName, rootPath)
	return url
}

func CreateFileInRepoApiUrl(options models.CreateFileInRepoApiUrlMutationOptions) string {
	gitUrl := GetEnv("GIT_URL", "https://github.com")
	url := fmt.Sprintf("%v/api/v1/repos/%v/%v/contents/%v", gitUrl, options.OwnerName, options.RepositoryName, options.FilePath)
	return url
}

func ModifyFilesInRepoApiUrl(options models.ModifyFilesInRepoApiUrlOptions) string {
	gitUrl := GetEnv("GIT_URL", "https://github.com")
	url := fmt.Sprintf("%v/api/v1/repos/%v/%v/contents/", gitUrl, options.OwnerName, options.RepositoryName)
	return url
}

func FindTechnologyDetectionFile(gitFiles []models.GitDirectory) (path string) {
	for _, gitFile := range gitFiles {
		if gitFile.Name == "package.json" {
			return "package.json"
		}
		if gitFile.Name == "requirements.txt" {
			return "requirements.txt"
		}
		if gitFile.Name == "go.mod" {
			return "go.mod"
		}
	}

	return ""
}

func FindRepoStack(techFileContent string) (stack models.RepoTechnology) {
	switch true {
	case strings.Contains(techFileContent, "react") ||
		strings.Contains(techFileContent, "next"):
		return models.React

	case strings.Contains(techFileContent, "vue"):
		return models.Vue

	case strings.Contains(techFileContent, "django") ||
		strings.Contains(techFileContent, "djangorestframework"):
		return models.Django

	case strings.Contains(techFileContent, "fastapi"):
		return models.FastAPI

	case strings.Contains(techFileContent, "nest") ||
		strings.Contains(techFileContent, "djangorestframework"):
		return models.Nest

	case strings.Contains(techFileContent, "go"):
		return models.Gin

	default:
		return models.Unknown
	}
}

func ExistsUpdateAIArg(payload models.RepoPushWebhookPayload) bool {
	exists := false
	for _, commit := range payload.Commits {
		if strings.Contains(commit.Message, string(configs.UpdateAIInstructsArgCommand)) {
			exists = true
		}
	}

	return exists
}

func ExtractBranchName(ref string) string {
	const prefix = "refs/heads/"
	if strings.HasPrefix(ref, prefix) {
		return strings.TrimPrefix(ref, prefix)
	}
	return ref
}
