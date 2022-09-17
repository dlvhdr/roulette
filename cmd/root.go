package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dlvhdr/roulette/cmd/app"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "roulette",
	Short:   "Pick randomly from an arbitrary list of items",
	Long:    `Provide a list of items and Roulette will pick one of them randomly.`,
	Example: `choose -t "What's for dinner?" 🍕,🍔,🥓,🥦,🌯,🥩`,
	RunE: func(cmd *cobra.Command, args []string) error {
		title, err := cmd.Flags().GetString("title")
		if err != nil {
			return err
		}

		options, err := cmd.Flags().GetStringSlice("options")
		if err != nil || len(options) == 0 {
			return fmt.Errorf(lipgloss.NewStyle().Foreground(
				lipgloss.Color("#f7768e"),
			).Render("Please provide at least 1 option...\n"))
		}

		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}
		p := tea.NewProgram(app.InitialModel(title, options, debug))
		if err := p.Start(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			return err
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringSliceP("options", "o", []string{}, "The list of options")
	rootCmd.Flags().StringP("title", "t", "", "Optional title")
	rootCmd.Flags().Bool("debug", false, "Show debug information")
}
