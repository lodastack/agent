package common

type AgentConfig struct {
	Listen       string   `toml:"listen" json:"listen"`
	IfacePrefix  []string `toml:"ifaceprefix" json:"ifaceprefix"`
	PluginsDir   string   `toml:"pluginsdir" json:"pluginsdir"`
	PluginsUser  string   `toml:"pluginsuser" json:"pluginsuser"`
	RegistryAddr string   `toml:registryaddr" json:"registryaddr"`
	Git          string   `toml:"git" json:"git"`
}

var Conf *AgentConfig

func InitCollectConfig(config *AgentConfig) {
	// Default use root exec plugins
	if config.PluginsUser == "" {
		config.PluginsUser = "root"
	}
	Conf = config
}
