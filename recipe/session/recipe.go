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

package session

import (
	defaultErrors "errors"
	"net/http"
	"strconv"

	"github.com/supertokens/supertokens-golang/recipe/openid"
	"github.com/supertokens/supertokens-golang/recipe/openid/openidmodels"
	"github.com/supertokens/supertokens-golang/recipe/session/api"
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessionwithjwt"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type Recipe struct {
	RecipeModule supertokens.RecipeModule
	Config       sessmodels.TypeNormalisedInput
	RecipeImpl   sessmodels.RecipeInterface
	OpenIdRecipe *openid.Recipe
	APIImpl      sessmodels.APIInterface

	claimsAddedByOtherRecipes          []*claims.TypeSessionClaim
	claimValidatorsAddedByOtherRecipes []claims.SessionClaimValidator
}

const RECIPE_ID = "session"

var singletonInstance *Recipe

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *sessmodels.TypeInput, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (Recipe, error) {
	r := &Recipe{
		claimsAddedByOtherRecipes:          []*claims.TypeSessionClaim{},
		claimValidatorsAddedByOtherRecipes: []claims.SessionClaimValidator{},
	}

	r.RecipeModule = supertokens.MakeRecipeModule(recipeId, appInfo, r.handleAPIRequest, r.getAllCORSHeaders, r.getAPIsHandled, nil, r.handleError, onSuperTokensAPIError)

	verifiedConfig, configError := validateAndNormaliseUserInput(appInfo, config)
	if configError != nil {
		return Recipe{}, configError
	}

	supertokens.LogDebugMessage("session init: AntiCsrf: " + verifiedConfig.AntiCsrf)
	if verifiedConfig.CookieDomain != nil {
		supertokens.LogDebugMessage("session init: CookieDomain: " + *verifiedConfig.CookieDomain)
	} else {
		supertokens.LogDebugMessage("session init: CookieDomain: nil")
	}
	supertokens.LogDebugMessage("session init: CookieSameSite: " + verifiedConfig.CookieSameSite)
	supertokens.LogDebugMessage("session init: CookieSecure: " + strconv.FormatBool(verifiedConfig.CookieSecure))
	supertokens.LogDebugMessage("session init: RefreshTokenPath: " + verifiedConfig.RefreshTokenPath.GetAsStringDangerous())
	supertokens.LogDebugMessage("session init: SessionExpiredStatusCode: " + strconv.Itoa(verifiedConfig.SessionExpiredStatusCode))

	r.Config = verifiedConfig
	r.APIImpl = verifiedConfig.Override.APIs(api.MakeAPIImplementation())

	querierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
	if err != nil {
		return Recipe{}, err
	}
	recipeImplementation := makeRecipeImplementation(*querierInstance, verifiedConfig)

	if verifiedConfig.Jwt.Enable {
		openIdRecipe, err := openid.MakeRecipe(recipeId, appInfo, &openidmodels.TypeInput{
			Issuer:   verifiedConfig.Jwt.Issuer,
			Override: verifiedConfig.Override.OpenIdFeature,
		}, onSuperTokensAPIError)
		if err != nil {
			return Recipe{}, err
		}
		r.RecipeImpl = verifiedConfig.Override.Functions(sessionwithjwt.MakeRecipeImplementation(recipeImplementation, openIdRecipe.RecipeImpl, verifiedConfig))
		r.OpenIdRecipe = &openIdRecipe
	} else {
		r.RecipeImpl = verifiedConfig.Override.Functions(recipeImplementation)
	}

	return *r, nil
}

func getRecipeInstanceOrThrowError() (*Recipe, error) {
	if singletonInstance != nil {
		return singletonInstance, nil
	}
	return nil, defaultErrors.New("Initialisation not done. Did you forget to call the init function?")
}

func GetRecipeInstanceOrThrowError() (*Recipe, error) {
	return getRecipeInstanceOrThrowError()
}

func recipeInit(config *sessmodels.TypeInput) supertokens.Recipe {
	return func(appInfo supertokens.NormalisedAppinfo, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (*supertokens.RecipeModule, error) {
		if singletonInstance == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config, onSuperTokensAPIError)
			if err != nil {
				return nil, err
			}
			singletonInstance = &recipe
			return &singletonInstance.RecipeModule, nil
		}
		return nil, defaultErrors.New("Session recipe has already been initialised. Please check your code for bugs.")
	}
}

// Implement RecipeModule

