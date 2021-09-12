package providers

// import (
// 	"encoding/json"
// 	"strings"

// 	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
// )

// type AppleConfig struct {
// 	ClientID              string
// 	ClientSecret          ClientSecret
// 	Scope                 []string
// 	AuthorisationRedirect *struct {
// 		Params map[string]interface{}
// 	}
// }

// type ClientSecret struct {
// 	KeyId      string
// 	PrivateKey string
// 	TeamId     string
// }

// const appleID = "apple"

// func Apple(config AppleConfig) models.TypeProvider {
// 	return models.TypeProvider{
// 		ID: appleID,
// 		Get: func(redirectURI, authCodeFromRequest *string) models.TypeProviderGetResponse {
// 			accessTokenAPIURL := "https://appleid.apple.com/auth/token"
// 			clientSecret, _ := getClientSecret(config.ClientID, config.ClientSecret.KeyId, config.ClientSecret.TeamId, config.ClientSecret.PrivateKey)
// 			accessTokenAPIParams := map[string]string{
// 				"client_id":     config.ClientID,
// 				"client_secret": clientSecret,
// 				"grant_type":    "authorization_code",
// 			}
// 			if authCodeFromRequest != nil {
// 				accessTokenAPIParams["code"] = *authCodeFromRequest
// 			}
// 			if redirectURI != nil {
// 				accessTokenAPIParams["redirect_uri"] = *redirectURI
// 			}

// 			authorisationRedirectURL := "https://appleid.apple.com/auth/authorize"
// 			scopes := []string{"name", "email"}
// 			if config.Scope != nil {
// 				scopes = append(scopes, config.Scope...)
// 			}

// 			var additionalParams map[string]interface{} = nil
// 			if config.AuthorisationRedirect != nil && config.AuthorisationRedirect.Params != nil {
// 				additionalParams = config.AuthorisationRedirect.Params
// 			}

// 			authorizationRedirectParams := map[string]interface{}{
// 				"scope":         strings.Join(scopes, " "),
// 				"response_mode": "form_post",
// 				"response_type": "code",
// 				"client_id":     config.ClientID,
// 			}
// 			for key, value := range additionalParams {
// 				authorizationRedirectParams[key] = value
// 			}

// 			return models.TypeProviderGetResponse{
// 				AccessTokenAPI: models.AccessTokenAPI{
// 					URL:    accessTokenAPIURL,
// 					Params: accessTokenAPIParams,
// 				},
// 				AuthorisationRedirect: models.AuthorisationRedirect{
// 					URL:    authorisationRedirectURL,
// 					Params: authorizationRedirectParams,
// 				},
// 				GetProfileInfo: func(authCodeResponse interface{}) (models.UserInfo, error) {
// 					authCodeResponseJson, err := json.Marshal(authCodeResponse)
// 					if err != nil {
// 						return models.UserInfo{}, err
// 					}
// 					var accessTokenAPIResponse appleGetProfileInfoInput
// 					err = json.Unmarshal(authCodeResponseJson, &accessTokenAPIResponse)
// 					if err != nil {
// 						return models.UserInfo{}, err
// 					}
// 					return models.UserInfo{}, nil
// 				},
// 			}
// 		},
// 	}
// }

// func getClientSecret(clientId, keyId, teamId, privateKey string) (string, error) {
// 	return "", nil
// }

// type appleGetProfileInfoInput struct {
// 	AccessToken  string `json:"access_token"`
// 	ExpiresIn    int    `json:"expires_in"`
// 	TokenType    string `json:"token_type"`
// 	RefreshToken string `json:"refresh_token"`
// 	IDToken      string `json:"id_token"`
// }
