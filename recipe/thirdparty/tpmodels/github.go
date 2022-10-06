package tpmodels

import (
	"github.com/supertokens/supertokens-golang/supertokens"
)

type TypeGithubInput struct {
	Config   []GithubConfig
	Override func(provider *GithubProvider) *GithubProvider
}

type GithubConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
}

type GithubProvider struct {
	GetConfig func(clientID *string, userContext supertokens.UserContext) (GithubConfig, error)
	*TypeProvider
}
