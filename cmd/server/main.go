package main

import (
	"bufio"
	"context"
	"fmt"
	"go-dash/internal/agent"
	"go-dash/internal/config"
	"go-dash/internal/model"
	"go-dash/internal/utils/logger"
	"html"
	"os"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading env", zap.Error(err))
	}

	cfg := config.ModelConfig{
		ModelBaseUrl: os.Getenv("MODEL_BASE_URL"),
		ModelName:    os.Getenv("MODEL_NAME"),
		ModelAPIKey:  os.Getenv("MODEL_API_KEY"),
	}
	if cfg.ModelAPIKey == "" || cfg.ModelBaseUrl == "" || cfg.ModelName == "" {
		logger.Error("Error loading env", zap.Error(err))
		os.Exit(-1)
	}

	cm, err := model.NewDeepSeekModel(ctx, cfg)
	if err != nil {
		logger.Error("Error creating DeepSeek model", zap.Error(err))
	}

	super, err := agent.NewSupervisor(ctx, cm, 5)
	if err != nil {
		logger.Error("Error creating supervisor", zap.Error(err))
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: super,
	})

	scanner := bufio.NewScanner(os.Stdin)

	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for {
		scanErr := scanner.Err()
		fmt.Print("\nYou: ")
		if !scanner.Scan() {
			if err := scanner.Err; err != nil {
				logger.Error("Error reading input", zap.Error(scanErr))
			}
			break
		}
		query := strings.TrimSpace(scanner.Text())
		if query == "" {
			continue
		}
		if query == "exit" || query == "quit" {
			break
		}

		//非流式s
		events := runner.Query(ctx, query)
		for {
			event, ok := events.Next()
			if !ok {
				break
			}
			if event.Err != nil {
				logger.Error("agent error", zap.Error(event.Err))
				break
			}
			if event.Output != nil && event.Output.MessageOutput != nil {
				mv := event.Output.MessageOutput
				if mv.Message != nil && mv.Message.Content != "" {
					fmt.Println(html.UnescapeString(mv.Message.Content))
				}
			}
		}
	}
}
