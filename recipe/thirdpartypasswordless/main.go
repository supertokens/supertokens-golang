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

func ThirdPartyManuallyCreateOrUpdateUserWithContext(thirdPartyID string, thirdPartyUserID string, email string, tenantId *string, userContext supertokens.UserContext) (tplmodels.ManuallyCreateOrUpdateUserResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return tplmodels.ManuallyCreateOrUpdateUserResponse{}, err
	}
	return (*instance.RecipeImpl.ThirdPartyManuallyCreateOrUpdateUser)(thirdPartyID, thirdPartyUserID, email, tenantId, userContext)
}

func ThirdPartyGetProviderWithContext(thirdPartyID string, tenantId *string, clientType *string, userContext supertokens.UserContext) (tpmodels.GetProviderResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return tpmodels.GetProviderResponse{}, err
	}
	return (*instance.RecipeImpl.ThirdPartyGetProvider)(thirdPartyID, tenantId, clientType, userContext)
}

func GetUserByThirdPartyInfoWithContext(thirdPartyID string, thirdPartyUserID string, tenantId *string, userContext supertokens.UserContext) (*tplmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByThirdPartyInfo)(thirdPartyID, thirdPartyUserID, tenantId, userContext)
}

func GetUserByIdWithContext(userID string, tenantId *string, userContext supertokens.UserContext) (*tplmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, tenantId, userContext)
}

func GetUsersByEmailWithContext(email string, tenantId *string, userContext supertokens.UserContext) ([]tplmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	return (*instance.RecipeImpl.GetUsersByEmail)(email, tenantId, userContext)
}

func CreateCodeWithEmailWithContext(email string, userInputCode *string, tenantId *string, userContext supertokens.UserContext) (plessmodels.CreateCodeResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.CreateCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.CreateCode)(&email, nil, userInputCode, tenantId, userContext)
}

func CreateCodeWithPhoneNumberWithContext(phoneNumber string, userInputCode *string, tenantId *string, userContext supertokens.UserContext) (plessmodels.CreateCodeResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.CreateCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.CreateCode)(nil, &phoneNumber, userInputCode, tenantId, userContext)
}

func CreateNewCodeForDeviceWithContext(deviceID string, userInputCode *string, tenantId *string, userContext supertokens.UserContext) (plessmodels.ResendCodeResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.ResendCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.CreateNewCodeForDevice)(deviceID, userInputCode, tenantId, userContext)
}

func ConsumeCodeWithUserInputCodeWithContext(deviceID string, userInputCode string, preAuthSessionID string, tenantId *string, userContext supertokens.UserContext) (tplmodels.ConsumeCodeResponse, error) {
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
	}, nil, preAuthSessionID, tenantId, userContext)
}

func ConsumeCodeWithLinkCodeWithContext(linkCode string, preAuthSessionID string, tenantId *string, userContext supertokens.UserContext) (tplmodels.ConsumeCodeResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return tplmodels.ConsumeCodeResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ConsumeCode)(nil, &linkCode, preAuthSessionID, tenantId, userContext)
}

func GetUserByIDWithContext(userID string, tenantId *string, userContext supertokens.UserContext) (*tplmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, tenantId, userContext)
}

func GetUserByPhoneNumberWithContext(phoneNumber string, tenantId *string, userContext supertokens.UserContext) (*tplmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.GetUserByPhoneNumber)(phoneNumber, tenantId, userContext)
}

func UpdatePasswordlessUserWithContext(userID string, email *string, phoneNumber *string, tenantId *string, userContext supertokens.UserContext) (plessmodels.UpdateUserResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.UpdateUserResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.UpdatePasswordlessUser)(userID, email, phoneNumber, tenantId, userContext)
}

func DeleteEmailForPasswordlessUserWithContext(userID string, tenantId *string, userContext supertokens.UserContext) (plessmodels.DeleteUserResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.DeleteUserResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.DeleteEmailForPasswordlessUser)(userID, tenantId, userContext)
}

func DeletePhoneNumberForUserWithContext(userID string, tenantId *string, userContext supertokens.UserContext) (plessmodels.DeleteUserResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.DeleteUserResponse{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.DeletePhoneNumberForUser)(userID, tenantId, userContext)
}

func RevokeAllCodesByEmailWithContext(email string, tenantId *string, userContext supertokens.UserContext) error {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.RevokeAllCodes)(&email, nil, tenantId, userContext)
}

