package trace

// Config trace config struct
type Config struct {
	Enable    bool     `toml:"enable" json:"enable"`
	Collector []string `toml:"collector" json:"collector"`
}
