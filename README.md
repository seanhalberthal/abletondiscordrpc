# Ableton Live Discord Rich Presence

A Go application that displays Discord Rich Presence status when Ableton Live 12 is running on macOS.

![Discord Rich Presence Example](https://img.shields.io/badge/Discord-Rich%20Presence-7289da?style=for-the-badge&logo=discord&logoColor=white)
![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![macOS](https://img.shields.io/badge/macOS-000000?style=for-the-badge&logo=apple&logoColor=white)

## Features

- 🎵 **Automatic Detection**: Monitors Ableton Live process automatically
- 📝 **Project Name Display**: Shows current project name in Discord status
- 🎛️ **Custom Status**: Set personalized status messages via CLI
- ⏱️ **Session Timing**: Displays how long you've been working
- 🔄 **Auto-Reconnect**: Handles Discord disconnections gracefully
- 🖥️ **macOS Native**: Uses AppleScript for seamless integration

## Installation

### Prerequisites

- macOS (Intel or Apple Silicon)
- Discord running on your system
- Ableton Live 12

> **Note**: Make sure "Display current activity as a status message" is enabled in Discord Settings > Activity Privacy

### Quick Install

```bash
# Open Terminal and run the following commands to install the app. If you're not sure where to run them, you can download the repo to your desktop first.

# Run this first
git clone https://github.com/seanhalberthal/abletondiscordrpc.git

# Then copy and run these 3 together to install
cd abletondiscordrpc
sudo cp releases/ableton-discord-rpc /usr/local/bin/
sudo chmod +x /usr/local/bin/ableton-discord-rpc

# Note: you may need to enter your password for sudo access
```

### Start Using It

```bash
# Run the app
ableton-discord-rpc

# That's it! Open Ableton Live and your Discord status will update automatically
```

## Usage

### How It Works

1. **Start the app** - Run `ableton-discord-rpc` in Terminal
2. **Open Ableton Live** - Your Discord status automatically updates
3. **Customize your vibe** - Type commands to set custom status messages

### Commands

While the app is running, type these commands:

```bash
status making beats          # Set custom status
status cooking up fire       # Change your vibe
status bangers comin up      # Whatever fits your mood
quit                         # Exit the app
```

### What Your Friends See

Your Discord status will show:

- 🎵 **Your project name** (e.g., "My New Track")
- 🔥 **Your custom status** (e.g., "making beats")
- **How long you've been working** (automatic timer)

## Technical Details

### Process Detection

- Uses `pgrep -f "Live"` to detect Ableton Live process
- Monitors process state every 15 seconds
- Automatically starts/stops session timing

### Project Name Extraction

- Retrieves window title using AppleScript
- Parses project name from "ProjectName - Ableton Live 12" format
- Falls back to "Untitled Project" if unavailable

### Discord Integration

- Uses the `rich-go` library for Discord RPC
- Implements automatic reconnection on connection failures
- Supports custom status messages via CLI commands

## Development

### Building

```bash
# Build for current platform
go build -o ableton-discord-rpc

# Build for macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o ableton-discord-rpc-intel

# Build for macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o ableton-discord-rpc-arm64
```

## Troubleshooting

### Common Issues

**Discord status not showing:**

- Ensure Discord is running and you're logged in
- Check that the Discord App ID assets are properly uploaded
- Verify Discord allows Rich Presence in settings

**Ableton not detected:**

- Make sure Ableton Live 12 is running
- Check that the process name contains "Live"
- Verify AppleScript permissions in System Preferences

**Project name shows as "Unknown Project":**

- Ensure Ableton Live window is active
- Check System Preferences > Security & Privacy > Accessibility
- Grant terminal/application AppleScript permissions

## License

This project is licensed under the MIT License.

## Acknowledgments

- [rich-go](https://github.com/hugolgst/rich-go) - Discord Rich Presence library
- [Ableton Live](https://www.ableton.com/live/) - Digital Audio Workstation
- [Discord](https://discord.com/) - Communication platform

---

Made with ❤️ for music producers using Ableton Live