func RevokeAllCodesByPhoneNumberWithContext(phoneNumber string, tenantId *string, userContext supertokens.UserContext) error {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.RevokeAllCodes)(nil, &phoneNumber, tenantId, userContext)
}

func RevokeCodeWithContext(codeID string, tenantId *string, userContext supertokens.UserContext) error {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.RevokeCode)(codeID, tenantId, userContext)
}

func ListCodesByEmailWithContext(email string, tenantId *string, userContext supertokens.UserContext) ([]plessmodels.DeviceType, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return []plessmodels.DeviceType{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ListCodesByEmail)(email, tenantId, userContext)
}

func ListCodesByPhoneNumberWithContext(phoneNumber string, tenantId *string, userContext supertokens.UserContext) ([]plessmodels.DeviceType, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return []plessmodels.DeviceType{}, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ListCodesByPhoneNumber)(phoneNumber, tenantId, userContext)
}

func ListCodesByDeviceIDWithContext(deviceID string, tenantId *string, userContext supertokens.UserContext) (*plessmodels.DeviceType, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ListCodesByDeviceID)(deviceID, tenantId, userContext)
}

func ListCodesByPreAuthSessionIDWithContext(preAuthSessionID string, tenantId *string, userContext supertokens.UserContext) (*plessmodels.DeviceType, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return (*instance.RecipeImpl.ListCodesByPreAuthSessionID)(preAuthSessionID, tenantId, userContext)
}

func CreateMagicLinkByEmailWithContext(email string, tenantId *string, userContext supertokens.UserContext) (string, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return "", err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return instance.passwordlessRecipe.CreateMagicLink(&email, nil, tenantId, userContext)
}

func CreateMagicLinkByPhoneNumberWithContext(phoneNumber string, tenantId *string, userContext supertokens.UserContext) (string, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return "", err
	}
	if userContext == nil {
		userContext = &map[string]interface{}{}
	}
	return instance.passwordlessRecipe.CreateMagicLink(nil, &phoneNumber, tenantId, userContext)
}

