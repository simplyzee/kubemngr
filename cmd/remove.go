/*
Copyright © 2019 Zee Ahmed <zee@simplyzee.dev>

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
}

// RemoveKubectlVersion - removes specific kubectl version from machine
func RemoveKubectlVersion(version string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	kubectlVersion := homeDir + "/.kubemngr/kubectl-" + version

	// Check if version to be removed exists
	_, err = os.Stat(kubectlVersion)
	if err == nil {
		fmt.Printf("Removing kubectl %s", version)
		os.Remove(kubectlVersion)
		return nil
	}

	fmt.Printf("kubectl %s is not installed", version)
	return nil
}
