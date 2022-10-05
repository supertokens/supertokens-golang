package tpmodels

import (
	"github.com/supertokens/supertokens-golang/supertokens"
)

type TypeFacebookInput struct {
	Config   []FacebookConfig
	Override func(provider *FacebookProvider) *FacebookProvider
}

type FacebookConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
}

type FacebookProvider struct {
	GetConfig func(clientID *string, userContext supertokens.UserContext) (FacebookConfig, error)
	*TypeProvider
}
