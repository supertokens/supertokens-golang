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
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func makeRecipeImplementation(querier supertokens.Querier) evmodels.RecipeInterface {
	createEmailVerificationToken := func(userID, email string, tenantId *string, userContext supertokens.UserContext) (evmodels.CreateEmailVerificationTokenResponse, error) {
		response, err := querier.SendPostRequest(supertokens.GetPathPrefixForTenantId(tenantId)+"/recipe/user/email/verify/token", map[string]interface{}{
			"userId": userID,
			"email":  email,
		})
		if err != nil {
			return evmodels.CreateEmailVerificationTokenResponse{}, err
		}
		status, ok := response["status"]
		if ok && status == "OK" {
			return evmodels.CreateEmailVerificationTokenResponse{
				OK: &struct{ Token string }{Token: response["token"].(string)},
			}, nil
		}

		return evmodels.CreateEmailVerificationTokenResponse{
			EmailAlreadyVerifiedError: &struct{}{},
		}, nil
	}

	verifyEmailUsingToken := func(token string, tenantId *string, userContext supertokens.UserContext) (evmodels.VerifyEmailUsingTokenResponse, error) {
		response, err := querier.SendPostRequest(supertokens.GetPathPrefixForTenantId(tenantId)+"/recipe/user/email/verify", map[string]interface{}{
			"method": "token",
			"token":  token,
		})
		if err != nil {
			return evmodels.VerifyEmailUsingTokenResponse{}, err
		}
		status, ok := response["status"]
		if ok && status == "OK" {
			return evmodels.VerifyEmailUsingTokenResponse{
				OK: &struct{ User evmodels.User }{User: evmodels.User{
					ID:    response["userId"].(string),
					Email: response["email"].(string),
				}},
			}, nil
		}
		return evmodels.VerifyEmailUsingTokenResponse{
			EmailVerificationInvalidTokenError: &struct{}{},
		}, nil
	}

	isEmailVerified := func(userID, email string, tenantId *string, userContext supertokens.UserContext) (bool, error) {
		response, err := querier.SendGetRequest(supertokens.GetPathPrefixForTenantId(tenantId)+"/recipe/user/email/verify", map[string]string{
			"userId": userID,
			"email":  email,
		})
		if err != nil {
			return false, err
		}
		return response["isVerified"].(bool), nil
	}

	revokeEmailVerificationTokens := func(userId string, email string, tenantId *string, userContext supertokens.UserContext) (evmodels.RevokeEmailVerificationTokensResponse, error) {
		_, err := querier.SendPostRequest(supertokens.GetPathPrefixForTenantId(tenantId)+"/recipe/user/email/verify/token/remove", map[string]interface{}{
			"userId": userId,
			"email":  email,
		})
		if err != nil {
			return evmodels.RevokeEmailVerificationTokensResponse{}, err
		}
		return evmodels.RevokeEmailVerificationTokensResponse{
			OK: &struct{}{},
		}, nil
	}

	unverifyEmail := func(userId string, email string, tenantId *string, userContext supertokens.UserContext) (evmodels.UnverifyEmailResponse, error) {
		_, err := querier.SendPostRequest(supertokens.GetPathPrefixForTenantId(tenantId)+"/recipe/user/email/verify/remove", map[string]interface{}{
			"userId": userId,
			"email":  email,
		})
		if err != nil {
			return evmodels.UnverifyEmailResponse{}, err
		}
		return evmodels.UnverifyEmailResponse{
			OK: &struct{}{},
		}, nil
	}
	return evmodels.RecipeInterface{
		CreateEmailVerificationToken:  &createEmailVerificationToken,
		VerifyEmailUsingToken:         &verifyEmailUsingToken,
		IsEmailVerified:               &isEmailVerified,
		RevokeEmailVerificationTokens: &revokeEmailVerificationTokens,
		UnverifyEmail:                 &unverifyEmail,
	}
}
