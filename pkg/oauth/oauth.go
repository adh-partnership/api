package oauth

import (
	"fmt"

	"golang.org/x/oauth2"

	"github.com/adh-partnership/api/pkg/config"
)

var (
	DiscordOAuthConfig *oauth2.Config
	OAuthConfig        *oauth2.Config
)

// Deprecated: use BuildWithConfig instead
func Build(clientid, clientsecret, redirectURL, authorize, token string) {
	OAuthConfig = buildv2(clientid, clientsecret, redirectURL, authorize, token, []string{"identify", "email"})
}

func buildv2(clientid, secret, redirectURL, authorize, token string, scopes []string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientid,
		ClientSecret: secret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authorize,
			TokenURL: token,
		},
	}
}

func BuildWithConfig(c *config.Config) {
	Build(
		c.OAuth.ClientID,
		c.OAuth.ClientSecret,
		fmt.Sprintf("%s%s", c.OAuth.MyBaseURL, "/v1/user/login/callback"),
		fmt.Sprintf("%s%s", c.OAuth.BaseURL, c.OAuth.Endpoints.Authorize),
		fmt.Sprintf("%s%s", c.OAuth.BaseURL, c.OAuth.Endpoints.Token),
	)

	DiscordOAuthConfig = buildv2(
		c.Discord.ClientID,
		c.Discord.ClientSecret,
		fmt.Sprintf("%s%s", c.OAuth.MyBaseURL, "/v1/user/discord/callback"),
		"https://discord.com/api/v10/oauth2/authorize",
		"https://discord.com/api/v10/oauth2/token",
		[]string{"identify"}, // We really only need their id
	)

	DiscordOAuthConfig.Endpoint.AuthStyle = oauth2.AuthStyleInParams
}
