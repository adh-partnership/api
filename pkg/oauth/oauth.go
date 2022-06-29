package oauth

import (
	"golang.org/x/oauth2"
)

var OAuthConfig *oauth2.Config

func Build(clientid, clientsecret, redirecturl, authorize, token string) {
	OAuthConfig = &oauth2.Config{
		ClientID:     clientid,
		ClientSecret: clientsecret,
		RedirectURL:  redirecturl,
		Scopes:       []string{"fullname", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authorize,
			TokenURL: token,
		},
	}
}
