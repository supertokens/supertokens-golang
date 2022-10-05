package tpmodels

import (
	"github.com/supertokens/supertokens-golang/supertokens"
)

type TypeAppleInput struct {
	Config   []AppleConfig
	Override func(provider *AppleProvider) *AppleProvider
}

type AppleConfig struct {
	ClientID     string
	ClientSecret AppleClientSecret
	Scope        []string
}

type AppleClientSecret struct {
	KeyId      string
	PrivateKey string
	TeamId     string
}

type AppleProvider struct {
	GetConfig func(clientID *string, userContext supertokens.UserContext) (AppleConfig, error)
	*TypeProvider
}
