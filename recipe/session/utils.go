package session

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func GetRidFromHeader(req *http.Request) *string {
	rid := req.Header.Get("rid")
	if rid == "" {
		return nil
	}
	return &rid
}

func GetRequiredClaimValidators(
	sessionContainer sessmodels.SessionContainer,
	overrideGlobalClaimValidators func(globalClaimValidators []claims.SessionClaimValidator, sessionContainer sessmodels.SessionContainer, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error),
	userContext supertokens.UserContext,
) ([]claims.SessionClaimValidator, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	claimValidatorsAddedByOtherRecipes := instance.getClaimValidatorsAddedByOtherRecipes()
	globalClaimValidators, err := (*instance.RecipeImpl.GetGlobalClaimValidators)(sessionContainer.GetUserIDWithContext(userContext), claimValidatorsAddedByOtherRecipes, sessionContainer.GetTenantIdWithContext(userContext), userContext)
	if err != nil {
		return nil, err
	}
	if overrideGlobalClaimValidators != nil {
		globalClaimValidators, err = overrideGlobalClaimValidators(globalClaimValidators, sessionContainer, userContext)
		if err != nil {
			return nil, err
		}
	}
	return globalClaimValidators, nil
}

func ValidateAndNormaliseUserInput(appInfo supertokens.NormalisedAppinfo, config *sessmodels.TypeInput) (sessmodels.TypeNormalisedInput, error) {
	var (
		cookieDomain      *string = nil
		olderCookieDomain *string = nil
		err               error
	)

	if config != nil && config.CookieDomain != nil {
		cookieDomain, err = normaliseSessionScopeOrThrowError(*config.CookieDomain)
		if err != nil {
			return sessmodels.TypeNormalisedInput{}, err
		}
	}

	if config != nil && config.OlderCookieDomain != nil {
		olderCookieDomain, err = normaliseSessionScopeOrThrowError(*config.OlderCookieDomain)
		if err != nil {
			return sessmodels.TypeNormalisedInput{}, err
		}
	}

	if config != nil && config.CookieSameSite != nil {
		// we have this block just to check if the user input is correct
		_, err = normaliseSameSiteOrThrowError(*config.CookieSameSite)
		if err != nil {
			return sessmodels.TypeNormalisedInput{}, err
		}
	}

	cookieSameSite := func(request *http.Request, userContext supertokens.UserContext) (string, error) {
		if config != nil && config.CookieSameSite != nil {
			return normaliseSameSiteOrThrowError(*config.CookieSameSite)
		}
		origin, err := appInfo.GetOrigin(request, userContext)
		if err != nil {
			return "", err
		}
		protocolOfWebsiteDomain, err := GetURLScheme(origin.GetAsStringDangerous())
		if err != nil {
			return "", err
		}

		protocolOfAPIDomain, err := GetURLScheme(appInfo.APIDomain.GetAsStringDangerous())
		if err != nil {
			return "", err
		}

		topLevelWebsiteDomain, err := appInfo.GetTopLevelWebsiteDomain(request, userContext)
		if err != nil {
			return "", err
		}

		if protocolOfAPIDomain != protocolOfWebsiteDomain || appInfo.TopLevelAPIDomain != topLevelWebsiteDomain {
			return CookieSameSite_NONE, nil
		}
		return CookieSameSite_LAX, nil
	}

	cookieSecure := false
	if config == nil || config.CookieSecure == nil {
		cookieSecure = strings.HasPrefix(appInfo.APIDomain.GetAsStringDangerous(), "https")
	} else {
		cookieSecure = *config.CookieSecure
	}

	sessionExpiredStatusCode := 401
	if config != nil && config.SessionExpiredStatusCode != nil {
		sessionExpiredStatusCode = *config.SessionExpiredStatusCode
	}

	invalidClaimStatusCode := 403
	if config != nil && config.InvalidClaimStatusCode != nil {
		invalidClaimStatusCode = *config.InvalidClaimStatusCode
	}

	if sessionExpiredStatusCode == invalidClaimStatusCode {
		return sessmodels.TypeNormalisedInput{}, errors.New("SessionExpiredStatusCode and InvalidClaimStatusCode cannot have the same value")
	}

	antiCsrfFunctionOrString := sessmodels.AntiCsrfFunctionOrString{
		FunctionValue: func(request *http.Request, userContext supertokens.UserContext) (string, error) {
			sameSite, err := cookieSameSite(request, userContext)
			if err != nil {
				return "", err
			}
			if sameSite == CookieSameSite_NONE {
				return AntiCSRF_VIA_CUSTOM_HEADER, nil
			}
			return AntiCSRF_NONE, nil
		},
	}
	if config != nil && config.AntiCsrf != nil {
		if *config.AntiCsrf != AntiCSRF_NONE && *config.AntiCsrf != AntiCSRF_VIA_CUSTOM_HEADER && *config.AntiCsrf != AntiCSRF_VIA_TOKEN {
			return sessmodels.TypeNormalisedInput{}, errors.New("antiCsrf config must be one of 'NONE' or 'VIA_CUSTOM_HEADER' or 'VIA_TOKEN'")
		}
		antiCsrfFunctionOrString = sessmodels.AntiCsrfFunctionOrString{
			StrValue: *config.AntiCsrf,
		}
	}

	if antiCsrfFunctionOrString.FunctionValue != nil && antiCsrfFunctionOrString.StrValue != "" {
		return sessmodels.TypeNormalisedInput{}, errors.New("should never come here")
	}
	if antiCsrfFunctionOrString.FunctionValue == nil && antiCsrfFunctionOrString.StrValue == "" {
		return sessmodels.TypeNormalisedInput{}, errors.New("should never come here")
	}

	errorHandlers := sessmodels.NormalisedErrorHandlers{
		OnTokenTheftDetected: func(sessionHandle string, userID string, req *http.Request, res http.ResponseWriter) error {
			recipeInstance, err := getRecipeInstanceOrThrowError()
			if err != nil {
				return err
			}
			return sendTokenTheftDetectedResponse(*recipeInstance, sessionHandle, userID, req, res)
		},
		OnTryRefreshToken: func(message string, req *http.Request, res http.ResponseWriter) error {
			recipeInstance, err := getRecipeInstanceOrThrowError()
			if err != nil {
				return err
			}
			return sendTryRefreshTokenResponse(*recipeInstance, message, req, res)
		},
		OnUnauthorised: func(message string, req *http.Request, res http.ResponseWriter) error {
			recipeInstance, err := getRecipeInstanceOrThrowError()
			if err != nil {
				return err
			}
			return sendUnauthorisedResponse(*recipeInstance, message, req, res)
		},
		OnInvalidClaim: func(validationErrors []claims.ClaimValidationError, req *http.Request, res http.ResponseWriter) error {
			recipeInstance, err := getRecipeInstanceOrThrowError()
			if err != nil {
				return err
			}
			return sendInvalidClaimResponse(*recipeInstance, validationErrors, req, res)
		},
		OnClearDuplicateSessionCookies: func(message string, req *http.Request, res http.ResponseWriter) error {
			return supertokens.Send200Response(res, message)
		},
	}

	if config != nil && config.ErrorHandlers != nil {
		if config.ErrorHandlers.OnTokenTheftDetected != nil {
			errorHandlers.OnTokenTheftDetected = config.ErrorHandlers.OnTokenTheftDetected
		}
		if config.ErrorHandlers.OnUnauthorised != nil {
			errorHandlers.OnUnauthorised = config.ErrorHandlers.OnUnauthorised
		}
		if config.ErrorHandlers.OnInvalidClaim != nil {
			errorHandlers.OnInvalidClaim = config.ErrorHandlers.OnInvalidClaim
		}
		if config.ErrorHandlers.OnTryRefreshToken != nil {
			errorHandlers.OnTryRefreshToken = config.ErrorHandlers.OnTryRefreshToken
		}
		if config.ErrorHandlers.OnClearDuplicateSessionCookies != nil {
			errorHandlers.OnClearDuplicateSessionCookies = config.ErrorHandlers.OnClearDuplicateSessionCookies
		}
	}

	refreshAPIPath, err := supertokens.NewNormalisedURLPath(RefreshAPIPath)
	if err != nil {
		return sessmodels.TypeNormalisedInput{}, err
	}

	if config == nil {
		config = &sessmodels.TypeInput{}
	}

	if config.GetTokenTransferMethod == nil {
		config.GetTokenTransferMethod = defaultGetTokenTransferMethod
	}

	useDynamicSigningKey := true

	if config.UseDynamicAccessTokenSigningKey != nil {
		useDynamicSigningKey = *config.UseDynamicAccessTokenSigningKey
	}

	var jwksRefreshIntervalSec uint64 = 4 * 3600 // 4 hours
	if config != nil && config.JWKSRefreshIntervalSec != nil {
		jwksRefreshIntervalSec = *config.JWKSRefreshIntervalSec
	}

	typeNormalisedInput := sessmodels.TypeNormalisedInput{
		RefreshTokenPath:         appInfo.APIBasePath.AppendPath(refreshAPIPath),
		CookieDomain:             cookieDomain,
		OlderCookieDomain:        olderCookieDomain,
		GetCookieSameSite:        cookieSameSite,
		CookieSecure:             cookieSecure,
		SessionExpiredStatusCode: sessionExpiredStatusCode,
		InvalidClaimStatusCode:   invalidClaimStatusCode,
		AntiCsrfFunctionOrString: antiCsrfFunctionOrString,
		ExposeAccessTokenToFrontendInCookieBasedAuth: config.ExposeAccessTokenToFrontendInCookieBasedAuth,
		UseDynamicAccessTokenSigningKey:              useDynamicSigningKey,
		JWKSRefreshIntervalSec:                       jwksRefreshIntervalSec,
		ErrorHandlers:                                errorHandlers,
		GetTokenTransferMethod:                       config.GetTokenTransferMethod,
		Override: sessmodels.OverrideStruct{
			Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
				return originalImplementation
			}, APIs: func(originalImplementation sessmodels.APIInterface) sessmodels.APIInterface {
				return originalImplementation
			},
			OpenIdFeature: nil},
	}

	if config != nil && config.Override != nil {
		if config.Override.Functions != nil {
			typeNormalisedInput.Override.Functions = config.Override.Functions
		}
		if config.Override.APIs != nil {
			typeNormalisedInput.Override.APIs = config.Override.APIs
		}
		typeNormalisedInput.Override.OpenIdFeature = config.Override.OpenIdFeature
	}

	return typeNormalisedInput, nil
}

