package logdna

import (
	"fmt"
	"os"
	"time"
)

const (
	defaultEnvAPIKey          = "LOGDNA_API_KEY" // #nosec G101
	defaultEndpoint           = "https://logs.logdna.com/logs/ingest"
	defaultTimeoput           = time.Second * 30
	defaultCheckpointSize     = 64
	defaultCheckpointInterval = time.Second * 1
)

var envAPIKey string

func init() {
	envAPIKey = os.Getenv(defaultEnvAPIKey)
}

// Config contains parameters for the LogDNA client.
type Config struct {
	APIKey string

	App      string
	Env      string
	Hostname string
	IP       string
	MacAddr  string

	Tags []string

	// IncludeStandardMeta bool
	// IndexMeta bool

	MinimumLevel   string
	Sync           bool
	Debug          bool
	NoRetry        bool
	Timeout        time.Duration
	CustomEndpoint string

	CheckpointSize     int
	CheckpointInterval time.Duration

	// used internally
	apikey   string
	endpoint string
	hostname string
	macaddr  string
	ip       string
	timeout  time.Duration
}

func (c Config) Validate() error {
	if len(c.App) > 32 {
		return fmt.Errorf("`App` length is [%d]. It must be lower than 32", len(c.App))
	}
	if len(c.Env) > 32 {
		return fmt.Errorf("`Env` length is [%d]. It must be lower than 32", len(c.Env))
	}
	if len(c.Hostname) > 32 {
		return fmt.Errorf("`Hostname` length is [%d]. It must be lower than 32", len(c.Hostname))
	}
	return nil
}

// Init intializes a config.
func (c *Config) Init() error {
	if err := c.Validate(); err != nil {
		return err
	}

	c.apikey = c.getAPIKey()
	c.endpoint = c.getEndpoint()
	c.hostname = c.getHostname()
	c.macaddr, c.ip = getMacAndIP()
	c.timeout = c.getTimeout()
	return nil
}

func (c Config) getAPIKey() string {
	if c.APIKey != "" {
		return c.APIKey
	}
	return envAPIKey
}

func (c Config) getEndpoint() string {
	if c.CustomEndpoint != "" {
		return c.CustomEndpoint
	}
	return defaultEndpoint
}

func (c Config) getHostname() string {
	if c.Hostname != "" {
		return c.Hostname
	}

	name, _ := os.Hostname()
	return name
}

func (c Config) getTimeout() time.Duration {
	if c.Timeout > 0 {
		return c.Timeout
	}
	return defaultTimeoput
}

func (c Config) getCheckpointSize() int {
	if c.CheckpointSize > 0 {
		return c.CheckpointSize
	}
	return defaultCheckpointSize
}

func (c Config) getCheckpointInterval() time.Duration {
	if c.CheckpointInterval > 0 {
		return c.CheckpointInterval
	}
	return defaultCheckpointInterval
}
