package discord

import (
	"fmt"
	"time"

	"github.com/hugolgst/rich-go/client"
)

// Client wraps the Discord RPC client with our application logic
type Client struct {
	appID     string
	connected bool
}

// NewClient creates a new Discord RPC client
func NewClient(appID string) *Client {
	return &Client{
		appID: appID,
	}
}

// Connect establishes connection to Discord
func (c *Client) Connect() error {
	err := client.Login(c.appID)
	if err != nil {
		return fmt.Errorf("failed to connect to Discord: %w", err)
	}
	c.connected = true
	return nil
}

// Disconnect closes the connection to Discord
func (c *Client) Disconnect() {
	if c.connected {
		client.Logout()
		c.connected = false
	}
}

// IsConnected returns whether the client is connected to Discord
func (c *Client) IsConnected() bool {
	return c.connected
}

// SetActivity updates the Discord Rich Presence activity
func (c *Client) SetActivity(details, state, largeText string, startTime *time.Time) error {
	if !c.connected {
		return fmt.Errorf("Discord client not connected")
	}

	activity := client.Activity{
		Details:    "🎵 " + details,
		State:      "🔥 " + state,
		LargeImage: "ableton-logo",
		LargeText:  largeText,
		SmallImage: "music-note",
		SmallText:  "Active",
	}

	if startTime != nil && !startTime.IsZero() {
		activity.Timestamps = &client.Timestamps{
			Start: startTime,
		}
	}

	return client.SetActivity(activity)
}

// SetWaitingActivity sets the activity when Ableton is not running
func (c *Client) SetWaitingActivity(appName string) error {
	if !c.connected {
		return fmt.Errorf("Discord client not connected")
	}

	now := time.Now()
	activity := client.Activity{
		State:      "Waiting for " + appName + "...",
		Details:    "Not currently making music",
		LargeImage: "ableton-logo",
		LargeText:  appName,
		Timestamps: &client.Timestamps{
			Start: &now,
		},
	}

	return client.SetActivity(activity)
}
