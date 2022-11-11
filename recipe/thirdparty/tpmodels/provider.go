package tpmodels

import "github.com/supertokens/supertokens-golang/supertokens"

type ProviderInput struct {
	ThirdPartyID string
	Config       ProviderConfig
	Override     func(provider *TypeProvider) *TypeProvider
}

type ProviderConfig struct {
	Clients []ProviderClientConfig

	AuthorizationEndpoint            string
	AuthorizationEndpointQueryParams map[string]interface{}
	TokenEndpoint                    string
	TokenParams                      map[string]interface{}
	ForcePKCE                        bool // Providers like twitter expects PKCE to be used along with secret
	UserInfoEndpoint                 string
	JwksURI                          string
	OIDCDiscoveryEndpoint            string
	UserInfoMap                      TypeUserInfoMap
	ValidateIdTokenPayload           func(idTokenPayload map[string]interface{}, clientConfig ProviderConfigForClient) (bool, error)
	TenantId                         string
}

type ProviderClientConfig struct {
	ClientType       string // optional
	ClientID         string
	ClientSecret     string
	Scope            []string
	AdditionalConfig map[string]interface{}
}

type ProviderConfigForClient struct {
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
	ValidateIdTokenPayload           func(idTokenPayload map[string]interface{}, clientConfig ProviderConfigForClient) (bool, error)
	TenantId                         string
}

type TypeProvider struct {
	ID string

	GetProviderConfig              func(tenantId *string, userContext supertokens.UserContext) (ProviderConfig, error)
	GetConfig                      func(clientType *string, input ProviderConfig, userContext supertokens.UserContext) (ProviderConfigForClient, error)
	GetAuthorisationRedirectURL    func(config ProviderConfigForClient, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (TypeAuthorisationRedirect, error)
	ExchangeAuthCodeForOAuthTokens func(config ProviderConfigForClient, redirectURIInfo TypeRedirectURIInfo, userContext supertokens.UserContext) (TypeOAuthTokens, error) // For apple, add userInfo from callbackInfo to oAuthTOkens
	GetUserInfo                    func(config ProviderConfigForClient, oAuthTokens TypeOAuthTokens, userContext supertokens.UserContext) (TypeUserInfo, error)
}
