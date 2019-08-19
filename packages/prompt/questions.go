package prompt

import (
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey"
)

// QuestionValidation provides validation options for Survey Question
// instances.
type QuestionValidation struct {
	Required  bool `yaml:"required"`
	MinLength int  `yaml:"minLength"`
	MaxLength int  `yaml:"maxLength"`
}

// QuestionOptions provides option storage for Survey Question
// instances.
type QuestionOptions struct {
	Message       string   `yaml:"message"`
	Default       string   `yaml:"default"`
	Help          string   `yaml:"help"`
	Options       []string `yaml:"options"`
	PageSize      int      `yaml:"pageSize"`
	VimMode       bool     `yaml:"vimMode"`
	Editor        string   `yaml:"editor"`
	HideDefault   bool     `yaml:"hideDefault"`
	AppendDefault bool     `yaml:"appendDefault"`
	FileName      string   `yaml:"fileName"`
}

// Question holds the Survey question configs.
type Question struct {
	Name     string             `yaml:"name"`
	Type     string             `yaml:"type"`
	Options  QuestionOptions    `yaml:"prompt"`
	Validate QuestionValidation `yaml:"validate"`
}

// IsValidPromptType checks that the given promptType is valid.
func IsValidPromptType(promptType string) bool {
	switch promptType {
	case
		"Input",
		"Multiline",
		"Password",
		"Confirm",
		"Select",
		"MultiSelect",
		"Editor":
		return true
	}

	return false
}

// AsSurveyQuestion coerces the question into a survey.Question type.
func (question Question) AsSurveyQuestion() *survey.Question {
	return &survey.Question{
		Name:      question.Name,
		Prompt:    question.Prompt(),
		Validate:  question.Validator(),
		Transform: question.Transformer(),
	}
}

// Transformer returns the survey.Transformer for the specific Question.
func (question Question) Transformer() survey.Transformer {
	return nil
}

// Validator returns the survey.Validator for the specific Question.
func (question Question) Validator() survey.Validator {
	return nil
}

// Prompt constructs a new survey.Prompt instance from the Question.
func (question Question) Prompt() survey.Prompt {
	if !IsValidPromptType(question.Type) {
		fmt.Printf("An invaid Type was specified for the `%s` Question.\n", question.Name)
	}

	switch question.Type {
	case "Input":
		return &survey.Input{
			Message: question.Options.Message,
			Help:    question.Options.Help,
			Default: question.Options.Default,
		}
	case "Multiline":
		return &survey.Multiline{
			Message: question.Options.Message,
			Help:    question.Options.Help,
			Default: question.Options.Default,
		}
	case "Password":
		return &survey.Password{
			Message: question.Options.Message,
			Help:    question.Options.Help,
		}
	case "Confirm":
		theDefault, err := strconv.ParseBool(question.Options.Default)

		if err != nil {
			panic(err)
		}

		return &survey.Confirm{
			Message: question.Options.Message,
			Help:    question.Options.Help,
			Default: theDefault,
		}
	case "Select":
		return &survey.Select{
			Message:  question.Options.Message,
			Help:     question.Options.Help,
			Default:  question.Options.Default,
			Options:  question.Options.Options,
			PageSize: question.Options.PageSize,
			VimMode:  question.Options.VimMode,
		}
	case "MultiSelect":
		return &survey.MultiSelect{
			Message:  question.Options.Message,
			Help:     question.Options.Help,
			Default:  []string{question.Options.Default},
			Options:  question.Options.Options,
			PageSize: question.Options.PageSize,
			VimMode:  question.Options.VimMode,
		}
	case "Editor":
		return &survey.Editor{
			Message:       question.Options.Message,
			Help:          question.Options.Help,
			Default:       question.Options.Default,
			Editor:        question.Options.Editor,
			HideDefault:   question.Options.HideDefault,
			AppendDefault: question.Options.AppendDefault,
			// FileName:      question.Options.FileName,
		}
	}

	return nil
}
