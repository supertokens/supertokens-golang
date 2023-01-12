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
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/emaildelivery/smtpService"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/smsdelivery/supertokensService"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/smsdelivery/twilioService"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init(config tplmodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func ThirdPartyManuallyCreateOrUpdateUserWithContext(thirdPartyID string, thirdPartyUserID string, email string, userContext supertokens.UserContext) (tplmodels.ManuallyCreateOrUpdateUserResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return tplmodels.ManuallyCreateOrUpdateUserResponse{}, err
	}
	return (*instance.RecipeImpl.ThirdPartyManuallyCreateOrUpdateUser)(thirdPartyID, thirdPartyUserID, email, userContext)
}

func ThirdPartyGetProviderWithContext(thirdPartyID string, tenantId *string, clientType *string, userContext supertokens.UserContext) (tpmodels.GetProviderResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return tpmodels.GetProviderResponse{}, err
	}
	return (*instance.RecipeImpl.ThirdPartyGetProvider)(thirdPartyID, tenantId, clientType, userContext)
}

func GetUserByThirdPartyInfoWithContext(thirdPartyID string, thirdPartyUserID string, userContext supertokens.UserContext) (*tplmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByThirdPartyInfo)(thirdPartyID, thirdPartyUserID, userContext)
}

func GetUserByIdWithContext(userID string, userContext supertokens.UserContext) (*tplmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, userContext)
}

func GetUsersByEmailWithContext(email string, userContext supertokens.UserContext) ([]tplmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUsersByEmail)(email, userContext)
}

func CreateCodeWithEmailWithContext(email string, userInputCode *string, userContext supertokens.UserContext) (plessmodels.CreateCodeResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.CreateCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.CreateCode)(&email, nil, userInputCode, userContext)
}

func CreateCodeWithPhoneNumberWithContext(phoneNumber string, userInputCode *string, userContext supertokens.UserContext) (plessmodels.CreateCodeResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.CreateCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.CreateCode)(nil, &phoneNumber, userInputCode, userContext)
}

func CreateNewCodeForDeviceWithContext(deviceID string, userInputCode *string, userContext supertokens.UserContext) (plessmodels.ResendCodeResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.ResendCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.CreateNewCodeForDevice)(deviceID, userInputCode, userContext)
}

func ConsumeCodeWithUserInputCodeWithContext(deviceID string, userInputCode string, preAuthSessionID string, userContext supertokens.UserContext) (tplmodels.ConsumeCodeResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
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
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return tplmodels.ConsumeCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ConsumeCode)(nil, &linkCode, preAuthSessionID, userContext)
}

func GetUserByIDWithContext(userID string, userContext supertokens.UserContext) (*tplmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, userContext)
}

func GetUserByPhoneNumberWithContext(phoneNumber string, userContext supertokens.UserContext) (*tplmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.GetUserByPhoneNumber)(phoneNumber, userContext)
}

func UpdatePasswordlessUserWithContext(userID string, email *string, phoneNumber *string, userContext supertokens.UserContext) (plessmodels.UpdateUserResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.UpdateUserResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.UpdatePasswordlessUser)(userID, email, phoneNumber, userContext)
}

func DeleteEmailForPasswordlessUserWithContext(userID string, userContext supertokens.UserContext) (plessmodels.DeleteUserResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.DeleteUserResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.DeleteEmailForPasswordlessUser)(userID, userContext)
}

func DeletePhoneNumberForUserWithContext(userID string, userContext supertokens.UserContext) (plessmodels.DeleteUserResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.DeleteUserResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.DeletePhoneNumberForUser)(userID, userContext)
}

func RevokeAllCodesByEmailWithContext(email string, userContext supertokens.UserContext) error {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.RevokeAllCodes)(&email, nil, userContext)
}

func RevokeAllCodesByPhoneNumberWithContext(phoneNumber string, userContext supertokens.UserContext) error {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.RevokeAllCodes)(nil, &phoneNumber, userContext)
}

func RevokeCodeWithContext(codeID string, userContext supertokens.UserContext) error {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.RevokeCode)(codeID, userContext)
}

func ListCodesByEmailWithContext(email string, userContext supertokens.UserContext) ([]plessmodels.DeviceType, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return []plessmodels.DeviceType{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ListCodesByEmail)(email, userContext)
}

func ListCodesByPhoneNumberWithContext(phoneNumber string, userContext supertokens.UserContext) ([]plessmodels.DeviceType, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return []plessmodels.DeviceType{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ListCodesByPhoneNumber)(phoneNumber, userContext)
}

func ListCodesByDeviceIDWithContext(deviceID string, userContext supertokens.UserContext) (*plessmodels.DeviceType, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ListCodesByDeviceID)(deviceID, userContext)
}

func ListCodesByPreAuthSessionIDWithContext(preAuthSessionID string, userContext supertokens.UserContext) (*plessmodels.DeviceType, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ListCodesByPreAuthSessionID)(preAuthSessionID, userContext)
}

func CreateMagicLinkByEmailWithContext(email string, userContext supertokens.UserContext) (string, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return "", err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return instance.passwordlessRecipe.CreateMagicLink(&email, nil, userContext)
}

func CreateMagicLinkByPhoneNumberWithContext(phoneNumber string, userContext supertokens.UserContext) (string, error) {
	instance, err := GetRecipeInstanceOrThrowError()
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
	instance, err := GetRecipeInstanceOrThrowError()
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
	instance, err := GetRecipeInstanceOrThrowError()
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
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	return (*instance.EmailDelivery.IngredientInterfaceImpl.SendEmail)(input, userContext)
}

func SendSmsWithContext(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
	instance, err := GetRecipeInstanceOrThrowError()
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

func ThirdPartyManuallyCreateOrUpdateUser(thirdPartyID string, thirdPartyUserID string, email string) (tplmodels.ManuallyCreateOrUpdateUserResponse, error) {
	return ThirdPartyManuallyCreateOrUpdateUserWithContext(thirdPartyID, thirdPartyUserID, email, &map[string]interface{}{})
}

func ThirdPartyGetProvider(thirdPartyID string, tenantId *string, clientType *string) (tpmodels.GetProviderResponse, error) {
	return ThirdPartyGetProviderWithContext(thirdPartyID, tenantId, clientType, &map[string]interface{}{})
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

func MakeSupertokensSMSService(apiKey string) *smsdelivery.SmsDeliveryInterface {
	return supertokensService.MakeSupertokensSMSService(apiKey)
}
