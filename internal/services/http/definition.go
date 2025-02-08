package http

import (
	"encoding/json"

	"github.com/creasty/defaults"

	"github.com/somatech1/mikros/components/definition"
)

type Definitions struct {
	DisableAuth          bool `toml:"disable_auth,omitempty" default:"false"`
	DisablePanicRecovery bool `toml:"disable_panic_recovery,omitempty" default:"false"`
	HideErrorDetails     bool `toml:"hide_error_details,omitempty"`
}

func newDefinitions(definitions *definition.Definitions) *Definitions {
	if currentDefs, ok := definitions.LoadService(definition.ServiceType_HTTP); ok {
		if b, err := json.Marshal(currentDefs); err == nil {
			var serviceDefs Definitions
			if err := json.Unmarshal(b, &serviceDefs); err != nil {
				return &serviceDefs
			}
		}
	}

	// Use the default values
	defs := &Definitions{}
	_ = defaults.Set(defs)

	return defs
}
