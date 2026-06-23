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
	md, err := html2md.ConvertString(htmlContent)
	if err != nil {
		return htmlContent
	}

	md = c.processImages(md)
	md = c.processSingleH1(md)

	return strings.TrimSpace(md)
}

func (c *HTMLToMarkdownConverter) processSingleH1(markdown string) string {
	h1Regex := regexp.MustCompile(`^#\s+.*$`)

	lines := strings.Split(markdown, "\n")
	var h1Count int
	var h1LineIndex int

	for i, line := range lines {
		if h1Regex.MatchString(line) {
			h1Count++
			h1LineIndex = i
		}
	}

	if h1Count == 1 {
		lines = append(lines[:h1LineIndex], lines[h1LineIndex+1:]...)
		return strings.Join(lines, "\n")
	}

	return markdown
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
