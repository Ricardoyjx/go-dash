package model

import (
	"context"

	"go-dash/internal/config"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var logger *zap.Logger
var ctx context.Context

func NewDeepSeekModel(ctx context.Context, config config.ModelConfig) (*deepseek.ChatModel, error) {
	logger, _ = zap.NewDevelopment()

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading env", zap.Error(err))
	}

	model, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  config.ModerlApiKey,
		Model:   config.ModelName,
		BaseURL: config.ModelBaseUrl,
	})

	if err != nil {
		logger.Error("Error creating DeepSeek model", zap.Error(err))
	}
	return model, err
}
