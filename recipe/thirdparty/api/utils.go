package api

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func findProvider(options tpmodels.APIOptions, thirdPartyId string, tenantId *string) (tpmodels.TypeProviderInterface, error) {
	if tenantId == nil {
		for _, provider := range options.Providers {
			if provider.GetId() == thirdPartyId {
				return provider, nil
			}
		}
		return nil, supertokens.BadInputError{Msg: "The third party provider " + thirdPartyId + " seems to be missing from the backend configs."}
	}

	var definedProvider tpmodels.TypeProviderInterface = nil
	for _, provider := range options.Providers {
		if provider.GetId() == thirdPartyId {
			definedProvider = provider
		}
	}

	result, err := supertokens.FetchTenantIDConfigMapping(thirdPartyId, *tenantId)
	if err != nil {
		return nil, err
	}

	if result.UnknownMappingError != nil {
		return nil, supertokens.BadInputError{Msg: "The tenantId " + *tenantId + " seems to be missing from the backend configs."}
	}

	if definedProvider == nil {
		definedProvider = createProvider(thirdPartyId, result.OK.Config)
		return definedProvider, nil
	}

	// TODO
	// for _, client := range result.OK.Config.Clients {
	//  check if it's a built-in provider and if it is, add or update config by clientType
	// 	definedProvider.AddOrUpdateClient(client)
	// }

	return nil, errors.New("needs implementation")
}

func createProvider(thirdPartyId string, config interface{}) tpmodels.TypeProviderInterface {
	// TODO impl
	switch thirdPartyId {
	case "active-directory":
	case "apple":
	case "discord":
	case "facebook":
	case "github":
	case "google":
	case "google-workspaces":
	case "okta":
	}
	return createCustomProvider(thirdPartyId, config)
}

func createCustomProvider(thirdPartyId string, config interface{}) tpmodels.TypeProviderInterface {
	// TODO impl
	return nil
}
