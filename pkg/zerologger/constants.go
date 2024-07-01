package zerologger

import "github.com/rs/zerolog"

const (
	// TraceLevel defines trace log level.
	TraceLevel = "trace"
	// DebugLevel defines debug log level.
	DebugLevel = "debug"
	// InfoLevel defines info log level.
	InfoLevel = "info"
	// WarnLevel defines warn log level.
	WarnLevel = "warn"
	// ErrorLevel defines error log level.
	ErrorLevel = "error"
	// FatalLevel defines fatal log level.
	FatalLevel = "fatal"
)

// ExtractZerologLevel получает уровень логгирования в совместимом с zerolog формате
func ExtractZerologLevel(level string) zerolog.Level {
	switch level {
	case TraceLevel:
		return zerolog.TraceLevel
	case DebugLevel:
		return zerolog.DebugLevel
	case InfoLevel:
		return zerolog.InfoLevel
	case WarnLevel:
		return zerolog.WarnLevel
	case ErrorLevel:
		return zerolog.ErrorLevel
	case FatalLevel:
		return zerolog.FatalLevel
	default:
		return zerolog.DebugLevel // TODO - может быть info???
	}
}
