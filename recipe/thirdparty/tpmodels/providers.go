/* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package tpmodels

import (
	"github.com/supertokens/supertokens-golang/supertokens"
)

type TypeGoogleInput struct {
	Config   []GoogleConfig
	Override func(provider GoogleProvider) GoogleProvider
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
}

type GoogleProvider struct {
	GetConfig func(clientID *string, userContext supertokens.UserContext) (GoogleConfig, error)

	GetAuthorisationRedirectURL    func(clientID *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (TypeAuthorisationRedirect, error)
	ExchangeAuthCodeForOAuthTokens func(clientID *string, redirectInfo TypeRedirectURIInfo, userContext supertokens.UserContext) (TypeOAuthTokens, error) // For apple, add userInfo from callbackInfo to oAuthTOkens
	GetUserInfo                    func(clientID *string, oAuthTokens TypeOAuthTokens, userContext supertokens.UserContext) (TypeUserInfo, error)
}

type TypeGoogleWorkspacesInput struct {
	Config   []GoogleWorkspacesConfig
	Override func(provider GoogleWorkspacesProvider) GoogleWorkspacesProvider
}

type GoogleWorkspacesConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
	Domain       *string
}

type GoogleWorkspacesProvider struct {
	GetConfig func(clientID *string, userContext supertokens.UserContext) (GoogleWorkspacesConfig, error)

	GetAuthorisationRedirectURL    func(clientID *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (TypeAuthorisationRedirect, error)
	ExchangeAuthCodeForOAuthTokens func(clientID *string, redirectInfo TypeRedirectURIInfo, userContext supertokens.UserContext) (TypeOAuthTokens, error) // For apple, add userInfo from callbackInfo to oAuthTOkens
	GetUserInfo                    func(clientID *string, oAuthTokens TypeOAuthTokens, userContext supertokens.UserContext) (TypeUserInfo, error)
}

type TypeGithubInput struct {
	Config   []GithubConfig
	Override func(provider GithubProvider) GithubProvider
}

type GithubConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
}

type GithubProvider struct {
	GetConfig func(clientID *string, userContext supertokens.UserContext) (GithubConfig, error)

	GetAuthorisationRedirectURL    func(clientID *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (TypeAuthorisationRedirect, error)
	ExchangeAuthCodeForOAuthTokens func(clientID *string, redirectInfo TypeRedirectURIInfo, userContext supertokens.UserContext) (TypeOAuthTokens, error) // For apple, add userInfo from callbackInfo to oAuthTOkens
	GetUserInfo                    func(clientID *string, oAuthTokens TypeOAuthTokens, userContext supertokens.UserContext) (TypeUserInfo, error)
}

type TypeDiscordInput struct {
	Config   []DiscordConfig
	Override func(provider DiscordProvider) DiscordProvider
}

type DiscordConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
}

type DiscordProvider struct {
	GetConfig func(clientID *string, userContext supertokens.UserContext) (DiscordConfig, error)

	GetAuthorisationRedirectURL    func(clientID *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (TypeAuthorisationRedirect, error)
	ExchangeAuthCodeForOAuthTokens func(clientID *string, redirectInfo TypeRedirectURIInfo, userContext supertokens.UserContext) (TypeOAuthTokens, error) // For apple, add userInfo from callbackInfo to oAuthTOkens
	GetUserInfo                    func(clientID *string, oAuthTokens TypeOAuthTokens, userContext supertokens.UserContext) (TypeUserInfo, error)
}

type TypeFacebookInput struct {
	Config   []FacebookConfig
	Override func(provider FacebookProvider) FacebookProvider
}

type FacebookConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
}

type FacebookProvider struct {
	GetConfig func(clientID *string, userContext supertokens.UserContext) (FacebookConfig, error)

	GetAuthorisationRedirectURL    func(clientID *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (TypeAuthorisationRedirect, error)
	ExchangeAuthCodeForOAuthTokens func(clientID *string, redirectInfo TypeRedirectURIInfo, userContext supertokens.UserContext) (TypeOAuthTokens, error) // For apple, add userInfo from callbackInfo to oAuthTOkens
	GetUserInfo                    func(clientID *string, oAuthTokens TypeOAuthTokens, userContext supertokens.UserContext) (TypeUserInfo, error)
}

type TypeAppleInput struct {
	Config   []AppleConfig
	Override func(provider AppleProvider) AppleProvider
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

	GetAuthorisationRedirectURL    func(clientID *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (TypeAuthorisationRedirect, error)
	ExchangeAuthCodeForOAuthTokens func(clientID *string, redirectInfo TypeRedirectURIInfo, userContext supertokens.UserContext) (TypeOAuthTokens, error) // For apple, add userInfo from callbackInfo to oAuthTOkens
	GetUserInfo                    func(clientID *string, oAuthTokens TypeOAuthTokens, userContext supertokens.UserContext) (TypeUserInfo, error)
}

type TypeOktaInput struct {
	Config   []OktaConfig
	Override func(provider OktaProvider) OktaProvider
}

type OktaConfig struct {
	ClientID     string
	ClientSecret string
	OktaDomain   string
	Scope        []string
}

type OktaProvider struct {
	GetConfig func(clientID *string, userContext supertokens.UserContext) (OktaConfig, error)

	GetAuthorisationRedirectURL    func(clientID *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (TypeAuthorisationRedirect, error)
	ExchangeAuthCodeForOAuthTokens func(clientID *string, callbackInfo TypeRedirectURIInfo, userContext supertokens.UserContext) (TypeOAuthTokens, error)
	GetUserInfo                    func(clientID *string, oAuthTokens TypeOAuthTokens, userContext supertokens.UserContext) (TypeUserInfo, error)
}
