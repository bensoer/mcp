package models

type FrontMatter struct {
	URI              string   `yaml:"uri"`
	Name             string   `yaml:"name"`
	Description      string   `yaml:"description"`
	Languages        []string `yaml:"languages"`
	FileTypes        []string `yaml:"file_types"`
	Priority         string   `yaml:"priority"`
	RelatedResources []string `yaml:"related_resources"`
}
