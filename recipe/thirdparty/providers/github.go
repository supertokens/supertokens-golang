package providers

import (
	"net/http"
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
				GetProfileInfo: func(authCodeResponse interface{}) (models.UserInfo, error) {
					accessTokenAPIResponse := authCodeResponse.(githubGetProfileInfoInput)
					accessToken := accessTokenAPIResponse.AccessToken
					authHeader := "Bearer " + accessToken
					response, err := getGithubAuthRequest(authHeader)
					if err != nil {
						return models.UserInfo{}, err
					}
					userInfo := response["data"].(map[string]interface{})
					emailsInfoResponse, err := getGithubEmailsInfo(authHeader)
					if err != nil {
						return models.UserInfo{}, err
					}
					emailsInfo := emailsInfoResponse["data"].([]interface{})
					ID := userInfo["id"].(string) // github userId will be a number
					// if user has choosen not to show their email publicly, userInfo here will
					// have email as null. So we instead get the info from the emails api and
					// use the email which is maked as primary one.
					var emailInfo map[string]interface{}
					for _, emailInfo := range emailsInfo {
						emailInfoMap := emailInfo.(map[string]interface{})
						if emailInfoMap["primary"].(bool) {
							emailInfo = emailInfoMap
						}
					}
					if emailInfo == nil {
						return models.UserInfo{
							ID: ID,
						}, nil
					}
					isVerified := false
					if emailInfo != nil {
						isVerified = emailInfo["verified"].(bool)
					}
					return models.UserInfo{
						ID: ID,
						Email: &models.EmailStruct{
							ID:         emailInfo["email"].(string),
							IsVerified: isVerified,
						},
					}, nil
				},
			}
		},
	}
}

func getGithubAuthRequest(authHeader string) (map[string]interface{}, error) {
	url := "https://api.github.com/user"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	return doGetRequest(req)
}

func getGithubEmailsInfo(authHeader string) (map[string]interface{}, error) {
	url := "https://api.github.com/user/emails"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	return doGetRequest(req)
}

type githubGetProfileInfoInput struct {
	AccessToken string
	ExpiresIn   int
	TokenType   string
}
