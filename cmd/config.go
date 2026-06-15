package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jakubowczarek/sir/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure repositories to monitor",
	Long:  `Add or remove repositories to monitor for PR review requests.`,
	RunE:  runConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func runConfig(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		cfg = &config.Config{}
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nCurrent repositories:")
		if len(cfg.Repositories) == 0 {
			fmt.Println("  (none)")
		} else {
			for i, repo := range cfg.Repositories {
				fmt.Printf("  %d. %s\n", i+1, repo)
			}
		}

		fmt.Println("\nOptions:")
		fmt.Println("  [a] Add repository")
		fmt.Println("  [r] Remove repository")
		fmt.Println("  [q] Quit")
		fmt.Print("\nChoice: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(strings.ToLower(choice))

		switch choice {
		case "a":
			fmt.Print("Enter repository (format: owner/repo): ")
			repo, _ := reader.ReadString('\n')
			repo = strings.TrimSpace(repo)
			if repo == "" {
				fmt.Println("Repository cannot be empty")
				continue
			}
			if !strings.Contains(repo, "/") {
				fmt.Println("Invalid format. Use: owner/repo")
				continue
			}
			cfg.Repositories = append(cfg.Repositories, repo)
			fmt.Printf("✓ Added %s\n", repo)

		case "r":
			if len(cfg.Repositories) == 0 {
				fmt.Println("No repositories to remove")
				continue
			}
			fmt.Print("Enter number to remove: ")
			var num int
			fmt.Fscanf(reader, "%d\n", &num)
			if num < 1 || num > len(cfg.Repositories) {
				fmt.Println("Invalid number")
				continue
			}
			removed := cfg.Repositories[num-1]
			cfg.Repositories = append(cfg.Repositories[:num-1], cfg.Repositories[num:]...)
			fmt.Printf("✓ Removed %s\n", removed)

		case "q":
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
			fmt.Println("\n✓ Configuration saved!")
			return nil

		default:
			fmt.Println("Invalid choice")
		}
	}
}
