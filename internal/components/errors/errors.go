package errors

import (
	"context"
	"encoding/json"

	errorsApi "github.com/somatech1/mikros/apis/errors"
	loggerApi "github.com/somatech1/mikros/apis/logger"
	"github.com/somatech1/mikros/components/logger"
	"google.golang.org/grpc/status"
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
	Kind        ErrorKind
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
		err.SublevelError = options.Error.Error()
	}

	return &ServiceError{
		err:    err,
		logger: options.Logger,
	}
}

func FromGRPCStatus(st *status.Status) error {
	msg := st.Message()

	var retErr Error
	err := json.Unmarshal([]byte(msg), &retErr)
	if err != nil {
		retErr = Error{
			Message: msg,
			Kind:    KindInternal,
		}
	}

	return &retErr
}

func (s *ServiceError) WithCode(code int32) errorsApi.Error {
	s.err.Code = code
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
		if s.err.SublevelError != "" {
			logFields = append(logFields, logger.String("error.message", s.err.SublevelError))
		}

		s.logger(ctx, s.err.Message, append(logFields, s.attributes...)...)
	}

	// And give back the proper error for the API
	return s.err
}

// withKind wraps an ErrorKind into a structured log Attribute.
func withKind(kind ErrorKind) loggerApi.Attribute {
	return logger.String("error.kind", string(kind))
}

// Error is the framework error type that a service handler should return to
// keep a standard error between services.
type Error struct {
	Code          int32     `json:"code"`
	ServiceName   string    `json:"service_name,omitempty"`
	Message       string    `json:"message,omitempty"`
	Destination   string    `json:"destination,omitempty"`
	Kind          ErrorKind `json:"kind"`
	SublevelError string    `json:"details,omitempty"`

	hideDetails bool
}

func (e *Error) Error() string {
	return e.String()
}

func (e *Error) String() string {
	out := Error{
		Code:        e.Code,
		Destination: e.Destination,
		Kind:        e.Kind,
		Message:     e.Message,
	}

	// The framework can be initialized disabling error message details at the
	// output to avoid showing internal information.
	if !e.hideDetails {
		out.SublevelError = e.SublevelError
		out.ServiceName = e.ServiceName
		out.Destination = e.Destination
	}

	b, _ := json.Marshal(out)
	return string(b)
}
