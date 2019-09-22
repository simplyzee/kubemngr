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
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	getter "github.com/hashicorp/go-getter"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "A tool manage different kubectl versions inside a workspace.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			err := DownloadKubectl(args[0])

			if err != nil {
				log.Fatal(err)
			}

		} else {
			fmt.Println("specify a kubectl version to install")
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

//DownloadKubectl - download user specified version of kubectl
func DownloadKubectl(version string) error {

	// TODO use tmp directory to download instead of kubemngr.
	// This was failing originally with the error: invalid cross-link device
	// filepath := "/tmp/"

	// TODO better sanity check for checking arg is valid
	if len(version) == 0 {
		log.Fatal(0)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	kubectl := fmt.Sprintf("%v/.kubemngr/kubectl-%v", homeDir, version)

	// Check if current version already exists
	if _, err = os.Stat(kubectl); err == nil {
		log.Fatalf("%s is already installed.", version)
	}

	// Create temp file of kubectl version in tmp directory
	out, err := os.Create(kubectl)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	uname := getOSInfo()
	// Compare system name to set value for building url to download kubectl binary
	if uname.Sysname != "Linux" && uname.Sysname != "Darwin" {
		log.Fatalf("Unsupported OS: %s\nCheck github.com/zee-ahmed/kubemngr for issues.", uname.Sysname)
	}
	if uname.Machine != "arm" && uname.Machine != "arm64" && uname.Machine != "x86_64" {
		log.Fatalf("Unsupported arch: %s\nCheck github.com/zee-ahmed/kubemngr for issues.", uname.Machine)
	}

	var sys = strings.ToLower(uname.Sysname)
	var machine string
	if uname.Machine == "x86_64" {
		machine = "amd64"
	} else {
		machine = strings.ToLower(uname.Machine)
	}

	// Check to make sure the file is a binary before moving the contents over to the user's home dir
	url := "https://storage.googleapis.com/kubernetes-release/release/%v/bin/%v/%v/kubectl"
	client := getter.Client{
		Src:              fmt.Sprintf(url, version, sys, machine),
		Dst:              kubectl,
		ProgressListener: defaultProgressBar,
	}
	fmt.Printf("Downloading %v\n", client.Src)
	err = client.Get()
	if err != nil {
		log.Fatal(err)
	}

	// elf - application/x-executable check
	mime, _, err := mimetype.DetectFile(kubectl)
	if mime != "application/octet-stream" {
		fmt.Printf("The downloaded binary is not in the expected format. Please check the version and try again.")
		os.Remove(kubectl)
		os.Exit(1)
	}

	// Set executable permissions on the kubectl binary
	if err := os.Chmod(kubectl, 0755); err != nil {
		log.Fatal(err)
	}

	return nil
}

type uname struct {
	Sysname string
	Machine string
}

func getOSInfo() uname {
	var utsname unix.Utsname

	if err := unix.Uname(&utsname); err != nil {
		fmt.Printf("Uname: %v", err)
	}

	return uname{
		Sysname: string(bytes.Trim(utsname.Sysname[:], "\x00")),
		Machine: string(bytes.Trim(utsname.Machine[:], "\x00")),
	}
}
