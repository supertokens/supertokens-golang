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

package supertokens

import (
	"errors"
	"net/http"
)

const RECIPE_ID = "accountlinking"

type AccountLinkingRecipe struct {
	RecipeModule RecipeModule
	Config       AccountLinkingTypeNormalisedInput
	RecipeImpl   AccountLinkingRecipeInterface
}

var singletonInstance *AccountLinkingRecipe

func makeAccountLinkingRecipe(recipeId string, appInfo NormalisedAppinfo, config *AccountLinkingTypeInput, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (AccountLinkingRecipe, error) {
	r := &AccountLinkingRecipe{}
	verifiedConfig := validateAndNormaliseAccountLinkingUserInput(appInfo, config)
	r.Config = verifiedConfig

	querierInstance, err := GetNewQuerierInstanceOrThrowError(recipeId)
	if err != nil {
		return AccountLinkingRecipe{}, err
	}
	recipeImplementation := makeRecipeImplementation(*querierInstance, verifiedConfig)
	r.RecipeImpl = verifiedConfig.Override.Functions(recipeImplementation)

	recipeModuleInstance := MakeRecipeModule(recipeId, appInfo, r.handleAPIRequest, r.getAllCORSHeaders, r.getAPIsHandled, nil, r.handleError, onSuperTokensAPIError)
	r.RecipeModule = recipeModuleInstance

	return *r, nil
}

func getAccountLinkingRecipeInstanceOrThrowError() (*AccountLinkingRecipe, error) {
	if singletonInstance != nil {
		return singletonInstance, nil
	}
	return nil, errors.New("Initialisation not done. Did you forget to call the init function?")
}

func accountLinkingRecipeInit(config *AccountLinkingTypeInput) Recipe {
	return func(appInfo NormalisedAppinfo, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (*RecipeModule, error) {
		if singletonInstance == nil {
			recipe, err := makeAccountLinkingRecipe(RECIPE_ID, appInfo, config, onSuperTokensAPIError)
			if err != nil {
				return nil, err
			}
			singletonInstance = &recipe

			return &singletonInstance.RecipeModule, nil
		}
		return nil, errors.New("Account linking recipe has already been initialised. Please check your code for bugs.")
	}
}

// implement RecipeModule

func (r *AccountLinkingRecipe) getAPIsHandled() ([]APIHandled, error) {
	return []APIHandled{}, nil
}

func (r *AccountLinkingRecipe) handleAPIRequest(id string, tenantId string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, _ NormalisedURLPath, _ string, userContext UserContext) error {
	return errors.New("should never come here")
}

func (r *AccountLinkingRecipe) getAllCORSHeaders() []string {
	return []string{}
}

func (r *AccountLinkingRecipe) handleError(err error, req *http.Request, res http.ResponseWriter, userContext UserContext) (bool, error) {
	return false, nil
}

func verifyEmailForRecipeUserIfLinkedAccountsAreVerified(user User, recipeUserId RecipeUserID, userContext UserContext) error {
	if InternalUseEmailVerificationRecipeProxyInstance == nil {
		// if email verification recipe is not initialised, then no op
		return nil
	}

	// This is just a helper function cause it's called in many places
	// like during sign up, sign in and post linking accounts.
	// This is not exposed to the developer as it's called in the relevant
	// recipe functions.
	// We do not do this in the core cause email verification is a different
	// recipe.
	// Finally, we only mark the email of this recipe user as verified and not
	// the other recipe users in the primary user (if this user's email is verified),
	// cause when those other users sign in, this function will be called for them anyway.

	if user.IsPrimaryUser {
		var recipeUserEmail *string = nil
		isAlreadyVerified := false
		for _, method := range user.LoginMethods {
			if method.RecipeUserID.GetAsString() == recipeUserId.GetAsString() {
				recipeUserEmail = method.Email
				isAlreadyVerified = method.Verified
			}
		}

		if recipeUserEmail != nil {
			if isAlreadyVerified {
				return nil
			}

			shouldVerifyEmail := false
			for _, method := range user.LoginMethods {
				if method.HasSameEmailAs(recipeUserEmail) && method.Verified {
					shouldVerifyEmail = true
				}
			}

			if shouldVerifyEmail {
				// While the token we create here is tenant specific, the verification status is not
				// So we can use any tenantId the user is associated with here as long as we use the
				// same in the verifyEmailUsingToken call
				token, err := InternalUseEmailVerificationRecipeProxyInstance.CreateEmailVerificationToken(recipeUserId, *recipeUserEmail, user.TenantIDs[0], userContext)
				if err != nil {
					return err
				}
				if token.OK != nil {
					// we purposely pass in false below cause we don't want account
					// linking to happen
					_, err := InternalUseEmailVerificationRecipeProxyInstance.VerifyEmailUsingToken(token.OK.Token, user.TenantIDs[0], false, userContext)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
