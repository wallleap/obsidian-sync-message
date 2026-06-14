package plugins

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

type BrowserRenderer struct {
	timeout   time.Duration
	userAgent string
}

func NewBrowserRenderer() *BrowserRenderer {
	return &BrowserRenderer{
		timeout:   30 * time.Second,
		userAgent: "",
	}
}

func (r *BrowserRenderer) SetUserAgent(ua string) {
	r.userAgent = ua
}

func findChromePath() string {
	paths := []string{
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
		"/usr/bin/google-chrome",
		"/usr/bin/chromium",
		"/usr/bin/chromium-browser",
	}
	
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func (r *BrowserRenderer) RenderURL(url string) (string, error) {
	var opts []chromedp.ExecAllocatorOption
	
	if chromePath := findChromePath(); chromePath != "" {
		os.Setenv("CHROME_PATH", chromePath)
		fmt.Printf("[DEBUG] Found Chrome at: %s\n", chromePath)
	} else {
		fmt.Println("[WARNING] Chrome not found, using default path")
	}
	
	if r.userAgent != "" {
		opts = append(opts, chromedp.UserAgent(r.userAgent))
	}
	
	opts = append(opts,
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var htmlContent string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		return "", fmt.Errorf("browser render failed: %w", err)
	}

	return htmlContent, nil
}

func (r *BrowserRenderer) RenderURLWithWait(url string, waitSelector string) (string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var htmlContent string

	tasks := []chromedp.Action{
		chromedp.Navigate(url),
	}

	if waitSelector != "" {
		tasks = append(tasks, chromedp.WaitReady(waitSelector, chromedp.ByQuery))
	} else {
		tasks = append(tasks, chromedp.WaitReady("body", chromedp.ByQuery))
	}

	tasks = append(tasks,
		chromedp.Sleep(2*time.Second),
		chromedp.OuterHTML("html", &htmlContent),
	)

	err := chromedp.Run(ctx, tasks...)

	if err != nil {
		return "", err
	}

	return htmlContent, nil
}

func (r *BrowserRenderer) RenderURLWithScript(url string, waitScript string) (string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var htmlContent string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Evaluate(waitScript, nil),
		chromedp.Sleep(1*time.Second),
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		return "", err
	}

	return htmlContent, nil
}

func (r *BrowserRenderer) GetElementText(url, selector string) (string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var text string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady(selector, chromedp.ByQuery),
		chromedp.Text(selector, &text),
	)

	if err != nil {
		return "", err
	}

	return text, nil
}

func (r *BrowserRenderer) GetElementsHTML(url, selector string) ([]string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var elements []*cdp.Node
	var htmls []string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Nodes(selector, &elements),
	)

	if err != nil {
		return nil, err
	}

	for _, el := range elements {
		if el.Children != nil {
			for _, child := range el.Children {
				if child.NodeName == "IMG" {
					for i := 0; i < len(child.Attributes); i += 2 {
						if i+1 < len(child.Attributes) {
							name := child.Attributes[i]
							value := child.Attributes[i+1]
							if name == "data-src" || name == "src" {
								htmls = append(htmls, "<img src=\""+value+"\">")
							}
						}
					}
				}
			}
		}
	}

	return htmls, nil
}
