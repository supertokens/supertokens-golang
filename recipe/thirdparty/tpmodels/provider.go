package tpmodels

import "github.com/supertokens/supertokens-golang/supertokens"

type ProviderInput struct {
	ThirdPartyID string
	Config       ProviderConfigInput
	Override     func(provider *TypeProvider) *TypeProvider
}

type ProviderConfigInput struct {
	Clients []ProviderClientConfigInput

	AuthorizationEndpoint            string
	AuthorizationEndpointQueryParams map[string]interface{}
	TokenEndpoint                    string
	TokenParams                      map[string]interface{}
	ForcePKCE                        bool // Providers like twitter expects PKCE to be used along with secret
	UserInfoEndpoint                 string
	JwksURI                          string
	OIDCDiscoveryEndpoint            string
	UserInfoMap                      TypeUserInfoMap
	ValidateIdTokenPayload           func(idTokenPayload map[string]interface{}, clientConfig ProviderClientConfig) (bool, error)
}

type ProviderClientConfigInput struct {
	ClientType       string // optional
	ClientID         string
	ClientSecret     string
	Scope            []string
	AdditionalConfig map[string]interface{}
}

type ProviderClientConfig struct {
	ClientType       string // optional
	ClientID         string
	ClientSecret     string
	Scope            []string
	AdditionalConfig map[string]interface{}

	AuthorizationEndpoint            string
	AuthorizationEndpointQueryParams map[string]interface{}
	TokenEndpoint                    string
	TokenParams                      map[string]interface{}
	ForcePKCE                        bool // Providers like twitter expects PKCE to be used along with secret
	UserInfoEndpoint                 string
	JwksURI                          string
	OIDCDiscoveryEndpoint            string
	UserInfoMap                      TypeUserInfoMap
	ValidateIdTokenPayload           func(idTokenPayload map[string]interface{}, clientConfig ProviderClientConfig) (bool, error)
}

type TypeProvider struct {
	ID string

	GetConfig                      func(clientType *string, tenantId *string, input ProviderConfigInput, userContext supertokens.UserContext) (ProviderClientConfig, error)
	GetAuthorisationRedirectURL    func(clientType *string, tenantId *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (TypeAuthorisationRedirect, error)
	ExchangeAuthCodeForOAuthTokens func(clientType *string, tenantId *string, redirectInfo TypeRedirectURIInfo, userContext supertokens.UserContext) (TypeOAuthTokens, error) // For apple, add userInfo from callbackInfo to oAuthTOkens
	GetUserInfo                    func(clientType *string, tenantId *string, oAuthTokens TypeOAuthTokens, userContext supertokens.UserContext) (TypeUserInfo, error)
}
