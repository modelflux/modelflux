# Model Actions

An abstracted layer above the models themselves, `model actions` allow the models to take specific
actions and generate specific outputs with some prompt engineering and logic.

## Types of actions

- CustomAction - A custom action node that allows the user to specify a system prompt and examples
  that can be customized to the user's preferences
  - Input type: string
  - Output type: string
  - Options: prompt
- SummarizeAction - Summarizes the input in x number of words
  - Input type: string
  - Output type: string
  - Options: number of words
- CategorizeAction - Categorizes a list of something based on a list of categories
  - Input type: []string
  - Output type: []string
  - Options: list of categories
