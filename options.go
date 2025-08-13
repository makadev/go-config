package config

// Options configures how the configuration should be loaded
type Options struct {
	// ConfigPaths specifies custom paths to search for configuration files.
	// If empty, defaults to ["config.local.yaml", "config.local.json", "config.yaml", "config.json"]
	ConfigPaths []string

	// EnvPrefix adds a prefix to all environment variable names
	EnvPrefix string

	// SkipEnv disables environment variable loading
	SkipEnv bool

	// SkipFiles disables configuration file loading
	SkipFiles bool

	// Dump Secret
	Secret bool

	// String to use for secret masking
	SecretWith string
}

func NewOptions() *Options {
	return &Options{
		ConfigPaths: DefaultConfigPaths,
		EnvPrefix:   "",
		SkipEnv:     false,
		SkipFiles:   false,
		Secret:      true,
		SecretWith:  "[REDACTED]",
	}
}

// DefaultConfigPaths defines the default order of configuration file loading
var DefaultConfigPaths = []string{
	"config.local.yaml",
	"config.local.json",
	"config.yaml",
	"config.json",
}
