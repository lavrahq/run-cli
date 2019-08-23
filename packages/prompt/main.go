package prompt

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lavrahq/cli/packages/when"
	"github.com/lavrahq/cli/util"
	"github.com/lavrahq/cli/util/cmdutil"
	"github.com/logrusorgru/aurora"
)

// GlobalAnswers is all prompt answers.
type GlobalAnswers map[string]AnswerMap

// AnswerMap is the answers map type
type AnswerMap map[string]interface{}

// Answers are the global stored answers
var Answers = make(GlobalAnswers)

// Prompt holds Prompt configuration.
type Prompt struct {
	Name      string `yaml:"name"`
	Answers   AnswerMap
	Questions []Question `yaml:"questions"`
}

// Answer holds the Prompt Answer configuration.
type Answer struct {
	prompt   Prompt
	question Question
	name     string
	value    string
}

// WriteAnswer writes the answers to the global AnswerMap var.
func (answer *Answer) WriteAnswer(name string, value interface{}) error {
	if Answers[answer.prompt.Name] == nil {
		Answers[answer.prompt.Name] = make(AnswerMap)
	}

	Answers[answer.prompt.Name][answer.name] = answer.question.Transformer()(value)
	Answers[answer.prompt.Name]["Raw"+answer.name] = value

	return nil
}

// Make creates a Prompt and returns it.
func Make(name string, questions []Question) Prompt {
	return Prompt{
		Name:      name,
		Questions: questions,
	}
}

// WhenEnvironment is the object passed into the When environment.
type WhenEnvironment struct {
	Answers AnswerMap
	Env     map[string]string
}

// Ask initializes the survey prompt, asking the questions provided.
func (p Prompt) Ask() AnswerMap {
	fmt.Println()
	fmt.Printf(" %s \n\n", aurora.Green(fmt.Sprintf("%s questions:", p.Name)))
	for _, e := range p.Questions {
		if when.ImplicitlyTrue(e.When) {
			err := survey.AskOne(e.Prompt(), &Answer{prompt: p, name: e.Name, question: e}, survey.WithValidator(e.CheckValid))
			cmdutil.CheckCommandError(err, fmt.Sprintf("asking question, %s", e.Name))

			continue
		}

		env := WhenEnvironment{
			Answers: Answers[p.Name],
			Env:     util.GetEnvMap(),
		}

		if when.True(e.When, env) {
			err := survey.AskOne(e.Prompt(), &Answer{prompt: p, name: e.Name, question: e}, survey.WithValidator(e.CheckValid))
			cmdutil.CheckCommandError(err, fmt.Sprintf("asking question, %s", e.Name))
		}
	}

	fmt.Println()

	return Answers[p.Name]
}
