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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

type Kubernetes []struct {
	TagName string `json:"tag_name"`
}

// listRemoteCmd represents the listRemote command
var listRemoteCmd = &cobra.Command{
	Use:   "listRemote",
	Short: "List available remote kubectl versions to download and install",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Getting available list of kubectl versions to install")

		k8s := Kubernetes{}
		_ = ListAvailableRemotes("https://api.github.com/repos/kubernetes/kubernetes/releases?per_page=100", &k8s)

		for _, releases := range k8s {
			version := releases.TagName

			// filter out  alpha and release candidatest and only show stable releases
			filterStable := strings.NewReplacer("-rc.1", "", "-beta.2", "", "-beta.1", "", "-alpha.3", "", "-alpha.2", "", "-alpha.1", "", "-rc.2", "", "-rc.3", "")
			stable := filterStable.Replace(version)

			fmt.Println(stable)
		}
	},
}

func init() {
	rootCmd.AddCommand(listRemoteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listRemoteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listRemoteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ListAvailableRemotes(url string, target interface{}) error {

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	jsonErr := json.Unmarshal(body, &target)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return jsonErr
}
