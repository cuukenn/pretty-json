# agents.md — pretty-json

本文件供编码代理（opencode 等）在修改本项目时参考。改动任何 Go 文件后请阅读对应源码，不要仅凭本文件臆断。

## 项目概述

一个基于 **Wails v2 + Go 后端 + Vue 前端**的 Windows 桌面 JSON 工具，提供 JSON 的格式化、压缩、校验、统计分析、JSON↔YAML/TOML 互转、JSONPath 查询、JSON 树浏览，以及从 JSON 生成 Go / TypeScript / Java (+Lombok) 代码，支持拖放打开文件与保存。

## 技术栈



* Go **1.25.1**，模块名 `pretty-json`

* Wails **v2.16.0**（`github.com/wailsapp/wails/v2`）

* `BurntSushi/toml v1.6.0`、`gopkg.in/yaml.v3 v3.0.1`（互转用）

* 前端：Vue（Wails Vue 模板），`frontend/` 由 npm 管理

## 目录与文件职责（Go 侧，根目录）



| 文件            | 职责                                                                                                  |
| ------------- | --------------------------------------------------------------------------------------------------- |
| `main.go`     | 入口；Wails 配置（标题 pretty-json、1280×860、启用拖放、绑定 `App`）                                                  |
| `app.go`      | `App` 结构体、`startup`、文件拖放；**单条目解析缓存&#x20;**`parseCache`（大文档反复解析复用，缓存解析树按只读对待）                        |
| `json.go`     | 解析 / 编码核心（`Obj.MarshalJSON` 保 key 顺序 + 大数精度）、`FormatJSON`/`MinifyJSON`/`ValidateJSON`/`AnalyzeJSON` |
| `convert.go`  | `JSONToYAML`/`YAMLToJSON`/`JSONToTOML`/`TOMLToJSON`                                                 |
| `jsonpath.go` | `JSONPathQuery`（含带点的引号 key、单值 vs 数组语义）                                                              |
| `tree.go`     | `BuildJSONTree` 构建 JSON 树                                                                           |
| `codegen.go`  | `GenerateCode`：按 `lang` 生成 Go / TS / Java（支持 Lombok）                                                |
| `files.go`    | `OpenJSONFile` / `SaveTextFile` 文件读写                                                                |
| `*.go_test`   | 各模块单元测试与基准（`BenchmarkEncodeValue`）                                                                  |

## 前端可调用的后端方法（Wails `Bind`）



* `FormatJSON(jsonStr) (string, error)` — 美化

* `MinifyJSON(jsonStr) (string, error)` — 压缩

* `ValidateJSON(jsonStr) (bool, string)` — 校验，返回是否合法 + 错误信息

* `AnalyzeJSON(jsonStr) AnalyzeResult` — 统计（键、值、层级、错误定位等）

* `JSONToYAML` / `YAMLToJSON` / `JSONToTOML` / `TOMLToJSON` `(string) (string, error)`

* `JSONPathQuery(jsonStr, path string) (*JSONPathResult, error)`

* `BuildJSONTree(jsonStr string) (*TreeNode, error)`

* `GenerateCode(jsonStr, lang string, opts map[string]any) (string, error)`

* `OpenJSONFile() (*FileData, error)` / `SaveTextFile(content, defaultName string) (string, error)`

## 构建 / 开发 / 测试



```
\# 单元测试

go test ./...

\# 基准

go test -bench=. ./...

\# 开发模式（热重载前端）

wails dev

\# 生产打包

wails build
```

前端（`frontend/`）：`npm install`、`npm run build`、`npm run dev`（由 `wails.json` 配置）。

## 编码约定（改动前必读对应源码）



* **保序与精度**：`json.go` 的编码刻意保持 key 顺序并保留大数精度（有 `TestFormatJSONPreservesBigNumberPrecision` / `TestEncodeParityWithLegacyPath` 守护），**不要引入会重排 key 或丢精度**的改动。

* **解析缓存**：`app.go` 的 `parseCache` 只缓存最近一次解析，且缓存树按只读对待 —— 新增需要解析 JSON 的方法时，优先复用缓存路径；若修改了解析结果，必须同步更新缓存逻辑。

* **互转顺序**：JSON 是唯一权威中间格式；YAML/TOML 转换以 JSON 为基准（先转 JSON 再转目标格式）。

* **错误处理**：方法返回 `(string, error)` / `(bool, string)` 时，把可读错误消息带给前端；空 / 非法输入需有明确行为（见各 `_test.go`）。

* **测试**：为新增 / 改动补 `_test.go` 用例，遵循现有 `TestXxx` 命名与表驱动风格。