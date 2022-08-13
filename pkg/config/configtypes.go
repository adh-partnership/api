package config

type Config struct {
	Server   ConfigServer   `yaml:"server"`
	Database ConfigDatabase `yaml:"database"`
	Discord  ConfigDiscord  `yaml:"discord"`
	Email    ConfigEmail    `yaml:"email"`
	RabbitMQ ConfigRabbitMQ `yaml:"rabbitmq"`
	Redis    ConfigRedis    `yaml:"redis"`
	Session  ConfigSession  `yaml:"session"`
	OAuth    ConfigOAuth    `yaml:"oauth"`
	VATUSA   ConfigVATUSA   `yaml:"vatusa"`
}

type ConfigServer struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type ConfigDiscord struct {
	Webhooks map[string]string `yaml:"webhooks"`
}

type ConfigDatabase struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Database    string `yaml:"database"`
	Automigrate bool   `yaml:"automigrate"`
}

type ConfigEmail struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}

type ConfigRabbitMQ struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type ConfigRedis struct {
	Password      string   `yaml:"password"`
	Database      int      `yaml:"database"`
	Address       string   `yaml:"address"`
	Sentinel      bool     `yaml:"sentinel"`
	MasterName    string   `yaml:"master_name"`
	SentinelAddrs []string `yaml:"sentinel_addrs"`
}

type ConfigSession struct {
	Cookie ConfigSessionCookie `yaml:"cookie"`
}

type ConfigSessionCookie struct {
	Name   string `yaml:"name"`
	Secret string `yaml:"secret"`
	Domain string `yaml:"domain"`
	Path   string `yaml:"path"`
	MaxAge int    `yaml:"max_age"`
}

type ConfigOAuth struct {
	BaseURL      string `yaml:"base_URL"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	MyBaseURL    string `yaml:"my_base_URL"`

	Endpoints ConfigOAuthEndpoints `yaml:"endpoints"`
}

type ConfigOAuthEndpoints struct {
	Authorize string `yaml:"authorize"`
	Token     string `yaml:"token"`
	UserInfo  string `yaml:"user"`
}

type ConfigVATUSA struct {
	APIKey string `yaml:"api_key"`
}
