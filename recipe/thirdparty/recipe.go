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

package thirdparty

import (
	"errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/api"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tperrors"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "thirdparty"

type Recipe struct {
	RecipeModule supertokens.RecipeModule
	Config       tpmodels.TypeNormalisedInput
	RecipeImpl   tpmodels.RecipeInterface
	APIImpl      tpmodels.APIInterface
	Providers    []tpmodels.ProviderInput
}

var singletonInstance *Recipe

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *tpmodels.TypeInput, emailDeliveryIngredient *emaildelivery.Ingredient, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (Recipe, error) {
	r := &Recipe{}

	r.RecipeModule = supertokens.MakeRecipeModule(recipeId, appInfo, r.handleAPIRequest, r.getAllCORSHeaders, r.getAPIsHandled, nil, r.handleError, onSuperTokensAPIError)

	querierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
	if err != nil {
		return Recipe{}, err
	}
	verifiedConfig, err := validateAndNormaliseUserInput(r, appInfo, config)
	if err != nil {
		return Recipe{}, err
	}
	r.Config = verifiedConfig
	r.APIImpl = verifiedConfig.Override.APIs(api.MakeAPIImplementation())
	r.RecipeImpl = verifiedConfig.Override.Functions(MakeRecipeImplementation(*querierInstance, config.SignInAndUpFeature.Providers))
	r.Providers = config.SignInAndUpFeature.Providers

	supertokens.AddPostInitCallback(func() error {
		evRecipe := emailverification.GetRecipeInstance()
		if evRecipe != nil {
			evRecipe.AddGetEmailForUserIdFunc(r.getEmailForUserId)
		}

		mtRecipe := multitenancy.GetRecipeInstance()
		if mtRecipe != nil {
			mtRecipe.SetStaticThirdPartyProviders(verifiedConfig.SignInAndUpFeature.Providers)
		}

		return nil
	})

	return *r, nil
}

func recipeInit(config *tpmodels.TypeInput) supertokens.Recipe {
	return func(appInfo supertokens.NormalisedAppinfo, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (*supertokens.RecipeModule, error) {
		if singletonInstance == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config, nil, onSuperTokensAPIError)
			if err != nil {
				return nil, err
			}
			singletonInstance = &recipe
			return &singletonInstance.RecipeModule, nil
		}
		return nil, errors.New("ThirdParty recipe has already been initialised. Please check your code for bugs.")
	}
}

func GetRecipeInstanceOrThrowError() (*Recipe, error) {
	if singletonInstance != nil {
		return singletonInstance, nil
	}
	return nil, errors.New("Initialisation not done. Did you forget to call the init function?")
}

// implement RecipeModule

func (r *Recipe) getAPIsHandled() ([]supertokens.APIHandled, error) {
	signInUpAPI, err := supertokens.NewNormalisedURLPath(SignInUpAPI)
	if err != nil {
		return nil, err
	}
	authorisationAPI, err := supertokens.NewNormalisedURLPath(AuthorisationAPI)
	if err != nil {
		return nil, err
	}
	appleRedirectHandlerAPI, err := supertokens.NewNormalisedURLPath(AppleRedirectHandlerAPI)
	if err != nil {
		return nil, err
	}
	return append([]supertokens.APIHandled{{
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: signInUpAPI,
		ID:                     SignInUpAPI,
		Disabled:               r.APIImpl.SignInUpPOST == nil,
	}, {
		Method:                 http.MethodGet,
		PathWithoutAPIBasePath: authorisationAPI,
		ID:                     AuthorisationAPI,
		Disabled:               r.APIImpl.AuthorisationUrlGET == nil,
	}, {
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: appleRedirectHandlerAPI,
		ID:                     AppleRedirectHandlerAPI,
		Disabled:               r.APIImpl.AppleRedirectHandlerPOST == nil,
	}}), nil
}

func (r *Recipe) handleAPIRequest(id string, tenantId string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, path supertokens.NormalisedURLPath, method string, userContext supertokens.UserContext) error {
	options := tpmodels.APIOptions{
		Config:               r.Config,
		OtherHandler:         theirHandler,
		RecipeID:             r.RecipeModule.GetRecipeID(),
		RecipeImplementation: r.RecipeImpl,
		Providers:            r.Providers,
		Req:                  req,
		Res:                  res,
		AppInfo:              r.RecipeModule.GetAppInfo(),
	}
	if id == SignInUpAPI {
		return api.SignInUpAPI(r.APIImpl, options, userContext)
	} else if id == AuthorisationAPI {
		return api.AuthorisationUrlAPI(r.APIImpl, options, userContext)
	} else if id == AppleRedirectHandlerAPI {
		return api.AppleRedirectHandler(r.APIImpl, options, userContext)
	}
	return errors.New("should never come here")
}

func (r *Recipe) getAllCORSHeaders() []string {
	return []string{}
}

func (r *Recipe) handleError(err error, req *http.Request, res http.ResponseWriter) (bool, error) {
	if errors.As(err, &tperrors.ClientTypeNotFoundError{}) {
		supertokens.SendNon200ResponseWithMessage(res, err.Error(), 400)
		return true, nil
	}
	return false, nil
}

func (r *Recipe) getEmailForUserId(userID string, userContext supertokens.UserContext) (evmodels.TypeEmailInfo, error) {
	userInfo, err := (*r.RecipeImpl.GetUserByID)(userID, userContext)
	if err != nil {
		return evmodels.TypeEmailInfo{}, err
	}
	if userInfo == nil {
		return evmodels.TypeEmailInfo{
			UnknownUserIDError: &struct{}{},
		}, nil
	}
	return evmodels.TypeEmailInfo{
		OK: &struct{ Email string }{
			Email: userInfo.Email,
		},
	}, nil
}

func ResetForTest() {
	singletonInstance = nil
}
