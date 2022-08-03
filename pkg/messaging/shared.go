package messaging

import "fmt"

type Config struct {
	User     string
	Password string
	Host     string
	Port     int
}

var config *Config

func generateDSN() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", config.User, config.Password, config.Host, config.Port)
}

func Setup(host string, port int, user string, password string) {
	config = &Config{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
	}
}
