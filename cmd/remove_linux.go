// +build linux

/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a kubectl version from machine",
	Run: func(cmd *cobra.Command, args []string) {
		err := RemoveKubectlVersion(args[0])
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// RemoveKubectlVersion - removes specific kubectl version from machine
func RemoveKubectlVersion(version string) error {
	// Get user home directory path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	kubectlVersion := homeDir + "/.kubemngr/kubectl-" + version

	// Check if version exists to remove it
	_, err = os.Stat(kubectlVersion)
	if err == nil {
		fmt.Printf("Removing kubectl %s", version)
		os.Remove(kubectlVersion)
		return nil
	} else {
		fmt.Printf("kubectl version %s doesn't exist", version)
		return nil
	}
}
