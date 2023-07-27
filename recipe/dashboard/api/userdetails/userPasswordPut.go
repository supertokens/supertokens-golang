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

	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type userPasswordPutResponse struct {
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

type userPasswordPutRequestBody struct {
	UserId      *string `json:"userId"`
	NewPassword *string `json:"newPassword`
}

func UserPasswordPut(apiInterface dashboardmodels.APIInterface, options dashboardmodels.APIOptions, userContext supertokens.UserContext) (userPasswordPutResponse, error) {
	body, err := supertokens.ReadFromRequest(options.Req)

	if err != nil {
		return userPasswordPutResponse{}, err
	}

	var readBody userPasswordPutRequestBody
	err = json.Unmarshal(body, &readBody)
	if err != nil {
		return userPasswordPutResponse{}, err
	}

	if readBody.UserId == nil {
		return userPasswordPutResponse{}, supertokens.BadInputError{
			Msg: "Required parameter 'userId' is missing",
		}
	}

	if readBody.NewPassword == nil {
		return userPasswordPutResponse{}, supertokens.BadInputError{
			Msg: "Required parameter 'newPassword' is missing",
		}
	}

	recipeToUse := "none"

	emailPasswordInstance := emailpassword.GetRecipeInstance()

	if emailPasswordInstance != nil {
		recipeToUse = "emailpassword"
	}

	if recipeToUse == "none" {
		tpepInstance := thirdpartyemailpassword.GetRecipeInstance()

		if tpepInstance != nil {
			recipeToUse = "thirdpartyemailpassword"
		}
	}

	if recipeToUse == "none" {
		// This means that neither emailpassword or thirdpartyemailpassword is initialised
		return userPasswordPutResponse{}, errors.New("Should never come here")
	}

	if recipeToUse == "emailpassword" {
		var passwordField epmodels.NormalisedFormField

		for _, value := range emailPasswordInstance.Config.SignUpFeature.FormFields {
			if value.ID == "password" {
				passwordField = value
			}
		}

		validationError := passwordField.Validate(*readBody.NewPassword)

		if validationError != nil {
			return userPasswordPutResponse{
				Status: "INVALID_PASSWORD_ERROR",
				Error:  *validationError,
			}, nil
		}

		passwordResetToken, resetTokenErr := emailpassword.CreateResetPasswordToken(*readBody.UserId, "public", userContext) // TODO multitenancy pass tenantId

		if resetTokenErr != nil {
			return userPasswordPutResponse{}, resetTokenErr
		}

		if passwordResetToken.UnknownUserIdError != nil {
			// Techincally it can but its an edge case so we assume that it wont
			return userPasswordPutResponse{}, errors.New("Should never come here")
		}

		passwordResetResponse, passwordResetErr := emailpassword.ResetPasswordUsingToken(passwordResetToken.OK.Token, *readBody.NewPassword, "public", userContext) // TODO multitenancy pass tenantId

		if passwordResetErr != nil {
			return userPasswordPutResponse{}, passwordResetErr
		}

		if passwordResetResponse.ResetPasswordInvalidTokenError != nil {
			return userPasswordPutResponse{}, errors.New("Should never come here")
		}

		return userPasswordPutResponse{
			Status: "OK",
		}, nil
	}

	var passwordField epmodels.TypeInputFormField

	for _, value := range thirdpartyemailpassword.GetRecipeInstance().Config.SignUpFeature.FormFields {
		if value.ID == "password" {
			passwordField = value
		}
	}

	validationError := passwordField.Validate(*readBody.NewPassword)

	if validationError != nil {
		return userPasswordPutResponse{
			Status: "INVALID_PASSWORD_ERROR",
			Error:  *validationError,
		}, nil
	}

	passwordResetToken, resetTokenErr := thirdpartyemailpassword.CreateResetPasswordToken(*readBody.UserId, userContext)

	if resetTokenErr != nil {
		return userPasswordPutResponse{}, resetTokenErr
	}

	if passwordResetToken.UnknownUserIdError != nil {
		// Techincally it can but its an edge case so we assume that it wont
		return userPasswordPutResponse{}, errors.New("Should never come here")
	}

	passwordResetResponse, passwordResetErr := thirdpartyemailpassword.ResetPasswordUsingToken(passwordResetToken.OK.Token, *readBody.NewPassword, userContext)

	if passwordResetErr != nil {
		return userPasswordPutResponse{}, passwordResetErr
	}

	if passwordResetResponse.ResetPasswordInvalidTokenError != nil {
		return userPasswordPutResponse{}, errors.New("Should never come here")
	}

	return userPasswordPutResponse{
		Status: "OK",
	}, nil
}
