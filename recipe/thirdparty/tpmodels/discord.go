package tpmodels

import (
	"github.com/supertokens/supertokens-golang/supertokens"
)

type TypeDiscordInput struct {
	Config   []DiscordConfig
	Override func(provider *DiscordProvider) *DiscordProvider
}

type DiscordConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
}

type DiscordProvider struct {
	GetConfig func(clientID *string, userContext supertokens.UserContext) (DiscordConfig, error)
	*TypeProvider
}
