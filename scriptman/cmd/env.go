// Copyright © 2017 Rodrigue Cloutier <rodcloutier at gmail dot com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/rodcloutier/scriptman"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Outputs the environment variable needed to use the scripts",
	Long: `Command that will parse the scriptman.yaml file and output the
environment needed to use the scripts
    `,
	Run: func(cmd *cobra.Command, args []string) {
		// Read a requirement in the local directory
		filename, _ := filepath.Abs("./scriptman.yml")
		yamlFile, err := ioutil.ReadFile(filename)

		if err != nil {
			panic(err)
		}

		var project scriptman.Project

		err = yaml.Unmarshal(yamlFile, &project)
		if err != nil {
			panic(err)
		}

		err = env(project)

		if err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(envCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// envCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// envCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func env(project scriptman.Project) error {
	paths := []string{}
	for _, req := range project.Requirements {
		paths = append(paths, req.FullDestination())
	}
	paths = append(paths, "$PATH")

	fmt.Println("export PATH=" + strings.Join(paths, ":"))

	return nil
}
