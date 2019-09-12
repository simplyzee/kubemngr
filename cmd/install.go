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
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/dustin/go-humanize"
	_ "github.com/google/go-github/v27/github"
	"github.com/spf13/cobra"
)

var sys, machine string

type WriteCounter struct {
	Total uint64
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "A tool manage different kubectl versions inside a workspace.",
	//RunE: func(cmd *cobra.Command, args []string) error {
	//	return errors.New("provide a kubectl version")
	//},
	Run: func(cmd *cobra.Command, args []string) {
		err := DownloadKubectl(args[0])

		if err != nil {
			log.Fatal(err)
		}
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

//DownloadKubectl - download user specified version of kubectl
func DownloadKubectl(version string) error {

	// TODO use tmp directory to download instead of kubemngr.
	// This was failing originally with the error: invalid cross-link device
	// filepath := "/tmp/"

	// TODO better sanity check for checking arg is valid
	if len(version) == 0 {
		log.Fatal(0)
	}

	// Get user home directory path
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
	out, err := os.Create(homeDir + "/.kubemngr/kubectl-" + version + ".tmp")
	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	// Get OS information to filter download type i.e linux / darwin
	uname := GetOSInfo()

	// TODO refactor me
	// doesn't work on OSX
	if ArrayToString(uname.Sysname) == "Linux" {
		sys = "linux"
	} else if ArrayToString(uname.Sysname) == "Darwin" {
		sys = "darwin"
	} else {
		sys = "UNKNOWN"
		fmt.Println("Unknown system")
	}

	if ArrayToString(uname.Machine) == "arm" {
		machine = "arm"
	} else if ArrayToString(uname.Machine) == "arm64" {
		machine = "arm64"
	} else if ArrayToString(uname.Machine) == "x86_64" {
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

	// Set executable permissions on the kubectl binary
	if err := os.Chmod(homeDir+"/.kubemngr/kubectl-"+version+".tmp", 0755); err != nil {
		log.Fatal(err)
	}

	// Rename the tmp file back to the original file and store it in the kubemngr directory
	currentFilePath := homeDir + "/.kubemngr/kubectl-" + version + ".tmp"
	newFilePath := homeDir + "/.kubemngr/kubectl-" + version

	err = os.Rename(currentFilePath, newFilePath)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// GetOSInfo - Get operating system information of machine
func GetOSInfo() syscall.Utsname {
	var uname syscall.Utsname

	if err := syscall.Uname(&uname); err != nil {
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

func ArrayToString(x [65]int8) string {
	var buf [65]byte
	for i, b := range x {
		buf[i] = byte(b)
	}
	str := string(buf[:])
	if i := strings.Index(str, "\x00"); i != -1 {
		str = str[:i]
	}
	return str
}
