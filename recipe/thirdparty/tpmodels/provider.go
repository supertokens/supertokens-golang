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
	ValidateIdTokenPayload           func(idTokenPayload map[string]interface{}, clientConfig ProviderConfigForClientType) error
	TenantId                         string
}

type ProviderClientConfig struct {
	ClientType       string // optional
	ClientID         string
	ClientSecret     string
	Scope            []string
	AdditionalConfig map[string]interface{}
}

type ProviderConfigForClientType struct {
	ClientID         string
	ClientSecret     string
	Scope            []string
	AdditionalConfig map[string]interface{}

	AuthorizationEndpoint            string
	AuthorizationEndpointQueryParams map[string]interface{}
	TokenEndpoint                    string
	TokenParams                      map[string]interface{}
	UserInfoEndpointQueryParams      map[string]interface{}
	UserInfoEndpointHeaders          map[string]interface{}
	ForcePKCE                        bool // Providers like twitter expects PKCE to be used along with secret
	UserInfoEndpoint                 string
	JwksURI                          string
	OIDCDiscoveryEndpoint            string
	UserInfoMap                      TypeUserInfoMap
	ValidateIdTokenPayload           func(idTokenPayload map[string]interface{}, clientConfig ProviderConfigForClientType) error
	TenantId                         string
}

type TypeProvider struct {
	ID string

	GetAllClientTypeConfigForTenant func(tenantId *string, userContext supertokens.UserContext) (ProviderConfig, error)
	GetConfigForClientType          func(clientType *string, input ProviderConfig, userContext supertokens.UserContext) (ProviderConfigForClientType, error)
	GetAuthorisationRedirectURL     func(config ProviderConfigForClientType, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (TypeAuthorisationRedirect, error)
	ExchangeAuthCodeForOAuthTokens  func(config ProviderConfigForClientType, redirectURIInfo TypeRedirectURIInfo, userContext supertokens.UserContext) (TypeOAuthTokens, error) // For apple, add userInfo from callbackInfo to oAuthTOkens
	GetUserInfo                     func(config ProviderConfigForClientType, oAuthTokens TypeOAuthTokens, userContext supertokens.UserContext) (TypeUserInfo, error)
}
