package models

/*

description: Rules on branch creation, branch naming conventions, pushing and PR practices
languages:
    - all
file_types:
    - "*.*"
priority: required
related_resources:
    - standards://git/commit-messages
    - standards://git/commit-staging
*/

/*

		URI:         "standards://python/logging",
		Name:        "Python Logging Standards",
		Title:       "Python Logging Standards",
		MIMEType:    "text/markdown",
		Description: "Python Logging Standards, Formats, Syntax Expections and Examples for AI Agents",
	},
*/

type FrontMatter struct {
	URI              string   `yaml:"uri"`
	Name             string   `yaml:"name"`
	Description      string   `yaml:"description"`
	Languages        []string `yaml:"languages"`
	FileTypes        []string `yaml:"file_types"`
	Priority         string   `yaml:"priority"`
	RelatedResources []string `yaml:"related_resources"`
}
