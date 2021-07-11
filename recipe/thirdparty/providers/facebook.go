package providers

import "github.com/supertokens/supertokens-golang/recipe/thirdparty/models"

type TypeThirdPartyProviderFacebookConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
}

const FacebookID = "facebook"

func FacebookGet(redirectURI string, authCodeFromRequest string) models.TypeProviderGetResponse {
	return models.TypeProviderGetResponse{}
}
