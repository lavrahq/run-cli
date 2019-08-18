package tmpl

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/lavrahq/cli/util/dir"
	. "github.com/lavrahq/cli/util/logs"
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

// Make initializes a Template given a dir and the template name
// or remote.
func Make(expandDir dir.Directory, template string) Template {
	var templateDir dir.Directory
	templateConfig := Template{
		From:      template,
		Directory: expandDir,
	}

	if templateConfig.IsTemplateAvailableLocally() {
		templateDir, _ = dir.Make(path.Join(GetLocalPath(), template))
	}

	safeRemote := templateConfig.GetSafeRemote()
	if !templateConfig.IsTemplateAvailableRemotely(safeRemote) {
		Log.Info("The remote template provided is not available.")
	}

	templateDir, _ = dir.Make(path.Join(GetCachePath(), templateConfig.GetLocalPathByRemote()))
	templateConfig.TemplateDirectory = templateDir

	return templateConfig
}

// CheckIfCoreRemote checks if the remote provided is a Lavra
// repository remote.
func (temp Template) CheckIfCoreRemote() bool {
	return getCountOfSlashesInRemote(temp.From) == 0
}

// CheckIfGithubRemote checks if the remoote provided is a Github
// repository remote.
func (temp Template) CheckIfGithubRemote() bool {
	return getCountOfSlashesInRemote(temp.From) == 1
}

// IsTemplateAvailableLocally checks whether or not a local template
// exists by the specific name.
func (temp Template) IsTemplateAvailableLocally() bool {
	tm := path.Join(GetLocalPath(), temp.From, "template.json")
	_, err := os.Stat(tm)

	return !os.IsNotExist(err)
}

// GetSafeRemote returns a string with the full Git URL for Lavra-owned
// projects as well Github-hosted projects, and finally the given remote
// string.
func (temp Template) GetSafeRemote() string {
	if temp.CheckIfCoreRemote() {
		s := []string{"https://github.com/lavrahq/cli-project-template-", temp.From, ".git"}

		return strings.Join(s, "")
	}

	if temp.CheckIfGithubRemote() {
		s := []string{"https://github.com/", temp.From, ".git"}

		return strings.Join(s, "")
	}

	return temp.From
}

// GetLocalPathByRemote returns the local path of the remote provided.
func (temp Template) GetLocalPathByRemote() string {
	h := md5.New()
	h.Write([]byte(temp.From))

	return hex.EncodeToString(h.Sum(nil))
}

// LoadRemoteTemplateIntoMemory loads the remote template into memory
// for
func (temp Template) LoadRemoteTemplateIntoMemory() (*git.Repository, error) {
	Log.Info("Received request to load template into memory.", zap.String("remote", temp.From))

	return git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:        temp.GetSafeRemote(),
		RemoteName: "origin",
	})
}

// IsTemplateAvailableRemotely returns a boolean stating whether the
// the remote template is available.
func (temp Template) IsTemplateAvailableRemotely(remote string) bool {
	_, err := temp.LoadRemoteTemplateIntoMemory()

	return err == nil
}

// EnsureTemplateIsFetched fetches the remote template, ensuring that the
// fetched version is the latest.
func (temp Template) EnsureTemplateIsFetched() (bool, error) {
	storePath := temp.TemplateDirectory.Path

	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		git.PlainClone(storePath, false, &git.CloneOptions{
			URL:               temp.GetSafeRemote(),
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})

		return true, nil
	}

	repo, err := git.PlainOpen(storePath)
	if err != nil {
		return false, err
	}

	w, err := repo.Worktree()
	if err != nil {
		return false, err
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
	})
	if err != nil {
		if err.Error() == "already up-to-date" {
			return true, nil
		}

		return false, nil
	}

	return true, nil
}
