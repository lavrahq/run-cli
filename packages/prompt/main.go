package prompt

import (
	"fmt"

	"github.com/AlecAivazis/survey"
	"github.com/logrusorgru/aurora"
)

// AnswerMap is the answers map type
type AnswerMap map[string]interface{}

var answers = make(AnswerMap)

// Prompt holds Prompt configuration.
type Prompt struct {
	Name      string     `yaml:"name"`
	Questions []Question `yaml:"questions"`
}

// Answers holds the Prompt Answer configuration.
type Answer struct {
	value string
}

// WriteAnswer writes the answers to the global AnswerMap var.
func (answer *Answer) WriteAnswer(name string, value interface{}) error {
	answers[name] = value

	return nil
}

// Make creates a Prompt and returns it.
func Make(name string, questions []Question) Prompt {
	return Prompt{
		Name:      name,
		Questions: questions,
	}
}

// Ask initializes the survey prompt, asking the questions provided.
func (p Prompt) Ask() AnswerMap {
	var questions []*survey.Question

	for _, e := range p.Questions {
		questions = append(questions, e.AsSurveyQuestion())
	}

	fmt.Println()
	fmt.Printf(" %s \n\n", aurora.Green("questions:"))
	survey.Ask(questions, &Answer{})

	return answers
}
