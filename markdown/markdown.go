package markdown

import (
	"github.com/charmbracelet/glamour"
)

func RenderMarkdown(content string) (string, error) {

	// 使用glamour渲染Markdown到终端
	r, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle("auto"),
		glamour.WithWordWrap(80),
	)

	out, err := r.Render(content)
	if err != nil {
		// 如果渲染失败，直接输出原始响应
		return content, err
	}

	return out, nil
}
