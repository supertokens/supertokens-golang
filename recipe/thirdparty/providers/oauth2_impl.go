package providers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/derekstavis/go-qs"
	"github.com/golang-jwt/jwt/v4"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const scopeParameter = "scope"
const scopeSeparator = " "

// typeCombinedOAuth2Config implements the core functionality of OAuth2. Other providers,

type typeCombinedOAuth2Config struct {
	ClientType       string
	ClientID         string
	ClientSecret     string
	Scope            []string
	AdditionalConfig map[string]interface{}

	AuthorizationEndpoint            string
	AuthorizationEndpointQueryParams map[string]interface{}
	TokenEndpoint                    string
	TokenParams                      map[string]interface{}
	ForcePKCE                        bool
	UserInfoEndpoint                 string
	JwksURI                          string
	OIDCDiscoveryEndpoint            string
	UserInfoMap                      tpmodels.TypeUserInfoMap
	ValidateIdTokenPayload           func(idTokenPayload map[string]interface{}, config *typeCombinedOAuth2Config) (bool, error)
}

func (config *typeCombinedOAuth2Config) GetAuthorisationRedirectURL(redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
	config.discoverEndpoints()

	queryParams := map[string]interface{}{
		scopeParameter:  strings.Join(config.Scope, scopeSeparator),
		"client_id":     config.ClientID,
		"redirect_uri":  redirectURIOnProviderDashboard,
		"response_type": "code",
	}
	var pkceCodeVerifier *string
	if config.ClientSecret == "" || config.ForcePKCE {
		challenge, verifier, err := generateCodeChallengeS256(32)
		if err != nil {
			return tpmodels.TypeAuthorisationRedirect{}, err
		}
		queryParams["code_challenge"] = challenge
		queryParams["code_challenge_method"] = "S256"
		pkceCodeVerifier = &verifier
	}

	for k, v := range config.AuthorizationEndpointQueryParams {
		if v == nil {
			delete(queryParams, k)
		} else {
			queryParams[k] = v
		}
	}

	url := config.AuthorizationEndpoint

	/* Transformation needed for dev keys BEGIN */
	if isUsingDevelopmentClientId(config.ClientID) {
		queryParams["client_id"] = getActualClientIdFromDevelopmentClientId(config.ClientID)
		queryParams["actual_redirect_uri"] = url
		url = DevOauthAuthorisationUrl
	}
	/* Transformation needed for dev keys END */

	queryParamsStr, err := qs.Marshal(queryParams)
	if err != nil {
		return tpmodels.TypeAuthorisationRedirect{}, err
	}

	return tpmodels.TypeAuthorisationRedirect{
		URLWithQueryParams: url + "?" + queryParamsStr,
		PKCECodeVerifier:   pkceCodeVerifier,
	}, nil
}

func (config *typeCombinedOAuth2Config) ExchangeAuthCodeForOAuthTokens(redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
	config.discoverEndpoints()

	tokenAPIURL := config.TokenEndpoint
	accessTokenAPIParams := map[string]interface{}{
		"client_id":    config.ClientID,
		"redirect_uri": redirectURIInfo.RedirectURIOnProviderDashboard,
		"code":         redirectURIInfo.RedirectURIQueryParams["code"].(string),
		"grant_type":   "authorization_code",
	}
	if config.ClientSecret != "" {
		accessTokenAPIParams["client_secret"] = config.ClientSecret
	}
	if redirectURIInfo.PKCECodeVerifier != nil {
		accessTokenAPIParams["code_verifier"] = *redirectURIInfo.PKCECodeVerifier
	}

	for k, v := range config.TokenParams {
		if v == nil {
			delete(accessTokenAPIParams, k)
		} else {
			accessTokenAPIParams[k] = v
		}
	}

	/* Transformation needed for dev keys BEGIN */
	if isUsingDevelopmentClientId(config.ClientID) {
		accessTokenAPIParams["client_id"] = getActualClientIdFromDevelopmentClientId(config.ClientID)
		accessTokenAPIParams["redirect_uri"] = DevOauthRedirectUrl
	}
	/* Transformation needed for dev keys END */

	oAuthTokens, err := doPostRequest(tokenAPIURL, accessTokenAPIParams, nil)
	if err != nil {
		return nil, err
	}

	return oAuthTokens, nil
}

