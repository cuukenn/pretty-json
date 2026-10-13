# pretty-json

一款基于 **Wails v2** 的 Windows 桌面 JSON 处理工具。内置 JSON 格式化、压缩、校验、统计分析、JSON ↔ YAML/TOML 互转、JSONPath 查询、JSON 树浏览，以及从 JSON 一键生成 Go / TypeScript / Java(+Lombok) 代码。

## 功能特性

- **格式化 / 压缩**：美化或 minify JSON，保持 key 顺序、保留大数精度
- **校验与分析**：合法性校验 + 错误定位；统计键值数、层级深度等结构信息
- **格式互转**：JSON → YAML、YAML → JSON、JSON → TOML、TOML → JSON（JSON 为权威中间格式）
- **JSONPath 查询**：支持标准 JSONPath，含带点号引号 key、单值与数组语义
- **JSON 树浏览**：将 JSON 构建为层级树，便于浏览大文档
- **代码生成**：从 JSON 生成 Go / TypeScript / Java（可选 Lombok）类型代码
- **文件拖放**：拖入 JSON 文件直接打开，可保存处理结果

## 技术栈

- Go 1.25 + Wails v2（`github.com/wailsapp/wails/v2`）
- 前端：Vue（Wails Vue 模板，`frontend/`，npm 管理）
- 依赖：`BurntSushi/toml`（TOML 解析）、`gopkg.in/yaml.v3`（YAML 解析）

## 构建 / 开发 / 测试

```sh
# 单元测试
go test ./...

# 基准测试
go test -bench=. ./...

# 开发模式（前端热重载）
wails dev

# 生产打包
wails build
```

前端依赖（`frontend/`）在 `wails dev` / `wails build` 时由 `wails.json` 自动触发 `npm install` / `npm run build`。

## 使用

启动应用后，粘贴 JSON 或拖入 JSON 文件，即可进行格式化、校验、互转、JSONPath 查询与代码生成。处理后的文本可保存到本地文件。

## 目录结构（Go 侧）

| 文件 | 职责 |
|---|---|
| `main.go` | 入口；Wails 应用配置与绑定 |
| `app.go` | `App` 结构体、启动逻辑、解析缓存 |
| `json.go` | JSON 解析/编码、格式化、压缩、校验、分析 |
| `convert.go` | JSON ↔ YAML / TOML 互转 |
| `jsonpath.go` | JSONPath 查询 |
| `tree.go` | JSON 树构建 |
| `codegen.go` | Go / TS / Java 代码生成 |
| `files.go` | 文件打开与保存 |

面向编码代理的详细约定见 [`agents.md`](./agents.md)。

## License

[Apache License 2.0](./LICENSE)
