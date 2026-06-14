package plugins

import (
	"regexp"
	"strings"

	html2md "github.com/JohannesKaufmann/html-to-markdown/v2"
)

type HTMLToMarkdownConverter struct {
	baseURL string
}

func NewHTMLToMarkdownConverter(baseURL string) *HTMLToMarkdownConverter {
	return &HTMLToMarkdownConverter{baseURL: baseURL}
}

func (c *HTMLToMarkdownConverter) Convert(htmlContent string) string {
	// 创建转换器
	md, err := html2md.ConvertString(htmlContent)
	if err != nil {
		return htmlContent
	}

	// 处理图片
	md = c.processImages(md)

	return strings.TrimSpace(md)
}

func (c *HTMLToMarkdownConverter) processImages(markdown string) string {
	// 匹配 Markdown 图片格式: ![alt](url)
	imgRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)

	md := imgRegex.ReplaceAllStringFunc(markdown, func(match string) string {
		parts := imgRegex.FindStringSubmatch(match)
		if len(parts) == 3 {
			alt := parts[1]
			url := parts[2]

			// 如果 URL 是相对路径，转换为绝对路径
			if c.baseURL != "" && !strings.HasPrefix(url, "http") {
				url = c.resolveURL(url)
			}

			return "![" + alt + "](" + url + ")"
		}
		return match
	})

	return md
}

func (c *HTMLToMarkdownConverter) resolveURL(href string) string {
	if href == "" {
		return ""
	}

	// 简单的相对路径转绝对路径
	base := strings.TrimSuffix(c.baseURL, "/")
	if strings.HasPrefix(href, "/") {
		return base + href
	}
	return base + "/" + href
}

func ExtractTitle(htmlContent string) string {
	re := regexp.MustCompile(`<title[^>]*>(.*?)</title>`)
	match := re.FindStringSubmatch(htmlContent)
	if len(match) > 1 {
		return strings.TrimSpace(StripTags(match[1]))
	}
	return ""
}

func StripTags(s string) string {
	re := regexp.MustCompile(`<[^>]+>`)
	return re.ReplaceAllString(s, "")
}
