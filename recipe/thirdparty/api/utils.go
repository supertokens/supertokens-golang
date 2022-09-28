package api

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func findProvider(options tpmodels.APIOptions, thirdPartyId string) (*tpmodels.TypeProvider, error) {
	providers := options.Providers

	for i := 0; i < len(providers); i++ {
		id := providers[i].ID

		if id != thirdPartyId {
			continue
		}

		// first if there is only one provider with thirdPartyId in the providers array,
		var otherProvidersWithSameId []tpmodels.TypeProvider = []tpmodels.TypeProvider{}
		for y := 0; y < len(providers); y++ {
			if providers[y].ID == id && &providers[y] != &providers[i] {
				otherProvidersWithSameId = append(otherProvidersWithSameId, providers[y])
			}
		}
		if len(otherProvidersWithSameId) == 0 {
			return &providers[i], nil
		}
	}

	return nil, supertokens.BadInputError{Msg: "The third party provider " + thirdPartyId + " seems to be missing from the backend configs. If it is configured, then please make sure that you are passing the correct clientId from the frontend."}
}
