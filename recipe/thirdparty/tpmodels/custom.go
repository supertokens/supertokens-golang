package tpmodels

import (
	"github.com/supertokens/supertokens-golang/supertokens"
)

type TypeCustomProviderBuilderInput struct {
	ID                            string
	AuthorizationURL              string
	ExchangeCodeForOAuthTokensURL string
	GetUserInfoURL                string
	DefaultScopes                 []string
}

type CustomProviderFunc func(input TypeCustomProviderInput) (TypeProvider, error)

type TypeCustomProviderInput struct {
	Config   []CustomProviderConfig
	Override func(provider *CustomProvider) *CustomProvider
}

type CustomProviderConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
	ExtraConfig  map[string]interface{}
}

type CustomProvider struct {
	GetConfig func(clientID *string, userContext supertokens.UserContext) (CustomProviderConfig, error)
	*TypeProvider
}
