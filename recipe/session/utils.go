package session

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/supertokens/supertokens-golang/recipe/session/api"
	"github.com/supertokens/supertokens-golang/recipe/session/models"
	"github.com/supertokens/supertokens-golang/supertokens"
	"golang.org/x/net/publicsuffix"
)

func validateAndNormaliseUserInput(appInfo supertokens.NormalisedAppinfo, config *models.TypeInput) (models.TypeNormalisedInput, error) {
	var (
		cookieDomain *string = nil
		err          error
	)

	if config != nil && config.CookieDomain != nil {
		cookieDomain, err = normaliseSessionScopeOrThrowError(*config.CookieDomain)
		if err != nil {
			return models.TypeNormalisedInput{}, err
		}
	}

	topLevelAPIDomain, err := GetTopLevelDomainForSameSiteResolution(appInfo.APIDomain.GetAsStringDangerous())
	if err != nil {
		return models.TypeNormalisedInput{}, err
	}
	topLevelWebsiteDomain, err := GetTopLevelDomainForSameSiteResolution(appInfo.WebsiteDomain.GetAsStringDangerous())
	if err != nil {
		return models.TypeNormalisedInput{}, err
	}

	cookieSameSite := cookieSameSite_LAX
	if topLevelAPIDomain != topLevelWebsiteDomain {
		cookieSameSite = cookieSameSite_NONE
	}

	if config == nil || config.CookieSameSite == nil {
		cookieSameSite = *config.CookieSameSite
	} else {
		cookieSameSite, err = normaliseSameSiteOrThrowError(*config.CookieSameSite)
		if err != nil {
			return models.TypeNormalisedInput{}, err
		}
	}
	cookieSecure := false
	if config != nil || config.CookieSecure != nil {
		cookieSecure = strings.HasPrefix(appInfo.APIDomain.GetAsStringDangerous(), "https")
	} else {
		cookieSecure = *config.CookieSecure
	}

	sessionExpiredStatusCode := 401
	if config != nil && config.SessionExpiredStatusCode != nil {
		sessionExpiredStatusCode = *config.SessionExpiredStatusCode
	}

	if config != nil && config.AntiCsrf != nil {
		if *config.AntiCsrf != antiCSRF_NONE && *config.AntiCsrf != antiCSRF_VIA_CUSTOM_HEADER && *config.AntiCsrf != antiCSRF_VIA_TOKEN {
			return models.TypeNormalisedInput{}, errors.New("antiCsrf config must be one of 'NONE' or 'VIA_CUSTOM_HEADER' or 'VIA_TOKEN'")
		}
	}

	antiCsrf := antiCSRF_NONE
	if config == nil || config.AntiCsrf == nil {
		if cookieSameSite == cookieSameSite_NONE {
			antiCsrf = antiCSRF_VIA_CUSTOM_HEADER
		} else {
			antiCsrf = antiCSRF_NONE
		}
	} else {
		antiCsrf = *config.AntiCsrf
	}

	errorHandlers := models.NormalisedErrorHandlers{
		OnTokenTheftDetected: func(sessionHandle string, userID string, req *http.Request, res http.ResponseWriter) error {
			recipeInstance, err := getRecipeInstanceOrThrowError()
			if err != nil {
				return err
			}
			return api.SendTokenTheftDetectedResponse(*recipeInstance, sessionHandle, userID, req, res)
		},
		OnTryRefreshToken: func(message string, req *http.Request, res http.ResponseWriter) error {
			recipeInstance, err := getRecipeInstanceOrThrowError()
			if err != nil {
				return err
			}
			return api.SendTryRefreshTokenResponse(*recipeInstance, message, req, res)
		},
		OnUnauthorised: func(message string, req *http.Request, res http.ResponseWriter) error {
			recipeInstance, err := getRecipeInstanceOrThrowError()
			if err != nil {
				return err
			}
			return api.SendUnauthorisedResponse(*recipeInstance, message, req, res)
		},
	}

	if config != nil && config.ErrorHandlers != nil {
		if config.ErrorHandlers.OnTokenTheftDetected != nil {
			errorHandlers.OnTokenTheftDetected = config.ErrorHandlers.OnTokenTheftDetected
		}
		if config.ErrorHandlers.OnUnauthorised != nil {
			errorHandlers.OnUnauthorised = config.ErrorHandlers.OnUnauthorised
		}
	}

	IsAnIPAPIDomain, err := supertokens.IsAnIPAddress(topLevelAPIDomain)
	if err != nil {
		return models.TypeNormalisedInput{}, err
	}
	IsAnIPWebsiteDomain, err := supertokens.IsAnIPAddress(topLevelWebsiteDomain)
	if err != nil {
		return models.TypeNormalisedInput{}, err
	}

	if cookieSameSite == cookieSameSite_NONE &&
		!cookieSecure &&
		!(topLevelAPIDomain == "localhost" || IsAnIPAPIDomain) &&
		!(topLevelWebsiteDomain == "localhost" || IsAnIPWebsiteDomain) {
		return models.TypeNormalisedInput{}, errors.New("Since your API and website domain are different, for sessions to work, please use https on your apiDomain and dont set cookieSecure to false.")
	}

	refreshAPIPath, err := supertokens.NewNormalisedURLPath(refreshAPIPath)
	if err != nil {
		return models.TypeNormalisedInput{}, err
	}

	typeNormalisedInput := models.TypeNormalisedInput{
		RefreshTokenPath:         appInfo.APIBasePath.AppendPath(*refreshAPIPath),
		CookieDomain:             cookieDomain,
		CookieSameSite:           cookieSameSite,
		CookieSecure:             cookieSecure,
		SessionExpiredStatusCode: sessionExpiredStatusCode,
		AntiCsrf:                 antiCsrf,
		Override: struct {
			Functions func(originalImplementation models.RecipeImplementation) models.RecipeImplementation
			APIs      func(originalImplementation models.APIImplementation) models.APIImplementation
		}{Functions: func(originalImplementation models.RecipeImplementation) models.RecipeImplementation {
			return originalImplementation
		}, APIs: func(originalImplementation models.APIImplementation) models.APIImplementation {
			return originalImplementation
		}},
	}

	if config != nil && config.Override != nil {
		if config.Override.Functions != nil {
			typeNormalisedInput.Override.Functions = config.Override.Functions
		}
		if config.Override.APIs != nil {
			typeNormalisedInput.Override.APIs = config.Override.APIs
		}
	}

	return typeNormalisedInput, nil
}
func normaliseSameSiteOrThrowError(sameSite string) (string, error) {
	sameSite = strings.TrimSpace(sameSite)
	sameSite = strings.ToLower(sameSite)
	if sameSite != cookieSameSite_STRICT && sameSite != cookieSameSite_LAX && sameSite != cookieSameSite_NONE {
		return "", errors.New(`cookie same site must be one of "strict", "lax", or "none"`)
	}
	return sameSite, nil
}

