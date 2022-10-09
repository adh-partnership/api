package config

type Config struct {
	Server   ConfigServer   `json:"server"`
	Database ConfigDatabase `json:"database"`
	Discord  ConfigDiscord  `json:"discord"`
	Email    ConfigEmail    `json:"email"`
	Session  ConfigSession  `json:"session"`
	OAuth    ConfigOAuth    `json:"oauth"`
	VATUSA   ConfigVATUSA   `json:"vatusa"`
	Facility ConfigFacility `json:"facility"`
	Storage  ConfigStorage  `json:"storage"`
}

type ConfigServer struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type ConfigDiscord struct {
	Webhooks map[string]string `json:"webhooks"`
}

type ConfigDatabase struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	User        string `json:"user"`
	Password    string `json:"password"`
	Database    string `json:"database"`
	Automigrate bool   `json:"automigrate"`
}

type ConfigEmail struct {
	Host     string   `json:"host"`
	Port     string   `json:"port"`
	User     string   `json:"user"`
	Password string   `json:"password"`
	From     string   `json:"from"`
	CC       []string `json:"cc"`
}

type ConfigSession struct {
	Cookie ConfigSessionCookie `json:"cookie"`
}

type ConfigSessionCookie struct {
	Name   string `json:"name"`
	Secret string `json:"secret"`
	Domain string `json:"domain"`
	Path   string `json:"path"`
	MaxAge int    `json:"max_age"`
}

type ConfigOAuth struct {
	BaseURL      string `json:"base_URL"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	MyBaseURL    string `json:"my_base_URL"`

	Endpoints ConfigOAuthEndpoints `json:"endpoints"`
}

type ConfigOAuthEndpoints struct {
	Authorize string `json:"authorize"`
	Token     string `json:"token"`
	UserInfo  string `json:"user"`
}

type ConfigVATUSA struct {
	Facility string `json:"facility"`
	APIKey   string `json:"api_key"`
}

type ConfigStorage struct {
	BaseURL   string `json:"base_URL"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
	Region    string `json:"region"`
	Endpoint  string `json:"endpoint"`
}

type ConfigFacility struct {
	Feedback ConfigFacilityFeedback `json:"feedback"`
	Stats    ConfigFacilityStats    `json:"stats"`
	Visiting ConfigFacilityVisiting `json:"visiting"`
}

type ConfigFacilityStats struct {
	Enabled                     bool     `json:"enabled"`
	DiscordBroadcast            bool     `json:"discord_broadcast"`
	DiscordBroadcastWebhookName string   `json:"discord_broadcast_webhook_name"`
	Prefixes                    []string `json:"prefixes"`
}

type ConfigFacilityVisiting struct {
	DiscordWebhookName string `json:"discord_webhook_name"`
	MinRating          string `json:"min_rating"`
}

type ConfigFacilityFeedback struct {
	DiscordWebhookName string `json:"discord_webhook_name"`
}
