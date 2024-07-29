package bead

type Bead struct {
	Name    string            `yaml:"name"`
	Enabled *bool             `yaml:"enabled"`
	Fields  map[string]string `yaml:"fields"`
}
