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
	"fmt"
	"license-generator/config"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply license template to source files",
	// TODO: write a detail description for this command
	Long: `Description`,
	RunE: func(cmd *cobra.Command, args []string) error {
		licenseName, err := cmd.Flags().GetString("name")
		if err != nil {
			return fmt.Errorf("license name is required")
		}

		destination, err := cmd.Flags().GetString("destination")
		if err != nil {
			return fmt.Errorf("destination is required")
		}

		extensions, err := cmd.Flags().GetStringArray("extension")
		if err != nil {
			return fmt.Errorf("error reading extensions")
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("an error occurred while reading user home directory")
		}

		license := strings.Join([]string{homeDir, config.LicenseLocation, licenseName}, "/")
		if !fileExists(license) {
			return fmt.Errorf("license `%s` does not exist.", licenseName)
		}

		if !dirExists(destination) {
			return fmt.Errorf("destination `%s` does not exist.", destination)
		}

		licenseContent, err := os.ReadFile(license)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", license, err)
		}

		count := 0
		err = filepath.WalkDir(destination, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("error accessing path %s: %w", path, err)
			}

			if !entry.IsDir() && atLeastEndsWith(entry.Name(), extensions) {
				existingContent, err := os.ReadFile(path)
				if err != nil {
					return fmt.Errorf("error reading file %s: %w", path, err)
				}

				if copyMarkAlreadyExists(licenseContent, existingContent) {
					return nil
				}

				licenseContentWithNewLine := licenseContent
				if existingContent[0] != '\n' {
					licenseContentWithNewLine = append(licenseContent, '\n')
				}

				updatedContent := append(licenseContentWithNewLine, existingContent...)
				err = os.WriteFile(path, updatedContent, 0644)
				if err != nil {
					return fmt.Errorf("error writing license for file %s: %w", path, err)
				}

				count++
			}

			return nil
		})

		fmt.Printf("Updated license in %d file(s).\n", count)

		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", destination, err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.Flags().StringP("name", "n", "", "name of license to apply")
	applyCmd.Flags().StringP("destination", "d", "", "where to apply license")
	applyCmd.Flags().StringArrayP("extension", "x", []string{}, "file extensions to be applied")
	applyCmd.MarkFlagRequired("name")
	applyCmd.MarkFlagRequired("destination")
	applyCmd.MarkFlagDirname("destination")
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func atLeastEndsWith(s string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}

	return false
}

func copyMarkAlreadyExists(copyMark, content []byte) bool {
	for i := 0; i < len(copyMark); i++ {
		if copyMark[i] != content[i] {
			return false
		}

	}

	return true
}
