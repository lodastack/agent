package outputs

type Config struct {
	Name       string   `toml:"name"`
	Servers    []string `toml:"servers"`
	BufferSize int      `toml:"buffersize"`
}
