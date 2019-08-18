package projects

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/lavrahq/cli/products"
	"github.com/lavrahq/cli/utilities/dir"
	"gopkg.in/yaml.v2"
)

// ProjectNetwork is a representation of a "network" entity within
// the deployment platform.
type ProjectNetwork struct {
	Name   string            `yaml:"name"`
	Labels map[string]string `yaml:"labels"`
}

// Project holds project-related configuration information.
type Project struct {
	Directory   dir.Directory
	Name        string                            `yaml:"name"`
	Description string                            `yaml:"description"`
	Products    map[string]products.ProductConfig `yaml:"products"`
	Network     ProjectNetwork                    `yaml:"network"`
}

// Untrack removes tracking of a project, leaving files intact.
// func (p Project) Untrack() error {
// 	return nil
// }

// Track logs a project from a dir and begins tracking it.
// func (p Project) Track(dir) error {

// }

// Load loads a project and returns a Project struct instance for manipulating
// the project.
func Load(dir dir.Directory) (Project, error) {
	var proj = Project{
		Directory: dir,
	}

	if !dir.IsProject() {
		return proj, errors.New("not a project directory")
	}

	fileData, err := ReadProjectFile(dir)
	if err != nil {
		return proj, err
	}

	if err := yaml.Unmarshal(fileData, &proj); err != nil {
		return proj, err
	}

	return proj, nil
}

// NewProject creates a project in the specified directory and adds it to
// the tracked projects if track is specified.
func NewProject(dir dir.Directory, track bool) (Project, error) {
	var proj = Project{
		Directory: dir,
	}

	if dir.IsProject() {
		proj, err := Load(dir)
		if err != nil {
			return proj, nil
		}

		return proj, nil
	}

	return proj, nil
}

// ReadProjectFile reads a project.yml file from a given directory.
func ReadProjectFile(dir dir.Directory) ([]byte, error) {
	if !dir.IsProject() {
		return nil, fmt.Errorf("expected project file at %s", dir.ProjectPath())
	}

	dat, err := ioutil.ReadFile(dir.ProjectPath())
	if err != nil {
		return nil, err
	}

	return dat, nil
}
