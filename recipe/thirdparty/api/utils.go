package api

import "github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"

func findRightProvider(providers []tpmodels.TypeProvider, thirdPartyId string, clientId *string) *tpmodels.TypeProvider {

	/* logic to find the right provider

	ClientID not passed
		Case 1: Single Config - Return the config
		Case 2: Multiple Config - Return the config with IsDefault = true

	ClientID passed
		Case 1: Single Config - Return the config irrespective of the clientID Match
		Case 2: Multiple Config
			Return the config with the matching clientID else return config with IsDefault = true

	*/

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
			return &providers[i]
		}

		if clientId == nil && providers[i].IsDefault {
			return &providers[i]
		}

		if clientId != nil && *clientId == providers[i].Get(nil, nil, &map[string]interface{}{}).GetClientId(&map[string]interface{}{}) {
			return &providers[i]
		}
	}

	return nil
}
