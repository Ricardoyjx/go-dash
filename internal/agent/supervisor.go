package agent

import (
	"context"

	"go-dash/internal/prompt"
	Logger "go-dash/internal/utils/logger"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"go.uber.org/zap"
)

func NewSupervisor(ctx context.Context, cm model.ToolCallingChatModel, maxIter int) (*adk.ChatModelAgent, error) {

	Logger.Info("Creating supervisor agent")

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Model:         cm,
		Name:          "Supervisor",
		Description:   "Supervisor agent",
		Instruction:   prompt.Coordinator,
		ToolsConfig:   adk.ToolsConfig{},
		MaxIterations: 5,
		Exit:          adk.ExitTool{},
	})
	if err != nil {
		Logger.Error("Failed to create supervisor agent", zap.Error(err))
	}

	return agent, err

}
