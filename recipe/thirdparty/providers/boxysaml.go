package providers

import (
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
)

const boxySamlID = "boxy-saml"

func BoxySaml(input tpmodels.ProviderInput) *tpmodels.TypeProvider {
	if input.Config.ThirdPartyId == "" {
		input.Config.ThirdPartyId = boxySamlID
	}
	if input.Config.Name == "" {
		input.Config.Name = "Boxy SAML"
	}

	if input.Config.UserInfoMap.FromUserInfoAPI.UserId == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.UserId = "id"
	}
	if input.Config.UserInfoMap.FromUserInfoAPI.Email == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.Email = "email"
	}

	return NewProvider(input)
}
