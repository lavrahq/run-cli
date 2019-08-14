package projects

import (
	"fmt"
)

// Project holds project-related configuration information.
type Project struct {
	Name string `yaml:"name"`
}

// NewProject creates a project in the specified directory and adds it to
// the projects paths.
func NewProject(dir string) {
	fmt.Println("the dir is", dir)
}