// TODO: implement test cases? - see node code (search for getTopLevelDomainForSameSiteResolution)
func GetTopLevelDomainForSameSiteResolution(URL string) (string, error) {
	urlObj, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	hostname := urlObj.Host
	isAnIP, err := supertokens.IsAnIPAddress(hostname)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(hostname, "localhost") || strings.HasPrefix(hostname, "localhost.org") || isAnIP {
		return "localhost", nil
	}
	parsedURL, err := publicsuffix.EffectiveTLDPlusOne(hostname)
	if err != nil {
		return "", errors.New("Please make sure that the apiDomain and websiteDomain have correct values")
	}
	return parsedURL, nil
}

func normaliseSessionScopeOrThrowError(sessionScope string) (*string, error) {
	sessionScope = strings.TrimSpace(sessionScope)
	sessionScope = strings.ToLower(sessionScope)

	sessionScope = strings.TrimPrefix(sessionScope, ".")

	if !strings.HasPrefix(sessionScope, "http://") && !strings.HasPrefix(sessionScope, "https://") {
		sessionScope = "http://" + sessionScope
	}

	urlObj, err := url.Parse(sessionScope)
	if err != nil {
		return nil, errors.New("Please provide a valid sessionScope")
	}

	sessionScope = urlObj.Host
	sessionScope = strings.TrimPrefix(sessionScope, ".")

	noDotNormalised := sessionScope

	isAnIP, err := supertokens.IsAnIPAddress(sessionScope)
	if err != nil {
		return nil, err
	}
	if sessionScope == "localhost" || isAnIP {
		noDotNormalised = sessionScope
	}
	if strings.HasPrefix(sessionScope, ".") {
		noDotNormalised = "." + sessionScope
	}
	return &noDotNormalised, nil
}

func getCurrTimeInMS() uint64 {
	return uint64(time.Now().UnixNano() / 1000000)
}

func attachCreateOrRefreshSessionResponseToRes(config models.TypeNormalisedInput, res http.ResponseWriter, response models.CreateOrRefreshAPIResponse) {
	accessToken := response.AccessToken
	refreshToken := response.RefreshToken
	idRefreshToken := response.IDRefreshToken
	setFrontTokenInHeaders(res, response.Session.UserID, response.AccessToken.Expiry, response.Session.UserDataInJWT)
	attachAccessTokenToCookie(config, res, accessToken.Token, accessToken.Expiry)
	attachRefreshTokenToCookie(config, res, refreshToken.Token, refreshToken.Expiry)
	setIDRefreshTokenInHeaderAndCookie(config, res, idRefreshToken.Token, idRefreshToken.Expiry)
	if response.AntiCsrfToken != nil {
		setAntiCsrfTokenInHeaders(res, *response.AntiCsrfToken)
	}
}
