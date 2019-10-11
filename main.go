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
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zee-ahmed/kubemngr/cmd"
)

var (
	clientVersion = "0.1.2"
	usrLocalBin   = "/usr/local/bin"
)

type paths []string

func (p paths) indexOf(element string) int {
	for k, v := range p {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func main() {
	// set kubemngr directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	directory := homeDir + "/.kubemngr"
	createDirectory(directory)
	binDirectory := homeDir + "/.local/bin"
	createDirectory(binDirectory)

	path, exists := os.LookupEnv("PATH")
	if !exists {
		log.Fatal("Can't access environment variable: PATH")
	}

	var paths paths = strings.Split(path, ":")
	pathsBeforeUsrLocalBin := paths[:paths.indexOf(usrLocalBin)]
	if pathsBeforeUsrLocalBin.indexOf(binDirectory) < 0 {
		fmt.Printf("PATH does not give precedent to %v/.local/bin. kubectl will be executed from /usr/local/bin unless PATH is amended.\n\n", homeDir)

		shell, exists := os.LookupEnv("SHELL")
		if !exists || shell != "/bin/zsh" && shell != "/bin/bash" {
			os.Exit(0)
		}

		fmt.Printf("\tDetected shell is %v, suggested amendment:\n\t", shell)
		if shell == "/bin/zsh" {
			fmt.Println(`echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc`)
		} else if shell == "/bin/bash" {
			fmt.Println(`echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc`)
		}

		os.Exit(0)
	}

	cmd.Execute(clientVersion)
}

func createDirectory(dirName string) bool {
	src, err := os.Stat(dirName)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirName, 0755)
		if errDir != nil {
			panic(err)
		}
		return true
	}

	if src.Mode().IsRegular() {
		fmt.Println(dirName, "already exist as a file!")
		return false
	}

	return false
}
