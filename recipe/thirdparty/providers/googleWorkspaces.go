package providers

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
)

func GoogleWorkspaces(input tpmodels.ProviderInput) tpmodels.TypeProvider {
	input.ThirdPartyID = "google-workspaces"

	input.Config.ValidateIdTokenPayload = func(idTokenPayload map[string]interface{}, clientConfig tpmodels.ProviderClientConfig) (bool, error) {
		return idTokenPayload["hd"] == clientConfig.AdditionalConfig["domain"], nil
	}

	return Google(input)
}
