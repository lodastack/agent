package member

// Config member config struct
type Config struct {
	Enable bool     `toml:"enable" json:"enable"`
	Key    string   `toml:"key" json:"key"`
	Nodes  []string `toml:"nodes" json:"nodes"`
}
