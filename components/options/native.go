package options

import (
	"github.com/somatech1/mikros/components/definition"
)

type NativeServiceOptions struct{}

func (n *NativeServiceOptions) Kind() definition.ServiceType {
	return definition.ServiceType_Native
}