var accessTokenCookiesExpiryDurationMillis uint64 = 31536000000

func normaliseSameSiteOrThrowError(sameSite string) (string, error) {
	sameSite = strings.TrimSpace(sameSite)
	sameSite = strings.ToLower(sameSite)
	if sameSite != CookieSameSite_STRICT && sameSite != CookieSameSite_LAX && sameSite != CookieSameSite_NONE {
		return "", errors.New(`cookie same site must be one of "strict", "lax", or "none"`)
	}
	return sameSite, nil
}

func GetURLScheme(URL string) (string, error) {
	urlObj, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	return urlObj.Scheme, nil
}

func normaliseSessionScopeOrThrowError(sessionScope string) (*string, error) {
	helper := func(scope string) (string, error) {
		scope = strings.TrimSpace(scope)
		scope = strings.ToLower(scope)

		scope = strings.TrimPrefix(scope, ".")

		if !strings.HasPrefix(scope, "http://") && !strings.HasPrefix(scope, "https://") {
			scope = "http://" + scope
		}

		parsedURL, err := url.Parse(scope)
		if err != nil {
			return "", errors.New("please provide a valid sessionScope")
		}

		hostname := parsedURL.Hostname()

		return hostname, nil
	}

	noDotNormalised, err := helper(sessionScope)
	if err != nil {
		return nil, err
	}

	isAnIP, err := supertokens.IsAnIPAddress(sessionScope)
	if err != nil {
		return nil, err
	}

	if noDotNormalised == "localhost" || isAnIP {
		return &noDotNormalised, nil
	}

	if strings.HasPrefix(sessionScope, ".") {
		noDotNormalised = "." + noDotNormalised
		return &noDotNormalised, nil
	}

	return &noDotNormalised, nil
}

