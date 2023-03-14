package output

type Output struct {
	Secrets []OutputSecret `yaml:"secrets"`
}

type OutputSecret struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace"`
	Type      string            `yaml:"type"`
	Data      map[string]string `yaml:"data"`
}
