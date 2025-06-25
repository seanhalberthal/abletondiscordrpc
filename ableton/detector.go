package ableton

import (
	"os/exec"
	"strings"
	"time"
)

// Info represents information about the current Ableton Live session
type Info struct {
	IsRunning   bool
	ProjectName string
	StartTime   time.Time
}

// Detector handles Ableton Live process and project detection
type Detector struct {
	lastSeenRunning  bool
	sessionStart     time.Time
	totalSessionTime time.Duration
	lastStopTime     time.Time
	sessionTimeout   time.Duration // Reset session after this much idle time
}

// NewDetector creates a new Ableton detector instance
func NewDetector() *Detector {
	return &Detector{
		sessionTimeout: 2 * time.Hour, // Reset session after 2 hours of inactivity
	}
}

// GetInfo returns current Ableton Live information
func (d *Detector) GetInfo() *Info {
	info := &Info{
		IsRunning:   false,
		ProjectName: "Unknown Project",
		StartTime:   d.sessionStart,
	}

	// Check if Ableton Live is running using pgrep
	cmd := exec.Command("pgrep", "-f", "Live")
	err := cmd.Run()
	if err != nil {
		// Process not found, update session tracking if it was running before
		if d.lastSeenRunning {
			// Add the time since session start to total session time
			if !d.sessionStart.IsZero() {
				d.totalSessionTime += time.Since(d.sessionStart)
			}
			d.lastStopTime = time.Now()
			d.lastSeenRunning = false
		}

		// Return info with preserved session start time for display
		if !d.sessionStart.IsZero() {
			info.StartTime = d.sessionStart
		}
		return info
	}

	info.IsRunning = true

	// Initialize session start time if this is the first time we see it running
	if !d.lastSeenRunning {
		// Check if session has timed out
		if !d.lastStopTime.IsZero() && time.Since(d.lastStopTime) > d.sessionTimeout {
			// Reset session after timeout
			d.sessionStart = time.Now()
			d.totalSessionTime = 0
		} else if d.sessionStart.IsZero() {
			// First time ever
			d.sessionStart = time.Now()
		} else {
			// Continue from previous session, accounting for accumulated time
			d.sessionStart = time.Now().Add(-d.totalSessionTime)
		}
		d.lastSeenRunning = true
	}
	info.StartTime = d.sessionStart

	// Get project name from window title
	windowTitle := d.getWindowTitle()
	if windowTitle != "" {
		info.ProjectName = d.parseProjectName(windowTitle)
	}

	return info
}

// getWindowTitle retrieves the current Ableton Live window title using AppleScript
func (d *Detector) getWindowTitle() string {
	// Try multiple AppleScript approaches
	scripts := []string{
		// Direct window title from Ableton Live application
		`tell application "Ableton Live 12 Suite" to get name of front window`,
		// Alternative approach using System Events
		`tell application "System Events" to get name of front window of application process "Live"`,
		// Fallback using window title from any Ableton process
		`tell application "System Events" to get name of front window of first application process whose name contains "Live"`,
	}

	for _, script := range scripts {
		cmd := exec.Command("osascript", "-e", script)
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			title := strings.TrimSpace(string(output))
			if title != "" && title != "missing value" {
				return title
			}
		}
	}

	return ""
}

// parseProjectName extracts the project name from the window title
func (d *Detector) parseProjectName(windowTitle string) string {
	if windowTitle == "" {
		return "Untitled Project"
	}

	// Ableton window titles are typically "ProjectName - Ableton Live 12 Suite"
	parts := strings.Split(windowTitle, " - ")
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}

	return windowTitle
}
