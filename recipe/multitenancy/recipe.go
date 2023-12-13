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

package multitenancy

import (
	"errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/multitenancy/api"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancyclaims"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "multitenancy"

type Recipe struct {
	RecipeModule              supertokens.RecipeModule
	Config                    multitenancymodels.TypeNormalisedInput
	RecipeImpl                multitenancymodels.RecipeInterface
	APIImpl                   multitenancymodels.APIInterface
	staticThirdPartyProviders []tpmodels.ProviderInput

	GetAllowedDomainsForTenantId func(tenantId string, userContext supertokens.UserContext) ([]string, error)
}

var singletonInstance *Recipe

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *multitenancymodels.TypeInput, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (*Recipe, error) {
	r := &Recipe{}
	verifiedConfig := validateAndNormaliseUserInput(config)
	r.Config = verifiedConfig

	querierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
	if err != nil {
		return nil, err
	}
	recipeImplementation := makeRecipeImplementation(*querierInstance, verifiedConfig, appInfo)
	r.RecipeImpl = verifiedConfig.Override.Functions(recipeImplementation)

	r.APIImpl = verifiedConfig.Override.APIs(api.MakeAPIImplementation())

	recipeModuleInstance := supertokens.MakeRecipeModule(recipeId, appInfo, r.handleAPIRequest, r.getAllCORSHeaders, r.getAPIsHandled, nil, r.handleError, onSuperTokensAPIError)
	r.RecipeModule = recipeModuleInstance

	r.staticThirdPartyProviders = []tpmodels.ProviderInput{}

	r.GetAllowedDomainsForTenantId = verifiedConfig.GetAllowedDomainsForTenantId

	return r, nil
}

func GetRecipeInstanceOrThrowError() (*Recipe, error) {
	if singletonInstance != nil {
		return singletonInstance, nil
	}

	return nil, errors.New("Initialisation not done. Did you forget to call the init function?")
}

func GetRecipeInstance() *Recipe {
	return singletonInstance
}

func recipeInit(config *multitenancymodels.TypeInput) supertokens.Recipe {
	return func(appInfo supertokens.NormalisedAppinfo, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (*supertokens.RecipeModule, error) {
		if singletonInstance == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config, onSuperTokensAPIError)
			if err != nil {
				return nil, err
			}

			if recipe.GetAllowedDomainsForTenantId != nil {
				supertokens.AddPostInitCallback(func() error {
					sessionRecipe, err := session.GetRecipeInstanceOrThrowError()

					if err != nil {
						return nil // skip adding claims if session recipe is not initialised
					}

					sessionRecipe.AddClaimFromOtherRecipe(multitenancyclaims.AllowedDomainsClaim)

					return nil
				})
			}

			singletonInstance = recipe
			return &singletonInstance.RecipeModule, nil
		}
		return nil, errors.New("Multitenancy recipe has already been initialised. Please check your code for bugs.")
	}
}

// implement RecipeModule

func (r *Recipe) getAPIsHandled() ([]supertokens.APIHandled, error) {
	loginMethodsAPI, err := supertokens.NewNormalisedURLPath(LoginMethodsAPI)
	if err != nil {
		return nil, err
	}

	return []supertokens.APIHandled{
		{
			Method:                 http.MethodGet,
			PathWithoutAPIBasePath: loginMethodsAPI,
			ID:                     LoginMethodsAPI,
			Disabled:               r.APIImpl.LoginMethodsGET == nil,
		},
	}, nil
}

func (r *Recipe) handleAPIRequest(id string, tenantId string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, _ supertokens.NormalisedURLPath, _ string, userContext supertokens.UserContext) error {
	options := multitenancymodels.APIOptions{
		RecipeImplementation:      r.RecipeImpl,
		Config:                    r.Config,
		RecipeID:                  RECIPE_ID,
		Req:                       req,
		Res:                       res,
		OtherHandler:              theirHandler,
		StaticThirdPartyProviders: r.staticThirdPartyProviders,
	}
	if id == LoginMethodsAPI {
		return api.LoginMethodsAPI(r.APIImpl, tenantId, options, userContext)
	}
	return errors.New("should never come here")
}

func (r *Recipe) getAllCORSHeaders() []string {
	return []string{}
}

func (r *Recipe) handleError(err error, req *http.Request, res http.ResponseWriter, userContext supertokens.UserContext) (bool, error) {
	return false, nil
}

func ResetForTest() {
	singletonInstance = nil
}

func (r *Recipe) SetStaticThirdPartyProviders(providers []tpmodels.ProviderInput) {
	// the `staticThirdPartyProviders` is always overwritten with the provider list from
	// the last thirdparty recipe that was initialised. In case of multitenancy, the
	// core is expected to contain the tenant configs for that tenant, so, it wouldn't
	// be too much of a problem about which static list we use from here.
	r.staticThirdPartyProviders = append([]tpmodels.ProviderInput{}, providers...)
}

// This function is called when the multitenancy package is imported. This is set so that
// the supertokens Init can create an instance of the multitenancy recipe automatically
// if the user has not explicitly created one.
func init() {
	supertokens.InternalUseDefaultMultitenancyRecipe = recipeInit(nil)

	// Create multitenancy claims when the module is imported
	multitenancyclaims.AllowedDomainsClaim, multitenancyclaims.AllowedDomainsClaimValidators = NewAllowedDomainsClaim()

	supertokens.GetTenantIdFuncFromUsingMultitenancyRecipe = func(tenantIdFromFrontend string, userContext supertokens.UserContext) (string, error) {
		mtRecipe := GetRecipeInstance()
		return (*mtRecipe.RecipeImpl.GetTenantId)(tenantIdFromFrontend, userContext)
	}
}
