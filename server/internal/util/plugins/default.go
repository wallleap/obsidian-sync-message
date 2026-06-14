package plugins

type DefaultPlugin struct{}

func (p *DefaultPlugin) Name() string {
	return "default"
}

func (p *DefaultPlugin) CanHandle(url string) bool {
	return true
}

func (p *DefaultPlugin) ExtractContent(htmlContent string) (string, string) {
	title := ExtractTitle(htmlContent)
	return title, htmlContent
}
