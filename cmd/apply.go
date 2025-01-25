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
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply copyright template to source files",
	// TODO: write a detail description for this command
	Long: `Description`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("\"stamp apply\" requires at least one argument to execute.")
		}

		templateName, err := cmd.Flags().GetString("name")
		if err != nil {
			return fmt.Errorf("template name is required")
		}

		templateContent, err := getTemplateContent(templateName)
		if err != nil {
			return err
		}

		extensions, err := cmd.Flags().GetStringArray("extensions")
		if err != nil {
			return fmt.Errorf("error reading extensions")
		}

		totalFilesApplied := 0
		for _, arg := range args {
			info, err := os.Stat(arg)
			if os.IsNotExist(err) {
				return err
			}

			if info.IsDir() {
				appliedFiles, err := applyToDirectory(arg, templateContent, extensions)
				if err != nil {
					return err
				}
				totalFilesApplied += appliedFiles
			} else {
				applied, err := applyToFile(arg, templateContent)
				if err != nil {
					return err
				}
				if applied {
					totalFilesApplied++
				}
			}
		}

		fmt.Printf("Stamped copyright in %d file(s).\n", totalFilesApplied)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.Flags().StringP("name", "n", "", "name of copyright to stamp")
	applyCmd.Flags().StringArrayP("extensions", "x", []string{}, "file extensions to be applied")
	applyCmd.MarkFlagRequired("name")
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
		if !strings.HasSuffix(s, suffix) {
			return false
		}
	}

	return true
}

func copyrightAlreadyExists(copyMark, content []byte) bool {
	for i := 0; i < len(copyMark); i++ {
		if copyMark[i] != content[i] {
			return false
		}

	}

	return true
}

func applyToDirectory(dirToApply string, templateContent []byte, extensions []string) (int, error) {
	if !dirExists(dirToApply) {
		return 0, fmt.Errorf("destination `%s` does not exist.", dirToApply)
	}

	count := 0
	err := filepath.WalkDir(dirToApply, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		if !entry.IsDir() && atLeastEndsWith(entry.Name(), extensions) {
			applied, err := applyToFile(path, templateContent)
			if err != nil {
				return err
			}
			if applied {
				count++
			}
		}

		return nil
	})

	if err != nil {
		return count, fmt.Errorf("error accessing path %s: %w", dirToApply, err)
	}

	return count, nil
}

func applyToFile(file string, templateContent []byte) (bool, error) {
	fileContent, err := os.ReadFile(file)
	if err != nil {
		return false, fmt.Errorf("error reading file %s: %w", file, err)
	}

	if copyrightAlreadyExists(templateContent, fileContent) {
		return false, nil
	}

	templateContentToApply := templateContent
	if fileContent[0] != '\n' {
		templateContentToApply = append(templateContent, '\n')
	}

	updatedFileContent := append(templateContentToApply, fileContent...)
	err = os.WriteFile(file, updatedFileContent, 0644)
	if err != nil {
		return false, fmt.Errorf("error stamping copyright for file %s: %w", file, err)
	}

	return true, nil
}

func getTemplateContent(templateName string) ([]byte, error) {
	templatesDir, err := internal.EnsureDataDirectoryExists()
	if err != nil {
		return nil, err
	}

	template := strings.Join([]string{templatesDir, templateName}, "/")
	if !fileExists(template) {
		return nil, fmt.Errorf("template `%s` does not exist.", templateName)
	}

	templateContent, err := os.ReadFile(template)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", template, err)
	}

	return templateContent, nil
}
