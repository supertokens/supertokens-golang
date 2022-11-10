package api

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func findProvider(options tpmodels.APIOptions, thirdPartyId string, tenantId *string) (tpmodels.TypeProvider, error) {

	for _, provider := range options.Providers {
		if provider.ID == thirdPartyId {
			return provider, nil
		}
	}

	if tenantId == nil {
		return tpmodels.TypeProvider{}, supertokens.BadInputError{Msg: "The third party provider " + thirdPartyId + " seems to be missing from the backend configs."}
	}

	// If tenantId is not nil, we need to create the provider on the fly,
	// so that the GetConfig function will make use of the core to fetch the config
	return createProvider(thirdPartyId), nil
}

func createProvider(thirdPartyId string) tpmodels.TypeProvider {
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

func createCustomProvider(thirdPartyId string) tpmodels.TypeProvider {
	// TODO impl
	return tpmodels.TypeProvider{}
}