func (config *typeCombinedOAuth2Config) GetUserInfo(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
	config.discoverEndpoints()

	accessToken, accessTokenOk := oAuthTokens["access_token"].(string)
	idToken, idTokenOk := oAuthTokens["id_token"].(string)

	rawUserInfoFromProvider := tpmodels.TypeRawUserInfoFromProvider{}

	if idTokenOk && config.JwksURI != "" {
		claims := jwt.MapClaims{}
		jwksURL := config.JwksURI
		jwks, err := getJWKSFromURL(jwksURL)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}
		token, err := jwt.ParseWithClaims(idToken, claims, jwks.Keyfunc)
		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}

		if !token.Valid {
			return tpmodels.TypeUserInfo{}, errors.New("invalid id_token supplied")
		}
		rawUserInfoFromProvider.FromIdTokenPayload = map[string]interface{}(claims)
	}
	if accessTokenOk && config.UserInfoEndpoint != "" {
		headers := map[string]string{
			"Authorization": "Bearer " + accessToken,
		}
		userInfoFromAccessToken, err := doGetRequest(config.UserInfoEndpoint, nil, headers)
		rawUserInfoFromProvider.FromUserInfoAPI = userInfoFromAccessToken.(map[string]interface{})

		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}
	}

	userInfoResult, err := config.getSupertokensUserInfoResultFromRawUserInfo(rawUserInfoFromProvider)
	if err != nil {
		return tpmodels.TypeUserInfo{}, err
	}

	return tpmodels.TypeUserInfo{
		ThirdPartyUserId:        userInfoResult.ThirdPartyUserId,
		Email:                   userInfoResult.EmailInfo,
		RawUserInfoFromProvider: rawUserInfoFromProvider,
	}, nil
}

func (config *typeCombinedOAuth2Config) discoverEndpoints() {
	if config.OIDCDiscoveryEndpoint != "" {
		oidcInfo, err := getOIDCDiscoveryInfo(config.OIDCDiscoveryEndpoint)

		if err == nil {
			if authURL, ok := oidcInfo["authorization_endpoint"].(string); ok {
				if config.AuthorizationEndpoint == "" {
					config.AuthorizationEndpoint = authURL
				}
			}

			if tokenURL, ok := oidcInfo["token_endpoint"].(string); ok {
				if config.TokenEndpoint == "" {
					config.TokenEndpoint = tokenURL
				}
			}

			if userInfoURL, ok := oidcInfo["userinfo_endpoint"].(string); ok {
				if config.UserInfoEndpoint == "" {
					config.UserInfoEndpoint = userInfoURL
				}
			}

			if jwksUri, ok := oidcInfo["jwks_uri"].(string); ok {
				config.JwksURI = jwksUri
			}
		}
	}
}

func (config *typeCombinedOAuth2Config) getSupertokensUserInfoResultFromRawUserInfo(rawUserInfoResponse tpmodels.TypeRawUserInfoFromProvider) (tpmodels.TypeSupertokensUserInfo, error) {
	var rawUserInfo map[string]interface{}

	if config.UserInfoMap.From == tpmodels.FromIdTokenPayload {
		if rawUserInfoResponse.FromIdTokenPayload == nil {
			return tpmodels.TypeSupertokensUserInfo{}, errors.New("rawUserInfoResponse.FromIdToken is not available")
		}
		rawUserInfo = rawUserInfoResponse.FromIdTokenPayload

		if config.ValidateIdTokenPayload != nil {
			valid, err := config.ValidateIdTokenPayload(rawUserInfo, config)
			if err != nil {
				return tpmodels.TypeSupertokensUserInfo{}, err
			}
			if !valid {
				return tpmodels.TypeSupertokensUserInfo{}, errors.New("id_token payload is invalid")
			}
		}
	} else {
		if rawUserInfoResponse.FromUserInfoAPI == nil {
			return tpmodels.TypeSupertokensUserInfo{}, errors.New("rawUserInfoResponse.FromAccessToken is not available")
		}
		rawUserInfo = rawUserInfoResponse.FromUserInfoAPI
	}

	result := tpmodels.TypeSupertokensUserInfo{}
	result.ThirdPartyUserId = fmt.Sprint(accessField(rawUserInfo, config.UserInfoMap.UserId))
	result.EmailInfo = &tpmodels.EmailStruct{
		ID: fmt.Sprint(accessField(rawUserInfo, config.UserInfoMap.Email)),
	}
	if emailVerified, ok := accessField(rawUserInfo, config.UserInfoMap.EmailVerified).(bool); ok {
		result.EmailInfo.IsVerified = emailVerified
	} else if emailVerified, ok := accessField(rawUserInfo, config.UserInfoMap.EmailVerified).(string); ok {
		result.EmailInfo.IsVerified = emailVerified == "true"
	}
	return result, nil
}
