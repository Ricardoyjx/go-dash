# go-dash

基于 **Go + Eino** 的多 Agent 项目骨架。

## 技术栈

- Go（模块名：`go-dash`）
- [Eino](https://github.com/cloudwego/eino)：字节跳动开源的 Go LLM 应用开发框架
- Eino ADK：Agent 构建、编排与运行（Runner / CheckPoint / Interrupt-Resume）

## 目录结构

``` bin
go-dash/
├── cmd/
│   └── server/                  # 程序入口（main.go）
├── internal/
│   ├── agent/
│   │   ├── subagents/           # 子 Agent 定义（Researcher / Coder / Reviewer ...）
│   │   ├── supervisor.go        # 协调者 / 主管 Agent
│   │   └── assemble.go          # 多 Agent 组装入口
│   ├── config/                  # 环境变量与配置加载
│   ├── memory/                  # CheckPointStore 等状态存储
│   ├── model/                   # ChatModel 初始化（OpenAI / Ark / DeepSeek ...）
│   ├── prompt/                  # 各 Agent 的 Instruction / 提示词
│   └── tools/                   # Agent 工具（web_search / knowledge_base / ...）
├── configs/                     # 配置文件（.env.example 等）
├── docs/                        # 架构与设计文档
├── README.md
└── go.mod
```

## 目录职责

| 目录 | 职责 |
| --- | --- |
| `cmd/server` | 程序入口，负责装配 Config → Model → Agent → Runner 并启动服务 |
| `internal/agent` | Agent 层：`subagents` 放各子 Agent，`supervisor.go` 放协调 Agent，`assemble.go` 统一组装 |
| `internal/tools` | 工具层，Agent 通过 ToolCall 调用（搜索、知识库、人工澄清等） |
| `internal/prompt` | 提示词层，集中管理各 Agent 的 Instruction，便于调整 |
| `internal/model` | 模型层，按 `MODEL_TYPE` 创建不同的 ChatModel |
| `internal/memory` | 状态层，提供 CheckPointStore 实现（内存 / 持久化） |
| `internal/config` | 配置层，从环境变量加载运行配置 |
| `configs` | 配置文件样例，如 `.env.example` |
| `docs` | 架构图、协作模式、扩展说明等文档 |

## 协作模式（规划）

骨架按 Eino 官方推荐的多 Agent 协作方式设计：

- **Supervisor 模式**：CoordinatorAgent 作为主管，动态路由任务给子 Agent；
- **Agent 独立执行**：子 Agent 通过 AgentTool / Supervisor 编排被调用，互不共享完整对话上下文；
- **Human-in-the-loop**：子 Agent 可通过 `ask_for_clarification` 类工具中断运行，等待用户补充信息后 Resume。

## 环境变量（规划）

| 变量 | 说明 | 默认值 |
| --- | --- | --- |
| `MODEL_TYPE` | 模型提供商：`openai` / `ark` | `openai` |
| `OPENAI_API_KEY` | OpenAI 或兼容接口的 Key | - |
| `OPENAI_MODEL` | OpenAI 模型名 | `gpt-4o-mini` |
| `OPENAI_BASE_URL` | OpenAI 兼容 BaseURL（可指向代理或网关） | - |
| `ARK_API_KEY` / `ARK_MODEL` / `ARK_BASE_URL` | 火山方舟配置 | - |
| `AGENT_MAX_ITERATIONS` | 单个 Agent 最大工具调用轮数 | `5` |

## 扩展指引

- **新增子 Agent**：在 `internal/agent/subagents` 下新增定义文件，并在 `assemble.go` 中注册；
- **新增工具**：在 `internal/tools` 下新增工具实现，挂载到对应 Agent 的 ToolsConfig；
- **调整行为**：修改 `internal/prompt` 中的 Instruction，无需改动代码逻辑；
- **切换模型**：在 `internal/model` 增加 provider 分支，通过 `MODEL_TYPE` 切换；
- **持久化**：将 `internal/memory` 中的 CheckPointStore 替换为 Redis 或本地磁盘实现；
- **对外服务**：在 `cmd/server` 中启动 HTTP / SSE 接口，流式透传 Agent 事件。

## 参考

- [Eino 官方文档](https://www.cloudwego.io/zh/docs/eino/)
- [Eino 官方多 Agent 示例](https://github.com/cloudwego/eino-examples/tree/main/adk/multiagent)
