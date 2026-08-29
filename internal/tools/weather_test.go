package tools

import (
	"context"
	"encoding/json"
	"testing"
)

// TestWeatherTool 直接调用 weather_search 工具，验证参数约束和返回结果。
func TestWeatherTool(t *testing.T) {
	ctx := context.Background()

	weatherTool, err := NewWeatherTool()
	if err != nil {
		t.Fatalf("NewWeatherTool failed: %v", err)
	}

	// 1. 查看自动生成的 ToolInfo
	info, err := weatherTool.Info(ctx)
	if err != nil {
		t.Fatalf("Info failed: %v", err)
	}
	t.Logf("tool name: %s", info.Name)
	t.Logf("tool desc: %s", info.Desc)

	params, err := info.ParamsOneOf.ToJSONSchema()
	if err != nil {
		t.Fatalf("ToJSONSchema failed: %v", err)
	}
	paramsJSON, _ := json.Marshal(params)
	t.Logf("tool params: %s", paramsJSON)

	// 2. 模拟 LLM 发起的调用（参数是 JSON 字符串）
	out, err := weatherTool.InvokableRun(ctx, `{"location":"Beijing"}`)
	if err != nil {
		t.Fatalf("InvokableRun failed: %v", err)
	}
	t.Logf("weather result: %s", out)
}
