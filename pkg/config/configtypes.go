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

package config

type Config struct {
	Database ConfigDatabase      `json:"database"`
	Discord  ConfigDiscord       `json:"discord"`
	Email    ConfigEmail         `json:"email"`
	Facility ConfigFacility      `json:"facility"`
	Features ConfigFeatures      `json:"features"`
	Groups   map[string][]string `json:"groups"`
	Metrics  ConfigMetrics       `json:"metrics"`
	OAuth    ConfigOAuth         `json:"oauth"`
	Server   ConfigServer        `json:"server"`
	Session  ConfigSession       `json:"session"`
	Storage  ConfigStorage       `json:"storage"`
	VATUSA   ConfigVATUSA        `json:"vatusa"`
}

type ConfigServer struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type ConfigDiscord struct {
	Webhooks     map[string]string `json:"webhooks"`
	ClientID     string            `json:"client_id"`
	ClientSecret string            `json:"client_secret"`
}

type ConfigDatabase struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	User        string `json:"user"`
	Password    string `json:"password"`
	Database    string `json:"database"`
	Automigrate bool   `json:"automigrate"`
	CACert      string `json:"ca_cert"`
}

type ConfigEmail struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	User        string `json:"user"`
	Password    string `json:"password"`
	From        string `json:"from"`
	TemplateDir string `json:"template_dir"`
}

type ConfigFeatures struct {
	StaffingRequest  bool `json:"staffing_request"`
	ControllerOnline bool `json:"controller_online"`
}

type ConfigSession struct {
	Cookie ConfigSessionCookie `json:"cookie"`
}

type ConfigSessionCookie struct {
	Name     string `json:"name"`
	Secret   string `json:"secret"`
	Domain   string `json:"domain"`
	Path     string `json:"path"`
	MaxAge   int    `json:"max_age"`
	Secure   bool   `json:"secure"`
	SameSite string `json:"same_site"`
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
	TestMode bool   `json:"test_mode"`
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
	Activity         ConfigFacilityActivity `json:"activity"`
	Stats            ConfigFacilityStats    `json:"stats"`
	Visiting         ConfigFacilityVisiting `json:"visiting"`
	TrainingRequests ConfigFacilityTraining `json:"training_requests"`
	FrontendURL      string                 `json:"frontend_url"`
}

type ConfigFacilityTraining struct {
	Enabled            bool                          `json:"enabled"`
	Discord            ConfigFacilityTrainingDiscord `json:"discord"`
	Positions          []string                      `json:"positions"`
	MaxRequestsPerUser int                           `json:"max_requests_per_user"`
	SendToDiscord      bool                          `json:"send_to_discord"`
}

type ConfigFacilityTrainingDiscord struct {
	TrainingStaff    string `json:"training_staff"`
	Scheduled        string `json:"scheduled"`
	ShowAllScheduled bool   `json:"show_all_scheduled"`
}

type ConfigFacilityStats struct {
	Enabled  bool     `json:"enabled"`
	Prefixes []string `json:"prefixes"`
}

type ConfigFacilityVisiting struct {
	MinRating    string `json:"min_rating"`
	SendWelcome  bool   `json:"send_welcome"`
	SendRemoval  bool   `json:"send_removal"`
	SendRejected bool   `json:"send_rejected"`
}

type ConfigFacilityActivity struct {
	Inactive ConfigFacilityActivityInactive `json:"inactive"`
	Warning  ConfigFacilityActivityWarning  `json:"warning"`
}

type ConfigFacilityActivityWarning struct {
	Enabled    bool  `json:"enabled"`
	DaysBefore int   `json:"days_before"`
	Months     []int `json:"months"`
}

type ConfigFacilityActivityInactive struct {
	Enabled  bool  `json:"enabled"`
	Period   int   `json:"period"`
	MinHours int   `json:"min_hours"`
	Months   []int `json:"months"`
}

type ConfigMetrics struct {
	Enabled bool   `json:"enabled"`
	Port    int    `json:"port"`
	Path    string `json:"path"`
}
