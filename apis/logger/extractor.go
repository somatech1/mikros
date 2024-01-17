package logger

import (
	"context"
)

// Extractor is an interface that a plugin can implement to provide an API
// allowing the service extract content from its context to add them into
// log messages.
type Extractor interface {
	Extract(ctx context.Context) []Attribute
}