func GetCurrTimeInMS() uint64 {
	return uint64(time.Now().UnixNano() / 1000000)
}

func SetAccessTokenInResponse(config sessmodels.TypeNormalisedInput, res http.ResponseWriter, accessToken string, frontToken string, tokenTransferMethod sessmodels.TokenTransferMethod, request *http.Request, userContext supertokens.UserContext) error {
	setFrontTokenInHeaders(res, frontToken)
	// We set the expiration to 1 year, because we can't really access the expiration of the refresh token everywhere we are setting it.
	// This should be safe to do, since this is only the validity of the cookie (set here or on the frontend) but we check the expiration of the JWT anyway.
	// Even if the token is expired the presence of the token indicates that the user could have a valid refresh
	// Some browsers now cap the maximum expiry at 400 days, so we set it to 1 year, which should suffice.
	setToken(config, res, sessmodels.AccessToken, accessToken, GetCurrTimeInMS()+accessTokenCookiesExpiryDurationMillis, tokenTransferMethod, request, userContext)

	if config.ExposeAccessTokenToFrontendInCookieBasedAuth && tokenTransferMethod == sessmodels.CookieTransferMethod {
		// We set the expiration to 1 year, because we can't really access the expiration of the refresh token everywhere we are setting it.
		// This should be safe to do, since this is only the validity of the cookie (set here or on the frontend) but we check the expiration of the JWT anyway.
		// Even if the token is expired the presence of the token indicates that the user could have a valid refresh
		// Some browsers now cap the maximum expiry at 400 days, so we set it to 1 year, which should suffice.
		setToken(config, res, sessmodels.AccessToken, accessToken, GetCurrTimeInMS()+accessTokenCookiesExpiryDurationMillis, sessmodels.HeaderTransferMethod, request, userContext)
	}
	return nil
}

