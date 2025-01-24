/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"license-generator/config"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all license templates",
	// TODO: write a detail description for this command
	Long: `Description`,
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("an error occurred while reading user home directory")
		}

		licenseDir := strings.Join([]string{homeDir, config.LicenseLocation}, "/")
		dirs, err := os.ReadDir(licenseDir)
		if err != nil {
			return fmt.Errorf("an error occurred while reading license directory")
		}
		if len(dirs) == 0 {
			fmt.Println("There is no licenses added.")
		} else {
			for i, dir := range dirs {
				fmt.Printf("%d. %s\n", i+1, dir.Name())
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
