package providers

import "github.com/supertokens/supertokens-golang/recipe/thirdparty/models"

type TypeThirdPartyProviderGithubConfig struct {
	ClientID                    string
	ClientSecret                string
	Scope                       []string
	AuthorisationRedirectParams map[string]string
}

const GithubID = "github"

func GithubGet(redirectURI string, authCodeFromRequest string) models.TypeProviderGetResponse {
	return models.TypeProviderGetResponse{}
}
