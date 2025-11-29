package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// FormatJSON 格式化 JSON（带 2 空格缩进）
func (a *App) FormatJSON(jsonStr string) (string, error) {
	if strings.TrimSpace(jsonStr) == "" {
		return "", fmt.Errorf("请输入JSON内容")
	}
	var jsonData interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	if err != nil {
		return "", fmt.Errorf("JSON格式错误: %v", err)
	}
	formattedJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("格式化失败: %v", err)
	}
	return string(formattedJSON), nil
}

// MinifyJSON 压缩 JSON
func (a *App) MinifyJSON(jsonStr string) (string, error) {
	if strings.TrimSpace(jsonStr) == "" {
		return "", fmt.Errorf("请输入JSON内容")
	}
	var jsonData interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	if err != nil {
		return "", fmt.Errorf("JSON格式错误: %v", err)
	}
	minifiedJSON, err := json.Marshal(jsonData)
	if err != nil {
		return "", fmt.Errorf("压缩失败: %v", err)
	}
	return string(minifiedJSON), nil
}

// ValidateJSON 校验 JSON 格式
func (a *App) ValidateJSON(jsonStr string) (bool, string) {
	if strings.TrimSpace(jsonStr) == "" {
		return false, "请输入JSON内容"
	}
	var jsonData interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	if err != nil {
		return false, fmt.Sprintf("校验失败: %v", err)
	}
	return true, "JSON格式正确 ✅"
}

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "pretty-json",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
