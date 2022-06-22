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

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/emaildelivery/smtpService"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/smsdelivery/supertokensService"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/smsdelivery/twilioService"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init(config tplmodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func ThirdPartySignInUpWithContext(thirdPartyID string, thirdPartyUserID string, email tplmodels.EmailStruct, userContext supertokens.UserContext) (tplmodels.ThirdPartySignInUp, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return tplmodels.ThirdPartySignInUp{}, err
	}
	return (*instance.RecipeImpl.ThirdPartySignInUp)(thirdPartyID, thirdPartyUserID, email, userContext)
}

func GetUserByThirdPartyInfoWithContext(thirdPartyID string, thirdPartyUserID string, userContext supertokens.UserContext) (*tplmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByThirdPartyInfo)(thirdPartyID, thirdPartyUserID, userContext)
}

func GetUserByIdWithContext(userID string, userContext supertokens.UserContext) (*tplmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, userContext)
}

func GetUsersByEmailWithContext(email string, userContext supertokens.UserContext) ([]tplmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUsersByEmail)(email, userContext)
}

func CreateEmailVerificationTokenWithContext(userID string, userContext supertokens.UserContext) (evmodels.CreateEmailVerificationTokenResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return evmodels.CreateEmailVerificationTokenResponse{}, err
	}
	email, err := instance.getEmailForUserIdForEmailVerification(userID, userContext)
	if err != nil {
		return evmodels.CreateEmailVerificationTokenResponse{}, err
	}
	return (*instance.EmailVerificationRecipe.RecipeImpl.CreateEmailVerificationToken)(userID, email, userContext)
}

func VerifyEmailUsingTokenWithContext(token string, userContext supertokens.UserContext) (*tplmodels.User, error) {
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
	email, err := instance.getEmailForUserIdForEmailVerification(userID, userContext)
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
	email, err := instance.getEmailForUserIdForEmailVerification(userID, userContext)
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
	email, err := instance.getEmailForUserIdForEmailVerification(userID, userContext)
	if err != nil {
		return evmodels.UnverifyEmailResponse{}, err
	}
	return (*instance.EmailVerificationRecipe.RecipeImpl.UnverifyEmail)(userID, email, userContext)
}

func CreateCodeWithEmailWithContext(email string, userInputCode *string, userContext supertokens.UserContext) (plessmodels.CreateCodeResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.CreateCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.CreateCode)(&email, nil, userInputCode, userContext)
}

func CreateCodeWithPhoneNumberWithContext(phoneNumber string, userInputCode *string, userContext supertokens.UserContext) (plessmodels.CreateCodeResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.CreateCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.CreateCode)(nil, &phoneNumber, userInputCode, userContext)
}

func CreateNewCodeForDeviceWithContext(deviceID string, userInputCode *string, userContext supertokens.UserContext) (plessmodels.ResendCodeResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.ResendCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.CreateNewCodeForDevice)(deviceID, userInputCode, userContext)
}

func ConsumeCodeWithUserInputCodeWithContext(deviceID string, userInputCode string, preAuthSessionID string, userContext supertokens.UserContext) (tplmodels.ConsumeCodeResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return tplmodels.ConsumeCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ConsumeCode)(&plessmodels.UserInputCodeWithDeviceID{
		Code:     userInputCode,
		DeviceID: deviceID,
	}, nil, preAuthSessionID, userContext)
}

func ConsumeCodeWithLinkCodeWithContext(linkCode string, preAuthSessionID string, userContext supertokens.UserContext) (tplmodels.ConsumeCodeResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return tplmodels.ConsumeCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ConsumeCode)(nil, &linkCode, preAuthSessionID, userContext)
}

func GetUserByIDWithContext(userID string, userContext supertokens.UserContext) (*tplmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, userContext)
}

func GetUserByPhoneNumberWithContext(phoneNumber string, userContext supertokens.UserContext) (*tplmodels.User, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.GetUserByPhoneNumber)(phoneNumber, userContext)
}

func UpdatePasswordlessUserWithContext(userID string, email *string, phoneNumber *string, userContext supertokens.UserContext) (plessmodels.UpdateUserResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.UpdateUserResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.UpdatePasswordlessUser)(userID, email, phoneNumber, userContext)
}

func DeleteEmailForPasswordlessUserWithContext(userID string, userContext supertokens.UserContext) (plessmodels.DeleteUserResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.DeleteUserResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.DeleteEmailForPasswordlessUser)(userID, userContext)
}

func DeletePhoneNumberForUserWithContext(userID string, userContext supertokens.UserContext) (plessmodels.DeleteUserResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.DeleteUserResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.DeletePhoneNumberForUser)(userID, userContext)
}

func RevokeAllCodesByEmailWithContext(email string, userContext supertokens.UserContext) error {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.RevokeAllCodes)(&email, nil, userContext)
}

func RevokeAllCodesByPhoneNumberWithContext(phoneNumber string, userContext supertokens.UserContext) error {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.RevokeAllCodes)(nil, &phoneNumber, userContext)
}

