/*
*  Copyright (C) 2025 Nguyen Thanh Duy
*
*  This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License
*  as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
*
*  This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY;
*  without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
*  See the GNU General Public License for more details.
*
*  You should have received a copy of the GNU General Public License along with this program.
*  If not, see <https://www.gnu.org/licenses/>.
 */

package cmd

import (
	"copy-stamp/config"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new license to collection",
	// TODO: write a detail description for this command
	Long: `Description`,
	RunE: func(cmd *cobra.Command, args []string) error {
		licenseName, err := cmd.Flags().GetString("name")
		if err != nil {
			return fmt.Errorf("license name is required")
		}

		location, err := cmd.Flags().GetString("location")
		if err != nil {
			return fmt.Errorf("license location is required")
		}

		if _, err := os.Stat(location); os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", location)
		}

		licenseName = processLicenseName(licenseName)
		// TODO: control the file permission
		content, err := os.ReadFile(location)
		if err != nil {
			return fmt.Errorf("an error occurred while reading file: %s", location)
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("an error occurred while reading user home directory")
		}

		// TODO: control the directory permission
		licenseDir := strings.Join([]string{homeDir, config.LicenseLocation}, "/")
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

	addCmd.Flags().StringP("name", "n", "", "name of the license")
	addCmd.Flags().StringP("location", "l", "", "location of license to add")
	addCmd.MarkFlagRequired("name")
	addCmd.MarkFlagRequired("location")
	addCmd.MarkFlagFilename("location")
}
