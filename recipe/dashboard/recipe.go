/* Copyright (c) 2022, VRAI Labs and/or its affiliates. All rights reserved.
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

package dashboard

import (
	"errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/dashboard/api"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/api/search"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/api/userdetails"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/constants"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "dashboard"

type Recipe struct {
	RecipeModule supertokens.RecipeModule
	Config       dashboardmodels.TypeNormalisedInput
	RecipeImpl   dashboardmodels.RecipeInterface
	APIImpl      dashboardmodels.APIInterface
}

var singletonInstance *Recipe

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *dashboardmodels.TypeInput, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (Recipe, error) {
	r := &Recipe{}
	verifiedConfig := validateAndNormaliseUserInput(appInfo, config)
	r.Config = verifiedConfig

	querierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
	if err != nil {
		return Recipe{}, err
	}

	recipeImplementation := makeRecipeImplementation(*querierInstance)
	r.RecipeImpl = verifiedConfig.Override.Functions(recipeImplementation)

	r.APIImpl = verifiedConfig.Override.APIs(api.MakeAPIImplementation())

	recipeModuleInstance := supertokens.MakeRecipeModule(recipeId, appInfo, r.handleAPIRequest, r.getAllCORSHeaders, r.getAPIsHandled, nil, r.handleError, onSuperTokensAPIError)
	r.RecipeModule = recipeModuleInstance

	return *r, nil
}

func recipeInit(config *dashboardmodels.TypeInput) supertokens.Recipe {
	return func(appInfo supertokens.NormalisedAppinfo, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (*supertokens.RecipeModule, error) {
		if singletonInstance == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config, onSuperTokensAPIError)
			if err != nil {
				return nil, err
			}
			singletonInstance = &recipe
			return &singletonInstance.RecipeModule, nil
		}
		return nil, errors.New("Dashboard recipe has already been initialised. Please check your code for bugs.")
	}
}

func (r *Recipe) getAPIsHandled() ([]supertokens.APIHandled, error) {
	dashboardAPI, err := supertokens.NewNormalisedURLPath(constants.DashboardAPI)
	if err != nil {
		return nil, err
	}
	dashboardApiBasePath, err := supertokens.NewNormalisedURLPath(constants.DashboardAPIBasePath)
	if err != nil {
		return nil, err
	}
	signInAPI, err := supertokens.NewNormalisedURLPath(constants.SignInAPI)
	if err != nil {
		return nil, err
	}
	signOutAPI, err := supertokens.NewNormalisedURLPath(constants.SignOutAPI)
	if err != nil {
		return nil, err
	}
	validateKeyAPI, err := supertokens.NewNormalisedURLPath(constants.ValidateKeyAPI)
	if err != nil {
		return nil, err
	}
	usersListGetAPI, err := supertokens.NewNormalisedURLPath(constants.UsersListGetAPI)
	if err != nil {
		return nil, err
	}
	usersCountAPI, err := supertokens.NewNormalisedURLPath(constants.UsersCountAPI)
	if err != nil {
		return nil, err
	}
	userAPI, err := supertokens.NewNormalisedURLPath(constants.UserAPI)
	if err != nil {
		return nil, err
	}
	userEmailVerifyAPI, err := supertokens.NewNormalisedURLPath(constants.UserEmailVerifyAPI)
	if err != nil {
		return nil, err
	}
	userMetaDataAPI, err := supertokens.NewNormalisedURLPath(constants.UserMetadataAPI)
	if err != nil {
		return nil, err
	}
	userSessionsAPI, err := supertokens.NewNormalisedURLPath(constants.UserSessionsAPI)
	if err != nil {
		return nil, err
	}
	userPasswordAPI, err := supertokens.NewNormalisedURLPath(constants.UserPasswordAPI)
	if err != nil {
		return nil, err
	}
	userEmailVerifyTokenAPI, err := supertokens.NewNormalisedURLPath(constants.UserEmailVerifyTokenAPI)
	if err != nil {
		return nil, err
	}
	searchTagsAPI, err := supertokens.NewNormalisedURLPath(constants.SearchTagsAPI)
	if err != nil {
		return nil, err
	}
	dashboardAnalyticsAPI, err := supertokens.NewNormalisedURLPath(constants.DashboardAnalyticsAPI)
	if err != nil {
		return nil, err
	}
	tenantsListAPI, err := supertokens.NewNormalisedURLPath(constants.TenantsListAPI)
	if err != nil {
		return nil, err
	}

	return []supertokens.APIHandled{
		{
			ID:                     constants.DashboardAPI,
			PathWithoutAPIBasePath: dashboardAPI,
			Method:                 http.MethodGet,
			Disabled:               false,
		},
		{
			ID:                     constants.SignInAPI,
			PathWithoutAPIBasePath: dashboardApiBasePath.AppendPath(signInAPI),
			Method:                 http.MethodPost,
			Disabled:               false,
		},
		{
			ID:                     constants.SignOutAPI,
			PathWithoutAPIBasePath: dashboardApiBasePath.AppendPath(signOutAPI),
			Method:                 http.MethodPost,
			Disabled:               false,
		},
		{
			ID:                     constants.ValidateKeyAPI,
			PathWithoutAPIBasePath: dashboardApiBasePath.AppendPath(validateKeyAPI),
			Method:                 http.MethodPost,
			Disabled:               false,
		},
		{
			ID:                     constants.UsersListGetAPI,
			PathWithoutAPIBasePath: dashboardApiBasePath.AppendPath(usersListGetAPI),
			Method:                 http.MethodGet,
			Disabled:               false,
		},
		{
			ID:                     constants.UsersCountAPI,
			PathWithoutAPIBasePath: dashboardApiBasePath.AppendPath(usersCountAPI),
			Method:                 http.MethodGet,
			Disabled:               false,
		},
		{
			ID:                     constants.UserAPI,
			PathWithoutAPIBasePath: dashboardApiBasePath.AppendPath(userAPI),
			Method:                 http.MethodGet,
			Disabled:               false,
		},
		{
			ID:                     constants.UserAPI,
			PathWithoutAPIBasePath: dashboardApiBasePath.AppendPath(userAPI),
			Method:                 http.MethodPost,
			Disabled:               false,
		},
		{
			ID:                     constants.UserAPI,
			PathWithoutAPIBasePath: dashboardApiBasePath.AppendPath(userAPI),
			Method:                 http.MethodPut,
			Disabled:               false,
		},
		{
			ID:                     constants.UserAPI,
			PathWithoutAPIBasePath: dashboardApiBasePath.AppendPath(userAPI),
			Method:                 http.MethodDelete,
			Disabled:               false,
		},
		{
			ID:                     constants.UserEmailVerifyAPI,
			PathWithoutAPIBasePath: dashboardApiBasePath.AppendPath(userEmailVerifyAPI),
			Method:                 http.MethodGet,
			Disabled:               false,
		},
		{
			ID:                     constants.UserEmailVerifyAPI,
			PathWithoutAPIBasePath: dashboardApiBasePath.AppendPath(userEmailVerifyAPI),
			Method:                 http.MethodPut,
			Disabled:               false,
		},
		{
			ID:                     constants.UserMetadataAPI,
			PathWithoutAPIBasePath: dashboardApiBasePath.AppendPath(userMetaDataAPI),
			Method:                 http.MethodGet,
			Disabled:               false,
		},
		{
			ID:                     constants.UserMetadataAPI,
			PathWithoutAPIBasePath: dashboardApiBasePath.AppendPath(userMetaDataAPI),
			Method:                 http.MethodPut,
			Disabled:               false,
		},
		{
			ID:                     constants.UserSessionsAPI,
			PathWithoutAPIBasePath: dashboardApiBasePath.AppendPath(userSessionsAPI),
			Method:                 http.MethodGet,
			Disabled:               false,
		},
		{
			ID:                     constants.UserSessionsAPI,
			PathWithoutAPIBasePath: dashboardApiBasePath.AppendPath(userSessionsAPI),
			Method:                 http.MethodPost,
			Disabled:               false,
		},
		{
			ID:                     constants.UserPasswordAPI,
			PathWithoutAPIBasePath: dashboardApiBasePath.AppendPath(userPasswordAPI),
			Method:                 http.MethodPut,
			Disabled:               false,
		},
		{
			ID:                     constants.UserEmailVerifyTokenAPI,
			PathWithoutAPIBasePath: dashboardApiBasePath.AppendPath(userEmailVerifyTokenAPI),
			Method:                 http.MethodPost,
			Disabled:               false,
		},
		{
			ID:                     constants.SearchTagsAPI,
			PathWithoutAPIBasePath: dashboardApiBasePath.AppendPath(searchTagsAPI),
			Method:                 http.MethodGet,
			Disabled:               false,
		},
		{
			ID:                     constants.DashboardAnalyticsAPI,
			PathWithoutAPIBasePath: dashboardApiBasePath.AppendPath(dashboardAnalyticsAPI),
			Method:                 http.MethodPost,
			Disabled:               false,
		},
		{
			ID:                     constants.TenantsListAPI,
			PathWithoutAPIBasePath: dashboardApiBasePath.AppendPath(tenantsListAPI),
			Method:                 http.MethodGet,
			Disabled:               false,
		},
	}, nil
}

func (r *Recipe) handleAPIRequest(id string, tenantId string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, _ supertokens.NormalisedURLPath, _ string, userContext supertokens.UserContext) error {
	options := dashboardmodels.APIOptions{
		Config:               r.Config,
		RecipeID:             r.RecipeModule.GetRecipeID(),
		RecipeImplementation: r.RecipeImpl,
		AppInfo:              r.RecipeModule.GetAppInfo(),
		Req:                  req,
		Res:                  res,
		OtherHandler:         theirHandler,
	}
	if id == dashboardAPI {
		return api.Dashboard(r.APIImpl, options, userContext)
	} else if id == validateKeyAPI {
		return api.ValidateKey(r.APIImpl, options, userContext)
	} else if id == signInAPI {
		return api.SignInPost(r.APIImpl, options, userContext)
	}

	// Do API key validation for the remaining APIs
	return apiKeyProtector(r.APIImpl, tenantId, options, userContext, func() (interface{}, error) {
		if id == usersListGetAPI {
			return api.UsersGet(r.APIImpl, tenantId, options, userContext)
		} else if id == usersCountAPI {
			return api.UsersCountGet(r.APIImpl, tenantId, options, userContext)
		} else if id == userAPI {
			if req.Method == http.MethodGet {
				return userdetails.UserGet(r.APIImpl, tenantId, options, userContext)
			}

			if req.Method == http.MethodPut {
				return userdetails.UserPut(r.APIImpl, tenantId, options, userContext)
			}

			if req.Method == http.MethodDelete {
				return userdetails.UserDelete(r.APIImpl, tenantId, options, userContext)
			}
		} else if id == userEmailVerifyAPI {
			if req.Method == http.MethodGet {
				return userdetails.UserEmailVerifyGet(r.APIImpl, tenantId, options, userContext)
			}

			if req.Method == http.MethodPut {
				return userdetails.UserEmailVerifyPut(r.APIImpl, tenantId, options, userContext)
			}
		} else if id == userSessionsAPI {
			if req.Method == http.MethodGet {
				return userdetails.UserSessionsGet(r.APIImpl, tenantId, options, userContext)
			}

			if req.Method == http.MethodPost {
				return userdetails.UserSessionsRevoke(r.APIImpl, tenantId, options, userContext)
			}
		} else if id == userMetaDataAPI {
			if req.Method == http.MethodGet {
				return userdetails.UserMetaDataGet(r.APIImpl, tenantId, options, userContext)
			}

			if req.Method == http.MethodPut {
				return userdetails.UserMetaDataPut(r.APIImpl, tenantId, options, userContext)
			}
		} else if id == userEmailVerifyTokenAPI {
			return userdetails.UserEmailVerifyTokenPost(r.APIImpl, tenantId, options, userContext)
		} else if id == userPasswordAPI {
			return userdetails.UserPasswordPut(r.APIImpl, tenantId, options, userContext)
		} else if id == searchTagsAPI {
			return search.SearchTagsGet(r.APIImpl, tenantId, options, userContext)
		} else if id == signOutAPI {
			return api.SignOutPost(r.APIImpl, tenantId, options, userContext)
		} else if id == dashboardAnalyticsAPI {
			return api.AnalyticsPost(r.APIImpl, tenantId, options, userContext)
		} else if id == tenantsListAPI {
			return api.TenantsListGet(r.APIImpl, tenantId, options, userContext)
		}
		return nil, errors.New("should never come here")
	})
}

func (r *Recipe) getAllCORSHeaders() []string {
	return []string{}
}

func (r *Recipe) handleError(err error, req *http.Request, res http.ResponseWriter) (bool, error) {
	return false, nil
}

func ResetForTest() {
	singletonInstance = nil
}
