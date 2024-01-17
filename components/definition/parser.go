package definition

import (
	"github.com/BurntSushi/toml"
)

// Parse is responsible for loading the service definitions file (service.toml)
// into a proper Definitions structure.
func Parse(path string) (*Definitions, error) {
	defs, err := New()
	if err != nil {
		return nil, err
	}

	if _, err := toml.DecodeFile(path, &defs); err != nil {
		return nil, err
	}

	return defs, nil
}

// ParseExternalDefinitions allows loading specific service definitions from its
// file using a custom target. This provides external features (plugins) to load
// their definitions from the same file into their own structures.
func ParseExternalDefinitions(path string, defs interface{}) error {
	if _, err := toml.DecodeFile(path, defs); err != nil {
		return err
	}

	return nil
}
