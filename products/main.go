package products

// ProductConfig is the configuration of a Product
type ProductConfig struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Version string `yaml:"version"`
}
