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
	_ "github.com/google/go-github/v27/github"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"os"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "A tool manage different kubectl versions inside a workspace.",
	//RunE: func(cmd *cobra.Command, args []string) error {
	//	return errors.New("provide a kubectl version")
	//},
	Run: func(cmd *cobra.Command, args []string) {
		DownloadKubectl(args[0])
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func DownloadKubectl(arg string) (int, error) {
	//TODO: Get OS information machine / architecture
	// Check if version already exists before downloading
	
	userInput := arg

	if len(userInput) == 0 {
		log.Fatal(userInput)
	}

	fmt.Printf("Downloading kubectl version %s", userInput)

	res, err := http.Get("https://storage.googleapis.com/kubernetes-release/release/" + userInput + "/bin/linux/amd64/kubectl")
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create(homeDir + "/.kubemngr/" + userInput)
	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Printf("Download Complete")
}