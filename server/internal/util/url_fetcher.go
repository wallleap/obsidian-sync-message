package util

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ob-sync-server/internal/util/plugins"
)

// 微信公众号专用 User-Agent
const wechatUserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.34(0x16082222) NetType/WIFI Language/zh_CN"

func FetchURLContent(urlStr string) (string, string, string, error) {
	var htmlContent string
	var useBrowser bool

	if strings.HasPrefix(urlStr, "https://mp.weixin.qq.com/s") {
		htmlContent = fetchWithWeChatUA(urlStr)
		if strings.Contains(htmlContent, "环境异常") || strings.Contains(htmlContent, "验证") || 
		   strings.Contains(htmlContent, "请先完成验证") {
			fmt.Println("[INFO] Trying browser rendering for WeChat article...")
			useBrowser = true
		} else {
			fmt.Println("[INFO] Using simple HTTP fetch for WeChat article")
		}
	} else {
		htmlContent = fetchWithHTTP(urlStr, "")
	}

	if useBrowser || htmlContent == "" {
		var err error
		fmt.Println("[INFO] Using browser renderer for:", urlStr)
		renderer := plugins.NewBrowserRenderer()
		
		if strings.HasPrefix(urlStr, "https://mp.weixin.qq.com/s") {
			fmt.Println("[INFO] Setting WeChat User-Agent for browser")
			renderer.SetUserAgent(wechatUserAgent)
		}
		
		htmlContent, err = renderer.RenderURL(urlStr)
		if err != nil {
			fmt.Printf("[ERROR] Browser render failed: %v\n", err)
			if htmlContent == "" {
				return "", "", "", fmt.Errorf("failed to fetch content with browser")
			}
		}
	}

	if htmlContent == "" {
		return "", "", "", fmt.Errorf("failed to fetch content")
	}

	pluginManager := plugins.NewPluginManager()
	plugin := pluginManager.GetHandler(urlStr)

	title, contentHTML := plugin.ExtractContent(htmlContent)

	converter := plugins.NewHTMLToMarkdownConverter(urlStr)
	markdown := converter.Convert(contentHTML)

	return title, markdown, htmlContent, nil
}

func fetchWithHTTP(urlStr string, userAgent string) string {
	var resp *http.Response
	var err error

	if userAgent != "" {
		// 使用自定义 User-Agent
		client := &http.Client{}
		req, _ := http.NewRequest("GET", urlStr, nil)
		req.Header.Set("User-Agent", userAgent)
		resp, err = client.Do(req)
	} else {
		resp, err = http.Get(urlStr)
	}

	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return string(body)
}

func fetchWithWeChatUA(urlStr string) string {
	return fetchWithHTTP(urlStr, wechatUserAgent)
}
