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

package tpmodels

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/supertokens"
)

type TypeRedirectURIQueryParams = map[string]interface{}
type TypeOAuthTokens = map[string]interface{}

type TypeRawUserInfoFromProvider struct {
	FromIdTokenPayload map[string]interface{}
	FromUserInfoAPI    map[string]interface{}
}

type TypeUserInfo struct {
	ThirdPartyUserId        string
	Email                   *EmailStruct
	RawUserInfoFromProvider TypeRawUserInfoFromProvider
}

type EmailStruct struct {
	ID         string `json:"id"`
	IsVerified bool   `json:"isVerified"`
}

type TypeAuthorisationRedirect struct {
	URLWithQueryParams string
	PKCECodeVerifier   *string
}

type TypeRedirectURIInfo struct {
	RedirectURIOnProviderDashboard string                     `json:"redirectURIOnProviderDashboard"`
	RedirectURIQueryParams         TypeRedirectURIQueryParams `json:"redirectURIQueryParams"`
	PKCECodeVerifier               *string                    `json:"pkceCodeVerifier"`
}

type TypeFrom string

const (
	FromIdTokenPayload TypeFrom = "idTokenPayload"
	FromUserInfoAPI    TypeFrom = "userInfoAPI"
)

type TypeUserInfoMapFields struct {
	UserId        string `json:"userId,omitempty"`
	Email         string `json:"email,omitempty"`
	EmailVerified string `json:"emailVerified,omitempty"`
}

type TypeUserInfoMap struct {
	FromIdTokenPayload TypeUserInfoMapFields `json:"fromIdTokenPayload,omitempty"`
	FromUserInfoAPI    TypeUserInfoMapFields `json:"fromUserInfoAPI,omitempty"`
}

type User struct {
	ID         string `json:"id"`
	TimeJoined uint64 `json:"timeJoined"`
	Email      string `json:"email"`
	ThirdParty struct {
		ID     string `json:"id"`
		UserID string `json:"userId"`
	} `json:"thirdParty"`
	TenantIds []string `json:"tenantIds"`
}

type TypeInputSignInAndUp struct {
	Providers []ProviderInput
}

type TypeNormalisedInputSignInAndUp struct {
	Providers []ProviderInput
}

// IsValidFrontendRedirectURIFunc decides whether the Apple callback may redirect to the
// frontendRedirectURI found in the state. frontendRedirectURI is the target after parsing,
// normalising and re-serialising it: the browser is sent to exactly this scheme, host, port, path
// and fragment. Only the query changes: the posted fields (code, state, id_token, ...) are set on
// it, replacing params with the same name, and it is re-encoded with its keys sorted. Returning
// false rejects the request with a 400, and returning an error ends the request with that error
// without redirecting.
//
// The scheme is lowercased, "." and ".." path segments are removed (also percent-encoded ones,
// such as %2e%2e) and special characters are escaped. For example, https://app.example.com/a/../cb
// and https://app.example.com/a/%2e%2e/cb are both passed as https://app.example.com/cb.
// The value is normalised less than in the Node.js SDK: the host's case, an explicit default port
// and an empty path are kept, and an internationalised host is not converted to its IDNA
// (punycode) form. For example, https://app.example.com is passed as https://app.example.com (not
// https://app.example.com/), https://App.example.com:443/cb is passed unchanged, and
// https://bücher.example/cb is passed as https://b%C3%BCcher.example/cb (the browser opens
// https://xn--bcher-kva.example/cb).
//
// So parse frontendRedirectURI (e.g. with url.Parse) and compare its parts with your allowed
// targets, instead of comparing the string: the scheme, the host (lowercased, or converted with
// idna.Lookup.ToASCII from golang.org/x/net/idna if you allow internationalised domains), the port
// (with the default filled in) and the path, and for intent: URLs also the fragment, which picks
// the receiving app.
//
// Treat frontendRedirectURI as untrusted: whoever starts the sign-in controls the whole state.
// Never use prefix or substring checks. For example, all of these pass a
// strings.HasPrefix(uri, "https://app.example.com") check:
//   - https://app.example.com.evil.test/
//   - https://app.example.com@evil.test/
//   - intent://callback#Intent;package=attacker.app;end (the fragment picks the receiving app)
//
// The redirect carries Apple's code and id_token. Any app on the device can register a custom
// scheme, so prefer verified https universal/app links over custom schemes.
type IsValidFrontendRedirectURIFunc func(frontendRedirectURI string, req *http.Request, userContext supertokens.UserContext) (bool, error)

type TypeInput struct {
	SignInAndUpFeature TypeInputSignInAndUp
	Override           *OverrideStruct
	// IsValidFrontendRedirectURI is optional. When set, it replaces the default check that
	// the Apple callback only redirects to the configured website origin. Use it if the
	// origin depends on the request, if one backend serves several frontends, or to allow
	// mobile deep links (https app links or custom schemes such as myapp://callback).
	// Relative URLs, http(s) URLs without a valid host, and javascript:, data:, vbscript:, file:,
	// blob:, about: and filesystem: URLs are always rejected without calling it.
	//
	// The URL comes from the request and is untrusted: parse it and compare its parts, never
	// compare it by prefix or substring. See IsValidFrontendRedirectURIFunc for how the URL is
	// normalised and for examples, and prefer verified https universal/app links over custom
	// schemes, since the redirect carries Apple's code and id_token.
	IsValidFrontendRedirectURI IsValidFrontendRedirectURIFunc
}

