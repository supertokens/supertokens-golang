package providers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
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

			authorizationRedirectParams := map[string]interface{}{
				"scope":                  strings.Join(scopes, " "),
				"access_type":            "offline",
				"include_granted_scopes": "true",
				"response_type":          "code",
				"client_id":              config.ClientID,
			}
			for key, value := range additionalParams {
				authorizationRedirectParams[key] = value
			}

			return models.TypeProviderGetResponse{
				AccessTokenAPI: models.AccessTokenAPI{
					URL:    accessTokenAPIURL,
					Params: accessTokenAPIParams,
				},
				AuthorisationRedirect: models.AuthorisationRedirect{
					URL:    authorisationRedirectURL,
					Params: authorizationRedirectParams,
				},
				GetProfileInfo: func(authCodeResponse interface{}) (models.UserInfo, error) {
					authCodeResponseJson, err := json.Marshal(authCodeResponse)
					if err != nil {
						return models.UserInfo{}, err
					}
					var accessTokenAPIResponse googleGetProfileInfoInput
					err = json.Unmarshal(authCodeResponseJson, &accessTokenAPIResponse)
					if err != nil {
						return models.UserInfo{}, err
					}
					accessToken := accessTokenAPIResponse.AccessToken
					authHeader := "Bearer " + accessToken
					response, err := getGoogleAuthRequest(authHeader)
					if err != nil {
						return models.UserInfo{}, err
					}
					userInfo := response.(map[string]interface{})
					ID := userInfo["id"].(string)
					email := userInfo["email"].(string)
					if email == "" {
						return models.UserInfo{
							ID: ID,
						}, nil
					}
					isVerified := userInfo["verified_email"].(bool)
					return models.UserInfo{
						ID: ID,
						Email: &models.EmailStruct{
							ID:         email,
							IsVerified: isVerified,
						},
					}, nil
				},
			}
		},
	}
}

func getGoogleAuthRequest(authHeader string) (interface{}, error) {
	params := map[string]string{
		"alt": "json",
	}
	paramsJson, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	url := "https://www.googleapis.com/oauth2/v1/userinfo"
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(paramsJson)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authHeader)
	return doGetRequest(req)
}

func doGetRequest(req *http.Request) (interface{}, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type googleGetProfileInfoInput struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}
