package plugins

import (
	"regexp"
	"strings"
)

type JuejinPlugin struct{}

func (p *JuejinPlugin) Name() string {
	return "juejin"
}

func (p *JuejinPlugin) CanHandle(url string) bool {
	return strings.HasPrefix(url, "https://juejin.cn/post")
}

func (p *JuejinPlugin) ExtractContent(htmlContent string) (string, string) {
	title := ExtractTitle(htmlContent)

	contentReg := regexp.MustCompile(`<article[^>]*>([\s\S]*?)</article>`)
	contentMatch := contentReg.FindStringSubmatch(htmlContent)

	if len(contentMatch) < 2 {
		contentReg = regexp.MustCompile(`<div class="article-content[^>]*>([\s\S]*?)</div>`)
		contentMatch = contentReg.FindStringSubmatch(htmlContent)
	}

	if len(contentMatch) < 2 {
		contentReg = regexp.MustCompile(`<div class="markdown-body[^>]*>([\s\S]*?)</div>`)
		contentMatch = contentReg.FindStringSubmatch(htmlContent)
	}

	if len(contentMatch) >= 2 {
		return title, contentMatch[1]
	}

	return title, htmlContent
}
