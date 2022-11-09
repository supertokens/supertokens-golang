package supertokens

import (
	"errors"
)

type CreateOrUpdateTenantIDConfigResponse struct {
	OK *struct {
		Created bool
		Updated bool
	}
}

type TenantConfig struct {
	ThirdPartyId string `json:"thirdPartyId"`

	Clients []TenantClient `json:"clients"`

	// Fields below are optional for built-in providers
	AuthorizationEndpoint            string                 `json:"authorizationEndpoint,omitempty"`
	AuthorizationEndpointQueryParams map[string]interface{} `json:"authorizationEndpointQueryParams,omitempty"`
	TokenEndpoint                    string                 `json:"tokenEndpoint,omitempty"`
	TokenParams                      map[string]interface{} `json:"tokenParams,omitempty"`
	ForcePKCE                        bool                   `json:"forcePKCE,omitempty"`
	UserInfoEndpoint                 string                 `json:"userInfoEndpoint,omitempty"`
	JwksURI                          string                 `json:"jwksURI,omitempty"`
	OIDCDiscoveryEndpoint            string                 `json:"oidcDiscoveryEndpoint,omitempty"`
	UserInfoMap                      struct {
		From               string `json:"from"`
		IdField            string `json:"idField"`
		EmailField         string `json:"emailField"`
		EmailVerifiedField string `json:"emailVerifiedField"`
	} `json:"userInfoMap,omitempty"`

	FrontendInfo struct {
		Name            string `json:"name"`
		ButtonStyle     string `json:"buttonStyle"`
		ButtonComponent string `json:"buttonComponent"`
	} `json:"frontendInfo"`
}

type TenantClient struct {
	ClientID         string                 `json:"clientId"`
	ClientSecret     string                 `json:"clientSecret"`
	Scope            []string               `json:"scope"`
	AdditionalConfig map[string]interface{} `json:"additionalConfig"`
}

func CreateOrUpdateTenantIDConfigMapping(thirdPartyId string, tenantId string, config TenantConfig) (CreateOrUpdateTenantIDConfigResponse, error) {
	// TODO impl
	return CreateOrUpdateTenantIDConfigResponse{}, errors.New("needs implementation")
}

type FetchTenantIDConfigResponse struct {
	OK *struct {
		Config TenantConfig
	}
	UnknownMappingError *struct{}
}

func FetchTenantIDConfigMapping(thirdPartyId string, tenantId string) (FetchTenantIDConfigResponse, error) {
	// TODO impl
	return FetchTenantIDConfigResponse{}, errors.New("needs implementation")
}

type DeleteTenantIDConfigResponse struct {
	OK *struct {
		DidMappingExist bool
	}
}

func DeleteTenantIDConfigMapping(thirdPartyId string, tenantId string) (DeleteTenantIDConfigResponse, error) {
	// TODO impl
	return DeleteTenantIDConfigResponse{}, errors.New("needs implementation")
}

type ListMultitenantConfigResponse struct {
	OK *struct {
		Configs []struct {
			ThirdPartyId string
			TenantId     string
			Config       TenantConfig
		}
	}
}

func ListMultitenantConfigMapping(thirdPartyProviderId *string) (ListMultitenantConfigResponse, error) {
	// TODO impl
	return ListMultitenantConfigResponse{}, errors.New("needs implementation")
}
