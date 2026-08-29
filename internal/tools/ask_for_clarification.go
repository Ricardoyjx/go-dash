package tools

import (
	"context"
	"log"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
)

type askForClarificationOptions struct {
	NewInput *string
}

// WithNewInput carries the user's answer into the resumed run.
func WithNewInput(input string) tool.Option {
	return tool.WrapImplSpecificOptFn(func(o *askForClarificationOptions) {
		o.NewInput = &input
	})
}

type AskForClarificationInput struct {
	Question string `json:"question" jsonschema_description:"question to ask the user for missing information"`
}

// NewAskForClarificationTool interrupts the agent and asks the user for input
// (human-in-the-loop) when the task lacks necessary context.
func NewAskForClarificationTool() tool.InvokableTool {
	t, err := utils.InferOptionableTool(
		"ask_for_clarification",
		"Call this tool when the request is ambiguous or lacks necessary information. It interrupts the agent to ask the user for the missing details.",
		func(_ context.Context, input *AskForClarificationInput, opts ...tool.Option) (string, error) {
			o := tool.GetImplSpecificOptions[askForClarificationOptions](nil, opts...)
			if o.NewInput == nil {
				return "", compose.NewInterruptAndRerunErr(input.Question)
			}
			out := *o.NewInput
			o.NewInput = nil
			return out, nil
		},
	)
	if err != nil {
		log.Fatalf("infer ask_for_clarification tool: %v", err)
	}
	return t
}