func PasswordlessSignInUpByEmailWithContext(email string, tenantId *string, userContext supertokens.UserContext) (struct {
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
	resp, err := instance.passwordlessRecipe.SignInUp(&email, nil, tenantId, userContext)
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

func PasswordlessSignInUpByPhoneNumberWithContext(phoneNumber string, tenantId *string, userContext supertokens.UserContext) (struct {
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
	resp, err := instance.passwordlessRecipe.SignInUp(nil, &phoneNumber, tenantId, userContext)
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

func ThirdPartyManuallyCreateOrUpdateUser(thirdPartyID string, thirdPartyUserID string, email string, tenantId *string) (tplmodels.ManuallyCreateOrUpdateUserResponse, error) {
	return ThirdPartyManuallyCreateOrUpdateUserWithContext(thirdPartyID, thirdPartyUserID, email, tenantId, &map[string]interface{}{})
}

func ThirdPartyGetProvider(thirdPartyID string, tenantId *string, clientType *string) (tpmodels.GetProviderResponse, error) {
	return ThirdPartyGetProviderWithContext(thirdPartyID, tenantId, clientType, &map[string]interface{}{})
}

func GetUserByThirdPartyInfo(thirdPartyID string, thirdPartyUserID string, tenantId *string) (*tplmodels.User, error) {
	return GetUserByThirdPartyInfoWithContext(thirdPartyID, thirdPartyUserID, tenantId, &map[string]interface{}{})
}

func GetUserById(userID string, tenantId *string) (*tplmodels.User, error) {
	return GetUserByIDWithContext(userID, tenantId, &map[string]interface{}{})
}

func GetUsersByEmail(email string, tenantId *string) ([]tplmodels.User, error) {
	return GetUsersByEmailWithContext(email, tenantId, &map[string]interface{}{})
}

func CreateCodeWithEmail(email string, userInputCode *string, tenantId *string) (plessmodels.CreateCodeResponse, error) {
	return CreateCodeWithEmailWithContext(email, userInputCode, tenantId, &map[string]interface{}{})
}

func CreateCodeWithPhoneNumber(phoneNumber string, userInputCode *string, tenantId *string) (plessmodels.CreateCodeResponse, error) {
	return CreateCodeWithPhoneNumberWithContext(phoneNumber, userInputCode, tenantId, &map[string]interface{}{})
}

func CreateNewCodeForDevice(deviceID string, userInputCode *string, tenantId *string) (plessmodels.ResendCodeResponse, error) {
	return CreateNewCodeForDeviceWithContext(deviceID, userInputCode, tenantId, &map[string]interface{}{})
}

func ConsumeCodeWithUserInputCode(deviceID string, userInputCode string, preAuthSessionID string, tenantId *string) (tplmodels.ConsumeCodeResponse, error) {
	return ConsumeCodeWithUserInputCodeWithContext(deviceID, userInputCode, preAuthSessionID, tenantId, &map[string]interface{}{})
}

func ConsumeCodeWithLinkCode(linkCode string, preAuthSessionID string, tenantId *string) (tplmodels.ConsumeCodeResponse, error) {
	return ConsumeCodeWithLinkCodeWithContext(linkCode, preAuthSessionID, tenantId, &map[string]interface{}{})
}

func GetUserByID(userID string, tenantId *string) (*tplmodels.User, error) {
	return GetUserByIDWithContext(userID, tenantId, &map[string]interface{}{})
}

func GetUserByPhoneNumber(phoneNumber string, tenantId *string) (*tplmodels.User, error) {
	return GetUserByPhoneNumberWithContext(phoneNumber, tenantId, &map[string]interface{}{})
}

func UpdatePasswordlessUser(userID string, email *string, phoneNumber *string, tenantId *string) (plessmodels.UpdateUserResponse, error) {
	return UpdatePasswordlessUserWithContext(userID, email, phoneNumber, tenantId, &map[string]interface{}{})
}

func DeleteEmailForPasswordlessUser(userID string, tenantId *string) (plessmodels.DeleteUserResponse, error) {
	return DeleteEmailForPasswordlessUserWithContext(userID, tenantId, &map[string]interface{}{})
}

func DeletePhoneNumberForUser(userID string, tenantId *string) (plessmodels.DeleteUserResponse, error) {
	return DeletePhoneNumberForUserWithContext(userID, tenantId, &map[string]interface{}{})
}

func RevokeAllCodesByEmail(email string, tenantId *string) error {
	return RevokeAllCodesByEmailWithContext(email, tenantId, &map[string]interface{}{})
}

func RevokeAllCodesByPhoneNumber(phoneNumber string, tenantId *string) error {
	return RevokeAllCodesByPhoneNumberWithContext(phoneNumber, tenantId, &map[string]interface{}{})
}

func RevokeCode(codeID string, tenantId *string) error {
	return RevokeCodeWithContext(codeID, tenantId, &map[string]interface{}{})
}

func ListCodesByEmail(email string, tenantId *string) ([]plessmodels.DeviceType, error) {
	return ListCodesByEmailWithContext(email, tenantId, &map[string]interface{}{})
}

func ListCodesByPhoneNumber(phoneNumber string, tenantId *string) ([]plessmodels.DeviceType, error) {
	return ListCodesByPhoneNumberWithContext(phoneNumber, tenantId, &map[string]interface{}{})
}

func ListCodesByDeviceID(deviceID string, tenantId *string) (*plessmodels.DeviceType, error) {
	return ListCodesByDeviceIDWithContext(deviceID, tenantId, &map[string]interface{}{})
}

func ListCodesByPreAuthSessionID(preAuthSessionID string, tenantId *string) (*plessmodels.DeviceType, error) {
	return ListCodesByPreAuthSessionIDWithContext(preAuthSessionID, tenantId, &map[string]interface{}{})
}

func CreateMagicLinkByEmail(email string, tenantId *string) (string, error) {
	return CreateMagicLinkByEmailWithContext(email, tenantId, &map[string]interface{}{})
}

func CreateMagicLinkByPhoneNumber(phoneNumber string, tenantId *string) (string, error) {
	return CreateMagicLinkByPhoneNumberWithContext(phoneNumber, tenantId, &map[string]interface{}{})
}

func PasswordlessSignInUpByEmail(email string, tenantId *string) (struct {
	PreAuthSessionID string
	CreatedNewUser   bool
	User             tplmodels.User
}, error) {
	return PasswordlessSignInUpByEmailWithContext(email, tenantId, &map[string]interface{}{})
}

func PasswordlessSignInUpByPhoneNumber(phoneNumber string, tenantId *string) (struct {
	PreAuthSessionID string
	CreatedNewUser   bool
	User             tplmodels.User
}, error) {
	return PasswordlessSignInUpByPhoneNumberWithContext(phoneNumber, tenantId, &map[string]interface{}{})
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
