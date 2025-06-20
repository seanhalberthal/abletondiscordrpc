package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"abletondiscordrpc/ableton"
	"abletondiscordrpc/config"
	"abletondiscordrpc/discord"
)

func main() {
	fmt.Println("Starting Ableton Live Discord Rich Presence...")

	// Initialize configuration
	cfg := config.DefaultConfig()

	// Initialize Discord client
	discordClient := discord.NewClient(cfg.DiscordAppID)
	err := discordClient.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to Discord: %v", err)
	}
	defer discordClient.Disconnect()

	fmt.Println("Connected to Discord!")
	fmt.Println("Commands:")
	fmt.Println("  status <message> - Set custom status")
	fmt.Println("  quit             - Exit the program")

	// Initialize Ableton detector
	detector := ableton.NewDetector()

	// Set up a graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Create a ticker for polling
	ticker := time.NewTicker(cfg.PollingInterval)
	defer ticker.Stop()

	// Start CLI input handler
	go handleCLIInput(cfg)

	// Main loop
	for {
		select {
		case <-ticker.C:
			updateDiscordPresence(discordClient, detector, cfg)
		case <-quit:
			fmt.Println("Shutting down...")
			return
		}
	}
}

// updateDiscordPresence updates the Discord Rich Presence based on Ableton status
func updateDiscordPresence(discordClient *discord.Client, detector *ableton.Detector, cfg *config.Config) {
	abletonInfo := detector.GetInfo()

	if !abletonInfo.IsRunning {
		err := discordClient.SetWaitingActivity(cfg.AbletonAppName, &abletonInfo.StartTime)
		if err != nil {
			// Attempt to reconnect on error
			discordClient.Disconnect()
			if reconnectErr := discordClient.Connect(); reconnectErr != nil {
				log.Printf("Failed to reconnect to Discord: %v", reconnectErr)
			}
		}
		return
	}

	// Ableton is running - show active status
	err := discordClient.SetActivity(
		abletonInfo.ProjectName,
		cfg.CustomStatus,
		cfg.AbletonAppName,
		&abletonInfo.StartTime,
	)
	if err != nil {
		// Attempt to reconnect on error
		discordClient.Disconnect()
		if reconnectErr := discordClient.Connect(); reconnectErr != nil {
			log.Printf("Failed to reconnect to Discord: %v", reconnectErr)
		}
	}
}

// handleCLIInput handles user input for setting custom status
func handleCLIInput(cfg *config.Config) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "quit" || input == "exit" {
			fmt.Println("Exiting...")
			os.Exit(0)
		}

		if strings.HasPrefix(input, "status ") {
			newStatus := strings.TrimSpace(input[7:]) // Remove "status " prefix
			if newStatus != "" {
				cfg.UpdateCustomStatus(newStatus)
				fmt.Printf("Status updated to: %s\n", cfg.CustomStatus)
			} else {
				fmt.Println("Please provide a status message. Usage: status <message>")
			}
			continue
		}

		fmt.Println("Commands:")
		fmt.Println("  status <message> - Set custom status")
		fmt.Println("  quit             - Exit the program")
	}
}
