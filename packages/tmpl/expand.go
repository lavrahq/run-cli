package tmpl

import (
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/lavrahq/cli/packages/prompt"
	"github.com/lavrahq/cli/packages/when"
	"github.com/lavrahq/cli/util"
	"github.com/lavrahq/cli/util/cmdutil"
	"github.com/otiai10/copy"
)

// WhenEnvironment provides the available fields for `when`
// evaluations.
type WhenEnvironment struct {
	Answers  prompt.AnswerMap
	Template TemplateManifest
	Vars     map[string]interface{}
	Env      map[string]string
}

func copyExpansion(temp Template, c Copy) {
	var into = ""

	if c.Into != "" {
		into = c.Into
	}

	spin := util.Spin(fmt.Sprintf(" + From /%s to /%s", c.From, c.Into))
	err := copy.Copy(
		path.Join(temp.TemplateDirectory.Path, "template", c.From),
		path.Join(temp.Directory.Path, into),
	)

	if err != nil {
		spin.Failed(err)
	}

	spin.Done()
}

func fillExpansion(temp Template, fill Fill, env WhenEnvironment) {
	filePath := path.Join(temp.Directory.Path, fill.File)

	spin := util.Spin(fmt.Sprintf(" + Filling /%s", fill.File))
	defer spin.Done()

	tmpl, err := template.ParseFiles(filePath)
	cmdutil.CheckCommandError(err, fmt.Sprintf("parse template, %s", fill.File))

	file, err := os.Create(filePath)
	cmdutil.CheckCommandError(err, fmt.Sprintf("init template, %s", fill.File))

	err = tmpl.Execute(file, env)
	cmdutil.CheckCommandError(err, fmt.Sprintf("execute template, %s", fill.File))

	file.Close()
}

// Copy expands the template into the template's directory.
func (temp Template) Copy() {
	copySpinner := util.Spin("Copying Files")
	copySpinner.Done()

	for _, c := range temp.Manifest.Copy {
		if when.ImplicitlyTrue(c.When) {
			copyExpansion(temp, c)

			continue
		}

		env := WhenEnvironment{
			Answers:  prompt.Answers[temp.Manifest.Name],
			Template: temp.Manifest,
			Env:      util.GetEnvMap(),
		}

		if when.True(c.When, env) {
			copyExpansion(temp, c)
		}
	}
}

// Fill fills the templates specified wtihin the template.
func (temp Template) Fill() {
	fillSpinner := util.Spin("Running Templates")
	fillSpinner.Done()

	env := WhenEnvironment{
		Answers:  prompt.Answers[temp.Manifest.Name],
		Template: temp.Manifest,
		Env:      util.GetEnvMap(),
	}

	for _, f := range temp.Manifest.Fill {
		env.Vars = f.Vars

		if when.ImplicitlyTrue(f.When) {
			fillExpansion(temp, f, env)

			continue
		}

		if when.True(f.When, env) {
			fillExpansion(temp, f, env)
		}
	}
}

// Prompt runs the manifest Prompt.
func (temp Template) Prompt() prompt.AnswerMap {
	return temp.Manifest.Prompt.Ask()
}
