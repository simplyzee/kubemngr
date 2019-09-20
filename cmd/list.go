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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
)

const (
	binaryListURL = "https://api.github.com/repos/kubernetes/kubernetes/releases?per_page=100"
)

var (
	remote bool
)

type kubectlVersion struct {
	Version version.Version
}

func (kc *kubectlVersion) UnmarshalJSON(b []byte) error {
	aux := &struct {
		TagName string `json:"tag_name"`
	}{}

	if err := json.Unmarshal(b, &aux); err != nil {
		log.Fatal(err)
	}
	version, err := version.NewVersion(aux.TagName)
	if err != nil {
		log.Fatal(err)
	}
	kc.Version = *version
	return nil
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed kubectl binary versions. For available versions, see --remote",
	Run: func(cmd *cobra.Command, args []string) {
		var versions []kubectlVersion
		if remote {
			fmt.Println("Fetching remote versions ...")
			versions = fetchRemoteVersions()
		} else {
			versions = fetchLocalVersions()

			if len(versions) > 0 {
				fmt.Println("Installed kubectl versions:")
			} else {
				fmt.Println("No versions installed. See 'kubemngr list --remote' for available versions.")
			}
		}

		re := regexp.MustCompile(`-rc.1|-beta.2|-beta.1|-alpha.3|-alpha.2|-alpha.1|-rc.2|-rc.3`)
		for _, version := range versions {
			if !re.MatchString(version.Version.String()) {
				fmt.Println(version.Version.Original())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&remote, "remote", false, "Get versions from remote")
}

// fetchLocalVersions - List available installed kubectl versions
func fetchLocalVersions() []kubectlVersion {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	kubectl, err := ioutil.ReadDir(homeDir + "/.kubemngr/")
	if err != nil {
		log.Fatal(err)
	}

	list := []kubectlVersion{}
	for _, files := range kubectl {
		file := files.Name()
		name, err := version.NewVersion(strings.Replace(file, "kubectl-", "", -1))
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, kubectlVersion{Version: *name})
	}

	return list
}

// fetchRemoteVersions lists Kubectl binaries available at the configured remote location
func fetchRemoteVersions() []kubectlVersion {
	res, err := http.Get(binaryListURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	list := []kubectlVersion{}
	jsonErr := json.Unmarshal(body, &list)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Version.GreaterThan(&list[j].Version)
	})

	return list
}
