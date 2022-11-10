package api

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func findProvider(options tpmodels.APIOptions, thirdPartyId string, tenantId *string) (tpmodels.TypeProvider, error) {

	for _, provider := range options.Providers {
		if provider.GetID() == thirdPartyId {
			return provider.Build(), nil
		}
	}

	if tenantId == nil {
		return tpmodels.TypeProvider{}, supertokens.BadInputError{Msg: "The third party provider " + thirdPartyId + " seems to be missing from the backend configs."}
	}

	// If tenantId is not nil, we need to create the provider on the fly,
	// so that the GetConfig function will make use of the core to fetch the config
	newProvider := createProvider(thirdPartyId)
	return newProvider.Build(), nil
}

func createProvider(thirdPartyId string) tpmodels.TypeProviderInterface {
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
	return createCustomProvider(thirdPartyId)
}

func createCustomProvider(thirdPartyId string) tpmodels.TypeProviderInterface {
	// TODO impl
	return nil
}
