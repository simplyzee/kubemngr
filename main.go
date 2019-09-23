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

	"github.com/zee-ahmed/kubemngr/cmd"
)

var clientVersion = "0.1.2"

func main() {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	directory := homeDir + "/.kubemngr"
	createDirectory(directory)

	binDirectory := homeDir + "/.local/bin"
	createDirectory(binDirectory)

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
