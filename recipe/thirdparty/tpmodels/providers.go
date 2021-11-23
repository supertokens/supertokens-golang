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

type GoogleConfig struct {
	ClientID              string
	ClientSecret          string
	Scope                 []string
	AuthorisationRedirect *struct {
		Params map[string]interface{}
	}
	IsDefault bool
}

type GoogleWorkspacesConfig struct {
	ClientID              string
	ClientSecret          string
	Scope                 []string
	Domain                *string
	AuthorisationRedirect *struct {
		Params map[string]interface{}
	}
	IsDefault bool
}

type GithubConfig struct {
	ClientID              string
	ClientSecret          string
	Scope                 []string
	AuthorisationRedirect *struct {
		Params map[string]interface{}
	}
	IsDefault bool
}

type DiscordConfig struct {
	ClientID              string
	ClientSecret          string
	Scope                 []string
	AuthorisationRedirect *struct {
		Params map[string]interface{}
	}
	IsDefault bool
}

type FacebookConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
	IsDefault    bool
}

type AppleConfig struct {
	ClientID              string
	ClientSecret          AppleClientSecret
	Scope                 []string
	AuthorisationRedirect *struct {
		Params map[string]interface{}
	}
	IsDefault bool
}

type AppleClientSecret struct {
	KeyId      string
	PrivateKey string
	TeamId     string
}
