package action

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/modelflux/cli/pkg/model"
)

type Action[ActionInput string | []string] interface {
	Run(input ActionInput, m model.Model) (ActionInput, error)
}

const LOOP_LIMIT = 100

type CustomAction struct {
	Prompt string
}

func (a *CustomAction) Run(input string, m model.Model) (string, error) {
	systemPropmt := "SYSTEM: " + a.Prompt + "\n"
	userPrompt := "USER: " + input + "\n"
	return m.Generate(systemPropmt + userPrompt)
}

type SummarizeAction struct {
	WordCount int64
}

func (a *SummarizeAction) Run(input string, m model.Model) (string, error) {
	systemPrompt := "SYSTEM: Summarize the following document concisely," +
		" capturing only the most essential points, key takeaways," +
		" and critical insights. Exclude unnecessary details, examples," +
		" and repetitions while maintaining clarity and completeness. Keep the summary under " +
		strconv.Itoa(int(a.WordCount)) + " words."
	return m.Generate(systemPrompt + "\n" + input)
}

type CategorizeAction struct {
	Categories []string
}

func (a *CategorizeAction) Run(input []string, m model.Model) ([]string, error) {
	systemPrompt := "SYSTEM: You help classify a list of items into different types." +
		" You will respond with the category ONLY and NOTHING else. Here is a list of" +
		" categories to classify: " + strings.Join(a.Categories, ", ")

	categorizedItems := make([]string, len(input))

	for i, item := range input {
		userPrompt := "USER: " + item + "\n"
		for j := range LOOP_LIMIT {
			category, err := m.Generate(systemPrompt + "\n" + userPrompt)
			if err != nil {
				return nil, err
			}
			if slices.Contains(a.Categories, category) {
				categorizedItems[i] = category
				break
			}
			if j == LOOP_LIMIT-1 {
				return nil, fmt.Errorf("Loop limit exceeded")
			}
		}
	}

	return categorizedItems, nil
}