func RevokeCodeWithContext(codeID string, userContext supertokens.UserContext) error {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.RevokeCode)(codeID, userContext)
}

func ListCodesByEmailWithContext(email string, userContext supertokens.UserContext) ([]plessmodels.DeviceType, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return []plessmodels.DeviceType{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ListCodesByEmail)(email, userContext)
}

func ListCodesByPhoneNumberWithContext(phoneNumber string, userContext supertokens.UserContext) ([]plessmodels.DeviceType, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return []plessmodels.DeviceType{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ListCodesByPhoneNumber)(phoneNumber, userContext)
}

func ListCodesByDeviceIDWithContext(deviceID string, userContext supertokens.UserContext) (*plessmodels.DeviceType, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ListCodesByDeviceID)(deviceID, userContext)
}

func ListCodesByPreAuthSessionIDWithContext(preAuthSessionID string, userContext supertokens.UserContext) (*plessmodels.DeviceType, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ListCodesByPreAuthSessionID)(preAuthSessionID, userContext)
}

func CreateMagicLinkByEmailWithContext(email string, userContext supertokens.UserContext) (string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return "", err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return instance.passwordlessRecipe.CreateMagicLink(&email, nil, userContext)
}

func CreateMagicLinkByPhoneNumberWithContext(phoneNumber string, userContext supertokens.UserContext) (string, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return "", err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return instance.passwordlessRecipe.CreateMagicLink(nil, &phoneNumber, userContext)
}

