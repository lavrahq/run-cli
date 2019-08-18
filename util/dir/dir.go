package dir

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// Directory is an instance of the directory manipulation utility.
type Directory struct {
	Path string
}

// Make creates an instance of the Direcotry
func Make(path string) (Directory, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return Directory{}, err
	}

	stat, err := os.Stat(absPath)
	dir := Directory{
		Path: absPath,
	}

	if dir.Exists() {
		if !stat.IsDir() {
			return dir, errors.New("path is not a directory")
		}
	}

	return dir, nil
}

// Exists returns true if the directory exists and false if not.
func (dir Directory) Exists() bool {
	_, err := os.Stat(dir.Path)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

// IsProject returns true/false depending on whether the directory
// is a project directory.
func (dir Directory) IsProject() bool {
	filePath := path.Join(dir.Path, "project.yml")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

// IsTemplate returns true/false depending on whether the directory
// is a template directory.
func (dir Directory) IsTemplate() bool {
	filePath := path.Join(dir.Path, "template.yml")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

// ReadFile reads the file at the path specified, returning the bytes
// or an error.
func (dir Directory) ReadFile(file string) ([]byte, error) {
	absPath, _ := filepath.Abs(file)

	return ioutil.ReadFile(absPath)
}

// ReadTemplateManifest r
func (dir Directory) ReadTemplateManifest() ([]byte, error) {
	if !dir.IsTemplate() {
		return []byte{}, errors.New("not a template directory")
	}

	return dir.ReadFile(dir.TemplatePath())
}

// ProjectPath returns the expected path of the project.yml file
// or an empty string if the project.yml file was not found.
func (dir Directory) ProjectPath() string {
	if dir.IsProject() {
		return path.Join(dir.Path, "project.yml")
	}

	return ""
}

// TemplatePath returns the expected path of the project.yml file
// or an empty string if the template.yml file was not found.
func (dir Directory) TemplatePath() string {
	if dir.IsProject() {
		return path.Join(dir.Path, "template.yml")
	}

	return ""
}

// IsFresh checks if the Path has been set on the Directory instance.
func (dir Directory) IsFresh() bool {
	return dir.Path == ""
}
