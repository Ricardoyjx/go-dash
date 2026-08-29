package subagents

import (
	"context"
	"go-dash/internal/prompt"
	tools "go-dash/internal/tools"
	Logger "go-dash/internal/utils/logger"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"go.uber.org/zap"
)

func NewWeatherAssistant(ctx context.Context, cm model.ToolCallingChatModel, MaxIter int) (adk.Agent, error) {
	Logger.Info("Creating weather assistant agent")

	webSearch, err := tools.NewWebSearchTool()
	if err != nil {
		return nil, err
	}

	weatherSearch, err := tools.NewWeatherTool()
	if err != nil {
		return nil, err
	}

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Model:       cm,
		Name:        "WeatherAssistant",
		Description: "Weather assistant agent",
		Instruction: prompt.Coordinator,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{
					webSearch, weatherSearch, tools.NewAskForClarificationTool(),
				},
			},
		},
		MaxIterations: MaxIter,
	})
	if err != nil {
		Logger.Error("Failed to create supervisor agent", zap.Error(err))
	}

	return agent, err
}
