package tools

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type WebSearchInput struct {
	Query string `json:"query" jsonschema_description:"search query"`
}

type WebSearchOutput struct {
	Results []string `json:"results"`
}

// NewWebSearchTool returns a web search tool.
// The implementation is a mock; replace it with a real search provider.
func NewWebSearchTool() (tool.InvokableTool, error) {
	return utils.InferTool(
		"web_search",
		"Search the web for the given query and return relevant snippets.",
		func(_ context.Context, input *WebSearchInput) (*WebSearchOutput, error) {
			if input.Query == "" {
				return nil, fmt.Errorf("query is required")
			}
			// TODO: call a real search provider (Tavily / SerpAPI / Bing Search ...).
			return &WebSearchOutput{
				Results: []string{"[mock result] " + input.Query},
			}, nil
		},
	)
}
