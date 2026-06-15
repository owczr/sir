package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jakubowczarek/sir/internal/config"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Configure GitHub authentication",
	Long:  `Set up your GitHub personal access token and optionally configure GitHub Enterprise host.`,
	RunE:  runAuth,
}

func init() {
	rootCmd.AddCommand(authCmd)
}

func runAuth(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	// Get GitHub host
	fmt.Print("GitHub host (press Enter for github.com): ")
	host, _ := reader.ReadString('\n')
	host = strings.TrimSpace(host)
	if host == "" {
		host = "github.com"
	}

	// Get GitHub username
	fmt.Print("GitHub username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	// Get GitHub token
	fmt.Print("GitHub personal access token: ")
	token, _ := reader.ReadString('\n')
	token = strings.TrimSpace(token)
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	// Load existing config or create new one
	cfg, err := config.Load()
	if err != nil {
		cfg = &config.Config{}
	}

	// Update auth settings
	cfg.GitHubHost = host
	cfg.GitHubUsername = username
	cfg.GitHubToken = token

	// Save config
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("\n✓ Authentication configured successfully!")
	return nil
}
