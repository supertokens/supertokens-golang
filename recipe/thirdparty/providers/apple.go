package providers

import (
	"reflect"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
)

type TypeThirdPartyProviderAppleConfig struct {
	ClientID              string
	ClientSecret          ClientSecret
	Scope                 []string
	AuthorisationRedirect *struct {
		Params map[string]interface{}
	}
}

type ClientSecret struct {
	KeyId      string
	PrivateKey string
	TeamId     string
}

const appleID = "apple"

func Apple(config TypeThirdPartyProviderAppleConfig) models.TypeProvider {
	return models.TypeProvider{
		ID: appleID,
		Get: func(redirectURI, authCodeFromRequest *string) models.TypeProviderGetResponse {
			accessTokenAPIURL := "https://appleid.apple.com/auth/token"
			clientSecret, _ := getClientSecret(config.ClientID, config.ClientSecret.KeyId, config.ClientSecret.TeamId, config.ClientSecret.PrivateKey)
			accessTokenAPIParams := map[string]string{
				"client_id":     config.ClientID,
				"client_secret": clientSecret,
				"grant_type":    "authorization_code",
			}
			if authCodeFromRequest != nil {
				accessTokenAPIParams["code"] = *authCodeFromRequest
			}
			if redirectURI != nil {
				accessTokenAPIParams["redirect_uri"] = *redirectURI
			}

			authorisationRedirectURL := "https://appleid.apple.com/auth/authorize"
			scopes := []string{"name", "email"}
			if config.Scope != nil {
				scopes = append(scopes, config.Scope...)
			}

			var additionalParams map[string]interface{} = nil
			if config.AuthorisationRedirect != nil && config.AuthorisationRedirect.Params != nil {
				additionalParams = config.AuthorisationRedirect.Params
			}

			authorizationRedirectParams := map[string]string{
				"scope":         strings.Join(scopes, " "),
				"response_mode": "form_post",
				"response_type": "code",
				"client_id":     config.ClientID,
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
				GetProfileInfo: func(authCodeResponse interface{}) (models.UserInfo, error){
					return models.UserInfo{}, nil
				},
			}
		},
	}
}

// TODO:
func getClientSecret(clientId, keyId, teamId, privateKey string) (string, error) {
	return "", nil
}

type appleGetProfileInfoInput struct {
	AccessToken  string
	ExpiresIn    int
	TokenType    string
	RefreshToken string
	IDToken      string
}
