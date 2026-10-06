package services

import (
	"airun/code-reviewer/configs"
	"airun/code-reviewer/models"
	"airun/code-reviewer/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var defaultTemperature = 0.2
var aiDefaultConfiguration = models.GetModelReviewPayload{
	Model:       configs.ACTIVE_AI_MODEL.Name,
	Temperature: &defaultTemperature,
}

func GetModelReviewMutation(options models.GetModelReviewMutationOptions) (*http.Response, error) {
	apiKey := utils.GetEnv("OPENROUTER_API_KEY", "http://127.0.0.1:11434")

	payload := models.GetModelReviewPayload{
		Model:       aiDefaultConfiguration.Model,
		Temperature: aiDefaultConfiguration.Temperature,
		Messages: []models.AIMessage{
			{
				Role:    "user",
				Content: utils.GetReviewPromptFrom(options.Files),
			},
			{
				Role:    "system",
				Content: configs.REVIEW_PROMPT,
			},
		},
		Reasoning: map[string]interface{}{
			"effort": "none",
		},
	}
	jsonPayload, err := json.Marshal(payload)
	fmt.Print(string(jsonPayload))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")
	response, err := HttpClient.Do(req)
	if err != nil {
		return &http.Response{}, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Error Code %v", response.Status)
	}
	return response, nil
}
