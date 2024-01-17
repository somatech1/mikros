package errors

import (
	"errors"
	"fmt"

	errorsApi "github.com/somatech1/mikros/apis/errors"
	loggerApi "github.com/somatech1/mikros/apis/logger"
)

// Internal error codes.
const (
	CodeInternal int32 = iota + 1
	CodeNotFound
	CodeInvalidArgument
	CodePreconditionFailed
	CodeNoPermission
	CodeRPC
)

// ErrorKind is an error representation of a mapped error.
type ErrorKind string

var (
	KindValidation   ErrorKind = "ValidationError"
	KindInternal     ErrorKind = "InternalError"
	KindNotFound     ErrorKind = "NotFoundError"
	KindPrecondition ErrorKind = "ConditionError"
	KindPermission   ErrorKind = "PermissionError"
	KindRPC          ErrorKind = "RPCError"
	KindCustom       ErrorKind = "CustomError"
)

type Factory struct {
	hideMessageDetails bool
	serviceName        string
	logger             loggerApi.Logger
}

type FactoryOptions struct {
	HideMessageDetails bool
	ServiceName        string
	Logger             loggerApi.Logger
}

// NewFactory creates a new Factory object.
func NewFactory(options FactoryOptions) *Factory {
	return &Factory{
		serviceName:        options.ServiceName,
		logger:             options.Logger,
		hideMessageDetails: options.HideMessageDetails,
	}
}

// RPC sets that the current error is related to an RPC call with another gRPC
// service (destination), However, it checks if the error is already a known
// error, so it does not override it.
func (f *Factory) RPC(err error, destination string) errorsApi.Error {
	return newServiceError(&serviceErrorOptions{
		HideDetails: f.hideMessageDetails,
		Code:        CodeRPC,
		Kind:        KindRPC,
		ServiceName: f.serviceName,
		Message:     "service RPC error",
		Destination: destination,
		Logger:      f.logger.Warn,
		Error:       err,
	})
}

// InvalidArgument sets that the current error is related to an argument that
// didn't follow validation rules.
func (f *Factory) InvalidArgument(err error) errorsApi.Error {
	return newServiceError(&serviceErrorOptions{
		HideDetails: f.hideMessageDetails,
		Code:        CodeInvalidArgument,
		Kind:        KindValidation,
		ServiceName: f.serviceName,
		Message:     "request validation failed",
		Logger:      f.logger.Warn,
		Error:       err,
	})
}

// FailedPrecondition sets that the current error is related to an internal
// condition which wasn't satisfied.
func (f *Factory) FailedPrecondition(message string) errorsApi.Error {
	return newServiceError(&serviceErrorOptions{
		HideDetails: f.hideMessageDetails,
		Code:        CodePreconditionFailed,
		Kind:        KindPrecondition,
		ServiceName: f.serviceName,
		Message:     "failed precondition",
		Logger:      f.logger.Warn,
		Error:       errors.New(message),
	})
}

// NotFound sets that the current error is related to some data not being found,
// probably in the database.
func (f *Factory) NotFound() errorsApi.Error {
	return newServiceError(&serviceErrorOptions{
		HideDetails: f.hideMessageDetails,
		Code:        CodeNotFound,
		Kind:        KindNotFound,
		ServiceName: f.serviceName,
		Message:     "not found",
		Logger:      f.logger.Warn,
	})
}

// Internal sets that the current error is related to an internal service
// error.
func (f *Factory) Internal(err error) errorsApi.Error {
	return newServiceError(&serviceErrorOptions{
		HideDetails: f.hideMessageDetails,
		Code:        CodeInternal,
		Kind:        KindInternal,
		ServiceName: f.serviceName,
		Message:     "got an internal error",
		Logger:      f.logger.Error,
		Error:       err,
	})
}

// PermissionDenied sets that the current error is related to a client trying
// to access a resource without having permission to do so.
func (f *Factory) PermissionDenied() errorsApi.Error {
	return newServiceError(&serviceErrorOptions{
		HideDetails: f.hideMessageDetails,
		Code:        CodeNoPermission,
		Kind:        KindPermission,
		ServiceName: f.serviceName,
		Message:     fmt.Sprintf("no permission to access %s", f.serviceName),
		Logger:      f.logger.Info,
	})
}

// Custom lets a service set a custom error kind for its errors. Internally, it
// will be treated as an Internal error.
func (f *Factory) Custom(msg string) errorsApi.Error {
	return newServiceError(&serviceErrorOptions{
		HideDetails: f.hideMessageDetails,
		Kind:        KindCustom,
		ServiceName: f.serviceName,
		Message:     msg,
		Logger:      f.logger.Info,
	})
}
