/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var licenseName string
var location string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new license to collection",
	// TODO: write a detail description for this command
	Long: `Description`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(location); os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", location)
		}

		licenseName = processLicenseName(licenseName)
		content, err := os.ReadFile(location)
		if err != nil {
			return fmt.Errorf("an error occurred while reading file: %s", location)
		}

		// TODO: control the file permission
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("an error occurred while reading user home directory")
		}

		// TODO: control the directory permission
		licenseDir := strings.Join([]string{homeDir, ".config", "license-generator", "licenses"}, "/")
		err = os.MkdirAll(licenseDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("an error occurred while configuring data directory")
		}

		// TODO: check if license name already exists
		finalLocation := strings.Join([]string{licenseDir, licenseName}, "/")
		err = os.WriteFile(finalLocation, content, os.ModePerm)
		if err != nil {
			return fmt.Errorf("an error occurred while writing to file: %s", finalLocation)
		}

		fmt.Printf("License %s added\n", licenseName)

		return nil
	},
}

func processLicenseName(name string) string {
	name = strings.Trim(name, " ")
	name = strings.Replace(name, " ", "-", -1)
	return name
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&licenseName, "name", "n", "", "name of the license")
	addCmd.Flags().StringVarP(&location, "location", "l", "", "location of license to add")
	addCmd.MarkFlagRequired("name")
	addCmd.MarkFlagRequired("location")
	addCmd.MarkFlagFilename("location")
}
