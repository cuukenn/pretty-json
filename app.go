package main

import (
	"context"
	"sync"
)

// App struct
type App struct {
	ctx context.Context

	// 单条目解析缓存：大文档常被 校验/格式化/互转/JSONPath/代码生成 反复解析。
	// 缓存的解析树按只读对待（已确认所有调用点不修改解析结果）。
	parseCacheMu  sync.Mutex
	parseCacheKey string
	parseCacheVal any
	parseCacheOK  bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.enableFileDrop()
}
