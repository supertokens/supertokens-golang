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

package userdetails

import (
	"errors"
	"fmt"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type userEmailVerifyTokenPost struct {
	Status string `json:"status"`
}

func UserEmailVerifyTokenPost(apiInterface dashboardmodels.APIInterface, options dashboardmodels.APIOptions)(userEmailVerifyTokenPost, error) {
	req := options.Req
	userId := req.URL.Query().Get("userId")

	if userId == "" {
		return userEmailVerifyTokenPost{}, supertokens.BadInputError {
			Msg: "Missing required parameter 'userId'",
		}
	}

	emailresponse, emailErr := emailverification.GetRecipeInstance().GetEmailForUserID(userId, supertokens.MakeDefaultUserContextFromAPI(options.Req))

	if emailErr != nil {
		return userEmailVerifyTokenPost{}, emailErr
	}

	if emailresponse.OK == nil {
		return userEmailVerifyTokenPost{}, errors.New("Should never come here")
	}

	emailVerificationToken, tokenError := emailverification.CreateEmailVerificationToken(userId, &emailresponse.OK.Email)

	if tokenError != nil {
		return userEmailVerifyTokenPost{}, tokenError
	}

	if emailVerificationToken.EmailAlreadyVerifiedError != nil {
		return userEmailVerifyTokenPost{
			Status: "EMAIL_ALREADY_VERIFIED_ERROR",
		}, nil
	}

	emailVerificationURL := fmt.Sprintf(
		"%s%s/verify-email?token=%s&rid=%s",
		options.AppInfo.WebsiteDomain.GetAsStringDangerous(),
		options.AppInfo.WebsiteBasePath.GetAsStringDangerous(),
		emailVerificationToken.OK.Token,
		options.RecipeID,
	)

	emailverification.SendEmail(emaildelivery.EmailType{
		EmailVerification: &emaildelivery.EmailVerificationType{
			User: emaildelivery.User{
				ID: userId,
				Email: emailresponse.OK.Email,
			},
			EmailVerifyLink: emailVerificationURL,
		},
	})

	return userEmailVerifyTokenPost{
		Status: "OK",
	}, nil
}