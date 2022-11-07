package api

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func findProvider(options tpmodels.APIOptions, thirdPartyId string) (tpmodels.TypeProvider, error) {
	for _, provider := range options.Providers {
		if provider.ID == thirdPartyId {
			return provider, nil
		}
	}
	return tpmodels.TypeProvider{}, supertokens.BadInputError{Msg: "The third party provider " + thirdPartyId + " seems to be missing from the backend configs."}
}
