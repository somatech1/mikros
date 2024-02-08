package options

import (
	"github.com/somatech1/mikros/components/definition"
)

type ScriptServiceOptions struct{}

func (s *ScriptServiceOptions) Kind() definition.ServiceType {
	return definition.ServiceType_Script
}
