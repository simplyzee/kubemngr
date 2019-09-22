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
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/h2non/filetype"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

var sys, machine string

// WriteCounter tracks the total number of bytes
type WriteCounter struct {
	Total uint64
}

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

	// Check if current version already exists
	_, err = os.Stat(homeDir + "/.kubemngr/kubectl-" + version)
	if err == nil {
		log.Printf("kubectl version %s already exists", version)
		return nil
	}

	// Create temp file of kubectl version in tmp directory
	out, err := os.Create(homeDir + "/.kubemngr/kubectl-" + version)
	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	// Get OS information to filter download type i.e linux / darwin
	uname := GetOSInfo()

	// Compare system name to set value for building url to download kubectl binary
	if arrayToString(uname.Sysname) == "Linux" {
		sys = "linux"
	} else if arrayToString(uname.Sysname) == "Darwin" {
		sys = "darwin"
	} else {
		sys = "UNKNOWN"
		fmt.Println("Unknown system")
	}

	if arrayToString(uname.Machine) == "arm" {
		machine = "arm"
	} else if arrayToString(uname.Machine) == "arm64" {
		machine = "arm64"
	} else if arrayToString(uname.Machine) == "x86_64" {
		machine = "amd64"
	} else {
		machine = "UNKNOWN"
		fmt.Println("Unknown machine")
	}

	url := "https://storage.googleapis.com/kubernetes-release/release/" + version + "/bin/" + sys + "/" + machine + "/kubectl"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	// Initialise WriteCounter and copy the contents of the response body to the tmp file
	counter := &WriteCounter{}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		log.Fatal(err)
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Println()

	// Check to make sure the file is a binary before moving the contents over to the user's home dir
	buf, _ := ioutil.ReadFile(homeDir + "/.kubemngr/kubectl-" + version)

	// elf - application/x-executable check
	if !filetype.IsArchive(buf) {
		fmt.Println("failed to download kubectl file. Are you sure you specified the right version?")
		os.Remove(homeDir + "/.kubemngr/kubectl-" + version)
		os.Exit(1)
	}

	// Set executable permissions on the kubectl binary
	if err := os.Chmod(homeDir+"/.kubemngr/kubectl-"+version, 0755); err != nil {
		log.Fatal(err)
	}

	// Rename the tmp file back to the original file and store it in the kubemngr directory
	currentFilePath := homeDir + "/.kubemngr/kubectl-" + version
	newFilePath := homeDir + "/.kubemngr/kubectl-" + version

	err = os.Rename(currentFilePath, newFilePath)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// GetOSInfo - Get operating system information of machine
func GetOSInfo() unix.Utsname {
	var uname unix.Utsname

	if err := unix.Uname(&uname); err != nil {
		fmt.Printf("Uname: %v", err)
	}

	return uname
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

// PrintProgress - Helper function to print progress of a download
func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 50))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}
