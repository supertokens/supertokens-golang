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
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/dashboard/api"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/passwordless"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless"
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type userPutResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type updateEmailResponse struct {
	Status string
	Error  string
}

type updatePhoneResponse struct {
	Status string
	Error  string
}

type userPutRequestBody struct {
	UserId    *string `json:"userId"`
	RecipeId  *string `json:"recipeId"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
}

func updateEmailForRecipeId(recipeId string, userId string, email string, tenantId string, userContext supertokens.UserContext) (updateEmailResponse, error) {
	if recipeId == "emailpassword" {
		var emailField epmodels.NormalisedFormField

		for _, value := range emailpassword.GetRecipeInstance().Config.SignUpFeature.FormFields {
			if value.ID == "email" {
				emailField = value
			}
		}

		validationError := emailField.Validate(email, tenantId)

		if validationError != nil {
			return updateEmailResponse{
				Status: "INVALID_EMAIL_ERROR",
				Error:  *validationError,
			}, nil
		}

		updateResponse, err := emailpassword.UpdateEmailOrPassword(userId, &email, nil, nil, nil, userContext)

		if err != nil {
			return updateEmailResponse{}, err
		}

		if updateResponse.EmailAlreadyExistsError != nil {
			return updateEmailResponse{
				Status: "EMAIL_ALREADY_EXISTS_ERROR",
			}, nil
		}

		if updateResponse.UnknownUserIdError != nil {
			return updateEmailResponse{}, errors.New("Should never come here")
		}

		return updateEmailResponse{
			Status: "OK",
		}, nil
	}

	if recipeId == "passwordless" {
		isValidEmail := true
		validationError := ""

		passwordlessConfig := passwordless.GetRecipeInstance().Config

		if passwordlessConfig.ContactMethodPhone.Enabled {
			validationResult := passwordless.DefaultValidateEmailAddress(email, tenantId)

			if validationResult != nil {
				isValidEmail = false
				validationError = *validationResult
			}
		} else if passwordlessConfig.ContactMethodEmail.Enabled {
			validationResult := passwordlessConfig.ContactMethodEmail.ValidateEmailAddress(email, tenantId)

			if validationResult != nil {
				isValidEmail = false
				validationError = *validationResult
			}
		} else {
			validationResult := passwordlessConfig.ContactMethodEmailOrPhone.ValidateEmailAddress(email, tenantId)

			if validationResult != nil {
				isValidEmail = false
				validationError = *validationResult
			}
		}

		if !isValidEmail {
			return updateEmailResponse{
				Status: "INVALID_EMAIL_ERROR",
				Error:  validationError,
			}, nil
		}

		updateResponse, updateErr := passwordless.UpdateUser(userId, &email, nil, userContext)

		if updateErr != nil {
			return updateEmailResponse{}, updateErr
		}

		if updateResponse.UnknownUserIdError != nil {
			return updateEmailResponse{}, errors.New("Should never come here")
		}

		if updateResponse.EmailAlreadyExistsError != nil {
			return updateEmailResponse{
				Status: "EMAIL_ALREADY_EXISTS_ERROR",
			}, nil
		}

		return updateEmailResponse{
			Status: "OK",
		}, nil
	}

	if recipeId == "thirdpartypasswordless" {
		isValidEmail := true
		validationError := ""

		passwordlessConfig := thirdpartypasswordless.GetRecipeInstance().Config

		if passwordlessConfig.ContactMethodPhone.Enabled {
			validationResult := passwordless.DefaultValidateEmailAddress(email, tenantId)

			if validationResult != nil {
				isValidEmail = false
				validationError = *validationResult
			}
		} else if passwordlessConfig.ContactMethodEmail.Enabled {
			validationResult := passwordlessConfig.ContactMethodEmail.ValidateEmailAddress(email, tenantId)

			if validationResult != nil {
				isValidEmail = false
				validationError = *validationResult
			}
		} else {
			validationResult := passwordlessConfig.ContactMethodEmailOrPhone.ValidateEmailAddress(email, tenantId)

			if validationResult != nil {
				isValidEmail = false
				validationError = *validationResult
			}
		}

		if !isValidEmail {
			return updateEmailResponse{
				Status: "INVALID_EMAIL_ERROR",
				Error:  validationError,
			}, nil
		}

		updateResponse, updateErr := thirdpartypasswordless.UpdatePasswordlessUser(userId, &email, nil, userContext)

		if updateErr != nil {
			return updateEmailResponse{}, updateErr
		}

		if updateResponse.UnknownUserIdError != nil {
			return updateEmailResponse{}, errors.New("Should never come here")
		}

		if updateResponse.EmailAlreadyExistsError != nil {
			return updateEmailResponse{
				Status: "EMAIL_ALREADY_EXISTS_ERROR",
			}, nil
		}

		return updateEmailResponse{
			Status: "OK",
		}, nil
	}

	return updateEmailResponse{}, errors.New("Should never come here")
}

func updatePhoneForRecipeId(recipeId string, userId string, phone string, tenantId string, userContext supertokens.UserContext) (updatePhoneResponse, error) {
	if recipeId == "passwordless" {
		isValidPhone := true
		validationError := ""

		passwordlessConfig := passwordless.GetRecipeInstance().Config

		if passwordlessConfig.ContactMethodEmail.Enabled {
			validationResult := passwordless.DefaultValidatePhoneNumber(phone, tenantId)

			if validationResult != nil {
				isValidPhone = false
				validationError = *validationResult
			}
		} else if passwordlessConfig.ContactMethodPhone.Enabled {
			validationResult := passwordlessConfig.ContactMethodPhone.ValidatePhoneNumber(phone, tenantId)

			if validationResult != nil {
				isValidPhone = false
				validationError = *validationResult
			}
		} else {
			validationResult := passwordlessConfig.ContactMethodEmailOrPhone.ValidatePhoneNumber(phone, tenantId)

			if validationResult != nil {
				isValidPhone = false
				validationError = *validationResult
			}
		}

		if !isValidPhone {
			return updatePhoneResponse{
				Status: "INVALID_PHONE_ERROR",
				Error:  validationError,
			}, nil
		}

		updateResponse, updateErr := passwordless.UpdateUser(userId, nil, &phone, userContext)

		if updateErr != nil {
			return updatePhoneResponse{}, updateErr
		}

		if updateResponse.UnknownUserIdError != nil {
			return updatePhoneResponse{}, errors.New("Should never come here")
		}

		if updateResponse.EmailAlreadyExistsError != nil {
			return updatePhoneResponse{
				Status: "PHONE_ALREADY_EXISTS_ERROR",
			}, nil
		}

		return updatePhoneResponse{
			Status: "OK",
		}, nil
	}

	if recipeId == "thirdpartypasswordless" {
		isValidPhone := true
		validationError := ""

		passwordlessConfig := thirdpartypasswordless.GetRecipeInstance().Config

		if passwordlessConfig.ContactMethodEmail.Enabled {
			validationResult := passwordless.DefaultValidatePhoneNumber(phone, tenantId)

			if validationResult != nil {
				isValidPhone = false
				validationError = *validationResult
			}
		} else if passwordlessConfig.ContactMethodPhone.Enabled {
			validationResult := passwordlessConfig.ContactMethodPhone.ValidatePhoneNumber(phone, tenantId)

			if validationResult != nil {
				isValidPhone = false
				validationError = *validationResult
			}
		} else {
			validationResult := passwordlessConfig.ContactMethodEmailOrPhone.ValidatePhoneNumber(phone, tenantId)

			if validationResult != nil {
				isValidPhone = false
				validationError = *validationResult
			}
		}

		if !isValidPhone {
			return updatePhoneResponse{
				Status: "INVALID_PHONE_ERROR",
				Error:  validationError,
			}, nil
		}

		updateResponse, updateErr := thirdpartypasswordless.UpdatePasswordlessUser(userId, nil, &phone, userContext)

		if updateErr != nil {
			return updatePhoneResponse{}, updateErr
		}

		if updateResponse.UnknownUserIdError != nil {
			return updatePhoneResponse{}, errors.New("Should never come here")
		}

		if updateResponse.EmailAlreadyExistsError != nil {
			return updatePhoneResponse{
				Status: "PHONE_ALREADY_EXISTS_ERROR",
			}, nil
		}

		return updatePhoneResponse{
			Status: "OK",
		}, nil
	}

	/**
	 * If it comes here then the user is a not a passwordless user in which case the UI should not have allowed this
	 */
	return updatePhoneResponse{}, errors.New("Should never come here")
}

func UserPut(apiInterface dashboardmodels.APIInterface, tenantId string, options dashboardmodels.APIOptions, userContext supertokens.UserContext) (userPutResponse, error) {
	body, err := supertokens.ReadFromRequest(options.Req)

	if err != nil {
		return userPutResponse{}, err
	}

	var readBody userPutRequestBody
	err = json.Unmarshal(body, &readBody)
	if err != nil {
		return userPutResponse{}, err
	}

	if readBody.UserId == nil {
		return userPutResponse{}, supertokens.BadInputError{
			Msg: "Required parameter 'userid' is missing",
		}
	}

	if readBody.RecipeId == nil {
		return userPutResponse{}, supertokens.BadInputError{
			Msg: "Required parameter 'recipeId' is missing",
		}
	}

	if readBody.FirstName == nil {
		return userPutResponse{}, supertokens.BadInputError{
			Msg: "Required parameter 'firstName' is missing",
		}
	}

	if readBody.LastName == nil {
		return userPutResponse{}, supertokens.BadInputError{
			Msg: "Required parameter 'lastName' is missing",
		}
	}

	if readBody.Email == nil {
		return userPutResponse{}, supertokens.BadInputError{
			Msg: "Required parameter 'email' is missing",
		}
	}

	if readBody.Phone == nil {
		return userPutResponse{}, supertokens.BadInputError{
			Msg: "Required parameter 'phone' is missing",
		}
	}

	_, recipeId := api.GetUserForRecipeId(*readBody.UserId, *readBody.RecipeId, userContext)

	if *readBody.FirstName != "" || *readBody.LastName != "" {
		isRecipeInitialised := false

		_, err = usermetadata.GetRecipeInstanceOrThrowError()

		if err == nil {
			isRecipeInitialised = true
		}

		// If the recipe is not initialised we consider updating the names as a no-op instead of throwing an error
		if isRecipeInitialised {
			metadataupdate := make(map[string]interface{})

			if strings.TrimSpace(*readBody.FirstName) != "" {
				metadataupdate["first_name"] = strings.TrimSpace(*readBody.FirstName)
			}

			if strings.TrimSpace(*readBody.LastName) != "" {
				metadataupdate["last_name"] = strings.TrimSpace(*readBody.LastName)
			}

			usermetadata.UpdateUserMetadata(*readBody.UserId, metadataupdate, userContext)
		}
	}

	if strings.TrimSpace(*readBody.Email) != "" {
		updateResponse, updateError := updateEmailForRecipeId(recipeId, *readBody.UserId, strings.TrimSpace(*readBody.Email), tenantId, userContext)

		if updateError != nil {
			return userPutResponse{}, updateError
		}

		if updateResponse.Status != "OK" {
			return userPutResponse{
				Status: updateResponse.Status,
				Error:  updateResponse.Error,
			}, nil
		}
	}

	if strings.TrimSpace(*readBody.Phone) != "" {
		updateResponse, updateError := updatePhoneForRecipeId(recipeId, *readBody.UserId, *readBody.Phone, tenantId, userContext)

		if updateError != nil {
			return userPutResponse{}, updateError
		}

		if updateResponse.Status != "OK" {
			return userPutResponse{
				Status: updateResponse.Status,
				Error:  updateResponse.Error,
			}, nil
		}
	}

	return userPutResponse{
		Status: "OK",
	}, nil
}
