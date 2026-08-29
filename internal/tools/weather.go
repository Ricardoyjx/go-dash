package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/cloudwego/eino/components/tool"
	einoUtils "github.com/cloudwego/eino/components/tool/utils"
)

// WeatherInput 是 weather_search 工具的入参。
// 参数约束直接写在 struct tag 上，InferTool 会根据这些 tag 自动生成 ToolInfo。
type WeatherInput struct {
	Location string `json:"location" jsonschema:"required" jsonschema_description:"城市或地点名称，例如 Beijing、Shanghai、Tokyo"`
	Date     string `json:"date,omitempty" jsonschema_description:"查询日期 YYYY-MM-DD，默认为今天"`
}

// WeatherOutput 是 weather_search 工具返回的结果，InferTool 会自动序列化为 JSON 字符串。
type WeatherOutput struct {
	Location  string  `json:"location"`
	Date      string  `json:"date"`
	TempC     float64 `json:"temp_c"`
	Condition string  `json:"condition"`
	Humidity  int     `json:"humidity"`
	WindKph   float64 `json:"wind_kph"`
	Source    string  `json:"source"`
}

// NewWeatherTool 按官方文档推荐的 InferTool 方式，把本地函数转换为 Eino Tool。
// 数据源使用 wttr.in（免费、无需 API Key）；请求失败时回退到模拟数据，方便离线试用。
func NewWeatherTool() (tool.InvokableTool, error) {
	return einoUtils.InferTool(
		"weather_search",
		"查询指定地点和日期的天气信息，包括温度、天气状况、湿度和风速。",
		fetchWeather,
	)
}

func fetchWeather(ctx context.Context, input *WeatherInput) (*WeatherOutput, error) {
	if input.Location == "" {
		return nil, fmt.Errorf("location is required")
	}

	date := input.Date
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	out, err := weatherFromWttrIn(ctx, input.Location)
	if err != nil {
		// 网络不可用或服务异常时的兜底：返回模拟数据，Source 标记为 mock。
		out = &WeatherOutput{
			Location:  input.Location,
			Date:      date,
			TempC:     22.5,
			Condition: "Sunny (mock)",
			Humidity:  60,
			WindKph:   12,
			Source:    "mock",
		}
		return out, nil
	}

	out.Location = input.Location
	out.Date = date
	out.Source = "wttr.in"
	return out, nil
}

// weatherFromWttrIn 调用 wttr.in 的 JSON 接口获取当前天气。
func weatherFromWttrIn(ctx context.Context, location string) (*WeatherOutput, error) {
	u := "https://wttr.in/" + url.PathEscape(location) + "?format=j1"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wttr.in returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}

	var data struct {
		CurrentCondition []struct {
			TempC         string `json:"temp_C"`
			Humidity      string `json:"humidity"`
			WindspeedKmph string `json:"windspeedKmph"`
			WeatherDesc   []struct {
				Value string `json:"value"`
			} `json:"weatherDesc"`
		} `json:"current_condition"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	if len(data.CurrentCondition) == 0 {
		return nil, fmt.Errorf("wttr.in response has no current_condition")
	}

	cur := data.CurrentCondition[0]
	out := &WeatherOutput{
		TempC:     parseFloat(cur.TempC),
		Humidity:  int(parseFloat(cur.Humidity)),
		WindKph:   parseFloat(cur.WindspeedKmph),
		Condition: "unknown",
	}
	if len(cur.WeatherDesc) > 0 {
		out.Condition = cur.WeatherDesc[0].Value
	}
	return out, nil
}

func parseFloat(s string) float64 {
	var f float64
	// wttr.in 字段都是数字字符串；解析失败按 0 处理，避免整个工具报错。
	_, _ = fmt.Sscanf(s, "%f", &f)
	return f
}
