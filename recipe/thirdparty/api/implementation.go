/* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"golang.org/x/net/idna"
)

func MakeAPIImplementation() tpmodels.APIInterface {

	authorisationUrlGET := func(provider *tpmodels.TypeProvider, redirectURIOnProviderDashboard string, tenantId string, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpmodels.AuthorisationUrlGETResponse, error) {
		authRedirect, err := provider.GetAuthorisationRedirectURL(redirectURIOnProviderDashboard, userContext)
		if err != nil {
			return tpmodels.AuthorisationUrlGETResponse{}, err
		}

		return tpmodels.AuthorisationUrlGETResponse{
			OK: &authRedirect,
		}, nil
	}

	signInUpPOST := func(provider *tpmodels.TypeProvider, input tpmodels.TypeSignInUpInput, tenantId string, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpmodels.SignInUpPOSTResponse, error) {
		var oAuthTokens map[string]interface{}
		var err error

		if input.RedirectURIInfo != nil {
			oAuthTokens, err = provider.ExchangeAuthCodeForOAuthTokens(*input.RedirectURIInfo, userContext)
			if err != nil {
				return tpmodels.SignInUpPOSTResponse{}, err
			}
		} else {
			oAuthTokens = *input.OAuthTokens
		}

		userInfo, err := provider.GetUserInfo(oAuthTokens, userContext)
		if err != nil {
			return tpmodels.SignInUpPOSTResponse{}, err
		}

		if userInfo.Email == nil && provider.Config.RequireEmail != nil && !*provider.Config.RequireEmail {
			userInfo.Email = &tpmodels.EmailStruct{
				ID:         provider.Config.GenerateFakeEmail(userInfo.ThirdPartyUserId, tenantId, userContext),
				IsVerified: true,
			}
		}

		emailInfo := userInfo.Email
		if emailInfo == nil {
			return tpmodels.SignInUpPOSTResponse{
				NoEmailGivenByProviderError: &struct{}{},
			}, nil
		}

		response, err := (*options.RecipeImplementation.SignInUp)(provider.ID, userInfo.ThirdPartyUserId, emailInfo.ID, oAuthTokens, userInfo.RawUserInfoFromProvider, tenantId, userContext)
		if err != nil {
			return tpmodels.SignInUpPOSTResponse{}, err
		}

		if emailInfo.IsVerified {
			evInstance := emailverification.GetRecipeInstance()
			if evInstance != nil {
				tokenResponse, err := (*evInstance.RecipeImpl.CreateEmailVerificationToken)(response.OK.User.ID, response.OK.User.Email, tenantId, userContext)
				if err != nil {
					return tpmodels.SignInUpPOSTResponse{}, err
				}
				if tokenResponse.OK != nil {
					_, err := (*evInstance.RecipeImpl.VerifyEmailUsingToken)(tokenResponse.OK.Token, tenantId, userContext)
					if err != nil {
						return tpmodels.SignInUpPOSTResponse{}, err
					}
				}
			}
		}

		session, err := session.CreateNewSession(options.Req, options.Res, tenantId, response.OK.User.ID, nil, nil, userContext)
		if err != nil {
			return tpmodels.SignInUpPOSTResponse{}, err
		}
		return tpmodels.SignInUpPOSTResponse{
			OK: &struct {
				CreatedNewUser          bool
				User                    tpmodels.User
				Session                 sessmodels.SessionContainer
				OAuthTokens             tpmodels.TypeOAuthTokens
				RawUserInfoFromProvider tpmodels.TypeRawUserInfoFromProvider
			}{
				CreatedNewUser:          response.OK.CreatedNewUser,
				User:                    response.OK.User,
				Session:                 session,
				OAuthTokens:             oAuthTokens,
				RawUserInfoFromProvider: userInfo.RawUserInfoFromProvider,
			},
		}, nil
	}

	appleRedirectHandlerPOST := func(formPostInfoFromProvider map[string]interface{}, options tpmodels.APIOptions, userContext supertokens.UserContext) error {
		state, ok := formPostInfoFromProvider["state"].(string)
		if !ok || state == "" {
			supertokens.LogDebugMessage("appleRedirectHandlerPOST: rejecting the callback because the state is missing or not a string")
			return supertokens.BadInputError{Msg: "Please provide the state in the request body"}
		}

		stateBytes, err := decodeAppleState(state)
		if errors.Is(err, errAppleStateNotBase64) {
			supertokens.LogDebugMessage("appleRedirectHandlerPOST: rejecting the callback because the state is not base64 encoded")
			return supertokens.BadInputError{Msg: "state must be base64 encoded"}
		}
		if err != nil {
			// The state only has base64 characters, but can't be decoded (e.g. because of its length)
			supertokens.LogDebugMessage("appleRedirectHandlerPOST: rejecting the callback because the state can not be decoded: " + err.Error())
			return supertokens.BadInputError{Msg: "state must be a base64 encoded JSON object"}
		}

		var stateObj map[string]interface{}
		err = json.Unmarshal(stateBytes, &stateObj)
		if err != nil || stateObj == nil {
			supertokens.LogDebugMessage("appleRedirectHandlerPOST: rejecting the callback because the state is not a JSON object")
			return supertokens.BadInputError{Msg: "state must be a base64 encoded JSON object"}
		}

		redirectURL, ok := stateObj["frontendRedirectURI"].(string)
		if !ok {
			supertokens.LogDebugMessage("appleRedirectHandlerPOST: rejecting the callback because frontendRedirectURI in the state is missing or not a string")
			return supertokens.BadInputError{Msg: "state must contain frontendRedirectURI as a string"}
		}

		parsedRedirectURL, err := url.Parse(redirectURL)
		if err != nil {
			supertokens.LogDebugMessage("appleRedirectHandlerPOST: rejecting the callback because frontendRedirectURI can not be parsed: " + err.Error())
			if options.Config.IsValidFrontendRedirectURI != nil {
				return supertokens.BadInputError{Msg: "frontendRedirectURI in state must be an absolute URL with an allowed scheme"}
			}
			return supertokens.BadInputError{Msg: "frontendRedirectURI in state must be an absolute http or https URL"}
		}

		// Browsers remove "." and ".." segments (also when percent-encoded) from the path and convert
		// internationalised hosts of http(s) URLs to their ASCII form when they open the URL, so we check
		// and redirect to the URL in that form: this way the checks below and IsValidFrontendRedirectURI
		// see the URL the browser opens.
		parsedRedirectURL = withASCIIHost(removeDotSegments(parsedRedirectURL))

		if options.Config.IsValidFrontendRedirectURI != nil {
			// The app decides which targets are allowed, so we also accept custom schemes (e.g.: for mobile deep
			// links) here, but never schemes that could run code or read local data in the browser.
			err := checkRedirectTargetForCustomCheck(parsedRedirectURL)
			if err != nil {
				supertokens.LogDebugMessage("appleRedirectHandlerPOST: rejecting frontendRedirectURI: " + err.Error())
				return supertokens.BadInputError{Msg: "frontendRedirectURI in state must be an absolute URL with an allowed scheme"}
			}

			// We pass the URL as it will be serialised in the Location header (before the posted fields
			// are added), so the app validates exactly where the browser will be sent.
			isValid, err := options.Config.IsValidFrontendRedirectURI(parsedRedirectURL.String(), options.Req, userContext)
			if err != nil {
				return err
			}
			if !isValid {
				supertokens.LogDebugMessage("appleRedirectHandlerPOST: rejecting frontendRedirectURI because IsValidFrontendRedirectURI returned false")
				return supertokens.BadInputError{Msg: "frontendRedirectURI in state is not allowed"}
			}
		} else {
			originMismatchError := supertokens.BadInputError{Msg: "frontendRedirectURI in state does not match the configured origin"}

			redirectOrigin, err := getComparableOrigin(parsedRedirectURL)
			if err != nil {
				supertokens.LogDebugMessage("appleRedirectHandlerPOST: rejecting frontendRedirectURI: " + err.Error() + ". To allow other targets (e.g. mobile deep links), use IsValidFrontendRedirectURI.")
				return supertokens.BadInputError{Msg: "frontendRedirectURI in state must be an absolute http or https URL"}
			}

			websiteOrigin, err := options.AppInfo.GetOrigin(options.Req, userContext)
			if err != nil {
				return err
			}
			parsedWebsiteOrigin, err := url.Parse(websiteOrigin.GetAsStringDangerous())
			if err != nil {
				return err
			}
			expectedOrigin, err := getComparableOrigin(parsedWebsiteOrigin)
			if err != nil {
				supertokens.LogDebugMessage("appleRedirectHandlerPOST: rejecting frontendRedirectURI because the configured origin can not be compared: " + err.Error() + ". Use IsValidFrontendRedirectURI to allow redirects in this case.")
				return originMismatchError
			}

			if redirectOrigin != expectedOrigin {
				supertokens.LogDebugMessage(fmt.Sprintf("appleRedirectHandlerPOST: rejecting frontendRedirectURI because its origin (%s) does not match the configured origin (%s). If the origin depends on the request, check what it resolves to for this request, or use IsValidFrontendRedirectURI.", redirectOrigin, expectedOrigin))
				return originMismatchError
			}
		}

		query := parsedRedirectURL.Query()

		// The posted fields replace any query params with the same name, so the frontend always reads
		// the values Apple sent
		for k, v := range formPostInfoFromProvider {
			query.Set(k, fmt.Sprint(v))
		}

		parsedRedirectURL.RawQuery = query.Encode()

		options.Res.Header().Set("Location", parsedRedirectURL.String())
		options.Res.WriteHeader(http.StatusSeeOther)

		return nil
	}

	return tpmodels.APIInterface{
		AuthorisationUrlGET:      &authorisationUrlGET,
		SignInUpPOST:             &signInUpPOST,
		AppleRedirectHandlerPOST: &appleRedirectHandlerPOST,
	}
}

// The Apple callback state is base64 encoded, either with the standard or the URL-safe alphabet,
// and with or without padding
var appleStateRegex = regexp.MustCompile(`^[A-Za-z0-9+/_-]+={0,2}$`)

// errAppleStateNotBase64 is returned by decodeAppleState for a state that contains characters
// other than the base64 ones
var errAppleStateNotBase64 = errors.New("the state is not base64 encoded")

func decodeAppleState(state string) ([]byte, error) {
	if !appleStateRegex.MatchString(state) {
		return nil, errAppleStateNotBase64
	}
	normalisedState := strings.NewReplacer("-", "+", "_", "/").Replace(strings.TrimRight(state, "="))
	return base64.RawStdEncoding.DecodeString(normalisedState)
}

// Schemes we never redirect to, because the browser would run code or expose local data with them
var disallowedRedirectSchemes = []string{"javascript", "data", "vbscript", "file", "blob", "about", "filesystem"}

// checkRedirectTargetForCustomCheck returns an error if the URL can't be passed to the app's
// IsValidFrontendRedirectURI: it must be absolute, http(s) URLs must have a valid host (see
// toASCIIHost), and the scheme must not be one of disallowedRedirectSchemes.
func checkRedirectTargetForCustomCheck(parsedURL *url.URL) error {
	// url.Parse lowercases the scheme
	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme == "" {
		return errors.New("the URL must be absolute")
	}
	if scheme == "http" || scheme == "https" {
		_, err := getComparableOrigin(parsedURL)
		return err
	}
	for _, disallowedScheme := range disallowedRedirectSchemes {
		if scheme == disallowedScheme {
			return errors.New("the URL uses a disallowed scheme: " + scheme)
		}
	}
	return nil
}

// removeDotSegments returns a copy of the URL with the "." and ".." segments removed from its path, as
// described in RFC 3986 section 5.2.4 (see removeDotSegmentsFromPath). The rest of the URL is kept as
// it is. URLs with an opaque path (e.g. myapp:callback) or a path that doesn't start with "/" are
// returned unchanged, because their path has no segments to remove.
//
// If the URL has no host, a resulting path that starts with "//" is prefixed with "/." (which the
// browser removes again), so that the serialised URL is not read as one with a host: e.g.
// myapp:/a/..//evil.test/x becomes myapp:/.//evil.test/x instead of myapp://evil.test/x.
func removeDotSegments(parsedURL *url.URL) *url.URL {
	result := *parsedURL
	if result.Opaque != "" {
		return &result
	}
	// We work on the escaped path, so the percent-encoded characters other than the dots of dot
	// segments are kept as they were (e.g. %2f does not become a separator)
	escapedPath := result.EscapedPath()
	if !strings.HasPrefix(escapedPath, "/") {
		return &result
	}
	escapedPath = removeDotSegmentsFromPath(escapedPath)
	if result.Host == "" && result.User == nil && strings.HasPrefix(escapedPath, "//") {
		escapedPath = "/." + escapedPath
	}
	path, err := url.PathUnescape(escapedPath)
	if err != nil {
		// Can't happen: removing segments from a valid escaped path leaves a valid escaped path
		return &result
	}
	result.Path = path
	result.RawPath = escapedPath
	return &result
}

// removeDotSegmentsFromPath removes the "." and ".." segments from an escaped path that starts with
// "/", following RFC 3986 section 5.2.4. Like browsers, it also treats the percent-encoded forms
// (%2e in either case, e.g. %2e%2e or .%2E) as dot segments. Empty segments are kept, so
// /..//evil.test/x becomes //evil.test/x.
func removeDotSegmentsFromPath(escapedPath string) string {
	segments := strings.Split(escapedPath[1:], "/")
	output := make([]string, 0, len(segments))
	for i, segment := range segments {
		isLastSegment := i == len(segments)-1
		switch strings.ToLower(segment) {
		case ".", "%2e":
			if isLastSegment {
				// A path ending with a dot segment ends with a "/"
				output = append(output, "")
			}
		case "..", ".%2e", "%2e.", "%2e%2e":
			if len(output) > 0 {
				output = output[:len(output)-1]
			}
			if isLastSegment {
				output = append(output, "")
			}
		default:
			output = append(output, segment)
		}
	}
	return "/" + strings.Join(output, "/")
}

// withASCIIHost returns the URL with the host converted to its IDNA A-label (punycode) form (see
// toASCIIHost), which is also the form browsers open, if it is an http(s) URL whose host contains
// non-ASCII characters (e.g. https://Bücher.example/cb becomes https://xn--bcher-kva.example/cb).
// ASCII hosts keep their case. Other URLs, and URLs whose host can't be converted (the checks
// reject them), are returned unchanged.
func withASCIIHost(parsedURL *url.URL) *url.URL {
	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		return parsedURL
	}
	host := parsedURL.Hostname()
	if isASCII(host) {
		return parsedURL
	}
	asciiHost, err := toASCIIHost(host)
	if err != nil {
		return parsedURL
	}
	result := *parsedURL
	result.Host = asciiHost
	if port := parsedURL.Port(); port != "" {
		result.Host += ":" + port
	}
	return &result
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}

// hostIDNAProfile converts hosts to their IDNA A-label (punycode) form the way browsers do (UTS #46
// with the options of the URL standard's "domain to ASCII"): it maps the host (e.g. lowercases it
// and maps full-width characters), and allows characters such as '_' and hyphens in any position.
var hostIDNAProfile = idna.New(
	idna.MapForLookup(),
	idna.Transitional(false),
	idna.BidiRule(),
	idna.StrictDomainName(false),
	idna.CheckHyphens(false),
)

// forbiddenDomainCodePoints are the ASCII characters (other than the C0 controls and DEL) that the
// URL standard doesn't allow in a domain. IDNA mapping can produce them from other characters, e.g.
// the full-width solidus (U+FF0F) becomes "/".
const forbiddenDomainCodePoints = " #%/:<>?@[\\]^|"

// toASCIIHost returns the host (a domain name or an IP address, without brackets and port) in a
// lowercase form that can be compared directly: domain names are converted to their IDNA A-label
// (punycode) form with hostIDNAProfile. It returns an error for domain names that can't be
// converted, that contain characters the URL standard doesn't allow in a domain after the
// conversion, or that have an empty label (other than after a trailing dot, e.g. a..example).
func toASCIIHost(host string) (string, error) {
	if net.ParseIP(host) != nil {
		return strings.ToLower(host), nil
	}
	asciiHost, err := hostIDNAProfile.ToASCII(host)
	if err != nil {
		return "", errors.New("the URL's host is not a valid domain name: " + err.Error())
	}
	for i := 0; i < len(asciiHost); i++ {
		if asciiHost[i] <= 0x20 || asciiHost[i] == 0x7f || strings.IndexByte(forbiddenDomainCodePoints, asciiHost[i]) != -1 {
			return "", fmt.Errorf("the URL's host is not a valid domain name: it contains %q", asciiHost[i])
		}
	}
	for _, label := range strings.Split(strings.TrimSuffix(asciiHost, "."), ".") {
		if label == "" {
			return "", errors.New("the URL's host is not a valid domain name: it has an empty label")
		}
	}
	return strings.ToLower(asciiHost), nil
}

type comparableOrigin struct {
	scheme string
	host   string
	port   string
}

func (o comparableOrigin) String() string {
	return o.scheme + "://" + o.host + ":" + o.port
}

// getComparableOrigin returns the origin of an absolute http(s) URL in a form that can be compared
// directly: the scheme is lowercased, the host is converted to its lowercase IDNA A-label (punycode)
// form (see toASCIIHost), and the default port is filled in. It returns an error for other URLs and
// for hosts toASCIIHost rejects.
func getComparableOrigin(parsedURL *url.URL) (comparableOrigin, error) {
	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		return comparableOrigin{}, errors.New("the URL must be an absolute http or https URL")
	}

	host := parsedURL.Hostname()
	if host == "" {
		return comparableOrigin{}, errors.New("the URL must have a host")
	}
	host, err := toASCIIHost(host)
	if err != nil {
		return comparableOrigin{}, err
	}

	port := parsedURL.Port()
	if port == "" {
		if scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}

	return comparableOrigin{
		scheme: scheme,
		host:   host,
		port:   port,
	}, nil
}
