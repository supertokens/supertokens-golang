/* Copyright (c) 2025, VRAI Labs and/or its affiliates. All rights reserved.
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

package webauthn

import (
	"encoding/json"
	"fmt"

	"github.com/goforj/godump"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/webauthn/webauthnmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeRecipeImplementation(querier supertokens.Querier, recipeEmailDelivery emaildelivery.Ingredient) webauthnmodels.RecipeInterface {

	getGeneratedOptions := func(
		webauthnGeneratedOptionsId string,
		tenantId string,
		userContext supertokens.UserContext,
	) (webauthnmodels.GetGeneratedOptionsResponse, error) {
		resp, err := querier.SendGetRequest(
			fmt.Sprintf("/%s/recipe/webauthn/options", tenantId),
			map[string]string{"webauthnGeneratedOptionsId": webauthnGeneratedOptionsId},
			userContext,
		)
		if err != nil {
			return webauthnmodels.GetGeneratedOptionsResponse{}, err
		}
		status := resp["status"].(string)
		if status == "OPTIONS_NOT_FOUND_ERROR" {
			return webauthnmodels.GetGeneratedOptionsResponse{
				OptionsNotFoundError: &struct{}{},
			}, nil
		}
		result := webauthnmodels.GetGeneratedOptionsResponse{}
		result.OK = &struct {
			WebauthnGeneratedOptionsId string
			CreatedAt                  int64
			ExpiresAt                  int64
			Email                      *string
			RelyingPartyId             string
			RelyingPartyName           string
			Origin                     string
			Challenge                  string
			Timeout                    int
			UserVerification           webauthnmodels.UserVerification
			UserPresence               bool
		}{
			WebauthnGeneratedOptionsId: resp["webauthnGeneratedOptionsId"].(string),
			CreatedAt:                  int64(resp["createdAt"].(float64)),
			ExpiresAt:                  int64(resp["expiresAt"].(float64)),
			RelyingPartyId:             resp["relyingPartyId"].(string),
			RelyingPartyName:           resp["relyingPartyName"].(string),
			Origin:                     resp["origin"].(string),
			Challenge:                  resp["challenge"].(string),
			Timeout:                    int(resp["timeout"].(float64)),
			UserVerification:           webauthnmodels.UserVerification(resp["userVerification"].(string)),
		}
		if userPresence, ok := resp["userPresence"].(bool); ok {
			result.OK.UserPresence = userPresence
		}
		if emailVal, ok := resp["email"]; ok && emailVal != nil {
			email := emailVal.(string)
			result.OK.Email = &email
		}
		return result, nil
	}

	registerOptions := func(
		email *string,
		recoverAccountToken *string,
		relyingPartyId string,
		relyingPartyName string,
		origin string,
		timeout *int,
		attestation *webauthnmodels.Attestation,
		residentKey *webauthnmodels.ResidentKey,
		userVerification *webauthnmodels.UserVerification,
		userPresence *bool,
		supportedAlgorithmIds []webauthnmodels.COSEAlgorithmIdentifier,
		tenantId string,
		userContext supertokens.UserContext,
	) (webauthnmodels.RegisterOptionsResponse, error) {
		body := map[string]interface{}{
			"relyingPartyId":   relyingPartyId,
			"relyingPartyName": relyingPartyName,
			"origin":           origin,
		}
		if email != nil {
			body["email"] = *email
		}
		if recoverAccountToken != nil {
			body["recoverAccountToken"] = *recoverAccountToken
		}
		if timeout != nil {
			body["timeout"] = *timeout
		}
		if attestation != nil {
			body["attestation"] = string(*attestation)
		}
		if residentKey != nil {
			body["residentKey"] = string(*residentKey)
		}
		if userVerification != nil {
			body["userVerification"] = string(*userVerification)
		}
		if userPresence != nil {
			body["userPresence"] = *userPresence
		}
		if supportedAlgorithmIds != nil {
			body["supportedAlgorithmIds"] = supportedAlgorithmIds
		}

		godump.Dump(body)

		resp, err := querier.SendPostRequest(
			fmt.Sprintf("/%s/recipe/webauthn/options/register", tenantId),
			body,
			userContext,
		)
		if err != nil {
			return webauthnmodels.RegisterOptionsResponse{}, err
		}
		godump.Dump(resp)
		status := resp["status"].(string)
		if status == "RECOVER_ACCOUNT_TOKEN_INVALID_ERROR" {
			return webauthnmodels.RegisterOptionsResponse{
				RecoverAccountTokenInvalidError: &struct{}{},
			}, nil
		}
		if status == "INVALID_EMAIL_ERROR" {
			return webauthnmodels.RegisterOptionsResponse{
				InvalidEmailError: &struct{ Err string }{Err: resp["err"].(string)},
			}, nil
		}
		if status == "INVALID_OPTIONS_ERROR" {
			return webauthnmodels.RegisterOptionsResponse{
				InvalidOptionsError: &struct{}{},
			}, nil
		}
		return buildRegisterOptionsResponse(resp)
	}

	signInOptions := func(
		relyingPartyId string,
		origin string,
		timeout *int,
		userVerification *webauthnmodels.UserVerification,
		userPresence *bool,
		tenantId string,
		userContext supertokens.UserContext,
	) (webauthnmodels.SignInOptionsResponse, error) {
		body := map[string]interface{}{
			"relyingPartyName": "something",
			"relyingPartyId":   relyingPartyId,
			"origin":           origin,
		}
		if timeout != nil {
			body["timeout"] = *timeout
		}
		if userVerification != nil {
			body["userVerification"] = string(*userVerification)
		}
		if userPresence != nil {
			body["userPresence"] = *userPresence
		}

		resp, err := querier.SendPostRequest(
			fmt.Sprintf("/%s/recipe/webauthn/options/signin", tenantId),
			body,
			userContext,
		)
		if err != nil {
			return webauthnmodels.SignInOptionsResponse{}, err
		}
		status := resp["status"].(string)
		if status == "INVALID_OPTIONS_ERROR" {
			return webauthnmodels.SignInOptionsResponse{
				InvalidOptionsError: &struct{}{},
			}, nil
		}

		if status != "OK" {
			return webauthnmodels.SignInOptionsResponse{}, fmt.Errorf("unexpected status: %s", status)
		}

		return webauthnmodels.SignInOptionsResponse{
			OK: &struct {
				WebauthnGeneratedOptionsId string
				CreatedAt                  int64
				ExpiresAt                  int64
				RpId                       string
				Challenge                  string
				Timeout                    int
				UserVerification           webauthnmodels.UserVerification
			}{
				WebauthnGeneratedOptionsId: resp["webauthnGeneratedOptionsId"].(string),
				CreatedAt:                  int64(resp["createdAt"].(float64)),
				ExpiresAt:                  int64(resp["expiresAt"].(float64)),
				RpId:                       resp["relyingPartyId"].(string),
				Challenge:                  resp["challenge"].(string),
				Timeout:                    int(resp["timeout"].(float64)),
				UserVerification:           webauthnmodels.UserVerification(resp["userVerification"].(string)),
			},
		}, nil
	}

	signUp := func(
		webauthnGeneratedOptionsId string,
		credential webauthnmodels.RegistrationPayload,
		tenantId string,
		userContext supertokens.UserContext,
	) (webauthnmodels.SignUpResponse, error) {
		credentialMap, err := credentialToMap(credential)
		if err != nil {
			return webauthnmodels.SignUpResponse{}, err
		}
		body := map[string]interface{}{
			"webauthnGeneratedOptionsId": webauthnGeneratedOptionsId,
			"credential":                 credentialMap,
		}
		resp, err := querier.SendPostRequest(
			fmt.Sprintf("/%s/recipe/webauthn/signup", tenantId),
			body,
			userContext,
		)
		if err != nil {
			return webauthnmodels.SignUpResponse{}, err
		}
		return parseSignUpResponse(resp)
	}

	signIn := func(
		webauthnGeneratedOptionsId string,
		credential webauthnmodels.AuthenticationPayload,
		tenantId string,
		userContext supertokens.UserContext,
	) (webauthnmodels.SignInResponse, error) {
		credentialMap, err := credentialToMap(credential)
		if err != nil {
			return webauthnmodels.SignInResponse{}, err
		}
		body := map[string]interface{}{
			"webauthnGeneratedOptionsId": webauthnGeneratedOptionsId,
			"credential":                 credentialMap,
		}
		resp, err := querier.SendPostRequest(
			fmt.Sprintf("/%s/recipe/webauthn/signin", tenantId),
			body,
			userContext,
		)
		if err != nil {
			return webauthnmodels.SignInResponse{}, err
		}
		return parseSignInResponse(resp)
	}

	verifyCredentials := func(
		webauthnGeneratedOptionsId string,
		credential webauthnmodels.AuthenticationPayload,
		tenantId string,
		userContext supertokens.UserContext,
	) (webauthnmodels.VerifyCredentialsResponse, error) {
		credentialMap, err := credentialToMap(credential)
		if err != nil {
			return webauthnmodels.VerifyCredentialsResponse{}, err
		}
		body := map[string]interface{}{
			"webauthnGeneratedOptionsId": webauthnGeneratedOptionsId,
			"credential":                 credentialMap,
		}
		resp, err := querier.SendPostRequest(
			fmt.Sprintf("/%s/recipe/webauthn/credential/verify", tenantId),
			body,
			userContext,
		)
		if err != nil {
			return webauthnmodels.VerifyCredentialsResponse{}, err
		}
		status := resp["status"].(string)
		switch status {
		case "OK":
			return webauthnmodels.VerifyCredentialsResponse{OK: &struct{}{}}, nil
		case "INVALID_CREDENTIALS_ERROR":
			return webauthnmodels.VerifyCredentialsResponse{InvalidCredentialsError: &struct{}{}}, nil
		case "INVALID_OPTIONS_ERROR":
			return webauthnmodels.VerifyCredentialsResponse{InvalidOptionsError: &struct{}{}}, nil
		case "INVALID_AUTHENTICATOR_ERROR":
			return webauthnmodels.VerifyCredentialsResponse{InvalidAuthenticatorError: &struct{ Reason string }{Reason: resp["reason"].(string)}}, nil
		case "CREDENTIAL_NOT_FOUND_ERROR":
			return webauthnmodels.VerifyCredentialsResponse{CredentialNotFoundError: &struct{}{}}, nil
		case "OPTIONS_NOT_FOUND_ERROR":
			return webauthnmodels.VerifyCredentialsResponse{OptionsNotFoundError: &struct{}{}}, nil
		}
		return webauthnmodels.VerifyCredentialsResponse{}, fmt.Errorf("unknown status: %s", status)
	}

	generateRecoverAccountToken := func(
		userId string,
		email string,
		tenantId string,
		userContext supertokens.UserContext,
	) (webauthnmodels.GenerateRecoverAccountTokenResponse, error) {
		resp, err := querier.SendPostRequest(
			fmt.Sprintf("/%s/recipe/webauthn/user/recover/token", tenantId),
			map[string]interface{}{
				"userId": userId,
				"email":  email,
			},
			userContext,
		)
		if err != nil {
			return webauthnmodels.GenerateRecoverAccountTokenResponse{}, err
		}
		status := resp["status"].(string)
		if status == "UNKNOWN_USER_ID_ERROR" {
			return webauthnmodels.GenerateRecoverAccountTokenResponse{
				UnknownUserIdError: &struct{}{},
			}, nil
		}
		return webauthnmodels.GenerateRecoverAccountTokenResponse{
			OK: &struct{ Token string }{Token: resp["token"].(string)},
		}, nil
	}

	consumeRecoverAccountToken := func(
		token string,
		tenantId string,
		userContext supertokens.UserContext,
	) (webauthnmodels.ConsumeRecoverAccountTokenResponse, error) {
		resp, err := querier.SendPostRequest(
			fmt.Sprintf("/%s/recipe/webauthn/user/recover/token/consume", tenantId),
			map[string]interface{}{"token": token},
			userContext,
		)
		if err != nil {
			return webauthnmodels.ConsumeRecoverAccountTokenResponse{}, err
		}
		status := resp["status"].(string)
		if status == "RECOVER_ACCOUNT_TOKEN_INVALID_ERROR" {
			return webauthnmodels.ConsumeRecoverAccountTokenResponse{
				RecoverAccountTokenInvalidError: &struct{}{},
			}, nil
		}
		return webauthnmodels.ConsumeRecoverAccountTokenResponse{
			OK: &struct {
				Email  string
				UserId string
			}{
				Email:  resp["email"].(string),
				UserId: resp["userId"].(string),
			},
		}, nil
	}

	registerCredential := func(
		recipeUserId string,
		webauthnGeneratedOptionsId string,
		credential webauthnmodels.RegistrationPayload,
		userContext supertokens.UserContext,
	) (webauthnmodels.RegisterCredentialResponse, error) {
		credentialMap, err := credentialToMap(credential)
		if err != nil {
			return webauthnmodels.RegisterCredentialResponse{}, err
		}
		body := map[string]interface{}{
			"recipeUserId":               recipeUserId,
			"webauthnGeneratedOptionsId": webauthnGeneratedOptionsId,
			"credential":                 credentialMap,
		}
		resp, err := querier.SendPostRequest(
			"/recipe/webauthn/user/credential/register",
			body,
			userContext,
		)
		if err != nil {
			return webauthnmodels.RegisterCredentialResponse{}, err
		}
		status := resp["status"].(string)
		switch status {
		case "OK":
			return webauthnmodels.RegisterCredentialResponse{OK: &struct{}{}}, nil
		case "INVALID_CREDENTIALS_ERROR":
			return webauthnmodels.RegisterCredentialResponse{InvalidCredentialsError: &struct{}{}}, nil
		case "INVALID_OPTIONS_ERROR":
			return webauthnmodels.RegisterCredentialResponse{InvalidOptionsError: &struct{}{}}, nil
		case "INVALID_AUTHENTICATOR_ERROR":
			return webauthnmodels.RegisterCredentialResponse{InvalidAuthenticatorError: &struct{ Reason string }{Reason: resp["reason"].(string)}}, nil
		case "OPTIONS_NOT_FOUND_ERROR":
			return webauthnmodels.RegisterCredentialResponse{OptionsNotFoundError: &struct{}{}}, nil
		}
		return webauthnmodels.RegisterCredentialResponse{}, fmt.Errorf("unknown status: %s", status)
	}

	getUserFromCredentialId := func(
		credentialId string,
		tenantId string,
		userContext supertokens.UserContext,
	) (webauthnmodels.GetUserFromCredentialIdResponse, error) {
		resp, err := querier.SendGetRequest(
			fmt.Sprintf("/%s/recipe/webauthn/user/credential", tenantId),
			map[string]string{"credentialId": credentialId},
			userContext,
		)
		if err != nil {
			return webauthnmodels.GetUserFromCredentialIdResponse{}, err
		}
		status := resp["status"].(string)
		if status == "CREDENTIAL_NOT_FOUND_ERROR" {
			return webauthnmodels.GetUserFromCredentialIdResponse{
				CredentialNotFoundError: &struct{}{},
			}, nil
		}
		user := parseUser(resp["user"].(map[string]interface{}))
		return webauthnmodels.GetUserFromCredentialIdResponse{
			OK: &struct {
				User         supertokens.User
				RecipeUserId string
			}{
				User:         user,
				RecipeUserId: resp["recipeUserId"].(string),
			},
		}, nil
	}

	getUserByEmail := func(
		email string,
		tenantId string,
		userContext supertokens.UserContext,
	) (*supertokens.User, error) {
		resp, err := querier.SendGetRequest(
			fmt.Sprintf("/%s/recipe/webauthn/user", tenantId),
			map[string]string{"email": email},
			userContext,
		)
		if err != nil {
			return nil, err
		}
		status := resp["status"].(string)
		if status != "OK" {
			return nil, nil
		}
		user := parseUser(resp["user"].(map[string]interface{}))
		return &user, nil
	}

	getUserByID := func(
		userID string,
		userContext supertokens.UserContext,
	) (*supertokens.User, error) {
		resp, err := querier.SendGetRequest(
			"/user/id",
			map[string]string{"userId": userID},
			userContext,
		)
		if err != nil {
			return nil, err
		}
		status := resp["status"].(string)
		if status != "OK" {
			return nil, nil
		}

		user := parseUser(resp["user"].(map[string]interface{}))
		return &user, nil
	}

	listCredentials := func(
		recipeUserId string,
		userContext supertokens.UserContext,
	) (webauthnmodels.ListCredentialsResponse, error) {
		resp, err := querier.SendGetRequest(
			"/recipe/webauthn/user/credential/list",
			map[string]string{"recipeUserId": recipeUserId},
			userContext,
		)
		if err != nil {
			return webauthnmodels.ListCredentialsResponse{}, err
		}
		var credentials []webauthnmodels.Credential
		if rawList, ok := resp["credentials"].([]interface{}); ok {
			for _, item := range rawList {
				if credMap, ok := item.(map[string]interface{}); ok {
					credentials = append(credentials, webauthnmodels.Credential{
						WebauthnCredentialId: credMap["webauthnCredentialId"].(string),
						RelyingPartyId:       credMap["relyingPartyId"].(string),
						RecipeUserId:         credMap["recipeUserId"].(string),
						CreatedAt:            int64(credMap["createdAt"].(float64)),
					})
				}
			}
		}
		return webauthnmodels.ListCredentialsResponse{
			OK: &struct {
				Credentials []webauthnmodels.Credential
			}{Credentials: credentials},
		}, nil
	}

	removeCredential := func(
		webauthnCredentialId string,
		recipeUserId string,
		userContext supertokens.UserContext,
	) (webauthnmodels.RemoveCredentialResponse, error) {
		resp, err := querier.SendDeleteRequest(
			"/recipe/webauthn/user/credential/remove",
			nil,
			map[string]string{
				"recipeUserId":         recipeUserId,
				"webauthnCredentialId": webauthnCredentialId,
			},
			userContext,
		)
		if err != nil {
			return webauthnmodels.RemoveCredentialResponse{}, err
		}
		status := resp["status"].(string)
		if status == "CREDENTIAL_NOT_FOUND_ERROR" {
			return webauthnmodels.RemoveCredentialResponse{
				CredentialNotFoundError: &struct{}{},
			}, nil
		}
		return webauthnmodels.RemoveCredentialResponse{OK: &struct{}{}}, nil
	}

	sendEmail := func(
		input emaildelivery.EmailType,
		userContext supertokens.UserContext,
	) error {
		return (*recipeEmailDelivery.IngredientInterfaceImpl.SendEmail)(input, userContext)
	}

	return webauthnmodels.RecipeInterface{
		GetGeneratedOptions:         &getGeneratedOptions,
		RegisterOptions:             &registerOptions,
		SignInOptions:               &signInOptions,
		SignUp:                      &signUp,
		SignIn:                      &signIn,
		VerifyCredentials:           &verifyCredentials,
		GenerateRecoverAccountToken: &generateRecoverAccountToken,
		ConsumeRecoverAccountToken:  &consumeRecoverAccountToken,
		RegisterCredential:          &registerCredential,
		GetUserFromCredentialId:     &getUserFromCredentialId,
		GetUserByEmail:              &getUserByEmail,
		GetUserByID:                 &getUserByID,
		ListCredentials:             &listCredentials,
		RemoveCredential:            &removeCredential,
		SendEmail:                   &sendEmail,
	}
}

func buildRegisterOptionsResponse(resp map[string]interface{}) (webauthnmodels.RegisterOptionsResponse, error) {
	ok := &struct {
		WebauthnGeneratedOptionsId string
		CreatedAt                  int64
		ExpiresAt                  int64
		RP                         struct {
			ID   string
			Name string
		}
		User struct {
			ID          string
			Name        string
			DisplayName string
		}
		Challenge          string
		Timeout            int
		ExcludeCredentials []struct {
			ID         string
			Transports []webauthnmodels.AuthenticatorTransport
			Type       string
		}
		Attestation      webauthnmodels.Attestation
		PubKeyCredParams []struct {
			Alg  int
			Type string
		}
		AuthenticatorSelection struct {
			RequireResidentKey bool
			ResidentKey        webauthnmodels.ResidentKey
			UserVerification   webauthnmodels.UserVerification
		}
	}{}

	ok.WebauthnGeneratedOptionsId = resp["webauthnGeneratedOptionsId"].(string)
	ok.CreatedAt = int64(resp["createdAt"].(float64))
	ok.ExpiresAt = int64(resp["expiresAt"].(float64))
	ok.Challenge = resp["challenge"].(string)
	ok.Timeout = int(resp["timeout"].(float64))
	ok.Attestation = webauthnmodels.Attestation(resp["attestation"].(string))

	if rp, ok2 := resp["rp"].(map[string]interface{}); ok2 {
		ok.RP.ID = rp["id"].(string)
		ok.RP.Name = rp["name"].(string)
	}
	if user, ok2 := resp["user"].(map[string]interface{}); ok2 {
		ok.User.ID = user["id"].(string)
		ok.User.Name = user["name"].(string)
		ok.User.DisplayName = user["displayName"].(string)
	}
	if authSel, ok2 := resp["authenticatorSelection"].(map[string]interface{}); ok2 {
		if rrk, ok3 := authSel["requireResidentKey"].(bool); ok3 {
			ok.AuthenticatorSelection.RequireResidentKey = rrk
		}
		if rk, ok3 := authSel["residentKey"].(string); ok3 {
			ok.AuthenticatorSelection.ResidentKey = webauthnmodels.ResidentKey(rk)
		}
		if uv, ok3 := authSel["userVerification"].(string); ok3 {
			ok.AuthenticatorSelection.UserVerification = webauthnmodels.UserVerification(uv)
		}
	}
	if excludeCredentials, ok2 := resp["excludeCredentials"].([]interface{}); ok2 {
		ok.ExcludeCredentials = make([]struct {
			ID         string
			Transports []webauthnmodels.AuthenticatorTransport
			Type       string
		}, 0, len(excludeCredentials))

		for _, raw := range excludeCredentials {
			credentialMap, ok3 := raw.(map[string]interface{})
			if !ok3 {
				continue
			}

			credential := struct {
				ID         string
				Transports []webauthnmodels.AuthenticatorTransport
				Type       string
			}{}

			if id, ok4 := credentialMap["id"].(string); ok4 {
				credential.ID = id
			}
			if typ, ok4 := credentialMap["type"].(string); ok4 {
				credential.Type = typ
			}
			if transports, ok4 := credentialMap["transports"].([]interface{}); ok4 {
				credential.Transports = make([]webauthnmodels.AuthenticatorTransport, 0, len(transports))
				for _, transportRaw := range transports {
					if transport, ok5 := transportRaw.(string); ok5 {
						credential.Transports = append(credential.Transports, webauthnmodels.AuthenticatorTransport(transport))
					}
				}
			}

			ok.ExcludeCredentials = append(ok.ExcludeCredentials, credential)
		}
	}
	if pubKeyCredParams, ok2 := resp["pubKeyCredParams"].([]interface{}); ok2 {
		ok.PubKeyCredParams = make([]struct {
			Alg  int
			Type string
		}, 0, len(pubKeyCredParams))

		for _, raw := range pubKeyCredParams {
			paramMap, ok3 := raw.(map[string]interface{})
			if !ok3 {
				continue
			}

			param := struct {
				Alg  int
				Type string
			}{}
			if alg, ok4 := paramMap["alg"].(float64); ok4 {
				param.Alg = int(alg)
			}
			if typ, ok4 := paramMap["type"].(string); ok4 {
				param.Type = typ
			}

			ok.PubKeyCredParams = append(ok.PubKeyCredParams, param)
		}
	}

	return webauthnmodels.RegisterOptionsResponse{OK: ok}, nil
}

func parseSignUpResponse(resp map[string]interface{}) (webauthnmodels.SignUpResponse, error) {
	status := resp["status"].(string)
	switch status {
	case "OK":
		user := parseUser(resp["user"].(map[string]interface{}))
		return webauthnmodels.SignUpResponse{
			OK: &struct {
				User         supertokens.User
				RecipeUserId string
			}{
				User:         user,
				RecipeUserId: resp["recipeUserId"].(string),
			},
		}, nil
	case "EMAIL_ALREADY_EXISTS_ERROR":
		return webauthnmodels.SignUpResponse{EmailAlreadyExistsError: &struct{}{}}, nil
	case "OPTIONS_NOT_FOUND_ERROR":
		return webauthnmodels.SignUpResponse{OptionsNotFoundError: &struct{}{}}, nil
	case "INVALID_OPTIONS_ERROR":
		return webauthnmodels.SignUpResponse{InvalidOptionsError: &struct{}{}}, nil
	case "INVALID_CREDENTIALS_ERROR":
		return webauthnmodels.SignUpResponse{InvalidCredentialsError: &struct{}{}}, nil
	case "INVALID_AUTHENTICATOR_ERROR":
		return webauthnmodels.SignUpResponse{InvalidAuthenticatorError: &struct{ Reason string }{Reason: resp["reason"].(string)}}, nil
	}
	return webauthnmodels.SignUpResponse{}, fmt.Errorf("unknown status: %s", status)
}

func parseSignInResponse(resp map[string]interface{}) (webauthnmodels.SignInResponse, error) {
	status := resp["status"].(string)
	switch status {
	case "OK":
		user := parseUser(resp["user"].(map[string]interface{}))
		return webauthnmodels.SignInResponse{
			OK: &struct {
				User         supertokens.User
				RecipeUserId string
			}{
				User:         user,
				RecipeUserId: resp["recipeUserId"].(string),
			},
		}, nil
	case "INVALID_CREDENTIALS_ERROR":
		return webauthnmodels.SignInResponse{InvalidCredentialsError: &struct{}{}}, nil
	case "INVALID_OPTIONS_ERROR":
		return webauthnmodels.SignInResponse{InvalidOptionsError: &struct{}{}}, nil
	case "INVALID_AUTHENTICATOR_ERROR":
		return webauthnmodels.SignInResponse{InvalidAuthenticatorError: &struct{ Reason string }{Reason: resp["reason"].(string)}}, nil
	case "CREDENTIAL_NOT_FOUND_ERROR":
		return webauthnmodels.SignInResponse{CredentialNotFoundError: &struct{}{}}, nil
	case "UNKNOWN_USER_ID_ERROR":
		return webauthnmodels.SignInResponse{UnknownUserIdError: &struct{}{}}, nil
	case "OPTIONS_NOT_FOUND_ERROR":
		return webauthnmodels.SignInResponse{OptionsNotFoundError: &struct{}{}}, nil
	}
	return webauthnmodels.SignInResponse{}, fmt.Errorf("unknown status: %s", status)
}

func parseUser(userMap map[string]interface{}) supertokens.User {
	userJSON, _ := json.Marshal(userMap)
	var user supertokens.User
	json.Unmarshal(userJSON, &user)
	return user
}

func credentialToMap(v interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	return m, err
}
