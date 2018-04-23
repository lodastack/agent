package member

// Config member config struct
type Config struct {
	Enable bool     `toml:"enable"`
	Key    string   `toml:"key"`
	Nodes  []string `toml:"nodes"`
}
