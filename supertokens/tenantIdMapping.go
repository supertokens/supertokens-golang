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

	AuthorizationEndpoint            string                 `json:"authorizationEndpoint"`
	AuthorizationEndpointQueryParams map[string]interface{} `json:"authorizationEndpointQueryParams"`
	TokenEndpoint                    string                 `json:"tokenEndpoint"`
	TokenParams                      map[string]interface{} `json:"tokenParams"`
	ForcePKCE                        bool                   `json:"forcePKCE"`
	UserInfoEndpoint                 string                 `json:"userInfoEndpoint"`
	JwksURI                          string                 `json:"jwksURI"`
	OIDCDiscoveryEndpoint            string                 `json:"oidcDiscoveryEndpoint"`

	UserInfoMap struct {
		From               string `json:"from"`
		IdField            string `json:"idField"`
		EmailField         string `json:"emailField"`
		EmailVerifiedField string `json:"emailVerifiedField"`
	} `json:"userInfoMap"`

	Clients []TenantClient `json:"clients"`

	Frontend struct {
		Name            string `json:"name"`
		ButtonStyle     string `json:"buttonStyle"`
		ButtonComponent string `json:"buttonComponent"`
	} `json:"frontend"`
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