func sendTryRefreshTokenResponse(recipeInstance Recipe, _ string, _ *http.Request, response http.ResponseWriter) error {
	return supertokens.SendNon200ResponseWithMessage(response, "try refresh token", recipeInstance.Config.SessionExpiredStatusCode)
}

func sendUnauthorisedResponse(recipeInstance Recipe, _ string, _ *http.Request, response http.ResponseWriter) error {
	return supertokens.SendNon200ResponseWithMessage(response, "unauthorised", recipeInstance.Config.SessionExpiredStatusCode)
}

func sendInvalidClaimResponse(recipeInstance Recipe, claimValidationErrors []claims.ClaimValidationError, _ *http.Request, response http.ResponseWriter) error {
	return supertokens.SendNon200Response(response, recipeInstance.Config.InvalidClaimStatusCode, map[string]interface{}{
		"message":               "invalid claim",
		"claimValidationErrors": claimValidationErrors,
	})
}

func sendTokenTheftDetectedResponse(recipeInstance Recipe, sessionHandle string, _ string, _ *http.Request, response http.ResponseWriter) error {
	_, err := (*recipeInstance.RecipeImpl.RevokeSession)(sessionHandle, &map[string]interface{}{})
	if err != nil {
		return err
	}
	return supertokens.SendNon200ResponseWithMessage(response, "token theft detected", recipeInstance.Config.SessionExpiredStatusCode)
}

func ValidateClaimsInPayload(claimValidators []claims.SessionClaimValidator, newAccessTokenPayload map[string]interface{}, userContext supertokens.UserContext) []claims.ClaimValidationError {
	validationErrors := []claims.ClaimValidationError{}

	for _, validator := range claimValidators {
		claimValidationResult := validator.Validate(newAccessTokenPayload, userContext)
		supertokens.LogDebugMessage(fmt.Sprint("validateClaimsInPayload ", validator.ID, " validation res ", claimValidationResult))
		if !claimValidationResult.IsValid {
			validationErrors = append(validationErrors, claims.ClaimValidationError{
				ID:     validator.ID,
				Reason: claimValidationResult.Reason,
			})
		}
	}
	return validationErrors
}

func defaultGetTokenTransferMethod(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
	// We allow fallback (checking headers then cookies) by default when validating

	if !forCreateNewSession {
		return sessmodels.AnyTransferMethod
	}

	// In create new session we respect the frontend preference by default
	authMode := GetAuthmodeFromHeader(req)
	if authMode == nil {
		return sessmodels.AnyTransferMethod
	}
	switch *authMode {
	case sessmodels.CookieTransferMethod:
		return sessmodels.CookieTransferMethod
	case sessmodels.HeaderTransferMethod:
		return sessmodels.HeaderTransferMethod
	default:
		return sessmodels.AnyTransferMethod
	}
}
