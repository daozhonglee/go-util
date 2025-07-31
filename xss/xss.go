// Package xss 提供XSS防护功能
package xss

import (
	"html"

	"github.com/microcosm-cc/bluemonday"
)

// Clean 清理HTML内容，移除潜在的XSS攻击代码，保留安全的HTML标签
func Clean(content string) string {
	p := bluemonday.UGCPolicy()
	afterHtml := p.Sanitize(content)
	afterHtml = html.UnescapeString(afterHtml)
	return afterHtml
}
