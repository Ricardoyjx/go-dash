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
	Logger.Info("Creating supervisor agent")

	webSearch, err := tools.NewWebSearchTool()

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Model:       cm,
		Name:        "Supervisor",
		Description: "Supervisor agent",
		Instruction: prompt.Coordinator,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{
					webSearch, tools.NewAskForClarificationTool()},
			},
		},
		MaxIterations: 5,
	})
	if err != nil {
		Logger.Error("Failed to create supervisor agent", zap.Error(err))
	}

	return agent, err
}
