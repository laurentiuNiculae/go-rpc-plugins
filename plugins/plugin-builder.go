package plugins

type PluginBuilder interface {
	Build(addr, port string) Plugin
}
