package errors

import (
	"context"
	"encoding/json"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errorsApi "github.com/somatech1/mikros/apis/errors"
	loggerApi "github.com/somatech1/mikros/apis/logger"
	"github.com/somatech1/mikros/components/logger"
	"github.com/somatech1/mikros/components/service"
)

// ServiceError is a structure that holds internal error details to improve
// error log description for the end-user, and it implements the errorApi.Error
// interface.
type ServiceError struct {
	err        *Error
	attributes []loggerApi.Attribute
	logger     func(ctx context.Context, msg string, attrs ...loggerApi.Attribute)
}

type serviceErrorOptions struct {
	HideDetails bool
	Code        int32
	Kind        errorsApi.Kind
	ServiceName string
	Message     string
	Destination string
	Logger      func(ctx context.Context, msg string, attrs ...loggerApi.Attribute)
	Error       error
}

func newServiceError(options *serviceErrorOptions) *ServiceError {
	err := &Error{
		hideDetails: options.HideDetails,
		Code:        options.Code,
		ServiceName: options.ServiceName,
		Message:     options.Message,
		Destination: options.Destination,
		Kind:        options.Kind,
	}

	if options.Error != nil {
		err.SubLevelError = options.Error.Error()
	}

	return &ServiceError{
		err:    err,
		logger: options.Logger,
	}
}

func FromGRPCStatus(st *status.Status, from, to service.Name) error {
	var (
		msg    = st.Message()
		retErr Error
	)

	if err := json.Unmarshal([]byte(msg), &retErr); err != nil {
		return newServiceError(&serviceErrorOptions{
			Destination: to.String(),
			Kind:        errorsApi.KindInternal,
			ServiceName: from.String(),
			Message:     "got an internal error",
			Error:       errors.New(msg),
		}).Submit(context.TODO())
	}

	// If we're dealing with a non mikros error, change it to an Internal
	// one so services can properly handle them.
	if st.Code() != codes.Unknown {
		retErr.Kind = errorsApi.KindInternal
		retErr.SubLevelError = msg
	}

	return &retErr
}

func (s *ServiceError) WithCode(code errorsApi.Code) errorsApi.Error {
	s.err.Code = code.ErrorCode()
	return s
}

func (s *ServiceError) WithAttributes(attrs ...loggerApi.Attribute) errorsApi.Error {
	s.attributes = attrs
	return s
}

func (s *ServiceError) Submit(ctx context.Context) error {
	// Display the error message onto the output
	if s.logger != nil {
		logFields := []loggerApi.Attribute{withKind(s.err.Kind)}
		if s.err.SubLevelError != "" {
			logFields = append(logFields, logger.String("error.message", s.err.SubLevelError))
		}

		s.logger(ctx, s.err.Message, append(logFields, s.attributes...)...)
	}

	// And give back the proper error for the API
	return s.err
}

func (s *ServiceError) Kind() errorsApi.Kind {
	return s.err.Kind
}

// withKind wraps a Kind into a structured log Attribute.
func withKind(kind errorsApi.Kind) loggerApi.Attribute {
	return logger.String("error.kind", string(kind))
}

// Error is the framework error type that a service handler should return to
// keep a standard error between services.
type Error struct {
	Code          int32          `json:"code"`
	ServiceName   string         `json:"service_name,omitempty"`
	Message       string         `json:"message,omitempty"`
	Destination   string         `json:"destination,omitempty"`
	Kind          errorsApi.Kind `json:"kind"`
	SubLevelError string         `json:"details,omitempty"`

	hideDetails bool
}

func (e *Error) Error() string {
	return e.String()
}

func (e *Error) String() string {
	out := Error{
		Code:    e.Code,
		Kind:    e.Kind,
		Message: e.Message,
	}

	// The framework can be initialized disabling error message details at the
	// output to avoid showing internal information.
	if !e.hideDetails {
		out.Destination = e.Destination
		out.SubLevelError = e.SubLevelError
		out.ServiceName = e.ServiceName
	}

	b, _ := json.Marshal(out)
	return string(b)
}
