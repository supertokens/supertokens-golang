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

package thirdpartyemailpassword

import (
	"errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/api"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/recipeimplementation"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "thirdpartyemailpassword"

type Recipe struct {
	RecipeModule        supertokens.RecipeModule
	Config              tpepmodels.TypeNormalisedInput
	emailPasswordRecipe *emailpassword.Recipe
	thirdPartyRecipe    *thirdparty.Recipe
	RecipeImpl          tpepmodels.RecipeInterface
	APIImpl             tpepmodels.APIInterface
	EmailDelivery       emaildelivery.Ingredient
}

var singletonInstance *Recipe

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *tpepmodels.TypeInput, emailVerificationInstance *emailverification.Recipe, thirdPartyInstance *thirdparty.Recipe, emailPasswordInstance *emailpassword.Recipe, emailDeliveryIngredient *emaildelivery.Ingredient, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (Recipe, error) {
	r := &Recipe{}
	r.RecipeModule = supertokens.MakeRecipeModule(recipeId, appInfo, r.handleAPIRequest, r.getAllCORSHeaders, r.getAPIsHandled, r.handleError, onSuperTokensAPIError)

	verifiedConfig, err := validateAndNormaliseUserInput(r, appInfo, config)
	if err != nil {
		return Recipe{}, err
	}
	r.Config = verifiedConfig
	{
		emailpasswordquerierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(emailpassword.RECIPE_ID)
		if err != nil {
			return Recipe{}, err
		}
		thirdpartyquerierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(thirdparty.RECIPE_ID)
		if err != nil {
			return Recipe{}, err
		}

		r.RecipeImpl = verifiedConfig.Override.Functions(recipeimplementation.MakeRecipeImplementation(*emailpasswordquerierInstance, thirdpartyquerierInstance))
	}
	r.APIImpl = verifiedConfig.Override.APIs(api.MakeAPIImplementation())

	var emailPasswordRecipe emailpassword.Recipe
	emailPasswordRecipeImpl := recipeimplementation.MakeEmailPasswordRecipeImplementation(r.RecipeImpl)

	if emailDeliveryIngredient != nil {
		r.EmailDelivery = *emailDeliveryIngredient
	} else {
		r.EmailDelivery = emaildelivery.MakeIngredient(verifiedConfig.GetEmailDeliveryConfig(r.RecipeImpl, emailPasswordRecipeImpl))
	}

	if emailPasswordInstance == nil {
		emailPasswordConfig := &epmodels.TypeInput{
			SignUpFeature:                  verifiedConfig.SignUpFeature,
			ResetPasswordUsingTokenFeature: verifiedConfig.ResetPasswordUsingTokenFeature,
			Override: &epmodels.OverrideStruct{
				Functions: func(_ epmodels.RecipeInterface) epmodels.RecipeInterface {
					return emailPasswordRecipeImpl
				},
				APIs: func(_ epmodels.APIInterface) epmodels.APIInterface {
					return api.GetEmailPasswordIterfaceImpl(r.APIImpl)
				},
			},
		}
		emailPasswordRecipe, err = emailpassword.MakeRecipe(recipeId, appInfo, emailPasswordConfig, &r.EmailDelivery, onSuperTokensAPIError)
		if err != nil {
			return Recipe{}, err
		}
		r.emailPasswordRecipe = &emailPasswordRecipe
	} else {
		r.emailPasswordRecipe = emailPasswordInstance
	}

	if len(verifiedConfig.Providers) > 0 {
		if thirdPartyInstance == nil {
			thirdPartyConfig := &tpmodels.TypeInput{
				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
					Providers: verifiedConfig.Providers,
				},
				Override: &tpmodels.OverrideStruct{
					Functions: func(_ tpmodels.RecipeInterface) tpmodels.RecipeInterface {
						return recipeimplementation.MakeThirdPartyRecipeImplementation(r.RecipeImpl)
					},
					APIs: func(_ tpmodels.APIInterface) tpmodels.APIInterface {
						return api.GetThirdPartyIterfaceImpl(r.APIImpl)
					},
				},
			}
			thirdPartyRecipeinstance, err := thirdparty.MakeRecipe(recipeId, appInfo, thirdPartyConfig, &r.EmailDelivery, onSuperTokensAPIError)
			if err != nil {
				return Recipe{}, err
			}
			r.thirdPartyRecipe = &thirdPartyRecipeinstance
		} else {
			r.thirdPartyRecipe = thirdPartyInstance
		}
	}

	return *r, nil
}

func recipeInit(config *tpepmodels.TypeInput) supertokens.Recipe {
	return func(appInfo supertokens.NormalisedAppinfo, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (*supertokens.RecipeModule, error) {
		if singletonInstance == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config, nil, nil, nil, nil, onSuperTokensAPIError)
			if err != nil {
				return nil, err
			}
			singletonInstance = &recipe
			return &singletonInstance.RecipeModule, nil
		}
		return nil, errors.New("ThirdPartyEmailPassword recipe has already been initialised. Please check your code for bugs.")
	}
}

func getRecipeInstanceOrThrowError() (*Recipe, error) {
	if singletonInstance != nil {
		return singletonInstance, nil
	}
	return nil, errors.New("Initialisation not done. Did you forget to call the init function?")
}

// implement RecipeModule

func (r *Recipe) getAPIsHandled() ([]supertokens.APIHandled, error) {
	emailpasswordAPIhandled, err := r.emailPasswordRecipe.RecipeModule.GetAPIsHandled()
	if err != nil {
		return nil, err
	}
	apisHandled := append(emailpasswordAPIhandled)
	if r.thirdPartyRecipe != nil {
		thirdpartyAPIhandled, err := r.thirdPartyRecipe.RecipeModule.GetAPIsHandled()
		if err != nil {
			return nil, err
		}
		apisHandled = append(apisHandled, thirdpartyAPIhandled...)
	}
	return apisHandled, nil
}

func (r *Recipe) handleAPIRequest(id string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, path supertokens.NormalisedURLPath, method string) error {
	ok, err := r.emailPasswordRecipe.RecipeModule.ReturnAPIIdIfCanHandleRequest(path, method)
	if err != nil {
		return err
	}
	if ok != nil {
		return r.emailPasswordRecipe.RecipeModule.HandleAPIRequest(id, req, res, theirHandler, path, method)
	}
	if r.thirdPartyRecipe != nil {
		ok, err := r.thirdPartyRecipe.RecipeModule.ReturnAPIIdIfCanHandleRequest(path, method)
		if err != nil {
			return err
		}
		if ok != nil {
			return r.thirdPartyRecipe.RecipeModule.HandleAPIRequest(id, req, res, theirHandler, path, method)
		}
	}
	return errors.New("should not come here")
}

func (r *Recipe) getAllCORSHeaders() []string {
	corsHeaders := r.emailPasswordRecipe.RecipeModule.GetAllCORSHeaders()
	if r.thirdPartyRecipe != nil {
		corsHeaders = append(corsHeaders, r.thirdPartyRecipe.RecipeModule.GetAllCORSHeaders()...)
	}
	return corsHeaders
}

func (r *Recipe) handleError(err error, req *http.Request, res http.ResponseWriter) (bool, error) {
	handleError, err := r.emailPasswordRecipe.RecipeModule.HandleError(err, req, res)
	if err != nil || handleError {
		return handleError, err
	}
	if r.thirdPartyRecipe != nil {
		handleError, err = r.thirdPartyRecipe.RecipeModule.HandleError(err, req, res)
		if err != nil || handleError {
			return handleError, err
		}
	}
	return false, err
}

func ResetForTest() {
	singletonInstance = nil
}
