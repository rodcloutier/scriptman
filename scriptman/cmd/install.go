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
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rodcloutier/scriptman"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		err = install(project)

		if err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func ensureDir(dir string) error {

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0777)
	}
	return nil
}

// TODO return the sha of the current commit
func clone(repository, fullDestination string) error {

	cwd, err := os.Getwd()
	defer os.Chdir(cwd)

	if err = os.Chdir(fullDestination); err != nil {
		log.Printf("Failed to chdir to %s", fullDestination)
		return err
	}

	// If already a git repository, do something

	cmd := exec.Command("git", "clone", repository, ".")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err = cmd.Run(); err != nil {
		log.Printf("Failed to clone repo %s", stderr.String())
		return err
	}

	// Checkout the file
	cmd = exec.Command("git", "checkout", "master")
	if err = cmd.Run(); err != nil {
		log.Printf("Failed to checkout master")
		return err
	}
	return err
}

func installRequirement(req scriptman.Requirement) error {

	var err error

	destination := req.FullDestination()
	if err = ensureDir(destination); err != nil {
		return err
	}

	repository := req.RepositoryURL()
	if err = clone(repository, destination); err != nil {
		return err
	}
	return nil
}

func install(project scriptman.Project) error {

	var err error
	for _, req := range project.Requirements {
		err = installRequirement(req)
		if err != nil {
			return err
		}
	}
	return nil
}
