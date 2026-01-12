package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/weprodev/wpd-message-gateway/config"
	"github.com/weprodev/wpd-message-gateway/contracts"
	"github.com/weprodev/wpd-message-gateway/manager"
)

// ANSI colors for prettier output
const (
	ColorReset  = "\033[0m"
	ColorBlue   = "\033[34m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorRed    = "\033[31m"
	ColorCyan   = "\033[36m"
)

func main() {
	printHeader()
	reader := bufio.NewReader(os.Stdin)

	// 1. Check/Setup .env
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		printWarn("No .env file found.")
		if promptBool(reader, "Create .env from .env.example? [Y/n]: ", true) {
			if err := copyFile(".env.example", ".env"); err != nil {
				printError("Failed to copy .env.example: %v", err)
				return
			}
			printSuccess("Created .env file.")
			fmt.Println("\nPlease open .env and configure your API keys, then run 'make sandbox' again.")
			return
		} else {
			printError("Cannot proceed without configuration.")
			return
		}
	}

	// 2. Load .env
	if err := loadEnvFile(".env"); err != nil {
		printError("Failed to load .env: %v", err)
		return
	}

	// 3. Load Config
	cfg, err := config.LoadFromEnv()
	if err != nil {
		printError("Failed to load config: %v", err)
		return
	}

	// 4. Validate Configuration Presence
	if len(cfg.EmailProviders) == 0 {
		printError("No email providers configured in .env.")
		fmt.Println("Please set at least one provider (e.g., MESSAGE_MAILGUN_API_KEY) in your .env file.")
		return
	}

	// 5. Init Manager
	mgr, err := manager.New(cfg)
	if err != nil {
		printError("Failed to initialize manager: %v", err)
		return
	}

	printSuccess(fmt.Sprintf("Loaded %d email provider(s)", len(cfg.EmailProviders)))

	// 6. Interactive Loop
	for {
		fmt.Println()
		printSection("Select Message Type")
		fmt.Println("1. Email")
		fmt.Println("2. SMS (Planned)")
		fmt.Println("3. Push (Planned)")
		fmt.Println("4. Chat (Planned)")
		fmt.Println("q. Quit")

		choice := prompt(reader, "Enter choice [1]: ", "1")

		switch choice {
		case "1":
			handleEmail(reader, mgr)
		case "2", "3", "4":
			printWarn("This provider type is not yet implemented/configured.")
		case "q", "exit":
			fmt.Println("Bye!")
			return
		default:
			printError("Invalid choice")
		}
	}
}

func handleEmail(reader *bufio.Reader, mgr *manager.Manager) {
	printSection("Email Sandbox")

	// Select Provider
	providers := mgr.AvailableEmailProviders()
	if len(providers) == 0 {
		printError("No email providers configured! Check your .env file.")
		return
	}

	fmt.Println("Available Providers:")
	for i, p := range providers {
		fmt.Printf("%d. %s\n", i+1, p)
	}

	pIndexStr := prompt(reader, "Select Provider [1]: ", "1")
	var providerName string
	// Simplified selection logic (assuming 1-based index)
	// In a real CLI lib we'd parse int, but keeping it KISS
	if pIndexStr == "1" && len(providers) >= 1 {
		providerName = providers[0]
	} else {
		// Fallback: try to find by name or index
		// For now simple assumption:
		printError("Invalid provider selection")
		return
	}

	// Inputs
	to := prompt(reader, "To (email): ", "")
	if to == "" {
		printError("Recipient required")
		return
	}

	subject := prompt(reader, "Subject [Test Email]: ", "Test Email")
	body := "<h1>Sandbox Test</h1><p>This is a test email sent from wpd-message-gateway sandbox.</p>"

	// Confirmation
	fmt.Printf("\n%sSending email to %s via %s...%s\n", ColorYellow, to, providerName, ColorReset)

	// Send
	start := time.Now()
	result, err := mgr.SendEmailWith(context.Background(), providerName, &contracts.Email{
		To:      []string{to},
		Subject: subject,
		HTML:    body,
	})

	if err != nil {
		printError("Failed to send: %v", err)
	} else {
		duration := time.Since(start)
		printSuccess("Email Sent Successfully!")
		fmt.Printf("ID: %s\n", result.ID)
		fmt.Printf("Time: %s\n", duration)
	}
}

// Helpers

func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			// Remove quotes if present
			val = strings.Trim(val, `"'`)
			_ = os.Setenv(key, val)
		}
	}
	return scanner.Err()
}

func prompt(r *bufio.Reader, label string, def string) string {
	fmt.Print(ColorCyan + label + ColorReset)
	input, _ := r.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return def
	}
	return input
}

func promptBool(r *bufio.Reader, label string, def bool) bool {
	// We handle the label manually to avoid double printing default
	fmt.Printf("%s%s%s", ColorCyan, label, ColorReset)

	input, _ := r.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "" {
		return def
	}
	return input == "y" || input == "yes"
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func printHeader() {
	fmt.Print(ColorBlue)
	fmt.Println("========================================")
	fmt.Println("   WPD Message Gateway - SANDBOX        ")
	fmt.Println("========================================")
	fmt.Print(ColorReset)
}

func printSection(title string) {
	fmt.Printf("\n%s--- %s ---%s\n", ColorBlue, title, ColorReset)
}

func printSuccess(msg string) {
	fmt.Printf("%s✅ %s%s\n", ColorGreen, msg, ColorReset)
}

func printWarn(msg string) {
	fmt.Printf("%s⚠️  %s%s\n", ColorYellow, msg, ColorReset)
}

func printError(format string, args ...interface{}) {
	fmt.Printf("%s❌ %s%s\n", ColorRed, fmt.Sprintf(format, args...), ColorReset)
}
