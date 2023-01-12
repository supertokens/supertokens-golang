package providers

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func oauth2_GetAuthorisationRedirectURL(config tpmodels.ProviderConfigForClientType, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
	queryParams := map[string]interface{}{
		"scope":         strings.Join(config.Scope, " "),
		"client_id":     config.ClientID,
		"redirect_uri":  redirectURIOnProviderDashboard,
		"response_type": "code",
	}
	var pkceCodeVerifier *string
	if config.ClientSecret == "" || (config.ForcePKCE != nil && *config.ForcePKCE) {
		challenge, verifier, err := generateCodeChallengeS256(64) // According to https://www.rfc-editor.org/rfc/rfc7636, length must be between 43 and 128
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

	authUrl := config.AuthorizationEndpoint

	/* Transformation needed for dev keys BEGIN */
	if isUsingDevelopmentClientId(config.ClientID) {
		queryParams["client_id"] = getActualClientIdFromDevelopmentClientId(config.ClientID)
		queryParams["actual_redirect_uri"] = authUrl
		authUrl = DevOauthAuthorisationUrl
	}
	/* Transformation needed for dev keys END */

	urlObj, err := url.Parse(authUrl)
	if err != nil {
		return tpmodels.TypeAuthorisationRedirect{}, err
	}

	queryParamsObj := urlObj.Query()
	for k, v := range queryParams {
		queryParamsObj.Add(k, fmt.Sprint(v))
	}
	urlObj.RawQuery = queryParamsObj.Encode()

	return tpmodels.TypeAuthorisationRedirect{
		URLWithQueryParams: urlObj.String(),
		PKCECodeVerifier:   pkceCodeVerifier,
	}, nil
}

func oauth2_ExchangeAuthCodeForOAuthTokens(config tpmodels.ProviderConfigForClientType, redirectURIInfo tpmodels.TypeRedirectURIInfo, userContext supertokens.UserContext) (tpmodels.TypeOAuthTokens, error) {
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

	for k, v := range config.TokenEndpointBodyParams {
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

func oauth2_GetUserInfo(config tpmodels.ProviderConfigForClientType, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
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
		if config.ValidateIdTokenPayload != nil {
			err := config.ValidateIdTokenPayload(rawUserInfoFromProvider.FromIdTokenPayload, config)
			if err != nil {
				return tpmodels.TypeUserInfo{}, err
			}
		}
	}

	if accessTokenOk && config.UserInfoEndpoint != "" {
		headers := map[string]string{
			"Authorization": "Bearer " + accessToken,
		}
		queryParams := map[string]interface{}{}

		for k, v := range config.UserInfoEndpointHeaders {
			if v == nil {
				delete(headers, k)
			} else {
				headers[k] = fmt.Sprint(v)
			}
		}

		for k, v := range config.UserInfoEndpointQueryParams {
			if v == nil {
				delete(queryParams, k)
			} else {
				queryParams[k] = v
			}
		}

		userInfoFromAccessToken, err := doGetRequest(config.UserInfoEndpoint, queryParams, headers)
		rawUserInfoFromProvider.FromUserInfoAPI = userInfoFromAccessToken.(map[string]interface{})

		if err != nil {
			return tpmodels.TypeUserInfo{}, err
		}
	}

	userInfoResult, err := oauth2_getSupertokensUserInfoResultFromRawUserInfo(config, rawUserInfoFromProvider)
	if err != nil {
		return tpmodels.TypeUserInfo{}, err
	}

	if config.TenantId != nil && *config.TenantId != tpmodels.DefaultTenantId {
		userInfoResult.ThirdPartyUserId += "|" + *config.TenantId
	}

	return tpmodels.TypeUserInfo{
		ThirdPartyUserId:        userInfoResult.ThirdPartyUserId,
		Email:                   userInfoResult.Email,
		RawUserInfoFromProvider: rawUserInfoFromProvider,
	}, nil
}

func oauth2_getSupertokensUserInfoResultFromRawUserInfo(config tpmodels.ProviderConfigForClientType, rawUserInfoResponse tpmodels.TypeRawUserInfoFromProvider) (tpmodels.TypeUserInfo, error) {
	result := tpmodels.TypeUserInfo{}

	if config.UserInfoMap.FromIdTokenPayload.UserId != "" {
		if userId, ok := accessField(rawUserInfoResponse.FromUserInfoAPI, config.UserInfoMap.FromUserInfoAPI.UserId); ok {
			result.ThirdPartyUserId = fmt.Sprint(userId)
		}
	}

	if config.UserInfoMap.FromIdTokenPayload.UserId != "" {
		if userId, ok := accessField(rawUserInfoResponse.FromIdTokenPayload, config.UserInfoMap.FromIdTokenPayload.UserId); ok {
			result.ThirdPartyUserId = fmt.Sprint(userId)
		}
	}

	if result.ThirdPartyUserId == "" {
		return result, errors.New("third party user id is missing")
	}

	var email string

	if config.UserInfoMap.FromUserInfoAPI.Email != "" {
		if emailVal, ok := accessField(rawUserInfoResponse.FromUserInfoAPI, config.UserInfoMap.FromUserInfoAPI.Email); ok {
			email = fmt.Sprint(emailVal)
		}
	}

	if config.UserInfoMap.FromIdTokenPayload.Email != "" {
		if emailVal, ok := accessField(rawUserInfoResponse.FromIdTokenPayload, config.UserInfoMap.FromIdTokenPayload.Email); ok {
			email = fmt.Sprint(emailVal)
		}
	}

	if email != "" {
		result.Email = &tpmodels.EmailStruct{
			ID:         email,
			IsVerified: false,
		}

		if config.UserInfoMap.FromUserInfoAPI.EmailVerified != "" {
			if emailVerifiedVal, ok := accessField(rawUserInfoResponse.FromUserInfoAPI, config.UserInfoMap.FromUserInfoAPI.EmailVerified); ok {
				if emailVerified, ok := emailVerifiedVal.(bool); ok {
					result.Email.IsVerified = emailVerified
				} else if emailVerified, ok := emailVerifiedVal.(string); ok {
					result.Email.IsVerified = emailVerified == "true"
				}
			}
		}

		if config.UserInfoMap.FromIdTokenPayload.EmailVerified != "" {
			if emailVerifiedVal, ok := accessField(rawUserInfoResponse.FromIdTokenPayload, config.UserInfoMap.FromIdTokenPayload.EmailVerified); ok {
				if emailVerified, ok := emailVerifiedVal.(bool); ok {
					result.Email.IsVerified = emailVerified
				} else if emailVerified, ok := emailVerifiedVal.(string); ok {
					result.Email.IsVerified = emailVerified == "true"
				}
			}
		}

	}

	return result, nil
}
