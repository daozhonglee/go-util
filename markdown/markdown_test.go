package markdown

import (
	"fmt"
	"testing"
)

func TestRenderMarkdown(t *testing.T) {
	content := "## 终端渲染示例\n> 引用内容\n`代码片段`"
	result, err := RenderMarkdown(content)
	if err != nil {
		t.Errorf("RenderMarkdown failed, err = %v", err)
		fmt.Println(err)
	}
	fmt.Println(result)
	// if result != content {
	// 	t.Errorf("RenderMarkdown failed, expected %s, got %s", content, result)
	// }
}
