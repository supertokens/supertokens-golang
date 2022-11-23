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

package thirdpartypasswordless

import (
	"errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/recipe/passwordless"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/api"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/recipeimplementation"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "thirdpartypasswordless"

type Recipe struct {
	RecipeModule       supertokens.RecipeModule
	Config             tplmodels.TypeNormalisedInput
	passwordlessRecipe *passwordless.Recipe
	thirdPartyRecipe   *thirdparty.Recipe
	RecipeImpl         tplmodels.RecipeInterface
	APIImpl            tplmodels.APIInterface
	EmailDelivery      emaildelivery.Ingredient
	SmsDelivery        smsdelivery.Ingredient
}

var singletonInstance *Recipe

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config tplmodels.TypeInput, thirdPartyInstance *thirdparty.Recipe, passwordlessInstance *passwordless.Recipe, emailDeliveryIngredient *emaildelivery.Ingredient, smsDeliveryIngredient *smsdelivery.Ingredient, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (Recipe, error) {
	r := &Recipe{}
	r.RecipeModule = supertokens.MakeRecipeModule(recipeId, appInfo, r.handleAPIRequest, r.getAllCORSHeaders, r.getAPIsHandled, nil, r.handleError, onSuperTokensAPIError)

	verifiedConfig, err := validateAndNormaliseUserInput(r, appInfo, config)
	if err != nil {
		return Recipe{}, err
	}
	r.Config = verifiedConfig
	{
		passwordlessquerierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(passwordless.RECIPE_ID)
		if err != nil {
			return Recipe{}, err
		}
		thirdpartyquerierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(thirdparty.RECIPE_ID)
		if err != nil {
			return Recipe{}, err
		}

		r.RecipeImpl = verifiedConfig.Override.Functions(recipeimplementation.MakeRecipeImplementation(*passwordlessquerierInstance, thirdpartyquerierInstance))
	}
	r.APIImpl = verifiedConfig.Override.APIs(api.MakeAPIImplementation())

	if emailDeliveryIngredient != nil {
		r.EmailDelivery = *emailDeliveryIngredient
	} else {
		r.EmailDelivery = emaildelivery.MakeIngredient(verifiedConfig.GetEmailDeliveryConfig())
	}

	if smsDeliveryIngredient != nil {
		r.SmsDelivery = *smsDeliveryIngredient
	} else {
		r.SmsDelivery = smsdelivery.MakeIngredient(verifiedConfig.GetSmsDeliveryConfig())
	}

	var passwordlessRecipe passwordless.Recipe
	if passwordlessInstance == nil {
		passwordlessConfig := plessmodels.TypeInput{
			ContactMethodPhone:        verifiedConfig.ContactMethodPhone,
			ContactMethodEmail:        verifiedConfig.ContactMethodEmail,
			ContactMethodEmailOrPhone: verifiedConfig.ContactMethodEmailOrPhone,
			FlowType:                  verifiedConfig.FlowType,
			GetCustomUserInputCode:    verifiedConfig.GetCustomUserInputCode,
			Override: &plessmodels.OverrideStruct{
				Functions: func(originalImplementation plessmodels.RecipeInterface) plessmodels.RecipeInterface {
					return recipeimplementation.MakePasswordlessRecipeImplementation(r.RecipeImpl)
				},
				APIs: func(originalImplementation plessmodels.APIInterface) plessmodels.APIInterface {
					return api.GetPasswordlessIterfaceImpl(r.APIImpl)
				},
			},
		}
		passwordlessRecipe, err = passwordless.MakeRecipe(recipeId, appInfo, passwordlessConfig, &r.EmailDelivery, &r.SmsDelivery, onSuperTokensAPIError)
		if err != nil {
			return Recipe{}, err
		}
		r.passwordlessRecipe = &passwordlessRecipe
	} else {
		r.passwordlessRecipe = passwordlessInstance
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

func recipeInit(config tplmodels.TypeInput) supertokens.Recipe {
	return func(appInfo supertokens.NormalisedAppinfo, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (*supertokens.RecipeModule, error) {
		if singletonInstance == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config, nil, nil, nil, nil, onSuperTokensAPIError)
			if err != nil {
				return nil, err
			}
			singletonInstance = &recipe
			return &singletonInstance.RecipeModule, nil
		}
		return nil, errors.New("ThirdPartyPasswordless recipe has already been initialised. Please check your code for bugs.")
	}
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

// implement RecipeModule

func (r *Recipe) getAPIsHandled() ([]supertokens.APIHandled, error) {
	passwordlessAPIhandled, err := r.passwordlessRecipe.RecipeModule.GetAPIsHandled()
	if err != nil {
		return nil, err
	}
	apisHandled := passwordlessAPIhandled
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
	ok, err := r.passwordlessRecipe.RecipeModule.ReturnAPIIdIfCanHandleRequest(path, method)
	if err != nil {
		return err
	}
	if ok != nil {
		return r.passwordlessRecipe.RecipeModule.HandleAPIRequest(id, req, res, theirHandler, path, method)
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
	corsHeaders := r.passwordlessRecipe.RecipeModule.GetAllCORSHeaders()
	if r.thirdPartyRecipe != nil {
		corsHeaders = append(corsHeaders, r.thirdPartyRecipe.RecipeModule.GetAllCORSHeaders()...)
	}
	return corsHeaders
}

func (r *Recipe) handleError(err error, req *http.Request, res http.ResponseWriter) (bool, error) {
	handleError, err := r.passwordlessRecipe.RecipeModule.HandleError(err, req, res)
	if err != nil || handleError {
		return handleError, err
	}
	if r.thirdPartyRecipe != nil {
		handleError, err = r.thirdPartyRecipe.RecipeModule.HandleError(err, req, res)
		if err != nil || handleError {
			return handleError, err
		}
	}
	return false, nil
}

func ResetForTest() {
	singletonInstance = nil
}
