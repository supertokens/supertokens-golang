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

func oauth2_GetAuthorisationRedirectURL(config tpmodels.ProviderConfigForClient, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
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

func oauth2_ExchangeAuthCodeForOAuthTokens(config tpmodels.ProviderConfigForClient, redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
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

func oauth2_GetUserInfo(config tpmodels.ProviderConfigForClient, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
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

	userInfoResult, err := oauth2_getSupertokensUserInfoResultFromRawUserInfo(config, rawUserInfoFromProvider)
	if err != nil {
		return tpmodels.TypeUserInfo{}, err
	}

	if config.TenantId != "" {
		userInfoResult.ThirdPartyUserId += "|" + config.TenantId // TODO delimiter
	}

	return tpmodels.TypeUserInfo{
		ThirdPartyUserId:        userInfoResult.ThirdPartyUserId,
		Email:                   userInfoResult.EmailInfo,
		RawUserInfoFromProvider: rawUserInfoFromProvider,
	}, nil
}

func oauth2_getSupertokensUserInfoResultFromRawUserInfo(config tpmodels.ProviderConfigForClient, rawUserInfoResponse tpmodels.TypeRawUserInfoFromProvider) (tpmodels.TypeSupertokensUserInfo, error) {

	if config.ValidateIdTokenPayload != nil {
		ok, err := config.ValidateIdTokenPayload(rawUserInfoResponse.FromIdTokenPayload, config)
		if err != nil {
			return tpmodels.TypeSupertokensUserInfo{}, err
		}
		if !ok {
			return tpmodels.TypeSupertokensUserInfo{}, errors.New("id_token payload validation failed")
		}
	}

	result := tpmodels.TypeSupertokensUserInfo{}
	if config.UserInfoMap.FromIdTokenPayload.UserId != "" {
		result.ThirdPartyUserId = fmt.Sprint(accessField(rawUserInfoResponse.FromIdTokenPayload, config.UserInfoMap.FromIdTokenPayload.UserId))
	} else if config.UserInfoMap.FromUserInfoAPI.UserId != "" {
		result.ThirdPartyUserId = fmt.Sprint(accessField(rawUserInfoResponse.FromUserInfoAPI, config.UserInfoMap.FromUserInfoAPI.UserId))
	} else {
		return tpmodels.TypeSupertokensUserInfo{}, errors.New("userId field is not specified in the UserInfoMap config")
	}

	if config.UserInfoMap.FromIdTokenPayload.Email != "" {
		result.EmailInfo = &tpmodels.EmailStruct{
			ID: fmt.Sprint(accessField(rawUserInfoResponse.FromIdTokenPayload, config.UserInfoMap.FromIdTokenPayload.Email)),
		}
	} else if config.UserInfoMap.FromUserInfoAPI.Email != "" {
		result.EmailInfo = &tpmodels.EmailStruct{
			ID: fmt.Sprint(accessField(rawUserInfoResponse.FromUserInfoAPI, config.UserInfoMap.FromUserInfoAPI.Email)),
		}
	} else {
		result.EmailInfo = nil
	}

	if result.EmailInfo != nil {
		if config.UserInfoMap.FromIdTokenPayload.EmailVerified != "" {
			if emailVerified, ok := accessField(rawUserInfoResponse.FromIdTokenPayload, config.UserInfoMap.FromIdTokenPayload.EmailVerified).(bool); ok {
				result.EmailInfo.IsVerified = emailVerified
			} else if emailVerified, ok := accessField(rawUserInfoResponse.FromIdTokenPayload, config.UserInfoMap.FromIdTokenPayload.EmailVerified).(string); ok {
				result.EmailInfo.IsVerified = emailVerified == "true"
			}
		} else if config.UserInfoMap.FromUserInfoAPI.EmailVerified != "" {
			if emailVerified, ok := accessField(rawUserInfoResponse.FromUserInfoAPI, config.UserInfoMap.FromUserInfoAPI.EmailVerified).(bool); ok {
				result.EmailInfo.IsVerified = emailVerified
			} else if emailVerified, ok := accessField(rawUserInfoResponse.FromUserInfoAPI, config.UserInfoMap.FromUserInfoAPI.EmailVerified).(string); ok {
				result.EmailInfo.IsVerified = emailVerified == "true"
			}
		}
	}

	return result, nil
}
