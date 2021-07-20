package providers

import (
	"reflect"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
)

type TypeThirdPartyProviderGithubConfig struct {
	ClientID              string
	ClientSecret          string
	Scope                 []string
	AuthorisationRedirect *struct {
		Params map[string]interface{}
	}
}

const githubID = "github"

func Github(config TypeThirdPartyProviderGithubConfig) models.TypeProvider {
	return models.TypeProvider{
		ID: githubID,
		Get: func(redirectURI, authCodeFromRequest *string) models.TypeProviderGetResponse {
			accessTokenAPIURL := "https://github.com/login/oauth/access_token"
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

			authorisationRedirectURL := "https://github.com/login/oauth/authorize"
			scopes := []string{"user"}
			if config.Scope != nil {
				scopes = append(scopes, config.Scope...)
			}

			var additionalParams map[string]interface{} = nil
			if config.AuthorisationRedirect != nil && config.AuthorisationRedirect.Params != nil {
				additionalParams = config.AuthorisationRedirect.Params
			}

			authorizationRedirectParams := map[string]string{
				"scope":     strings.Join(scopes, " "),
				"client_id": config.ClientID,
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
				// TODO:
				GetProfileInfo: func(authCodeResponse interface{}) models.UserInfo {
					return models.UserInfo{}
				},
			}
		},
	}
}

type githubGetProfileInfoInput struct {
	AccessToken string
	ExpiresIn   int
	TokenType   string
}