func (r *Recipe) getAPIsHandled() ([]supertokens.APIHandled, error) {
	refreshAPIPathNormalised, err := supertokens.NewNormalisedURLPath(refreshAPIPath)
	if err != nil {
		return nil, err
	}
	signoutAPIPathNormalised, err := supertokens.NewNormalisedURLPath(signoutAPIPath)
	if err != nil {
		return nil, err
	}
	resp := []supertokens.APIHandled{{
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: refreshAPIPathNormalised,
		ID:                     refreshAPIPath,
		Disabled:               r.APIImpl.RefreshPOST == nil,
	}, {
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: signoutAPIPathNormalised,
		ID:                     signoutAPIPath,
		Disabled:               r.APIImpl.SignOutPOST == nil,
	}}

	if r.OpenIdRecipe != nil {
		jwtAPIs, err := r.OpenIdRecipe.RecipeModule.GetAPIsHandled()
		if err != nil {
			return nil, err
		}
		resp = append(resp, jwtAPIs...)
	}

	return resp, nil
}

func (r *Recipe) handleAPIRequest(id string, req *http.Request, res http.ResponseWriter, theirhandler http.HandlerFunc, path supertokens.NormalisedURLPath, method string) error {
	options := sessmodels.APIOptions{
		Config:               r.Config,
		RecipeID:             r.RecipeModule.GetRecipeID(),
		RecipeImplementation: r.RecipeImpl,
		Req:                  req,
		Res:                  res,
		OtherHandler:         theirhandler,

		ClaimValidatorsAddedByOtherRecipes: r.getClaimValidatorsAddedByOtherRecipes(),
	}
	if id == refreshAPIPath {
		return api.HandleRefreshAPI(r.APIImpl, options)
	} else if id == signoutAPIPath {
		return api.SignOutAPI(r.APIImpl, options)
	} else if r.OpenIdRecipe != nil {
		return r.OpenIdRecipe.RecipeModule.HandleAPIRequest(id, req, res, theirhandler, path, method)
	}
	return nil
}

func (r *Recipe) getAllCORSHeaders() []string {
	resp := getCORSAllowedHeaders()
	if r.OpenIdRecipe != nil {
		resp = append(resp, r.OpenIdRecipe.RecipeModule.GetAllCORSHeaders()...)
	}
	return resp
}

func (r *Recipe) handleError(err error, req *http.Request, res http.ResponseWriter) (bool, error) {
	if defaultErrors.As(err, &errors.UnauthorizedError{}) {
		supertokens.LogDebugMessage("errorHandler: returning UNAUTHORISED")
		unauthErr := err.(errors.UnauthorizedError)
		return true, r.Config.ErrorHandlers.OnUnauthorised(err.Error(), unauthErr.ClearCookies == nil || *unauthErr.ClearCookies, req, res)
	} else if defaultErrors.As(err, &errors.TryRefreshTokenError{}) {
		supertokens.LogDebugMessage("errorHandler: returning TRY_REFRESH_TOKEN")
		return true, r.Config.ErrorHandlers.OnTryRefreshToken(err.Error(), req, res)
	} else if defaultErrors.As(err, &errors.TokenTheftDetectedError{}) {
		supertokens.LogDebugMessage("errorHandler: returning TOKEN_THEFT_DETECTED")
		errs := err.(errors.TokenTheftDetectedError)
		return true, r.Config.ErrorHandlers.OnTokenTheftDetected(errs.Payload.SessionHandle, errs.Payload.UserID, req, res)
	} else if defaultErrors.As(err, &errors.InvalidClaimError{}) {
		supertokens.LogDebugMessage("errorHandler: returning INVALID_CLAIMS")
		errs := err.(errors.InvalidClaimError)
		return true, r.Config.ErrorHandlers.OnInvalidClaim(errs.InvalidClaims, req, res)
	} else if r.OpenIdRecipe != nil {
		return r.OpenIdRecipe.RecipeModule.HandleError(err, req, res)
	}
	return false, nil
}

// Claim functions
func (r *Recipe) AddClaimFromOtherRecipe(claim *claims.TypeSessionClaim) error {
	for _, existingClaim := range r.claimsAddedByOtherRecipes {
		if claim.Key == existingClaim.Key {
			return defaultErrors.New("claim already added by other recipe")
		}
	}
	r.claimsAddedByOtherRecipes = append(r.claimsAddedByOtherRecipes, claim)
	return nil
}

func (r *Recipe) getClaimsAddedByOtherRecipes() []*claims.TypeSessionClaim {
	return r.claimsAddedByOtherRecipes
}

func (r *Recipe) AddClaimValidatorFromOtherRecipe(validator claims.SessionClaimValidator) error {
	r.claimValidatorsAddedByOtherRecipes = append(r.claimValidatorsAddedByOtherRecipes, validator)
	return nil
}

func (r *Recipe) getClaimValidatorsAddedByOtherRecipes() []claims.SessionClaimValidator {
	return r.claimValidatorsAddedByOtherRecipes
}

func ResetForTest() {
	singletonInstance = nil
}
