package tpmodels

import (
	"github.com/supertokens/supertokens-golang/supertokens"
)

type TypeGoogleInput struct {
	Config   []GoogleConfig
	Override func(provider *GoogleProvider) *GoogleProvider
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
}

type GoogleProvider struct {
	GetConfig func(clientID *string, userContext supertokens.UserContext) (GoogleConfig, error)
	TypeProvider
}
