package providers

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type CustomConfig struct {
	Clients []CustomClientConfig

	AuthorizationEndpoint            string
	AuthorizationEndpointQueryParams map[string]interface{}
	TokenEndpoint                    string
	TokenParams                      map[string]interface{}
	ForcePKCE                        bool // Providers like twitter expects PKCE to be used along with secret
	UserInfoEndpoint                 string
	JwksURI                          string
	OIDCDiscoveryEndpoint            string
	UserInfoMap                      tpmodels.TypeUserInfoMap
	ValidateIdTokenPayload           func(idTokenPayload map[string]interface{}, clientConfig CustomClientConfig) (bool, error)
}

type CustomClientConfig struct {
	ClientType       string // optional
	ClientID         string
	ClientSecret     string
	Scope            []string
	AdditionalConfig map[string]interface{}
}

type TypeCustomProviderImplementation struct {
	GetConfig func(clientType *string, tenantId *string, userContext supertokens.UserContext) (CustomClientConfig, error)
	*tpmodels.TypeProvider
}

type CustomProvider struct {
	ThirdPartyID string
	Config       CustomConfig
	Override     func(provider *TypeCustomProviderImplementation) *TypeCustomProviderImplementation
}

func (input CustomProvider) Build() *tpmodels.TypeProvider {
	customImpl := input.buildInternal()
	if input.Override != nil {
		customImpl = input.Override(customImpl)
	}
	return customImpl.TypeProvider
}

func (input CustomProvider) buildInternal() *TypeCustomProviderImplementation {
	return nil // TODO impl
}
