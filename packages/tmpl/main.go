package tmpl

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/lavrahq/cli/packages/fs"
	"github.com/lavrahq/cli/packages/prompt"
	"github.com/lavrahq/cli/util"
	"github.com/lavrahq/cli/util/cmdutil"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"gopkg.in/yaml.v2"
)

// Copy is an instance of each expansion configuration
// object. Expansions are used to determine the files/directories to
// expand (copy) form the template into the project.
type Copy struct {
	Into string `yaml:"into"`
	From string `yaml:"from"`
	When string `yaml:"when"`
}

// Fill is an instance of each template fill configuration
// object. Templates are used to specify templates the files that
// should be ran through the template engine for variable replacement.
type Fill struct {
	File  string                 `yaml:"file"`
	When  string                 `yaml:"when"`
	Needs []string               `yaml:"needs"`
	Vars  map[string]interface{} `yaml:"vars"`
}

// TemplateManifest is an instance of the template configuration
// file.
type TemplateManifest struct {
	Name        string        `yaml:"name"`
	Author      string        `yaml:"author"`
	Description string        `yaml:"description"`
	Prompt      prompt.Prompt `yaml:"prompt"`
	Copy        []Copy        `yaml:"copy"`
	Fill        []Fill        `yaml:"fill"`
}

// Template holds information related to the template being
// rendered.
type Template struct {
	From              string
	Directory         fs.Directory
	TemplateDirectory fs.Directory
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
func Make(expandDir fs.Directory, template string) Template {
	var templateDir fs.Directory
	templateConfig := Template{
		From:      template,
		Directory: expandDir,
	}

	if templateConfig.IsTemplateAvailableLocally() {
		templateDir, _ = fs.MakeDirectory(path.Join(GetLocalPath(), template))
	}

	safeRemote := templateConfig.GetSafeRemote()
	if !templateConfig.IsTemplateAvailableRemotely(safeRemote) {
		cmdutil.ExitWithMessage("The remote template provided is not available.")
	}

	templateDir, _ = fs.MakeDirectory(path.Join(GetCachePath(), templateConfig.GetLocalPathByRemote()))
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
func (temp Template) EnsureTemplateIsFetched() error {
	storePath := temp.TemplateDirectory.Path

	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		git.PlainClone(storePath, false, &git.CloneOptions{
			URL:               temp.GetSafeRemote(),
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})

		return nil
	}

	repo, err := git.PlainOpen(storePath)
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
	})
	if err != nil {
		if err.Error() == "already up-to-date" {
			return nil
		}

		return nil
	}

	return nil
}

// CachedPath returns the path to the locally cached template.
func (temp Template) CachedPath() string {
	return path.Join(GetCachePath(), temp.GetLocalPathByRemote(), "template.yml")
}

// LoadManifest loads the template.yml into the Manifest of the Template.
func (temp Template) LoadManifest() Template {
	progress := util.Spin("Fetching template manifest")

	bytes, err := ioutil.ReadFile(temp.CachedPath())
	cmdutil.CheckCommandError(err, "loading manifest")

	err = yaml.Unmarshal(bytes, &temp.Manifest)
	cmdutil.CheckCommandError(err, "converting manifest")

	temp.Manifest.Prompt.Name = temp.Manifest.Name

	progress.Done()

	return temp
}
