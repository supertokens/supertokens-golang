package providers

import (
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
)

type TypeThirdPartyProviderFacebookConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
}

const facebookID = "facebook"

func Facebook(config TypeThirdPartyProviderFacebookConfig) models.TypeProvider {
	return models.TypeProvider{
		ID: facebookID,
		Get: func(redirectURI, authCodeFromRequest *string) models.TypeProviderGetResponse {
			accessTokenAPIURL := "https://graph.facebook.com/v9.0/oauth/access_token"
			accessTokenAPIParams := map[string]string{
				"client_id":     config.ClientID,
				"client_secret": config.ClientSecret,
			}
			if authCodeFromRequest != nil {
				accessTokenAPIParams["code"] = *authCodeFromRequest
			}
			if redirectURI != nil {
				accessTokenAPIParams["redirect_uri"] = *redirectURI
			}

			authorisationRedirectURL := "https://www.facebook.com/v9.0/dialog/oauth"
			scopes := []string{"email"}
			if config.Scope != nil {
				scopes = append(scopes, config.Scope...)
			}

			authorizationRedirectParams := map[string]string{
				"scope":         strings.Join(scopes, " "),
				"response_type": "code",
				"client_id":     config.ClientID,
			}

			return models.TypeProviderGetResponse{
				AccessTokenAPI: models.URLParams{
					URL:    accessTokenAPIURL,
					Params: accessTokenAPIParams,
				},
				AuthorisationRedirect: models.URLParams{
					URL:    authorisationRedirectURL,
					Params: authorizationRedirectParams,
				},
				// TODO:
				GetProfileInfo: func(authCodeResponse interface{}) models.UserInfo {
					return models.UserInfo{}
				},
			}
		},
	}
}

type facebookGetProfileInfoInput struct {
	AccessToken string
	ExpiresIn   int
	TokenType   string
}
