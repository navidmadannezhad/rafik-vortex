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

var aiDefaultConfiguration = models.AiConfiuguration{
	Model: configs.ACTIVE_AI_MODEL.Name,
	Options: map[string]interface{}{
		"temperature": 0.2,
		"num_predict": 256,
	},
}

func GetModelReviewMutation(options models.GetModelReviewMutationOptions) (*http.Response, error) {
	modelUrl := utils.GetEnv("MODEL_URL", "http://127.0.0.1:11434")
	url := fmt.Sprintf("%v/api/generate", modelUrl)

	payload := models.AiConfiuguration{
		Model:   aiDefaultConfiguration.Model,
		Options: aiDefaultConfiguration.Options,
		Prompt:  utils.GetReviewPromptFrom(options.Files),
		System:  configs.REVIEW_PROMPT,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, err
	}
	response, err := HttpClient.Do(req)
	if err != nil {
		return &http.Response{}, err
	}
	return response, nil
}