type TypeNormalisedInput struct {
	SignInAndUpFeature         TypeNormalisedInputSignInAndUp
	Override                   OverrideStruct
	IsValidFrontendRedirectURI IsValidFrontendRedirectURIFunc
}

type OverrideStruct struct {
	Functions func(originalImplementation RecipeInterface) RecipeInterface
	APIs      func(originalImplementation APIInterface) APIInterface
}

type ProviderInput struct {
	Config                             ProviderConfig
	IncludeInNonPublicTenantsByDefault bool
	Override                           func(originalImplementation *TypeProvider) *TypeProvider
}

type ProviderConfig struct {
	ThirdPartyId string `json:"thirdPartyId"`
	Name         string `json:"name,omitempty"`

	Clients []ProviderClientConfig `json:"clients,omitempty"`

	// Fields below are optional for built-in providers
	AuthorizationEndpoint            string                 `json:"authorizationEndpoint,omitempty"`
	AuthorizationEndpointQueryParams map[string]interface{} `json:"authorizationEndpointQueryParams,omitempty"`
	TokenEndpoint                    string                 `json:"tokenEndpoint,omitempty"`
	TokenEndpointBodyParams          map[string]interface{} `json:"tokenEndpointBodyParams,omitempty"`
	UserInfoEndpoint                 string                 `json:"userInfoEndpoint,omitempty"`
	UserInfoEndpointQueryParams      map[string]interface{} `json:"userInfoEndpointQueryParams,omitempty"`
	UserInfoEndpointHeaders          map[string]interface{} `json:"userInfoEndpointHeaders,omitempty"`
	JwksURI                          string                 `json:"jwksURI,omitempty"`
	OIDCDiscoveryEndpoint            string                 `json:"oidcDiscoveryEndpoint,omitempty"`
	UserInfoMap                      TypeUserInfoMap        `json:"userInfoMap,omitempty"`
	RequireEmail                     *bool                  `json:"requireEmail,omitempty"`

	ValidateIdTokenPayload func(idTokenPayload map[string]interface{}, clientConfig ProviderConfigForClientType, userContext supertokens.UserContext) error `json:"-"`
	ValidateAccessToken    func(accessToken string, clientConfig ProviderConfigForClientType, userContext supertokens.UserContext) error                    `json:"-"`
	GenerateFakeEmail      func(thirdPartyUserId string, tenantId string, userContext supertokens.UserContext) string                                       `json:"-"`
}

type ProviderClientConfig struct {
	ClientType       string                 `json:"clientType,omitempty"` // optional
	ClientID         string                 `json:"clientId"`
	ClientSecret     string                 `json:"clientSecret,omitempty"`
	Scope            []string               `json:"scope,omitempty"`
	ForcePKCE        *bool                  `json:"forcePKCE,omitempty"`
	AdditionalConfig map[string]interface{} `json:"additionalConfig,omitempty"`
}

type ProviderConfigForClientType struct {
	Name string

	ClientID         string
	ClientSecret     string
	Scope            []string
	AdditionalConfig map[string]interface{}

	AuthorizationEndpoint            string
	AuthorizationEndpointQueryParams map[string]interface{}
	TokenEndpoint                    string
	TokenEndpointBodyParams          map[string]interface{}
	ForcePKCE                        *bool // Providers like twitter expects PKCE to be used along with secret
	UserInfoEndpoint                 string
	UserInfoEndpointQueryParams      map[string]interface{}
	UserInfoEndpointHeaders          map[string]interface{}
	JwksURI                          string
	OIDCDiscoveryEndpoint            string
	UserInfoMap                      TypeUserInfoMap
	ValidateIdTokenPayload           func(idTokenPayload map[string]interface{}, clientConfig ProviderConfigForClientType, userContext supertokens.UserContext) error
	ValidateAccessToken              func(accessToken string, clientConfig ProviderConfigForClientType, userContext supertokens.UserContext) error

	RequireEmail      *bool
	GenerateFakeEmail func(thirdPartyUserId string, tenantId string, userContext supertokens.UserContext) string
}

type TypeProvider struct {
	ID     string
	Config ProviderConfigForClientType

	GetConfigForClientType         func(clientType *string, userContext supertokens.UserContext) (ProviderConfigForClientType, error)
	GetAuthorisationRedirectURL    func(redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (TypeAuthorisationRedirect, error)
	ExchangeAuthCodeForOAuthTokens func(redirectURIInfo TypeRedirectURIInfo, userContext supertokens.UserContext) (TypeOAuthTokens, error) // For apple, add userInfo from callbackInfo to oAuthTOkens
	GetUserInfo                    func(oAuthTokens TypeOAuthTokens, userContext supertokens.UserContext) (TypeUserInfo, error)
}
