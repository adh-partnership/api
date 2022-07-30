package logger

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

var (
	Logger = logrus.New()
	Format string
)

func NewLogger(format string) {
	Format = format
	if format == "json" {
		Logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		Logger.SetFormatter(&nested.Formatter{
			HideKeys:        true,
			TimestampFormat: "2006-01-02T15:04:05Z07:00",
			FieldsOrder:     []string{"component", "category"},
			ShowFullLevel:   true,
		})
	}
}

func IsValidFormat(format string) bool {
	return format == "json" || format == "text"
}

func IsValidLogLevel(level string) bool {
	_, err := ParseLogLevel(level)
	return err == nil
}

func ParseLogLevel(level string) (logrus.Level, error) {
	return logrus.ParseLevel(level)
}
