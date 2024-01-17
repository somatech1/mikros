package definition

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func validateVersion(_ context.Context, fl validator.FieldLevel) bool {
	return ValidateVersion(fl.Field().String())
}

// validateCollectorName checks if a name corresponds to a valid tracing
// collector name, i.e., low letters or digits in snake case format.
func validateCollectorName(_ context.Context, fl validator.FieldLevel) bool {
	if name := fl.Field().String(); name != "" {
		for _, c := range name {
			isChar := unicode.IsLetter(c) && unicode.IsLower(c)

			if !isChar && !unicode.IsNumber(c) && c != '_' {
				return false
			}
		}

		return true
	}

	// We don't accept empty names.
	return false

}

// ValidateVersion is a helper function to validate the version format used by
// services.
func ValidateVersion(input string) bool {
	return regexp.MustCompile("^v[0-9]{1,2}(|[.][0-9]{1,2})(|[.][0-9]{1,2})$").MatchString(input)
}

// validateServiceType validates if valid service type was used inside the
// settings file. It also supports the notation 'type:port', where one can
// set a custom server port for the specific service type.
func validateServiceType(ctx context.Context, fl validator.FieldLevel) bool {
	if serviceType := fl.Field().String(); serviceType != "" {
		supportedTypes, ok := ctx.Value(serviceTypeCtx{}).([]string)
		if !ok {
			return false
		}

		if strings.Contains(serviceType, ":") {
			parts := strings.Split(serviceType, ":")
			if len(parts) > 1 {
				// The server port was defined and, we must validate it.
				if !validatePort(parts[1]) {
					return false
				}
			}

			serviceType = parts[0]
		}

		for _, t := range supportedTypes {
			if serviceType == t {
				return true
			}
		}
	}

	return false
}

func validatePort(port string) bool {
	_, err := strconv.ParseInt(port, 10, 32)
	return err == nil
}
