package providers

import (
	"reflect"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
)

type TypeThirdPartyProviderGoogleConfig struct {
	ClientID              string
	ClientSecret          string
	Scope                 []string
	AuthorisationRedirect *struct {
		Params map[string]interface{}
	}
}

const googleID = "google"

func Google(config TypeThirdPartyProviderGoogleConfig) models.TypeProvider {
	return models.TypeProvider{
		ID: googleID,
		Get: func(redirectURI, authCodeFromRequest *string) models.TypeProviderGetResponse {
			accessTokenAPIURL := "https://accounts.google.com/o/oauth2/token"
			accessTokenAPIParams := map[string]string{
				"client_id":     config.ClientID,
				"client_secret": config.ClientSecret,
				"grant_type":    "authorization_code",
			}
			if authCodeFromRequest != nil {
				accessTokenAPIParams["code"] = *authCodeFromRequest
			}
			if redirectURI != nil {
				accessTokenAPIParams["redirect_uri"] = *redirectURI
			}

			authorisationRedirectURL := "https://accounts.google.com/o/oauth2/v2/auth"
			scopes := []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"}
			if config.Scope != nil {
				scopes = append(scopes, config.Scope...)
			}

			var additionalParams map[string]interface{} = nil
			if config.AuthorisationRedirect != nil && config.AuthorisationRedirect.Params != nil {
				additionalParams = config.AuthorisationRedirect.Params
			}

			authorizationRedirectParams := map[string]string{
				"scope":                  strings.Join(scopes, " "),
				"access_type":            "offline",
				"include_granted_scopes": "true",
				"response_type":          "code",
				"client_id":              config.ClientID,
			}
			for key, value := range additionalParams {
				if reflect.ValueOf(value).Kind() == reflect.String {
					authorizationRedirectParams[key] = value.(string)
				}
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
				GetProfileInfo: func(authCodeResponse interface{}) models.UserInfo {
					// accessTokenAPIResponse := authCodeResponse.(googleGetProfileInfoInput)
					// accessToken := accessTokenAPIResponse.AccessToken
					// authHeader := "Bearer " + accessToken
					return models.UserInfo{}
				},
			}
		},
	}
}

type googleGetProfileInfoInput struct {
	AccessToken  string
	ExpiresIn    int
	TokenType    string
	Scope        string
	RefreshToken string
}
