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
	lastSeenRunning bool
	sessionStart    time.Time
}

// NewDetector creates a new Ableton detector instance
func NewDetector() *Detector {
	return &Detector{}
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
		// Process not found, reset session if it was running before
		if d.lastSeenRunning {
			d.sessionStart = time.Time{}
			d.lastSeenRunning = false
		}
		return info
	}

	info.IsRunning = true

	// Initialize session start time if this is the first time we see it running
	if !d.lastSeenRunning {
		d.sessionStart = time.Now()
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
