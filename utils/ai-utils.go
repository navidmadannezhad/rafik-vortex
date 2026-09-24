package utils

import (
	"airun/code-reviewer/configs"
	"airun/code-reviewer/models"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
)

func GetReviewPromptFrom(files []models.FileDiff) string {
	var diffText string

	for _, file := range files {
		path := file.Path
		if path == "" {
			path = "(unknown)"
		}
		diffText += fmt.Sprintf("### File: %s\n%s\n", path, file.Hunks)
	}

	maxChars := configs.ACTIVE_AI_MODEL.ModelMaximumCharacters
	runes := []rune(diffText)
	if len(runes) > maxChars {
		truncated := runes[:maxChars]
		last := -1
		for i := len(truncated) - 1; i >= 0; i-- {
			if truncated[i] == '\n' {
				last = i
				break
			}
		}
		if last > 0 {
			truncated = truncated[:last]
		}
		diffText = string(truncated) + "\n... [truncated]"
	}

	return fmt.Sprintf("Review the following code changes:\n\n%s", diffText)
}

func ParseAIReviewComments(response string) []models.AIReviewComment {
	var raw []models.AIReviewComment
	var arrayRE = regexp.MustCompile(`(?s)\[.*?\]`)

	// First try: parse the whole thing as JSON.
	if err := json.Unmarshal([]byte(response), &raw); err != nil {
		// Fallback: non-greedy match for first complete JSON array.
		match := arrayRE.FindString(response)
		if match == "" {
			return nil
		}
		if err := json.Unmarshal([]byte(match), &raw); err != nil {
			log.Println("warning: LLM response contained no parseable JSON array")
			return nil
		}
	}

	out := make([]models.AIReviewComment, 0, len(raw))
	for _, c := range raw {
		if c.Comment == "" {
			continue
		}
		out = append(out, c)
	}
	return out
}

func CommentListHasCritical(comments []models.AIReviewComment) bool {
	critical := false
	for _, comment := range comments {
		if comment.Severity == "critical" {
			critical = true
			break
		}
	}
	return critical
}
