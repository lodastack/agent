package common

import (
	"io/ioutil"
	"os/user"
	"testing"
)

func Test_InitCollectConfig(t *testing.T) {
	config := MustConfig()

	if Conf.Listen != config.Listen {
		t.Fatalf("test config fatal: %s - %s", Conf.Listen, config.Listen)
	}
	if Conf.IfacePrefix[0] != config.IfacePrefix[0] {
		t.Fatalf("test config fatal: %s - %s", Conf.IfacePrefix[0], config.IfacePrefix[0])
	}
	if Conf.PluginsDir != config.PluginsDir {
		t.Fatalf("test config fatal: %s - %s", Conf.PluginsDir, config.PluginsDir)
	}
	if Conf.PluginsUser != config.PluginsUser {
		t.Fatalf("test config fatal: %s - %s", Conf.PluginsUser, config.PluginsUser)
	}
	if Conf.RegistryAddr != config.RegistryAddr {
		t.Fatalf("test config fatal: %s - %s", Conf.RegistryAddr, config.RegistryAddr)
	}
	if Conf.Git != config.Git {
		t.Fatalf("test config fatal: %s - %s", Conf.Git, config.Git)
	}
}

func MustConfig() *AgentConfig {
	pluginsDir, err := ioutil.TempDir("", "install-config-test-")
	if err != nil {
		panic("create tmp file fatal")
	}
	user, err := user.Current()
	if err != nil {
		panic("get current user fatal")
	}

	config := &AgentConfig{
		Listen:       "localhost:0",
		IfacePrefix:  []string{"eth"},
		PluginsDir:   pluginsDir,
		PluginsUser:  user.Username,
		RegistryAddr: "registry.test.com",
		Git:          "git@git.test.com:%s.git",
	}
	InitCollectConfig(config)
	return config
}
