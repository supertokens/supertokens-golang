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
	"encoding/json"
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

type userEmailverifyTokenPostRequestBody struct {
	UserId *string `json:"userId"`
}

func UserEmailVerifyTokenPost(apiInterface dashboardmodels.APIInterface, options dashboardmodels.APIOptions, userContext supertokens.UserContext) (userEmailVerifyTokenPost, error) {
	body, err := supertokens.ReadFromRequest(options.Req)

	if err != nil {
		return userEmailVerifyTokenPost{}, err
	}

	var readBody userEmailverifyTokenPostRequestBody
	err = json.Unmarshal(body, &readBody)
	if err != nil {
		return userEmailVerifyTokenPost{}, err
	}

	if readBody.UserId == nil {
		return userEmailVerifyTokenPost{}, supertokens.BadInputError{
			Msg: "Required parameter 'userId' is missing",
		}
	}

	emailresponse, emailErr := emailverification.GetRecipeInstance().GetEmailForUserID(*readBody.UserId, userContext)

	if emailErr != nil {
		return userEmailVerifyTokenPost{}, emailErr
	}

	if emailresponse.OK == nil {
		return userEmailVerifyTokenPost{}, errors.New("Should never come here")
	}

	emailVerificationToken, tokenError := emailverification.CreateEmailVerificationToken(*readBody.UserId, &emailresponse.OK.Email)

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
				ID:    *readBody.UserId,
				Email: emailresponse.OK.Email,
			},
			EmailVerifyLink: emailVerificationURL,
		},
	})

	return userEmailVerifyTokenPost{
		Status: "OK",
	}, nil
}
