package utils

import (
	"airun/code-reviewer/models"
	"fmt"
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
