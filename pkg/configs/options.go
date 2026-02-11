package configs

// Options holds the options required to run the converters.
type Options struct {
	ProxyBufferHeuristic bool `yaml:"proxy_buffer_heuristic,omitempty" json:"proxy_buffer_heuristic,omitempty"`
	DisablePlugins       bool `yaml:"disable_plugins,omitempty"        json:"disable_plugins,omitempty"`
}

// NewOptions returns new instance of Options when invoked.
func NewOptions() *Options {
	return &Options{}
}
