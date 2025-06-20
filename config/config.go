package config

import "time"

// Config holds application configuration
type Config struct {
	DiscordAppID    string
	AbletonAppName  string
	PollingInterval time.Duration
	CustomStatus    string
}

// DefaultConfig returns the default application configuration
func DefaultConfig() *Config {
	return &Config{
		DiscordAppID:    "1385679105969885184",
		AbletonAppName:  "Ableton Live 12",
		PollingInterval: 15 * time.Second,
		CustomStatus:    "Making music",
	}
}

// UpdateCustomStatus updates the custom status message
func (c *Config) UpdateCustomStatus(status string) {
	c.CustomStatus = status
}
