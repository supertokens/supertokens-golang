package providers

import "github.com/supertokens/supertokens-golang/recipe/thirdparty/models"

type TypeThirdPartyProviderGoogleConfig struct {
	ClientID                    string
	ClientSecret                string
	Scope                       []string
	AuthorisationRedirectParams map[string]string
}

const GoogleID = "google"

func GoogleGet(redirectURI string, authCodeFromRequest string) models.TypeProviderGetResponse {
	return models.TypeProviderGetResponse{}
}
