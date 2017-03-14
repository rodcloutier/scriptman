package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/rodcloutier/scriptman/pkg"
)

type Requirement struct {
	Package     string `yaml:"package"`     // full url of package
	Version     string `yaml:"version"`     // sha or tag, absent defaults to master
	Destination string `yaml:"destination"` // defaults to 'vendor'?
	Repository  string `yaml:"repo"`        // the repository to use for cloning
}

type Project struct {
	Requirements []Requirement `yaml:"requirements"`
}

// UnmarshalYAML implements default value at load time for the Requirement type
func (r *Requirement) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawRequirement Requirement
	raw := rawRequirement{
		Destination: "vendor",
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}
	*r = Requirement(raw)
	return nil
}

func (r *Requirement) FullDestination() string {
	dir, _ := filepath.Abs(filepath.Join(r.Destination, r.Package))
	return dir
}

func (r *Requirement) RepositoryURL() string {

	// TODO only if req.Repository is a file path
	repo, _ := filepath.Abs(r.Repository)

	return repo
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

func installRequirement(req Requirement) error {

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

func install(project Project) error {

	var err error
	for _, req := range project.Requirements {
		err = installRequirement(req)
		if err != nil {
			return err
		}
	}
	return nil
}

func env(project Project) error {
	paths := []string{}
	for _, req := range project.Requirements {
		paths = append(paths, req.FullDestination())
	}
	paths = append(paths, "$PATH")

	fmt.Println("export PATH=" + strings.Join(paths, ":"))

	return nil
}

func main() {

	// Read a requirement in the local directory
	filename, _ := filepath.Abs("./scriptman.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var project Project

	err = yaml.Unmarshal(yamlFile, &project)
	if err != nil {
		panic(err)
	}

	env(project)

	if err != nil {
		panic(err)
	}
}
