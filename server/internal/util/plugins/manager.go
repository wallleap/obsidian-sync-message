package plugins

var defaultPlugins = []URLFetcherPlugin{
	&WeChatPlugin{},
	&JuejinPlugin{},
}

type PluginManager struct {
	plugins      []URLFetcherPlugin
	defaultPlugin URLFetcherPlugin
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins:      append([]URLFetcherPlugin{}, defaultPlugins...),
		defaultPlugin: &DefaultPlugin{},
	}
}

func (m *PluginManager) RegisterPlugin(plugin URLFetcherPlugin) {
	m.plugins = append(m.plugins, plugin)
}

func (m *PluginManager) GetHandler(url string) URLFetcherPlugin {
	for _, plugin := range m.plugins {
		if plugin.CanHandle(url) {
			return plugin
		}
	}
	return m.defaultPlugin
}
