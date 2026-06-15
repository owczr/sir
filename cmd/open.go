package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/jakubowczarek/sir/internal/config"
	"github.com/jakubowczarek/sir/internal/github"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open a PR in your browser",
	Long:  `List PRs waiting for your review and open the selected one in your default browser.`,
	RunE:  runOpen,
}

func init() {
	rootCmd.AddCommand(openCmd)
}

func runOpen(cmd *cobra.Command, args []string) error {
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

	// Collect all PRs
	var allPRs []github.PullRequest
	for _, repo := range cfg.Repositories {
		prs, err := client.GetReviewRequests(repo, cfg.GitHubUsername)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to check %s: %v\n", repo, err)
			continue
		}
		allPRs = append(allPRs, prs...)
	}

	if len(allPRs) == 0 {
		fmt.Println("Nothing to review!")
		return nil
	}

	// Display PRs
	fmt.Println("PRs waiting for your review:")
	for i, pr := range allPRs {
		fmt.Printf("  %d. [%s] %s\n", i+1, pr.Repository, pr.Title)
	}

	// Get user selection
	fmt.Print("\nEnter PR number to open (or 0 to cancel): ")
	var choice int
	_, err = fmt.Scanf("%d", &choice)
	if err != nil || choice < 0 || choice > len(allPRs) {
		return fmt.Errorf("invalid selection")
	}

	if choice == 0 {
		fmt.Println("Cancelled")
		return nil
	}

	selectedPR := allPRs[choice-1]
	
	// Open in browser
	if err := openBrowser(selectedPR.URL); err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}

	fmt.Printf("✓ Opened PR #%d in browser\n", selectedPR.Number)
	return nil
}

func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}
