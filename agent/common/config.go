package common

import (
	"fmt"
)

type AgentConfig struct {
	Listen       string   `toml:"listen"`
	IfacePrefix  []string `toml:"ifaceprefix"`
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
	Conf = config
	fmt.Println("load config:", Conf)
}