func PasswordlessSignInUpByEmailWithContext(email string, userContext supertokens.UserContext) (struct {
	PreAuthSessionID string
	CreatedNewUser   bool
	User             tplmodels.User
}, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return struct {
			PreAuthSessionID string
			CreatedNewUser   bool
			User             tplmodels.User
		}{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	resp, err := instance.passwordlessRecipe.SignInUp(&email, nil, userContext)
	if err != nil {
		return struct {
			PreAuthSessionID string
			CreatedNewUser   bool
			User             tplmodels.User
		}{}, err
	}
	return struct {
		PreAuthSessionID string
		CreatedNewUser   bool
		User             tplmodels.User
	}{
		PreAuthSessionID: resp.PreAuthSessionID,
		CreatedNewUser:   resp.CreatedNewUser,
		User: tplmodels.User{
			ID:          resp.User.ID,
			TimeJoined:  resp.User.TimeJoined,
			Email:       resp.User.Email,
			PhoneNumber: resp.User.PhoneNumber,
			ThirdParty:  nil,
		},
	}, nil
}

func PasswordlessSignInUpByPhoneNumberWithContext(phoneNumber string, userContext supertokens.UserContext) (struct {
	PreAuthSessionID string
	CreatedNewUser   bool
	User             tplmodels.User
}, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return struct {
			PreAuthSessionID string
			CreatedNewUser   bool
			User             tplmodels.User
		}{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	resp, err := instance.passwordlessRecipe.SignInUp(nil, &phoneNumber, userContext)
	if err != nil {
		return struct {
			PreAuthSessionID string
			CreatedNewUser   bool
			User             tplmodels.User
		}{}, err
	}
	return struct {
		PreAuthSessionID string
		CreatedNewUser   bool
		User             tplmodels.User
	}{
		PreAuthSessionID: resp.PreAuthSessionID,
		CreatedNewUser:   resp.CreatedNewUser,
		User: tplmodels.User{
			ID:          resp.User.ID,
			TimeJoined:  resp.User.TimeJoined,
			Email:       resp.User.Email,
			PhoneNumber: resp.User.PhoneNumber,
			ThirdParty:  nil,
		},
	}, nil
}

func SendEmailWithContext(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	return (*instance.EmailDelivery.IngredientInterfaceImpl.SendEmail)(input, userContext)
}

func SendSmsWithContext(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	return (*instance.SmsDelivery.IngredientInterfaceImpl.SendSms)(input, userContext)
}

func SendEmail(input emaildelivery.EmailType) error {
	return SendEmailWithContext(input, &map[string]interface{}{})
}

func SendSms(input smsdelivery.SmsType) error {
	return SendSmsWithContext(input, &map[string]interface{}{})
}

func ThirdPartySignInUp(thirdPartyID string, thirdPartyUserID string, email tplmodels.EmailStruct) (tplmodels.ThirdPartySignInUp, error) {
	return ThirdPartySignInUpWithContext(thirdPartyID, thirdPartyUserID, email, &map[string]interface{}{})
}

func GetUserByThirdPartyInfo(thirdPartyID string, thirdPartyUserID string) (*tplmodels.User, error) {
	return GetUserByThirdPartyInfoWithContext(thirdPartyID, thirdPartyUserID, &map[string]interface{}{})
}

func GetUserById(userID string) (*tplmodels.User, error) {
	return GetUserByIDWithContext(userID, &map[string]interface{}{})
}

func GetUsersByEmail(email string) ([]tplmodels.User, error) {
	return GetUsersByEmailWithContext(email, &map[string]interface{}{})
}

func CreateEmailVerificationToken(userID string) (evmodels.CreateEmailVerificationTokenResponse, error) {
	return CreateEmailVerificationTokenWithContext(userID, &map[string]interface{}{})
}

func VerifyEmailUsingToken(token string) (*tplmodels.User, error) {
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

func CreateCodeWithEmail(email string, userInputCode *string) (plessmodels.CreateCodeResponse, error) {
	return CreateCodeWithEmailWithContext(email, userInputCode, &map[string]interface{}{})
}

func CreateCodeWithPhoneNumber(phoneNumber string, userInputCode *string) (plessmodels.CreateCodeResponse, error) {
	return CreateCodeWithPhoneNumberWithContext(phoneNumber, userInputCode, &map[string]interface{}{})
}

func CreateNewCodeForDevice(deviceID string, userInputCode *string) (plessmodels.ResendCodeResponse, error) {
	return CreateNewCodeForDeviceWithContext(deviceID, userInputCode, &map[string]interface{}{})
}

func ConsumeCodeWithUserInputCode(deviceID string, userInputCode string, preAuthSessionID string) (tplmodels.ConsumeCodeResponse, error) {
	return ConsumeCodeWithUserInputCodeWithContext(deviceID, userInputCode, preAuthSessionID, &map[string]interface{}{})
}

func ConsumeCodeWithLinkCode(linkCode string, preAuthSessionID string) (tplmodels.ConsumeCodeResponse, error) {
	return ConsumeCodeWithLinkCodeWithContext(linkCode, preAuthSessionID, &map[string]interface{}{})
}

func GetUserByID(userID string) (*tplmodels.User, error) {
	return GetUserByIDWithContext(userID, &map[string]interface{}{})
}

func GetUserByPhoneNumber(phoneNumber string) (*tplmodels.User, error) {
	return GetUserByPhoneNumberWithContext(phoneNumber, &map[string]interface{}{})
}

func UpdatePasswordlessUser(userID string, email *string, phoneNumber *string) (plessmodels.UpdateUserResponse, error) {
	return UpdatePasswordlessUserWithContext(userID, email, phoneNumber, &map[string]interface{}{})
}

func DeleteEmailForPasswordlessUser(userID string) (plessmodels.DeleteUserResponse, error) {
	return DeleteEmailForPasswordlessUserWithContext(userID, &map[string]interface{}{})
}

func DeletePhoneNumberForUser(userID string) (plessmodels.DeleteUserResponse, error) {
	return DeletePhoneNumberForUserWithContext(userID, &map[string]interface{}{})
}

func RevokeAllCodesByEmail(email string) error {
	return RevokeAllCodesByEmailWithContext(email, &map[string]interface{}{})
}

func RevokeAllCodesByPhoneNumber(phoneNumber string) error {
	return RevokeAllCodesByPhoneNumberWithContext(phoneNumber, &map[string]interface{}{})
}

func RevokeCode(codeID string) error {
	return RevokeCodeWithContext(codeID, &map[string]interface{}{})
}

func ListCodesByEmail(email string) ([]plessmodels.DeviceType, error) {
	return ListCodesByEmailWithContext(email, &map[string]interface{}{})
}

func ListCodesByPhoneNumber(phoneNumber string) ([]plessmodels.DeviceType, error) {
	return ListCodesByPhoneNumberWithContext(phoneNumber, &map[string]interface{}{})
}

func ListCodesByDeviceID(deviceID string) (*plessmodels.DeviceType, error) {
	return ListCodesByDeviceIDWithContext(deviceID, &map[string]interface{}{})
}

func ListCodesByPreAuthSessionID(preAuthSessionID string) (*plessmodels.DeviceType, error) {
	return ListCodesByPreAuthSessionIDWithContext(preAuthSessionID, &map[string]interface{}{})
}

func CreateMagicLinkByEmail(email string) (string, error) {
	return CreateMagicLinkByEmailWithContext(email, &map[string]interface{}{})
}

func CreateMagicLinkByPhoneNumber(phoneNumber string) (string, error) {
	return CreateMagicLinkByPhoneNumberWithContext(phoneNumber, &map[string]interface{}{})
}

func PasswordlessSignInUpByEmail(email string) (struct {
	PreAuthSessionID string
	CreatedNewUser   bool
	User             tplmodels.User
}, error) {
	return PasswordlessSignInUpByEmailWithContext(email, &map[string]interface{}{})
}

func PasswordlessSignInUpByPhoneNumber(phoneNumber string) (struct {
	PreAuthSessionID string
	CreatedNewUser   bool
	User             tplmodels.User
}, error) {
	return PasswordlessSignInUpByPhoneNumberWithContext(phoneNumber, &map[string]interface{}{})
}

func MakeSMTPService(config emaildelivery.SMTPServiceConfig) *emaildelivery.EmailDeliveryInterface {
	return smtpService.MakeSMTPService(config)
}

func MakeTwilioService(config smsdelivery.TwilioServiceConfig) (*smsdelivery.SmsDeliveryInterface, error) {
	return twilioService.MakeTwilioService(config)
}

func MakeSupertokensService(apiKey string) *smsdelivery.SmsDeliveryInterface {
	return supertokensService.MakeSupertokensService(apiKey)
}
