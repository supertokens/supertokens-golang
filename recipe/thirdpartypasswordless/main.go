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

func ThirdPartyManuallyCreateOrUpdateUser(tenantId string, thirdPartyID string, thirdPartyUserID string, email string, userContext ...supertokens.UserContext) (tplmodels.ManuallyCreateOrUpdateUserResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return tplmodels.ManuallyCreateOrUpdateUserResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.ThirdPartyManuallyCreateOrUpdateUser)(thirdPartyID, thirdPartyUserID, email, tenantId, userContext[0])
}

func ThirdPartyGetProvider(tenantId string, thirdPartyID string, clientType *string, userContext ...supertokens.UserContext) (*tpmodels.TypeProvider, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.ThirdPartyGetProvider)(thirdPartyID, clientType, tenantId, userContext[0])
}

func GetUserByThirdPartyInfo(tenantId string, thirdPartyID string, thirdPartyUserID string, userContext ...supertokens.UserContext) (*tplmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetUserByThirdPartyInfo)(thirdPartyID, thirdPartyUserID, tenantId, userContext[0])
}

func GetUserById(userID string, userContext ...supertokens.UserContext) (*tplmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, userContext[0])
}

func GetUsersByEmail(tenantId string, email string, userContext ...supertokens.UserContext) ([]tplmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetUsersByEmail)(email, tenantId, userContext[0])
}

func CreateCodeWithEmail(tenantId string, email string, userInputCode *string, userContext ...supertokens.UserContext) (plessmodels.CreateCodeResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.CreateCodeResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.CreateCode)(&email, nil, userInputCode, tenantId, userContext[0])
}

func CreateCodeWithPhoneNumber(tenantId string, phoneNumber string, userInputCode *string, userContext ...supertokens.UserContext) (plessmodels.CreateCodeResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.CreateCodeResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.CreateCode)(nil, &phoneNumber, userInputCode, tenantId, userContext[0])
}

func CreateNewCodeForDevice(tenantId string, deviceID string, userInputCode *string, userContext ...supertokens.UserContext) (plessmodels.ResendCodeResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.ResendCodeResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.CreateNewCodeForDevice)(deviceID, userInputCode, tenantId, userContext[0])
}

func ConsumeCodeWithUserInputCode(tenantId string, deviceID string, userInputCode string, preAuthSessionID string, userContext ...supertokens.UserContext) (tplmodels.ConsumeCodeResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return tplmodels.ConsumeCodeResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.ConsumeCode)(&plessmodels.UserInputCodeWithDeviceID{
		Code:     userInputCode,
		DeviceID: deviceID,
	}, nil, preAuthSessionID, tenantId, userContext[0])
}

func ConsumeCodeWithLinkCode(tenantId string, linkCode string, preAuthSessionID string, userContext ...supertokens.UserContext) (tplmodels.ConsumeCodeResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return tplmodels.ConsumeCodeResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.ConsumeCode)(nil, &linkCode, preAuthSessionID, tenantId, userContext[0])
}

func GetUserByID(userID string, userContext ...supertokens.UserContext) (*tplmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetUserByID)(userID, userContext[0])
}

func GetUserByPhoneNumber(tenantId string, phoneNumber string, userContext ...supertokens.UserContext) (*tplmodels.User, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetUserByPhoneNumber)(phoneNumber, tenantId, userContext[0])
}

func UpdatePasswordlessUser(userID string, email *string, phoneNumber *string, userContext ...supertokens.UserContext) (plessmodels.UpdateUserResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.UpdateUserResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.UpdatePasswordlessUser)(userID, email, phoneNumber, userContext[0])
}

func DeleteEmailForPasswordlessUser(userID string, userContext ...supertokens.UserContext) (plessmodels.DeleteUserResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.DeleteUserResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.DeleteEmailForPasswordlessUser)(userID, userContext[0])
}

func DeletePhoneNumberForUser(userID string, userContext ...supertokens.UserContext) (plessmodels.DeleteUserResponse, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return plessmodels.DeleteUserResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.DeletePhoneNumberForUser)(userID, userContext[0])
}

func RevokeAllCodesByEmail(tenantId string, email string, userContext ...supertokens.UserContext) error {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.RevokeAllCodes)(&email, nil, tenantId, userContext[0])
}

func RevokeAllCodesByPhoneNumber(tenantId string, phoneNumber string, userContext ...supertokens.UserContext) error {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.RevokeAllCodes)(nil, &phoneNumber, tenantId, userContext[0])
}

func RevokeCode(tenantId string, codeID string, userContext ...supertokens.UserContext) error {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.RevokeCode)(codeID, tenantId, userContext[0])
}

func ListCodesByEmail(tenantId string, email string, userContext ...supertokens.UserContext) ([]plessmodels.DeviceType, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return []plessmodels.DeviceType{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.ListCodesByEmail)(email, tenantId, userContext[0])
}

func ListCodesByPhoneNumber(tenantId string, phoneNumber string, userContext ...supertokens.UserContext) ([]plessmodels.DeviceType, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return []plessmodels.DeviceType{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.ListCodesByPhoneNumber)(phoneNumber, tenantId, userContext[0])
}

func ListCodesByDeviceID(tenantId string, deviceID string, userContext ...supertokens.UserContext) (*plessmodels.DeviceType, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.ListCodesByDeviceID)(deviceID, tenantId, userContext[0])
}

func ListCodesByPreAuthSessionID(tenantId string, preAuthSessionID string, userContext ...supertokens.UserContext) (*plessmodels.DeviceType, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return nil, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.ListCodesByPreAuthSessionID)(preAuthSessionID, tenantId, userContext[0])
}

func CreateMagicLinkByEmail(tenantId string, email string, userContext ...supertokens.UserContext) (string, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return "", err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return instance.passwordlessRecipe.CreateMagicLink(&email, nil, tenantId, userContext[0])
}

func CreateMagicLinkByPhoneNumber(tenantId string, phoneNumber string, userContext ...supertokens.UserContext) (string, error) {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return "", err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return instance.passwordlessRecipe.CreateMagicLink(nil, &phoneNumber, tenantId, userContext[0])
}

func PasswordlessSignInUpByEmail(tenantId string, email string, userContext ...supertokens.UserContext) (struct {
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
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	resp, err := instance.passwordlessRecipe.SignInUp(&email, nil, tenantId, userContext[0])
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

func PasswordlessSignInUpByPhoneNumber(tenantId string, phoneNumber string, userContext ...supertokens.UserContext) (struct {
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
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	resp, err := instance.passwordlessRecipe.SignInUp(nil, &phoneNumber, tenantId, userContext[0])
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

func SendEmail(input emaildelivery.EmailType, userContext ...supertokens.UserContext) error {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.EmailDelivery.IngredientInterfaceImpl.SendEmail)(input, userContext[0])
}

func SendSms(input smsdelivery.SmsType, userContext ...supertokens.UserContext) error {
	instance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		return err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.SmsDelivery.IngredientInterfaceImpl.SendSms)(input, userContext[0])
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
