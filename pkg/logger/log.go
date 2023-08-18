/*
 * Copyright ADH Partnership
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

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
