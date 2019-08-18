package tmpl

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/lavrahq/cli/utilities/dir"
	"github.com/lavrahq/cli/utilities/logs"
	"go.uber.org/zap"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// TemplateSurveyValidator holds options available for validating
// survey prompts.
type TemplateSurveyValidator struct {
	Required bool `yaml:"required"`
}

// TemplateSurveyPrompt holds options available for various input
// prompt configs.
type TemplateSurveyPrompt struct {
	Message string `yaml:"message"`
}

// TemplateSurvey holds the template survey input.
type TemplateSurvey struct {
	Name     string                  `yaml:"name"`
	Type     string                  `yaml:"type"`
	Prompt   TemplateSurveyPrompt    `yaml:"prompt"`
	Validate TemplateSurveyValidator `yaml:"validate"`
}

// TemplateManifest is an instance of the template configuration
// file.
type TemplateManifest struct {
	Name        string           `yaml:"name"`
	Author      string           `yaml:"author"`
	Description string           `yaml:"description"`
	Survey      []TemplateSurvey `yaml:"survey"`
}

// Template holds information related to the template being
// rendered.
type Template struct {
	From              string
	Directory         dir.Directory
	TemplateDirectory dir.Directory
	Manifest          TemplateManifest
}

// getCountOfSlashesInRemote returns the number of forward
// slashes within a remote string.
func getCountOfSlashesInRemote(remote string) int {
	matchSlashes := regexp.MustCompile("/")
	slashes := matchSlashes.FindAllStringIndex(remote, -1)

	return len(slashes)
}

// CheckIfCoreRemote checks if the remote provided is a Lavra
// repository remote.
func CheckIfCoreRemote(remote string) bool {
	return getCountOfSlashesInRemote(remote) == 0
}

// CheckIfGithubRemote checks if the remoote provided is a Github
// repository remote.
func CheckIfGithubRemote(remote string) bool {
	return getCountOfSlashesInRemote(remote) == 1
}

// IsTemplateAvailableLocally checks whether or not a local template
// exists by the specific name.
func IsTemplateAvailableLocally(name string) bool {
	tm := path.Join(GetLocalPath(), name, "template.json")
	_, err := os.Stat(tm)

	return !os.IsNotExist(err)
}

// Make initializes a Template given a dir and the template name
// or remote.
func Make(expandDir dir.Directory, template string) Template {
	var templateDir dir.Directory
	templateConfig := Template{
		From:      template,
		Directory: expandDir,
	}

	if IsTemplateAvailableLocally(template) {
		templateDir, _ = dir.Make(path.Join(GetLocalPath(), template))
	}

	safeRemote := GetSafeRemote(template)
	if !IsTemplateAvailableRemotely(safeRemote) {
		logs.Log.Info("The remote template provided is not available.")
	}

	templateDir, _ = dir.Make(path.Join(GetCachePath(), GetLocalPathByRemote(safeRemote)))
	templateConfig.TemplateDirectory = templateDir

	return templateConfig
}

// GetSafeRemote returns a string with the full Git URL for Lavra-owned
// projects as well Github-hosted projects, and finally the given remote
// string.
func GetSafeRemote(remote string) string {
	if CheckIfCoreRemote(remote) {
		s := []string{"https://github.com/lavrahq/cli-project-template-", remote, ".git"}

		return strings.Join(s, "")
	}

	if CheckIfGithubRemote(remote) {
		s := []string{"https://github.com/", remote, ".git"}

		return strings.Join(s, "")
	}

	return remote
}

// GetLocalPathByRemote returns the local path of the remote provided.
func GetLocalPathByRemote(remote string) string {
	h := md5.New()
	h.Write([]byte(remote))

	return hex.EncodeToString(h.Sum(nil))
}

// LoadRemoteTemplateIntoMemory loads the remote template into memory
// for
func LoadRemoteTemplateIntoMemory(remote string) (*git.Repository, error) {
	logs.Log.Info(
		"Received request to load template into memory.",
		zap.String("remote", remote),
	)

	return git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:        GetSafeRemote(remote),
		RemoteName: "origin",
	})
}

// IsTemplateAvailableRemotely returns a boolean stating whether the
// the remote template is available.
func IsTemplateAvailableRemotely(remote string) bool {
	_, err := LoadRemoteTemplateIntoMemory(remote)

	return err == nil
}

// EnsureTemplateIsFetched fetches the remote template, ensuring that the
// fetched version is the latest.
func (tmpl Template) EnsureTemplateIsFetched() string {
	storePath := tmpl.TemplateDirectory.Path

	git.PlainClone(storePath, false, &git.CloneOptions{
		URL:               GetSafeRemote(tmpl.From),
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	return ""
}
