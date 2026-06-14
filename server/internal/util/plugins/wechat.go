package plugins

import (
	"fmt"
	"regexp"
	"strings"
)

type WeChatPlugin struct{}

func (p *WeChatPlugin) Name() string {
	return "wechat"
}

func (p *WeChatPlugin) CanHandle(url string) bool {
	return strings.HasPrefix(url, "https://mp.weixin.qq.com/s")
}

func (p *WeChatPlugin) ExtractContent(htmlContent string) (string, string) {
	title := p.extractTitle(htmlContent)
	fmt.Printf("[WeChat Debug] Extracted title: %s\n", title)

	contentReg := regexp.MustCompile(`<div id="js_content"[^>]*>([\s\S]*?)</div>`)
	contentMatch := contentReg.FindStringSubmatch(htmlContent)

	if len(contentMatch) < 2 {
		contentReg = regexp.MustCompile(`<div class="rich_media_content[^>]*>([\s\S]*?)</div>`)
		contentMatch = contentReg.FindStringSubmatch(htmlContent)
	}

	if len(contentMatch) < 2 {
		contentReg = regexp.MustCompile(`<section[^>]*>([\s\S]*?)</section>`)
		contentMatch = contentReg.FindStringSubmatch(htmlContent)
	}

	if len(contentMatch) < 2 {
		contentReg = regexp.MustCompile(`<article[^>]*>([\s\S]*?)</article>`)
		contentMatch = contentReg.FindStringSubmatch(htmlContent)
	}

	if len(contentMatch) >= 2 {
		content := contentMatch[1]
		content = p.replaceDataSrc(content)
		return title, content
	}

	return title, p.replaceDataSrc(htmlContent)
}

func (p *WeChatPlugin) replaceDataSrc(content string) string {
	dataSrcReg := regexp.MustCompile(`<img[^>]*data-src=["']([^"']+)["'][^>]*>`)
	content = dataSrcReg.ReplaceAllString(content, `<img src="$1">`)

	dataSrcSetReg := regexp.MustCompile(`<img[^>]*data-srcset=["']([^"']+)["'][^>]*>`)
	content = dataSrcSetReg.ReplaceAllString(content, `<img src="$1">`)

	return content
}

// 微信公众号专用标题提取
func (p *WeChatPlugin) extractTitle(htmlContent string) string {
	// 方法1: 提取 og:title meta 标签
	ogTitleReg := regexp.MustCompile(`<meta[^>]*property=["']og:title["'][^>]*content=["']([^"']+)["']`)
	match := ogTitleReg.FindStringSubmatch(htmlContent)
	if len(match) > 1 {
		title := strings.TrimSpace(match[1])
		fmt.Printf("[WeChat Debug] Method 1 (og:title): %s\n", title)
		return title
	}

	// 方法1b: 另一种 og:title 格式
	ogTitleReg2 := regexp.MustCompile(`<meta[^>]*content=["']([^"']+)["'][^>]*property=["']og:title["']`)
	match = ogTitleReg2.FindStringSubmatch(htmlContent)
	if len(match) > 1 {
		title := strings.TrimSpace(match[1])
		fmt.Printf("[WeChat Debug] Method 1b (og:title alt): %s\n", title)
		return title
	}

	// 方法2: 提取 h1.rich_media_title
	h1Reg := regexp.MustCompile(`<h1[^>]*class=["'][^"']*rich_media_title[^"']*["'][^>]*>([\s\S]*?)</h1>`)
	match = h1Reg.FindStringSubmatch(htmlContent)
	if len(match) > 1 {
		title := strings.TrimSpace(StripTags(match[1]))
		fmt.Printf("[WeChat Debug] Method 2 (h1.rich_media_title): %s\n", title)
		return title
	}

	// 方法3: 提取任意 h1 标签
	h1AnyReg := regexp.MustCompile(`<h1[^>]*>([\s\S]*?)</h1>`)
	match = h1AnyReg.FindStringSubmatch(htmlContent)
	if len(match) > 1 {
		title := strings.TrimSpace(StripTags(match[1]))
		fmt.Printf("[WeChat Debug] Method 3 (any h1): %s\n", title)
		return title
	}

	// 方法4: 提取 title 标签
	titleReg := regexp.MustCompile(`<title[^>]*>(.*?)</title>`)
	match = titleReg.FindStringSubmatch(htmlContent)
	if len(match) > 1 {
		title := strings.TrimSpace(match[1])
		// 微信的 title 通常是 "标题 - 公众号名 - 微信" 格式
		if idx := strings.Index(title, " - "); idx > 0 {
			title = title[:idx]
		}
		fmt.Printf("[WeChat Debug] Method 4 (title tag): %s\n", title)
		return title
	}

	fmt.Printf("[WeChat Debug] No title found!\n")
	return ""
}
