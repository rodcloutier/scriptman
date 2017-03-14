package scriptman

import (
	"path/filepath"
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
