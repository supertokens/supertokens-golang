package tpmodels

import (
	"github.com/supertokens/supertokens-golang/supertokens"
)

type GetAuthorisationRedirectURLFunc func(clientType *string, tenantId *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (TypeAuthorisationRedirect, error)
type ExchangeAuthCodeForOAuthTokensFunc func(clientType *string, tenantId *string, redirectInfo TypeRedirectURIInfo, userContext supertokens.UserContext) (TypeOAuthTokens, error) // For apple, add userInfo from callbackInfo to oAuthTOkens
type GetUserInfoFunc func(clientType *string, tenantId *string, oAuthTokens TypeOAuthTokens, userContext supertokens.UserContext) (TypeUserInfo, error)

type TypeProviderInterface interface {
	GetId() string
	GetAuthorisationRedirectURL(clientType *string, tenantId *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (TypeAuthorisationRedirect, error)
	ExchangeAuthCodeForOAuthTokens(clientType *string, tenantId *string, redirectInfo TypeRedirectURIInfo, userContext supertokens.UserContext) (TypeOAuthTokens, error) // For apple, add userInfo from callbackInfo to oAuthTOkens
	GetUserInfo(clientType *string, tenantId *string, oAuthTokens TypeOAuthTokens, userContext supertokens.UserContext) (TypeUserInfo, error)
}

type TypeProvider struct {
	ID string

	GetAuthorisationRedirectURLImpl    GetAuthorisationRedirectURLFunc
	ExchangeAuthCodeForOAuthTokensImpl ExchangeAuthCodeForOAuthTokensFunc
	GetUserInfoImpl                    GetUserInfoFunc
}

func (p TypeProvider) GetId() string {
	return p.ID
}

func (p TypeProvider) GetAuthorisationRedirectURL(clientType *string, tenantId *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (TypeAuthorisationRedirect, error) {
	return p.GetAuthorisationRedirectURLImpl(clientType, tenantId, redirectURIOnProviderDashboard, userContext)
}

func (p TypeProvider) ExchangeAuthCodeForOAuthTokens(clientType *string, tenantId *string, redirectInfo TypeRedirectURIInfo, userContext supertokens.UserContext) (TypeOAuthTokens, error) {
	return p.ExchangeAuthCodeForOAuthTokensImpl(clientType, tenantId, redirectInfo, userContext)
}
