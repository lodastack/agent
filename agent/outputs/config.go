package outputs

type Config struct {
	Name       string   `toml:"name" json:"name"`
	Servers    []string `toml:"servers" json:"servers"`
	BufferSize int      `toml:"buffersize" json:"buffersize"`
}
