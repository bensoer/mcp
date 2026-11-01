package models

// RelatedResource is another document to read only when When matches the current step.
type RelatedResource struct {
	URI  string `yaml:"uri"`
	When string `yaml:"when"`
}

type FrontMatter struct {
	URI              string            `yaml:"uri"`
	Name             string            `yaml:"name"`
	Description      string            `yaml:"description"`
	Languages        []string          `yaml:"languages"`
	FileTypes        []string          `yaml:"file_types"`
	Priority         string            `yaml:"priority"`
	RelatedResources []RelatedResource `yaml:"related_resources"`
}
