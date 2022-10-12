package tpmodels

import (
	"github.com/supertokens/supertokens-golang/supertokens"
)

type TypeOktaInput struct {
	Config   []OktaConfig
	Override func(provider *OktaProvider) *OktaProvider
}

type OktaConfig struct {
	ClientID     string
	ClientSecret string
	OktaDomain   string
	Scope        []string
}

type OktaProvider struct {
	GetConfig   func(clientID *string, userContext supertokens.UserContext) (OktaConfig, error)
	GetTenantID func(clientID *string, userContext supertokens.UserContext) (string, error)
	TypeProvider
}
