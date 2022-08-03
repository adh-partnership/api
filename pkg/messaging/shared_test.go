package messaging

import (
	"reflect"
	"testing"
)

func TestSetup(t *testing.T) {
	tests := []struct {
		name     string
		args     *Config
		expected *Config
	}{
		{
			name: "valid",
			args: &Config{
				User:     "user",
				Password: "password",
				Host:     "host",
				Port:     123,
			},
			expected: &Config{
				User:     "user",
				Password: "password",
				Host:     "host",
				Port:     123,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Setup(tt.args.Host, tt.args.Port, tt.args.User, tt.args.Password)
			if !reflect.DeepEqual(config, tt.expected) {
				t.Errorf("Setup() = %+v, want %+v", config, tt.expected)
			}
		})
	}
}

func TestGenerateDSN(t *testing.T) {
	tests := []struct {
		name     string
		args     *Config
		expected string
	}{
		{
			name: "valid",
			args: &Config{
				User:     "user",
				Password: "password",
				Host:     "host",
				Port:     123,
			},
			expected: "amqp://user:password@host:123/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Setup(tt.args.Host, tt.args.Port, tt.args.User, tt.args.Password)
			if got := generateDSN(); got != tt.expected {
				t.Errorf("generateDSN() = %v, want %v", got, tt.expected)
			}
		})
	}
}
