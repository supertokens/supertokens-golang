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

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/emaildelivery/smtpService"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/providers"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type signInUpResponse struct {
	CreatedNewUser bool
	User           tpmodels.User
}

func Init(config *tpmodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func SignInUpWithContext(thirdPartyID string, thirdPartyUserID string, email tpmodels.EmailStruct, userContext supertokens.UserContext) (tpmodels.SignInUpResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return tpmodels.SignInUpResponse{}, err
	}
	return (*instance.RecipeImpl.SignInUp)(thirdPartyID, thirdPartyUserID, email, userContext)
}

func GetUserByIDWithContext(userID string, userContext supertokens.UserContext) (*tpmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, userContext)
}

func GetUsersByEmailWithContext(email string, userContext supertokens.UserContext) ([]tpmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return []tpmodels.User{}, err
	}
	return (*instance.RecipeImpl.GetUsersByEmail)(email, userContext)
}

func GetUserByThirdPartyInfoWithContext(thirdPartyID, thirdPartyUserID string, userContext supertokens.UserContext) (*tpmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByThirdPartyInfo)(thirdPartyID, thirdPartyUserID, userContext)
}

func CreateEmailVerificationTokenWithContext(userID string, userContext supertokens.UserContext) (evmodels.CreateEmailVerificationTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.CreateEmailVerificationTokenResponse{}, err
	}
	email, err := instance.getEmailForUserId(userID, userContext)
	if err != nil {
		return evmodels.CreateEmailVerificationTokenResponse{}, err
	}
	return (*instance.EmailVerificationRecipe.RecipeImpl.CreateEmailVerificationToken)(userID, email, userContext)
}

func VerifyEmailUsingTokenWithContext(token string, userContext supertokens.UserContext) (*tpmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	response, err := (*instance.EmailVerificationRecipe.RecipeImpl.VerifyEmailUsingToken)(token, userContext)
	if err != nil {
		return nil, err
	}
	if response.EmailVerificationInvalidTokenError != nil {
		return nil, errors.New("email verification token is invalid")
	}
	return (*instance.RecipeImpl.GetUserByID)(response.OK.User.ID, userContext)
}

func IsEmailVerifiedWithContext(userID string, userContext supertokens.UserContext) (bool, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return false, err
	}
	email, err := instance.getEmailForUserId(userID, userContext)
	if err != nil {
		return false, err
	}
	return (*instance.EmailVerificationRecipe.RecipeImpl.IsEmailVerified)(userID, email, userContext)
}

func RevokeEmailVerificationTokensWithContext(userID string, userContext supertokens.UserContext) (evmodels.RevokeEmailVerificationTokensResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.RevokeEmailVerificationTokensResponse{}, err
	}
	email, err := instance.getEmailForUserId(userID, userContext)
	if err != nil {
		return evmodels.RevokeEmailVerificationTokensResponse{}, err
	}
	return (*instance.EmailVerificationRecipe.RecipeImpl.RevokeEmailVerificationTokens)(userID, email, userContext)
}

func UnverifyEmailWithContext(userID string, userContext supertokens.UserContext) (evmodels.UnverifyEmailResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.UnverifyEmailResponse{}, err
	}
	email, err := instance.getEmailForUserId(userID, userContext)
	if err != nil {
		return evmodels.UnverifyEmailResponse{}, err
	}
	return (*instance.EmailVerificationRecipe.RecipeImpl.UnverifyEmail)(userID, email, userContext)
}

func SendEmailWithContext(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	return (*instance.EmailDelivery.IngredientInterfaceImpl.SendEmail)(input, userContext)
}

func SignInUp(thirdPartyID string, thirdPartyUserID string, email tpmodels.EmailStruct) (tpmodels.SignInUpResponse, error) {
	return SignInUpWithContext(thirdPartyID, thirdPartyUserID, email, &map[string]interface{}{})
}

func GetUserByID(userID string) (*tpmodels.User, error) {
	return GetUserByIDWithContext(userID, &map[string]interface{}{})
}

func GetUsersByEmail(email string) ([]tpmodels.User, error) {
	return GetUsersByEmailWithContext(email, &map[string]interface{}{})
}

func GetUserByThirdPartyInfo(thirdPartyID, thirdPartyUserID string) (*tpmodels.User, error) {
	return GetUserByThirdPartyInfoWithContext(thirdPartyID, thirdPartyUserID, &map[string]interface{}{})
}

func CreateEmailVerificationToken(userID string) (evmodels.CreateEmailVerificationTokenResponse, error) {
	return CreateEmailVerificationTokenWithContext(userID, &map[string]interface{}{})
}

func VerifyEmailUsingToken(token string) (*tpmodels.User, error) {
	return VerifyEmailUsingTokenWithContext(token, &map[string]interface{}{})
}

func IsEmailVerified(userID string) (bool, error) {
	return IsEmailVerifiedWithContext(userID, &map[string]interface{}{})
}

func RevokeEmailVerificationTokens(userID string) (evmodels.RevokeEmailVerificationTokensResponse, error) {
	return RevokeEmailVerificationTokensWithContext(userID, &map[string]interface{}{})
}

func UnverifyEmail(userID string) (evmodels.UnverifyEmailResponse, error) {
	return UnverifyEmailWithContext(userID, &map[string]interface{}{})
}

func Apple(config tpmodels.AppleConfig) tpmodels.TypeProvider {
	return providers.Apple(config)
}

func Facebook(config tpmodels.FacebookConfig) tpmodels.TypeProvider {
	return providers.Facebook(config)
}

func Github(config tpmodels.GithubConfig) tpmodels.TypeProvider {
	return providers.Github(config)
}

func Discord(config tpmodels.DiscordConfig) tpmodels.TypeProvider {
	return providers.Discord(config)
}

func GoogleWorkspaces(config tpmodels.GoogleWorkspacesConfig) tpmodels.TypeProvider {
	return providers.GoogleWorkspaces(config)
}

func Google(config tpmodels.GoogleConfig) tpmodels.TypeProvider {
	return providers.Google(config)
}

func SendEmail(input emaildelivery.EmailType) error {
	return SendEmailWithContext(input, &map[string]interface{}{})
}

func MakeSMTPService(config emaildelivery.SMTPServiceConfig) *emaildelivery.EmailDeliveryInterface {
	return smtpService.MakeSMTPService(config)
}
