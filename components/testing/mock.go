package testing

import (
	"testing"

	"go.uber.org/mock/gomock"
)

// MockController gives access to a gomock.Controller object that enables
// creating a service mock to be used inside tests. This API can be called
// without a Testing object.
func MockController(t *testing.T) *gomock.Controller {
	return gomock.NewController(t)
}

// MockAny gives access to a gomock.Matcher object which can replace any
// kind of argument. This API can be called without a Testing object.
func MockAny() gomock.Matcher {
	return gomock.Any()
}
