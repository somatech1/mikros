package options

// Internal feature names
const (
	prefix = "mikros_framework-"

	// Internal features

	HttpFeatureName = prefix + "http"

	// These HTTP features plugins don't exist here, but to be supported by
	// internal services, they must have these names.

	HttpCorsFeatureName        = prefix + "http_cors"
	HttpAuthFeatureName        = prefix + "http_auth"
	TracingFeatureName         = prefix + "tracing"
	TrackerFeatureName         = prefix + "tracker"
	LoggerExtractorFeatureName = prefix + "logger_extractor"
)
