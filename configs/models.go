package configs

import "airun/code-reviewer/models"

var aiModels = map[string]models.AiModelProperties{
	"qwen2.5-coder:7b": {
		Name:                   "qwen2.5-coder:7b",
		ModelMaximumCharacters: 12000,
	},
	"smollm2": {
		Name:                   "smollm2",
		ModelMaximumCharacters: 8192,
	},
}

var ACTIVE_AI_MODEL = aiModels["smollm2"]
