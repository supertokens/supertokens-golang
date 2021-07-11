package providers

import "github.com/supertokens/supertokens-golang/recipe/thirdparty/models"

type TypeThirdPartyProviderAppleConfig struct {
	ClientID                   string
	ClientSecret                ClientSecret
	Scope                       []string
	AuthorisationRedirectParams map[string]string
}

type ClientSecret struct {
	KeyId      string
	PrivateKey string
	TeamId     string
}

const AppleID = "apple"

func AppleGetClientSecret(clientId string, keyId string, teamId string, privateKey string) string {
	return ""
}

func AppleGet(redirectURI string, authCodeFromRequest string) models.TypeProviderGetResponse {
	return models.TypeProviderGetResponse{}
}
