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

package emailverification

import (
	"errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/api"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evclaims"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "emailverification"

type Recipe struct {
	RecipeModule  supertokens.RecipeModule
	Config        evmodels.TypeNormalisedInput
	RecipeImpl    evmodels.RecipeInterface
	APIImpl       evmodels.APIInterface
	EmailDelivery emaildelivery.Ingredient

	GetEmailForUserID        evmodels.TypeGetEmailForUserID
	AddGetEmailForUserIdFunc func(function evmodels.TypeGetEmailForUserID)
}

var singletonInstance *Recipe

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config evmodels.TypeInput, emailDeliveryIngredient *emaildelivery.Ingredient, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (Recipe, error) {
	getEmailForUserIdFuncsFromOtherRecipes := []evmodels.TypeGetEmailForUserID{}

	r := &Recipe{}
	verifiedConfig, err := validateAndNormaliseUserInput(appInfo, config)
	if err != nil {
		return Recipe{}, err
	}
	r.Config = verifiedConfig
	r.APIImpl = verifiedConfig.Override.APIs(api.MakeAPIImplementation())

	querierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
	if err != nil {
		return Recipe{}, err
	}
	recipeImplementation := makeRecipeImplementation(*querierInstance)
	r.RecipeImpl = verifiedConfig.Override.Functions(recipeImplementation)

	recipeModuleInstance := supertokens.MakeRecipeModule(recipeId, appInfo, r.handleAPIRequest, r.getAllCORSHeaders, r.getAPIsHandled, r.handleError, onSuperTokensAPIError)
	r.RecipeModule = recipeModuleInstance

	if emailDeliveryIngredient != nil {
		r.EmailDelivery = *emailDeliveryIngredient
	} else {
		r.EmailDelivery = emaildelivery.MakeIngredient(verifiedConfig.GetEmailDeliveryConfig())
	}

	r.GetEmailForUserID = func(userID string, userContext supertokens.UserContext) (evmodels.TypeEmailInfo, error) {
		if r.Config.GetEmailForUserID != nil {
			emailRes, err := r.Config.GetEmailForUserID(userID, userContext)
			if err != nil {
				return evmodels.TypeEmailInfo{}, err
			}
			if emailRes.UnknownUserIDError != nil {
				return emailRes, nil
			}
		}

		var err error
		var emailRes evmodels.TypeEmailInfo
		for _, getEmailForUserIdFunc := range getEmailForUserIdFuncsFromOtherRecipes {
			emailRes, err = getEmailForUserIdFunc(userID, userContext)
			if err != nil {
				return emailRes, err
			}
			if emailRes.UnknownUserIDError != nil {
				return emailRes, nil
			}
		}
		return evmodels.TypeEmailInfo{
			UnknownUserIDError: &struct{}{},
		}, nil
	}

	r.AddGetEmailForUserIdFunc = func(function evmodels.TypeGetEmailForUserID) {
		getEmailForUserIdFuncsFromOtherRecipes = append(getEmailForUserIdFuncsFromOtherRecipes, function)
	}

	return *r, nil
}

func getRecipeInstanceOrThrowError() (*Recipe, error) {
	if singletonInstance != nil {
		return singletonInstance, nil
	}
	return nil, errors.New("Initialisation not done. Did you forget to call the init function?")
}

func GetRecipeInstance() *Recipe {
	return singletonInstance
}

func recipeInit(config evmodels.TypeInput) supertokens.Recipe {
	return func(appInfo supertokens.NormalisedAppinfo, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (*supertokens.RecipeModule, error) {
		if singletonInstance == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config, nil, onSuperTokensAPIError)
			if err != nil {
				return nil, err
			}
			singletonInstance = &recipe

			supertokens.AddPostInitCallback(func() error {
				sessionRecipe, err := session.GetRecipeInstanceOrThrowError()
				if err != nil {
					return err
				}

				sessionRecipe.AddClaimFromOtherRecipe(evclaims.EmailVerificationClaim)

				if config.Mode == evmodels.ModeRequired {
					sessionRecipe.AddClaimValidatorFromOtherRecipe(
						evclaims.EmailVerificationClaimValidators.IsVerified(nil),
					)
				}
				return nil
			})
			return &singletonInstance.RecipeModule, nil
		}
		return nil, errors.New("Emailverification recipe has already been initialised. Please check your code for bugs.")
	}
}

// implement RecipeModule

func (r *Recipe) getAPIsHandled() ([]supertokens.APIHandled, error) {
	generateEmailVerifyTokenAPINormalised, err := supertokens.NewNormalisedURLPath(generateEmailVerifyTokenAPI)
	if err != nil {
		return nil, err
	}
	emailVerifyAPINormalised, err := supertokens.NewNormalisedURLPath(emailVerifyAPI)
	if err != nil {
		return nil, err
	}

	return []supertokens.APIHandled{{
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: generateEmailVerifyTokenAPINormalised,
		ID:                     generateEmailVerifyTokenAPI,
		Disabled:               r.APIImpl.GenerateEmailVerifyTokenPOST == nil,
	}, {
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: emailVerifyAPINormalised,
		ID:                     emailVerifyAPI,
		Disabled:               r.APIImpl.VerifyEmailPOST == nil,
	}, {
		Method:                 http.MethodGet,
		PathWithoutAPIBasePath: emailVerifyAPINormalised,
		ID:                     emailVerifyAPI,
		Disabled:               r.APIImpl.IsEmailVerifiedGET == nil,
	}}, nil
}

func (r *Recipe) handleAPIRequest(id string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, _ supertokens.NormalisedURLPath, _ string) error {
	options := evmodels.APIOptions{
		Config:               r.Config,
		RecipeID:             r.RecipeModule.GetRecipeID(),
		RecipeImplementation: r.RecipeImpl,
		Req:                  req,
		Res:                  res,
		OtherHandler:         theirHandler,
		EmailDelivery:        r.EmailDelivery,
		GetEmailForUserID:    r.GetEmailForUserID,
	}
	if id == generateEmailVerifyTokenAPI {
		return api.GenerateEmailVerifyToken(r.APIImpl, options)
	} else {
		return api.EmailVerify(r.APIImpl, options)
	}
}

func (r *Recipe) getAllCORSHeaders() []string {
	return []string{}
}

func (r *Recipe) handleError(err error, req *http.Request, res http.ResponseWriter) (bool, error) {
	return false, nil
}

func ResetForTest() {
	singletonInstance = nil
	EmailVerificationEmailSentForTest = false
	EmailVerificationDataForTest = struct {
		User                    evmodels.User
		EmailVerifyURLWithToken string
		UserContext             supertokens.UserContext
	}{}
}
