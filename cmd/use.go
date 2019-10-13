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

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Use a specific version of one of the downloaded kubectl binaries",
	Run: func(cmd *cobra.Command, args []string) {
		err := UseKubectlBinary(args[0])
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}

// UseKubectlBinary - sets kubectl to the version specified
func UseKubectlBinary(version string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	kubectlVersion := homeDir + "/.kubemngr/kubectl-" + version
	kubectlLink := homeDir + "/.local/bin/kubectl"

	_, err = os.Stat(kubectlVersion)
	if os.IsNotExist(err) {
		log.Printf("kubectl %s does not exist", version)
		os.Exit(1)
	}

	if _, err := os.Lstat(kubectlLink); err == nil {
		os.Remove(kubectlLink)
	}

	err = os.Symlink(kubectlVersion, kubectlLink)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("kubectl version set to %s", version)

	return nil
}
