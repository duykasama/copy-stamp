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
	"copy-stamp/internal"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new copyright template to collection",
	// TODO: write a detail description for this command
	Long: `Description`,
	RunE: func(cmd *cobra.Command, args []string) error {
		templateName, err := cmd.Flags().GetString("name")
		if err != nil {
			return fmt.Errorf("error reading template name")
		}

		location, err := cmd.Flags().GetString("location")
		if err != nil {
			return fmt.Errorf("template location is required")
		}

		if _, err := os.Stat(location); os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", location)
		}

		templateName = processTemplateName(templateName)
		// TODO: control the file permission
		templateContent, err := os.ReadFile(location)
		if err != nil {
			return fmt.Errorf("an error occurred while reading file: %s", location)
		}

		templatesDir, err := internal.EnsureDataDirectoryExists()
		if err != nil {
			return err
		}

		if err != nil {
			return fmt.Errorf("an error occurred while configuring data directory")
		}

		// TODO: check if license name already exists
		templateFile := strings.Join([]string{templatesDir, templateName}, "/")
		err = os.WriteFile(templateFile, templateContent, 0644)
		if err != nil {
			return fmt.Errorf("an error occurred while writing to file: %s", templateFile)
		}

		fmt.Printf("Template %s added\n", templateName)

		return nil
	},
}

func processTemplateName(name string) string {
	name = strings.Trim(name, " ")
	name = strings.Replace(name, " ", "-", -1)
	return name
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringP("name", "n", "", "name of the template")
	addCmd.Flags().StringP("location", "l", "", "location of template to add")
	addCmd.MarkFlagRequired("name")
	addCmd.MarkFlagRequired("location")
	addCmd.MarkFlagFilename("location")
}
