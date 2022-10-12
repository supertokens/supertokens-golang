package tpmodels

import (
	"github.com/supertokens/supertokens-golang/supertokens"
)

type TypeGoogleWorkspacesInput struct {
	Config   []GoogleWorkspacesConfig
	Override func(provider *GoogleWorkspacesProvider) *GoogleWorkspacesProvider
}

type GoogleWorkspacesConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
	Domain       *string
}

type GoogleWorkspacesProvider struct {
	GetConfig   func(clientID *string, userContext supertokens.UserContext) (GoogleWorkspacesConfig, error)
	GetTenantID func(clientID *string, userContext supertokens.UserContext) (string, error)
	*TypeProvider
}
