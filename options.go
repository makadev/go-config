package config

// Options configures how the configuration should be loaded
type Options struct {
	// ConfigPaths specifies custom paths to search for configuration files.
	// If empty, defaults to ["config.local.yaml", "config.local.json", "config.yaml", "config.json"]
	ConfigPaths []string

	// EnvPrefix adds a prefix to all environment variable names
	EnvPrefix string

	// AutoEnv enables automatic environment variable loading by mangling the config key
	// f.e. "server.port" -> EnvPrefix + "SERVER_PORT"
	AutoEnv bool

	// SkipEnv disables environment variable loading
	SkipEnv bool

	// SkipFiles disables configuration file loading
	SkipFiles bool

	// Dump Secret
	Secret bool

	// String to use for secret masking
	SecretWith string

	// Configuration Tags to check
	ConfigTags []string
}

func NewOptions() *Options {
	return &Options{
		ConfigPaths: DefaultConfigPaths,
		EnvPrefix:   "",
		SkipEnv:     false,
		SkipFiles:   false,
		Secret:      true,
		SecretWith:  "[REDACTED]",
		ConfigTags:  DefaultConfigTag,
	}
}

var DefaultConfigTag = []string{
	"config",
	"yaml",
	"json",
}

// DefaultConfigPaths defines the default order of configuration file loading
var DefaultConfigPaths = []string{
	"config.local.yaml",
	"config.local.json",
	"config.yaml",
	"config.json",
}
