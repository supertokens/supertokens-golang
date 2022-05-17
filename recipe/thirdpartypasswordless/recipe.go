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

	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
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
	RecipeModule            supertokens.RecipeModule
	Config                  tplmodels.TypeNormalisedInput
	EmailVerificationRecipe emailverification.Recipe
	passwordlessRecipe      *passwordless.Recipe
	thirdPartyRecipe        *thirdparty.Recipe
	RecipeImpl              tplmodels.RecipeInterface
	APIImpl                 tplmodels.APIInterface
}

var singletonInstance *Recipe

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config tplmodels.TypeInput, emailVerificationInstance *emailverification.Recipe, thirdPartyInstance *thirdparty.Recipe, passwordlessInstance *passwordless.Recipe, onGeneralError func(err error, req *http.Request, res http.ResponseWriter)) (Recipe, error) {
	r := &Recipe{}
	r.RecipeModule = supertokens.MakeRecipeModule(recipeId, appInfo, r.handleAPIRequest, r.getAllCORSHeaders, r.getAPIsHandled, r.handleError, onGeneralError)

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

	if emailVerificationInstance == nil {
		// we override the recipe function for email verification
		// to return true for is email verified for passwordless users.
		var apiInterfaceFunc func(originalImplementation evmodels.APIInterface) evmodels.APIInterface
		var recipeInterfaceFunc func(originalImplementation evmodels.RecipeInterface) evmodels.RecipeInterface
		if config.Override != nil && config.Override.EmailVerificationFeature != nil {
			apiInterfaceFunc = config.Override.EmailVerificationFeature.APIs
			recipeInterfaceFunc = config.Override.EmailVerificationFeature.Functions
		}
		verifiedConfig.EmailVerificationFeature.Override = &evmodels.OverrideStruct{
			APIs: apiInterfaceFunc,
			Functions: func(originalImplementation evmodels.RecipeInterface) evmodels.RecipeInterface {
				ogIsEmailVerified := *originalImplementation.IsEmailVerified
				(*originalImplementation.IsEmailVerified) = func(userID, email string, userContext supertokens.UserContext) (bool, error) {
					user, err := (*(*r).RecipeImpl.GetUserByID)(userID, userContext)
					if err != nil {
						return false, err
					}

					if user == nil || user.ThirdParty != nil {
						return ogIsEmailVerified(userID, email, userContext)
					}
					// this is a passwordless user, so we always want
					// to return that their info / email is verified
					return true, nil
				}

				ogCreateEmailVerificationToken := *originalImplementation.CreateEmailVerificationToken
				(*originalImplementation.CreateEmailVerificationToken) = func(userID, email string, userContext supertokens.UserContext) (evmodels.CreateEmailVerificationTokenResponse, error) {
					user, err := (*(*r).RecipeImpl.GetUserByID)(userID, userContext)
					if err != nil {
						return evmodels.CreateEmailVerificationTokenResponse{}, err
					}

					if user == nil || user.ThirdParty != nil {
						return ogCreateEmailVerificationToken(userID, email, userContext)
					}

					return evmodels.CreateEmailVerificationTokenResponse{
						EmailAlreadyVerifiedError: &struct{}{},
					}, nil
				}

				if recipeInterfaceFunc != nil {
					// we call the user's function override as well post our change.
					return recipeInterfaceFunc(originalImplementation)
				} else {
					return originalImplementation
				}
			},
		}

		// TODO: do not pass nil to emaildelivery ingredient
		emailVerificationRecipe, err := emailverification.MakeRecipe(recipeId, appInfo, verifiedConfig.EmailVerificationFeature, nil, onGeneralError)
		if err != nil {
			return Recipe{}, err
		}
		r.EmailVerificationRecipe = emailVerificationRecipe

	} else {
		r.EmailVerificationRecipe = *emailVerificationInstance
	}

	var passwordlessRecipe passwordless.Recipe
	if passwordlessInstance == nil {
		passwordlessConfig := plessmodels.TypeInput{
			ContactMethodPhone:        verifiedConfig.ContactMethodPhone,
			ContactMethodEmail:        verifiedConfig.ContactMethodEmail,
			ContactMethodEmailOrPhone: verifiedConfig.ContactMethodEmailOrPhone,
			FlowType:                  verifiedConfig.FlowType,
			GetLinkDomainAndPath:      verifiedConfig.GetLinkDomainAndPath,
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
		passwordlessRecipe, err = passwordless.MakeRecipe(recipeId, appInfo, passwordlessConfig, nil, onGeneralError)
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
					EmailVerificationFeature: nil,
				},
			}
			thirdPartyRecipeinstance, err := thirdparty.MakeRecipe(recipeId, appInfo, thirdPartyConfig, &r.EmailVerificationRecipe, onGeneralError)
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
	return func(appInfo supertokens.NormalisedAppinfo, onGeneralError func(err error, req *http.Request, res http.ResponseWriter)) (*supertokens.RecipeModule, error) {
		if singletonInstance == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config, nil, nil, nil, onGeneralError)
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

// implement RecipeModule

func (r *Recipe) getAPIsHandled() ([]supertokens.APIHandled, error) {
	passwordlessAPIhandled, err := r.passwordlessRecipe.RecipeModule.GetAPIsHandled()
	if err != nil {
		return nil, err
	}
	emailverificationAPIhandled, err := r.EmailVerificationRecipe.RecipeModule.GetAPIsHandled()
	if err != nil {
		return nil, err
	}
	apisHandled := append(passwordlessAPIhandled, emailverificationAPIhandled...)
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
	return r.EmailVerificationRecipe.RecipeModule.HandleAPIRequest(id, req, res, theirHandler, path, method)
}

func (r *Recipe) getAllCORSHeaders() []string {
	corsHeaders := append(r.EmailVerificationRecipe.RecipeModule.GetAllCORSHeaders(), r.passwordlessRecipe.RecipeModule.GetAllCORSHeaders()...)
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
	return r.EmailVerificationRecipe.RecipeModule.HandleError(err, req, res)
}

func (r *Recipe) getEmailForUserIdForEmailVerification(userID string, userContext supertokens.UserContext) (string, error) {
	userInfo, err := (*r.RecipeImpl.GetUserByID)(userID, userContext)
	if err != nil {
		return "", err
	}
	if userInfo == nil {
		return "", errors.New("Unknown User ID provided")
	}
	if userInfo.ThirdParty == nil {
		// this is a passwordless user.. so we always return some random email,
		// and in the function for isEmailVerified, we will check if the user
		// is a passwordless user, and if they are, we will return true in there
		return "_____supertokens_passwordless_user@supertokens.com", nil
	}
	return *userInfo.Email, nil
}

func ResetForTest() {
	singletonInstance = nil
}
