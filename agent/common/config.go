package common

type AgentConfig struct {
	Listen       string   `toml:"listen"`
	IfacePrefix  []string `toml:"ifaceprefix"`
	NTPServer    string   `toml:"ntpserver"`
	PluginsDir   string   `toml:"pluginsdir"`
	PluginsUser  string   `toml:"pluginsuser"`
	RegistryAddr string   `toml:registryaddr"`
	Git          string   `toml:"git"`
}

var Conf *AgentConfig

func InitCollectConfig(config *AgentConfig) {
	// Default use root exec plugins
	if config.PluginsUser == "" {
		config.PluginsUser = "root"
	}
	if config.NTPServer == "" {
		config.NTPServer = "133.100.11.8"
	}
	Conf = config
}
