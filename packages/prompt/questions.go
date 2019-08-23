package prompt

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/gosimple/slug"
	"github.com/lavrahq/cli/util/cmdutil"
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
	Name      string             `yaml:"name"`
	Type      string             `yaml:"type"`
	Options   QuestionOptions    `yaml:"prompt"`
	Validate  QuestionValidation `yaml:"validate"`
	Transform string             `yaml:"transform"`
	When      string             `yaml:"when"`
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

// IsValidTransformerType checks that the given promptType is valid.
func IsValidTransformerType(transformerType string) bool {
	switch transformerType {
	case
		"Title",
		"ToLower",
		"Slug":
		return true
	}

	return false
}

// IsValidValidatorType checks that the given promptType is valid.
func IsValidValidatorType(validatorType string) bool {
	switch validatorType {
	case
		"Required",
		"MinLength",
		"MaxLength":
		return true
	}

	return false
}

// AsSurveyQuestion coerces the question into a survey.Question type.
func (question Question) AsSurveyQuestion() survey.Question {
	return survey.Question{
		Name:      question.Name,
		Prompt:    question.Prompt(),
		Validate:  question.CheckValid,
		Transform: question.Transformer(),
	}
}

func toSlug(ans interface{}) interface{} {
	return slug.Make(ans.(string))
}

// Transformer returns the survey.Transformer for the specific Question.
func (question Question) Transformer() survey.Transformer {
	if !IsValidPromptType(question.Type) {
		cmdutil.ExitWithMessage(fmt.Sprintf("An invaid Transformer was specified for the `%s` Question.\n", question.Name))
	}

	switch question.Transform {
	case "Title":
		return survey.Title
	case "ToLower":
		return survey.ToLower
	case "Slug":
		return toSlug
	}

	return func(ans interface{}) interface{} {
		return ans
	}
}

// CheckValid checks if the question's answer is valid according to the
// specifications for the specified Question.
func (question Question) CheckValid(ans interface{}) error {
	// since we are validating an Input, the assertion will always succeed
	if question.Validate.Required {
		if str, ok := ans.(string); !ok || len(str) == 0 {
			return errors.New("this response is required")
		}
	}

	if question.Validate.MinLength > 1 {
		if str, ok := ans.(string); !ok || len(str) < question.Validate.MinLength {
			return fmt.Errorf("this response must be %d or more characters", question.Validate.MinLength)
		}
	}

	if question.Validate.MaxLength > 0 {
		if str, ok := ans.(string); !ok || len(str) > question.Validate.MaxLength {
			return fmt.Errorf("this response must have %d or less characters", question.Validate.MaxLength)
		}
	}

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
		var theDefault = true

		if question.Options.Default == "" {
			theDefault, _ = strconv.ParseBool(question.Options.Default)
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
