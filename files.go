package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// FileData 文件内容与元信息
type FileData struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

var fileFilters = []wailsruntime.FileFilter{
	{DisplayName: "数据 / 代码文件", Pattern: "*.json;*.yaml;*.yml;*.toml;*.txt;*.go;*.ts;*.java"},
	{DisplayName: "所有文件", Pattern: "*"},
}

// OpenJSONFile 弹出打开对话框选择数据文件；用户取消时返回 nil
func (a *App) OpenJSONFile() (*FileData, error) {
	path, err := wailsruntime.OpenFileDialog(a.ctx, wailsruntime.OpenDialogOptions{
		Title:   "打开 JSON / YAML / TOML 文件",
		Filters: fileFilters,
	})
	if err != nil {
		return nil, err
	}
	if path == "" {
		return nil, nil // 用户取消
	}
	return readFileData(path)
}

// maxFileBytes 打开/拖入文件的大小上限，防止超大文件撑爆编辑器与桥接
const maxFileBytes = 20 << 20 // 20 MB

// formatSize 人类可读的文件大小
func formatSize(n int64) string {
	switch {
	case n >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(n)/(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(n)/(1<<10))
	}
	return fmt.Sprintf("%d B", n)
}

// readFileData 读取文件并去除 UTF-8 BOM；超过大小上限时返回错误
func readFileData(path string) (*FileData, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}
	if info.Size() > maxFileBytes {
		return nil, fmt.Errorf("文件过大（%s），打开上限为 %d MB", formatSize(info.Size()), maxFileBytes>>20)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}
	content := strings.TrimPrefix(string(data), "\ufeff")
	return &FileData{Name: filepath.Base(path), Path: path, Content: content}, nil
}

// SaveTextFile 弹出保存对话框写入内容；返回保存路径（空字符串表示用户取消）
func (a *App) SaveTextFile(content, defaultName string) (string, error) {
	if strings.TrimSpace(defaultName) == "" {
		defaultName = "pretty-json.json"
	}
	path, err := wailsruntime.SaveFileDialog(a.ctx, wailsruntime.SaveDialogOptions{
		Title:           "另存为",
		DefaultFilename: defaultName,
		Filters:         fileFilters,
	})
	if err != nil {
		return "", err
	}
	if path == "" {
		return "", nil // 用户取消
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("写入文件失败: %v", err)
	}
	return path, nil
}

// enableFileDrop 注册拖放：拖入文件后向前端发送 file:drop 事件（只处理第一个文件）
func (a *App) enableFileDrop() {
	wailsruntime.OnFileDrop(a.ctx, func(x, y int, paths []string) {
		for _, p := range paths {
			info, err := os.Stat(p)
			if err != nil || info.IsDir() {
				continue
			}
			fd, err := readFileData(p)
			if err != nil {
				wailsruntime.EventsEmit(a.ctx, "file:drop:error", err.Error())
				return
			}
			wailsruntime.EventsEmit(a.ctx, "file:drop", fd)
			return
		}
	})
}
