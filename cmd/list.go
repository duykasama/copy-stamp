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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all copyright templates",
	// TODO: write a detail description for this command
	Long: `Description`,
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("an error occurred while reading user home directory")
		}

		templateDir := strings.Join([]string{homeDir, config.TemplatesLocation}, "/")
		dirs, err := os.ReadDir(templateDir)
		if err != nil {
			return fmt.Errorf("an error occurred while reading templates directory")
		}
		if len(dirs) == 0 {
			fmt.Println("There is no templates added.")
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
