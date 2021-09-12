package session

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
	"time"

	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
)

const (
	accessTokenCookieKey  = "sAccessToken"
	refreshTokenCookieKey = "sRefreshToken"

	// there are two of them because one is used by the server to check if the user is logged in and the other is checked by the frontend to see if the user is logged in.
	idRefreshTokenCookieKey = "sIdRefreshToken"
	idRefreshTokenHeaderKey = "id-refresh-token"

	antiCsrfHeaderKey = "anti-csrf"
	ridHeaderKey      = "rid"

	frontTokenHeaderKey = "front-token"

	frontendSDKNameHeaderKey    = "supertokens-sdk-name"
	frontendSDKVersionHeaderKey = "supertokens-sdk-version"
)

type TokenInfo struct {
	Uid string      `json:"uid"`
	Ate uint64      `json:"ate"`
	Up  interface{} `json:"up"`
}

func clearSessionFromCookie(config sessmodels.TypeNormalisedInput, res http.ResponseWriter) {
	setCookie(config, res, accessTokenCookieKey, "", 0, "accessTokenPath")
	setCookie(config, res, refreshTokenCookieKey, "", 0, "refreshTokenPath")
	setCookie(config, res, idRefreshTokenCookieKey, "", 0, "accessTokenPath")
	setHeader(res, idRefreshTokenHeaderKey, "remove", false)
	setHeader(res, "Access-Control-Expose-Headers", idRefreshTokenHeaderKey, true)
}

func attachAccessTokenToCookie(config sessmodels.TypeNormalisedInput, res http.ResponseWriter, token string, expiry uint64) {
	setCookie(config, res, accessTokenCookieKey, token, expiry, "accessTokenPath")
}

func attachRefreshTokenToCookie(config sessmodels.TypeNormalisedInput, res http.ResponseWriter, token string, expiry uint64) {
	setCookie(config, res, refreshTokenCookieKey, token, expiry, "refreshTokenPath")
}

func getAccessTokenFromCookie(req *http.Request) *string {
	return getCookieValue(req, accessTokenCookieKey)
}

func getRefreshTokenFromCookie(req *http.Request) *string {
	return getCookieValue(req, refreshTokenCookieKey)
}

func getAntiCsrfTokenFromHeaders(req *http.Request) *string {
	return getHeader(req, antiCsrfHeaderKey)
}

func getRidFromHeader(req *http.Request) *string {
	return getHeader(req, ridHeaderKey)
}

func getIDRefreshTokenFromCookie(req *http.Request) *string {
	return getCookieValue(req, idRefreshTokenCookieKey)
}

func setAntiCsrfTokenInHeaders(res http.ResponseWriter, antiCsrfToken string) {
	setHeader(res, antiCsrfHeaderKey, antiCsrfToken, false)
	setHeader(res, "Access-Control-Expose-Headers", antiCsrfHeaderKey, true)
}

func setIDRefreshTokenInHeaderAndCookie(config sessmodels.TypeNormalisedInput, res http.ResponseWriter, idRefreshToken string, expiry uint64) {
	setHeader(res, idRefreshTokenHeaderKey, idRefreshToken+";"+fmt.Sprint(expiry), false)
	setHeader(res, "Access-Control-Expose-Headers", idRefreshTokenHeaderKey, true)

	setCookie(config, res, idRefreshTokenCookieKey, idRefreshToken, expiry, "accessTokenPath")
}

func setFrontTokenInHeaders(res http.ResponseWriter, userId string, atExpiry uint64, jwtPayload interface{}) {
	tokenInfo := &TokenInfo{
		Uid: userId,
		Ate: atExpiry,
		Up:  jwtPayload,
	}
	parsed, _ := json.Marshal(tokenInfo)
	data := []byte(parsed)
	setHeader(res, frontTokenHeaderKey, base64.StdEncoding.EncodeToString(data), false)
	setHeader(res, "Access-Control-Expose-Headers", frontTokenHeaderKey, true)
}

func getCORSAllowedHeaders() []string {
	return []string{
		antiCsrfHeaderKey, ridHeaderKey,
	}
}

func setHeader(res http.ResponseWriter, key, value string, allowDuplicateKey bool) {
	existingValue := res.Header().Get(strings.ToLower(key))
	if existingValue == "" {
		res.Header().Set(key, value)
	} else if allowDuplicateKey {
		res.Header().Set(key, existingValue+", "+value)
	} else {
		res.Header().Set(key, value)
	}
}

func setCookie(config sessmodels.TypeNormalisedInput, res http.ResponseWriter, name string, value string, expires uint64, pathType string) {
	var domain string
	if config.CookieDomain != nil {
		domain = *config.CookieDomain
	}
	secure := config.CookieSecure
	sameSite := config.CookieSameSite

	path := ""
	if pathType == "refreshTokenPath" {
		path = config.RefreshTokenPath.GetAsStringDangerous()
	} else if pathType == "accessTokenPath" {
		path = "/"
	}

	var sameSiteField = http.SameSiteNoneMode
	if sameSite == "lax" {
		sameSiteField = http.SameSiteLaxMode
	} else if sameSite == "strict" {
		sameSiteField = http.SameSiteStrictMode
	}

	httpOnly := true

	if domain != "" {
		cookie := &http.Cookie{
			Name:     name,
			Value:    url.QueryEscape(value),
			Domain:   domain,
			Secure:   secure,
			HttpOnly: httpOnly,
			Expires:  time.Unix(int64(expires/1000), 0),
			Path:     path,
			SameSite: sameSiteField,
		}
		setCookieValue(res, cookie)
	} else {
		cookie := &http.Cookie{
			Name:     name,
			Value:    url.QueryEscape(value),
			Secure:   secure,
			HttpOnly: httpOnly,
			Expires:  time.Unix(int64(expires/1000), 0),
			Path:     path,
			SameSite: sameSiteField,
		}
		setCookieValue(res, cookie)
	}
}

func getHeader(request *http.Request, key string) *string {
	value := request.Header.Get(key)
	if value == "" {
		return nil
	}
	return &value
}

func getCookieValue(request *http.Request, key string) *string {
	cookies := request.Cookies()
	if len(cookies) == 0 {
		return nil
	}

	for _, cookie := range cookies {
		if cookie.Name == key {
			val, err := url.QueryUnescape(cookie.Value)
			if err != nil {
				return nil
			}
			return &val
		}
	}
	return nil
}

// setCookieValue replaces cookie.go SetCookie, it replaces the cookie values instead of appending them
func setCookieValue(w http.ResponseWriter, cookie *http.Cookie) {
	cookieHeader := w.Header().Values("Set-Cookie")
	if len(cookieHeader) == 0 {
		w.Header().Set("Set-Cookie", cookie.String())
		return
	}
	existingCookies := make(map[string]string, len(cookieHeader))
	// map existing cookies by cookie name
	for _, ch := range cookieHeader {
		existingCookies[getCookieName(ch)] = ch
	}
	// replace if already existing
	existingCookies[getCookieName(cookie.String())] = cookie.String()
	// clear previous cookies from the headers
	w.Header().Del("Set-Cookie")
	// and add them back
	for _, ck := range existingCookies {
		w.Header().Add("Set-Cookie", ck)
	}
}

func getCookieName(cookie string) string {
	parts := strings.Split(textproto.TrimString(cookie), ";")
	if len(parts) == 1 && parts[0] == "" {
		return ""
	}
	parts[0] = textproto.TrimString(parts[0])
	kv := strings.Split(parts[0], "=")
	if len(kv) == 0 {
		return ""
	}
	return kv[0]
}
