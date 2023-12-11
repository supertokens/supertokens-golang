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

package supertokens

import (
	"net/http"
)

type NormalisedAppinfo struct {
	AppName                  string
	GetOrigin                func(request *http.Request, userContext UserContext) (NormalisedURLDomain, error)
	APIDomain                NormalisedURLDomain
	TopLevelAPIDomain        string
	GetTopLevelWebsiteDomain func(request *http.Request, userContext UserContext) (string, error)
	APIBasePath              NormalisedURLPath
	APIGatewayPath           NormalisedURLPath
	WebsiteBasePath          NormalisedURLPath
}

type AppInfo struct {
	AppName         string
	WebsiteDomain   string
	Origin          string
	GetOrigin       func(request *http.Request, userContext UserContext) (string, error)
	APIDomain       string
	WebsiteBasePath *string
	APIBasePath     *string
	APIGatewayPath  *string
}

type Recipe func(appInfo NormalisedAppinfo, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (*RecipeModule, error)

type TypeInput struct {
	Supertokens           *ConnectionInfo
	AppInfo               AppInfo
	RecipeList            []Recipe
	Telemetry             *bool
	Debug                 bool
	OnSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)
}

type ConnectionInfo struct {
	ConnectionURI      string
	APIKey             string
	NetworkInterceptor func(*http.Request, UserContext) *http.Request
}

type APIHandled struct {
	PathWithoutAPIBasePath NormalisedURLPath
	Method                 string
	ID                     string
	Disabled               bool
}

type UserContext = *map[string]interface{}

type GeneralErrorResponse struct {
	Message string
}
type ThirdParty struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
}

type RecipeID string

const (
	EmailPasswordRID RecipeID = "emailpassword"
	ThirdPartyRID    RecipeID = "thirdparty"
	PasswordlessRID  RecipeID = "passwordless"
)

type LoginMethods struct {
	RecipeLevelUser
	Verified                bool                              `json:"verified"`
	HasSameEmailAs          func(email string) bool           `json:"-"`
	HasSamePhoneNumberAs    func(phoneNumber string) bool     `json:"-"`
	HasSameThirdPartyInfoAs func(thirdParty *ThirdParty) bool `json:"-"`
}

type User struct {
	ID            string         `json:"id"`
	TimeJoined    int64          `json:"timeJoined"`
	IsPrimaryUser bool           `json:"isPrimaryUser"`
	TenantIDs     []string       `json:"tenantIds"`
	Emails        []string       `json:"emails"`
	PhoneNumbers  []string       `json:"phoneNumbers"`
	ThirdParty    []ThirdParty   `json:"thirdParty"`
	LoginMethods  []LoginMethods `json:"loginMethods"`
}

type AccountInfo struct {
	Email       *string     `json:"email,omitempty"`
	PhoneNumber *string     `json:"phoneNumber,omitempty"`
	ThirdParty  *ThirdParty `json:"thirdParty,omitempty"`
}

type AccountInfoWithRecipeID struct {
	RecipeID RecipeID `json:"recipeId"`
	AccountInfo
}

type RecipeLevelUser struct {
	TenantIDs    []string     `json:"tenantIds"`
	TimeJoined   int64        `json:"timeJoined"`
	RecipeUserID RecipeUserID `json:"recipeUserId"`
	AccountInfoWithRecipeID
}

type AccountInfoWithRecipeIdAndWithRecipeUserId struct {
	RecipeUserId *RecipeUserID
	AccountInfoWithRecipeID
}

type ShouldDoAutomaticAccountLinkingResponse struct {
	ShouldAutomaticallyLink   bool
	ShouldRequireVerification bool
}

type AccountLinkingTypeInput struct {
	OnAccountLinked                 func(user User, newAccountUser RecipeLevelUser, userContext UserContext) error
	ShouldDoAutomaticAccountLinking func(newAccountInfo AccountInfoWithRecipeIdAndWithRecipeUserId, user *User, tenantID string, userContext UserContext) (ShouldDoAutomaticAccountLinkingResponse, error)
	Override                        *AccountLinkingOverrideStruct
}

type AccountLinkingTypeNormalisedInput struct {
	OnAccountLinked                 func(user User, newAccountUser RecipeLevelUser, userContext UserContext) error
	ShouldDoAutomaticAccountLinking func(newAccountInfo AccountInfoWithRecipeIdAndWithRecipeUserId, user *User, tenantID string, userContext UserContext) (ShouldDoAutomaticAccountLinkingResponse, error)
	Override                        AccountLinkingOverrideStruct
}

type AccountLinkingOverrideStruct struct {
	Functions func(originalImplementation AccountLinkingRecipeInterface) AccountLinkingRecipeInterface
}
