package config

import (
	"sync"

	"github.com/lodastack/agent/agent/common"
	"github.com/lodastack/agent/agent/outputs"

	"github.com/BurntSushi/toml"
)

const (
	//APP NAME
	AppName = "Agent"
	//Usage
	Usage = "Agent Usage"
	//Vresion Num
	Version = "0.0.1"
	//Author Nmae
	Author = "devlopers@LodaStack"
	//Email Address
	Email = "devlopers@lodastack.com"
)

const (
	//PID FILE
	PID = "/var/run/loda-agent.pid"
)

var (
	mux = new(sync.RWMutex)
	C   = new(Config)
)

type Config struct {
	Agent  common.AgentConfig `toml:"agent"`
	Output outputs.Config     `toml:"output"`
	Log    LogConfig          `toml:"log"`
}

type LogConfig struct {
	Dir           string `toml:"logdir"`
	Level         string `toml:"loglevel"`
	Logrotatenum  int    `toml:"logrotatenum"`
	Logrotatesize uint64 `toml:"logrotatesize"`
}

func ParseConfig(path string) error {
	mux.Lock()
	defer mux.Unlock()

	if _, err := toml.DecodeFile(path, &C); err != nil {
		return err
	}
	return nil
}

func GetConfig() *Config {
	mux.RLock()
	defer mux.RUnlock()
	return C
}
