package plugins

type URLFetcherPlugin interface {
	CanHandle(url string) bool
	ExtractContent(htmlContent string) (title string, contentHTML string)
	Name() string
}
