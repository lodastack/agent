package config

import (
	"sync"

	"github.com/lodastack/agent/agent/common"
	"github.com/lodastack/agent/agent/outputs"
	"github.com/lodastack/agent/member"
	"github.com/lodastack/agent/trace"

	"github.com/BurntSushi/toml"
)

const (
	//APP NAME
	AppName = "Monitor Agent"
	//Usage
	Usage = "Agent Usage"
	//Author Nmae
	Author = "devlopers@LodaStack"
	//Email Address
	Email = "devlopers@lodastack.com"
)

var (
	//Vresion Num
	Version = ""
	//Vresion Commit
	Commit = ""
	//Vresion Branch
	Branch = ""
	//Build Time
	BuildTime = ""
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
	Agent  common.AgentConfig `toml:"agent" json:"agent"`
	Output outputs.Config     `toml:"output" json:"output"`
	Trace  trace.Config       `toml:"trace" json:"trace"`
	Member member.Config      `toml:"member" json:"member"`
	Log    LogConfig          `toml:"log" json:"log"`
}

type LogConfig struct {
	Dir           string `toml:"logdir" json:"logdir"`
	Level         string `toml:"loglevel" json:"loglevel"`
	Logrotatenum  int    `toml:"logrotatenum" json:"logrotatenum"`
	Logrotatesize uint64 `toml:"logrotatesize" json:"logrotatesize"`
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
