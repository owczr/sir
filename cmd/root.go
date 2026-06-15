package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jakubowczarek/sir/internal/config"
	"github.com/jakubowczarek/sir/internal/github"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sir",
	Short: "Should I Review? - Check for GitHub PR review requests",
	Long:  `A simple CLI tool to check if you have any pending GitHub PR review requests.`,
	RunE:  runCheck,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Ensure config directory exists
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "sir")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not create config directory: %v\n", err)
	}
}

func runCheck(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.GitHubToken == "" {
		return fmt.Errorf("GitHub token not configured. Run 'sir auth' to set up authentication")
	}

	if len(cfg.Repositories) == 0 {
		return fmt.Errorf("no repositories configured. Run 'sir config' to add repositories")
	}

	client := github.NewClient(cfg.GitHubToken, cfg.GitHubHost)

	totalPRs := 0
	results := make(map[string]int)

	for _, repo := range cfg.Repositories {
		prs, err := client.GetReviewRequests(repo, cfg.GitHubUsername)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to check %s: %v\n", repo, err)
			continue
		}
		if len(prs) > 0 {
			results[repo] = len(prs)
			totalPRs += len(prs)
		}
	}

	if totalPRs == 0 {
		fmt.Println("Nothing to review!")
		return nil
	}

	fmt.Println("PRs waiting for your review:")
	for repo, count := range results {
		fmt.Printf("  %s: %d PR(s)\n", repo, count)
	}

	return nil
}
