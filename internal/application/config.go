package application

// Config ...
type Config struct {
	BindPort   string `toml:"bind_port"`
	PythonPort string `toml:"python_port"`
}

// NewConfig - helper function
func NewConfig() *Config {
	return &Config{
		BindPort:   ":8080",
		PythonPort: ":7100",
	}
}
