package trace

// Config trace config struct
type Config struct {
	Enable    bool     `toml:"enable"`
	Collector []string `toml:"collector"`
}
