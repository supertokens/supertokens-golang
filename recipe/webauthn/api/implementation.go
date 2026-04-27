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

package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/webauthn/webauthnmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() webauthnmodels.APIInterface {

	registerOptionsPOST := func(
		email *string,
		recoverAccountToken *string,
		tenantId string,
		options webauthnmodels.APIOptions,
		userContext supertokens.UserContext,
	) (webauthnmodels.RegisterOptionsPOSTResponse, error) {
		relyingPartyId, err := options.Config.GetRelyingPartyId(tenantId, options.Req, userContext)
		if err != nil {
			return webauthnmodels.RegisterOptionsPOSTResponse{}, err
		}
		relyingPartyName, err := options.Config.GetRelyingPartyName(tenantId, options.Req, userContext)
		if err != nil {
			return webauthnmodels.RegisterOptionsPOSTResponse{}, err
		}
		origin, err := options.Config.GetOrigin(tenantId, options.Req, userContext)
		if err != nil {
			return webauthnmodels.RegisterOptionsPOSTResponse{}, err
		}

		defaultTimeout := webauthnmodels.DefaultRegisterOptionsTimeout
		defaultAttestation := webauthnmodels.DefaultRegisterOptionsAttestation
		defaultResidentKey := webauthnmodels.DefaultRegisterOptionsResidentKey
		defaultUserVerification := webauthnmodels.DefaultRegisterOptionsUserVerification
		defaultUserPresence := webauthnmodels.DefaultRegisterOptionsUserPresence
		defaultAlgIDs := append(
			[]webauthnmodels.COSEAlgorithmIdentifier{},
			webauthnmodels.DefaultRegisterOptionsSupportedAlgorithmIDs...,
		)

		resp, err := (*options.RecipeImplementation.RegisterOptions)(
			email,
			recoverAccountToken,
			relyingPartyId,
			relyingPartyName,
			origin,
			&defaultTimeout,
			&defaultAttestation,
			&defaultResidentKey,
			&defaultUserVerification,
			&defaultUserPresence,
			defaultAlgIDs,
			tenantId,
			userContext,
		)
		if err != nil {
			return webauthnmodels.RegisterOptionsPOSTResponse{}, err
		}
		if resp.RecoverAccountTokenInvalidError != nil {
			return webauthnmodels.RegisterOptionsPOSTResponse{
				RecoverAccountTokenInvalidError: resp.RecoverAccountTokenInvalidError,
			}, nil
		}
		if resp.InvalidEmailError != nil {
			return webauthnmodels.RegisterOptionsPOSTResponse{
				InvalidEmailError: resp.InvalidEmailError,
			}, nil
		}
		if resp.InvalidOptionsError != nil {
			return webauthnmodels.RegisterOptionsPOSTResponse{
				InvalidOptionsError: resp.InvalidOptionsError,
			}, nil
		}
		return webauthnmodels.RegisterOptionsPOSTResponse{OK: resp.OK}, nil
	}

	signInOptionsPOST := func(
		tenantId string,
		options webauthnmodels.APIOptions,
		userContext supertokens.UserContext,
	) (webauthnmodels.SignInOptionsPOSTResponse, error) {
		relyingPartyId, err := options.Config.GetRelyingPartyId(tenantId, options.Req, userContext)
		if err != nil {
			return webauthnmodels.SignInOptionsPOSTResponse{}, err
		}
		origin, err := options.Config.GetOrigin(tenantId, options.Req, userContext)
		if err != nil {
			return webauthnmodels.SignInOptionsPOSTResponse{}, err
		}

		defaultTimeout := webauthnmodels.DefaultSignInOptionsTimeout
		defaultUserVerification := webauthnmodels.DefaultSignInOptionsUserVerification
		defaultUserPresence := webauthnmodels.DefaultSignInOptionsUserPresence

		resp, err := (*options.RecipeImplementation.SignInOptions)(
			relyingPartyId,
			origin,
			&defaultTimeout,
			&defaultUserVerification,
			&defaultUserPresence,
			tenantId,
			userContext,
		)
		if err != nil {
			return webauthnmodels.SignInOptionsPOSTResponse{}, err
		}
		if resp.InvalidOptionsError != nil {
			return webauthnmodels.SignInOptionsPOSTResponse{
				InvalidOptionsError: resp.InvalidOptionsError,
			}, nil
		}
		return webauthnmodels.SignInOptionsPOSTResponse{OK: resp.OK}, nil
	}

	signUpPOST := func(
		webauthnGeneratedOptionsId string,
		credential webauthnmodels.RegistrationPayload,
		sess *sessmodels.SessionContainer,
		tenantId string,
		options webauthnmodels.APIOptions,
		userContext supertokens.UserContext,
	) (webauthnmodels.SignUpPOSTResponse, error) {
		resp, err := (*options.RecipeImplementation.SignUp)(
			webauthnGeneratedOptionsId,
			credential,
			tenantId,
			userContext,
		)
		if err != nil {
			return webauthnmodels.SignUpPOSTResponse{}, err
		}
		if resp.EmailAlreadyExistsError != nil {
			return webauthnmodels.SignUpPOSTResponse{EmailAlreadyExistsError: resp.EmailAlreadyExistsError}, nil
		}
		if resp.OptionsNotFoundError != nil {
			return webauthnmodels.SignUpPOSTResponse{OptionsNotFoundError: resp.OptionsNotFoundError}, nil
		}
		if resp.InvalidOptionsError != nil {
			return webauthnmodels.SignUpPOSTResponse{InvalidOptionsError: resp.InvalidOptionsError}, nil
		}
		if resp.InvalidCredentialsError != nil {
			return webauthnmodels.SignUpPOSTResponse{InvalidCredentialsError: resp.InvalidCredentialsError}, nil
		}
		if resp.InvalidAuthenticatorError != nil {
			return webauthnmodels.SignUpPOSTResponse{InvalidAuthenticatorError: resp.InvalidAuthenticatorError}, nil
		}

		newSession, err := session.CreateNewSession(
			options.Req, options.Res, tenantId, resp.OK.RecipeUserId,
			map[string]interface{}{}, map[string]interface{}{}, userContext,
		)
		if err != nil {
			return webauthnmodels.SignUpPOSTResponse{}, err
		}

		return webauthnmodels.SignUpPOSTResponse{
			OK: &struct {
				User         supertokens.User
				RecipeUserId string
			}{
				User:         resp.OK.User,
				RecipeUserId: resp.OK.RecipeUserId,
			},
			Session: &newSession,
		}, nil
	}

	signInPOST := func(
		webauthnGeneratedOptionsId string,
		credential webauthnmodels.AuthenticationPayload,
		sess *sessmodels.SessionContainer,
		tenantId string,
		options webauthnmodels.APIOptions,
		userContext supertokens.UserContext,
	) (webauthnmodels.SignInPOSTResponse, error) {
		resp, err := (*options.RecipeImplementation.SignIn)(
			webauthnGeneratedOptionsId,
			credential,
			tenantId,
			userContext,
		)
		if err != nil {
			return webauthnmodels.SignInPOSTResponse{}, err
		}
		if resp.InvalidCredentialsError != nil {
			return webauthnmodels.SignInPOSTResponse{InvalidCredentialsError: resp.InvalidCredentialsError}, nil
		}
		if resp.InvalidOptionsError != nil {
			return webauthnmodels.SignInPOSTResponse{InvalidOptionsError: resp.InvalidOptionsError}, nil
		}
		if resp.InvalidAuthenticatorError != nil {
			return webauthnmodels.SignInPOSTResponse{InvalidAuthenticatorError: resp.InvalidAuthenticatorError}, nil
		}
		if resp.CredentialNotFoundError != nil {
			return webauthnmodels.SignInPOSTResponse{CredentialNotFoundError: resp.CredentialNotFoundError}, nil
		}
		if resp.UnknownUserIdError != nil {
			return webauthnmodels.SignInPOSTResponse{UnknownUserIdError: resp.UnknownUserIdError}, nil
		}
		if resp.OptionsNotFoundError != nil {
			return webauthnmodels.SignInPOSTResponse{OptionsNotFoundError: resp.OptionsNotFoundError}, nil
		}

		newSession, err := session.CreateNewSession(
			options.Req, options.Res, tenantId, resp.OK.RecipeUserId,
			map[string]interface{}{}, map[string]interface{}{}, userContext,
		)
		if err != nil {
			return webauthnmodels.SignInPOSTResponse{}, err
		}

		return webauthnmodels.SignInPOSTResponse{
			OK: &struct {
				User         supertokens.User
				Session      sessmodels.SessionContainer
				RecipeUserId string
			}{
				User:         resp.OK.User,
				Session:      newSession,
				RecipeUserId: resp.OK.RecipeUserId,
			},
		}, nil
	}

	generateRecoverAccountTokenPOST := func(
		email string,
		tenantId string,
		options webauthnmodels.APIOptions,
		userContext supertokens.UserContext,
	) (webauthnmodels.GenerateRecoverAccountTokenPOSTResponse, error) {
		user, err := (*options.RecipeImplementation.GetUserByEmail)(email, tenantId, userContext)
		if err != nil {
			return webauthnmodels.GenerateRecoverAccountTokenPOSTResponse{}, err
		}
		if user == nil {
			supertokens.LogDebugMessage(fmt.Sprintf("Webauthn recover account email not sent; no user found for email: %s", email))
			return webauthnmodels.GenerateRecoverAccountTokenPOSTResponse{OK: &struct{}{}}, nil
		}

		resp, err := (*options.RecipeImplementation.GenerateRecoverAccountToken)(
			user.ID, email, tenantId, userContext,
		)
		if err != nil {
			return webauthnmodels.GenerateRecoverAccountTokenPOSTResponse{}, err
		}
		if resp.UnknownUserIdError != nil {
			supertokens.LogDebugMessage(fmt.Sprintf("Webauthn recover account email not sent; unknown user ID for email: %s", email))
			return webauthnmodels.GenerateRecoverAccountTokenPOSTResponse{OK: &struct{}{}}, nil
		}

		recoverAccountLink, err := getRecoverAccountLink(options.AppInfo, resp.OK.Token, tenantId, options.Req, userContext)
		if err != nil {
			return webauthnmodels.GenerateRecoverAccountTokenPOSTResponse{}, err
		}

		supertokens.LogDebugMessage(fmt.Sprintf("Sending webauthn recover account email to %s", email))
		sendErr := (*options.EmailDelivery.IngredientInterfaceImpl.SendEmail)(
			emaildelivery.EmailType{
				WebauthnRecoverAccount: &emaildelivery.WebauthnRecoverAccountType{
					User:               emaildelivery.User{ID: user.ID, Email: email},
					RecoverAccountLink: recoverAccountLink,
					TenantId:           tenantId,
				},
			},
			userContext,
		)
		if sendErr != nil {
			return webauthnmodels.GenerateRecoverAccountTokenPOSTResponse{}, sendErr
		}

		return webauthnmodels.GenerateRecoverAccountTokenPOSTResponse{OK: &struct{}{}}, nil
	}

	consumeRecoverAccountTokenPOST := func(
		token string,
		tenantId string,
		options webauthnmodels.APIOptions,
		userContext supertokens.UserContext,
	) (webauthnmodels.ConsumeRecoverAccountTokenPOSTResponse, error) {
		resp, err := (*options.RecipeImplementation.ConsumeRecoverAccountToken)(token, tenantId, userContext)
		if err != nil {
			return webauthnmodels.ConsumeRecoverAccountTokenPOSTResponse{}, err
		}
		if resp.RecoverAccountTokenInvalidError != nil {
			return webauthnmodels.ConsumeRecoverAccountTokenPOSTResponse{
				RecoverAccountTokenInvalidError: resp.RecoverAccountTokenInvalidError,
			}, nil
		}
		return webauthnmodels.ConsumeRecoverAccountTokenPOSTResponse{OK: resp.OK}, nil
	}

	recoverAccountPOST := func(
		webauthnGeneratedOptionsId string,
		credential webauthnmodels.RegistrationPayload,
		token string,
		tenantId string,
		options webauthnmodels.APIOptions,
		userContext supertokens.UserContext,
	) (webauthnmodels.RecoverAccountPOSTResponse, error) {
		consumeResp, err := (*options.RecipeImplementation.ConsumeRecoverAccountToken)(token, tenantId, userContext)
		if err != nil {
			return webauthnmodels.RecoverAccountPOSTResponse{}, err
		}
		if consumeResp.RecoverAccountTokenInvalidError != nil {
			return webauthnmodels.RecoverAccountPOSTResponse{
				RecoverAccountTokenInvalidError: consumeResp.RecoverAccountTokenInvalidError,
			}, nil
		}

		registerResp, err := (*options.RecipeImplementation.RegisterCredential)(
			consumeResp.OK.UserId,
			webauthnGeneratedOptionsId,
			credential,
			userContext,
		)
		if err != nil {
			return webauthnmodels.RecoverAccountPOSTResponse{}, err
		}
		if registerResp.InvalidCredentialsError != nil {
			return webauthnmodels.RecoverAccountPOSTResponse{InvalidCredentialsError: registerResp.InvalidCredentialsError}, nil
		}
		if registerResp.OptionsNotFoundError != nil {
			return webauthnmodels.RecoverAccountPOSTResponse{OptionsNotFoundError: registerResp.OptionsNotFoundError}, nil
		}
		if registerResp.InvalidOptionsError != nil {
			return webauthnmodels.RecoverAccountPOSTResponse{InvalidOptionsError: registerResp.InvalidOptionsError}, nil
		}
		if registerResp.InvalidAuthenticatorError != nil {
			return webauthnmodels.RecoverAccountPOSTResponse{InvalidAuthenticatorError: registerResp.InvalidAuthenticatorError}, nil
		}

		return webauthnmodels.RecoverAccountPOSTResponse{OK: &struct{}{}}, nil
	}

	emailExistsGET := func(
		email string,
		tenantId string,
		options webauthnmodels.APIOptions,
		userContext supertokens.UserContext,
	) (webauthnmodels.EmailExistsGETResponse, error) {
		user, err := (*options.RecipeImplementation.GetUserByEmail)(email, tenantId, userContext)
		if err != nil {
			return webauthnmodels.EmailExistsGETResponse{}, err
		}
		return webauthnmodels.EmailExistsGETResponse{
			OK: &struct{ Exists bool }{Exists: user != nil},
		}, nil
	}

	listCredentialsGET := func(
		sess sessmodels.SessionContainer,
		tenantId string,
		options webauthnmodels.APIOptions,
		userContext supertokens.UserContext,
	) (webauthnmodels.ListCredentialsGETResponse, error) {
		userId := sess.GetUserIDWithContext(userContext)
		user, err := (*options.RecipeImplementation.GetUserByID)(userId, userContext)
		if err != nil {
			return webauthnmodels.ListCredentialsGETResponse{}, err
		}
		if user == nil {
			return webauthnmodels.ListCredentialsGETResponse{}, fmt.Errorf("user not found")
		}

		var allCredentials []struct {
			WebauthnCredentialId string
			RelyingPartyId       string
			RecipeUserId         string
			CreatedAt            int64
		}

		for _, lm := range user.LoginMethods {
			if lm.RecipeID != "webauthn" {
				continue
			}
			resp, err := (*options.RecipeImplementation.ListCredentials)(lm.RecipeUserID, userContext)
			if err != nil {
				return webauthnmodels.ListCredentialsGETResponse{}, err
			}
			if resp.OK != nil {
				for _, cred := range resp.OK.Credentials {
					allCredentials = append(allCredentials, struct {
						WebauthnCredentialId string
						RelyingPartyId       string
						RecipeUserId         string
						CreatedAt            int64
					}{
						WebauthnCredentialId: cred.WebauthnCredentialId,
						RelyingPartyId:       cred.RelyingPartyId,
						RecipeUserId:         cred.RecipeUserId,
						CreatedAt:            cred.CreatedAt,
					})
				}
			}
		}

		return webauthnmodels.ListCredentialsGETResponse{
			OK: &struct {
				Credentials []struct {
					WebauthnCredentialId string
					RelyingPartyId       string
					RecipeUserId         string
					CreatedAt            int64
				}
			}{Credentials: allCredentials},
		}, nil
	}

	registerCredentialPOST := func(
		webauthnGeneratedOptionsId string,
		credential webauthnmodels.RegistrationPayload,
		sess sessmodels.SessionContainer,
		tenantId string,
		options webauthnmodels.APIOptions,
		userContext supertokens.UserContext,
	) (webauthnmodels.RegisterCredentialPOSTResponse, error) {
		recipeUserId := getSessionRecipeUserID(sess, userContext)
		resp, err := (*options.RecipeImplementation.RegisterCredential)(recipeUserId, webauthnGeneratedOptionsId, credential, userContext)
		if err != nil {
			return webauthnmodels.RegisterCredentialPOSTResponse{}, err
		}
		if resp.InvalidCredentialsError != nil {
			return webauthnmodels.RegisterCredentialPOSTResponse{InvalidCredentialsError: resp.InvalidCredentialsError}, nil
		}
		if resp.OptionsNotFoundError != nil {
			return webauthnmodels.RegisterCredentialPOSTResponse{OptionsNotFoundError: resp.OptionsNotFoundError}, nil
		}
		if resp.InvalidOptionsError != nil {
			return webauthnmodels.RegisterCredentialPOSTResponse{InvalidOptionsError: resp.InvalidOptionsError}, nil
		}
		if resp.InvalidAuthenticatorError != nil {
			return webauthnmodels.RegisterCredentialPOSTResponse{InvalidAuthenticatorError: resp.InvalidAuthenticatorError}, nil
		}
		return webauthnmodels.RegisterCredentialPOSTResponse{OK: &struct{}{}}, nil
	}

	removeCredentialPOST := func(
		webauthnCredentialId string,
		sess sessmodels.SessionContainer,
		tenantId string,
		options webauthnmodels.APIOptions,
		userContext supertokens.UserContext,
	) (webauthnmodels.RemoveCredentialPOSTResponse, error) {
		userId := sess.GetUserIDWithContext(userContext)
		user, err := (*options.RecipeImplementation.GetUserByID)(userId, userContext)
		if err != nil {
			return webauthnmodels.RemoveCredentialPOSTResponse{}, err
		}
		if user == nil {
			return webauthnmodels.RemoveCredentialPOSTResponse{}, fmt.Errorf("user not found")
		}

		userFromCredential, err := (*options.RecipeImplementation.GetUserFromCredentialId)(webauthnCredentialId, tenantId, userContext)
		if err != nil {
			return webauthnmodels.RemoveCredentialPOSTResponse{}, err
		}
		if userFromCredential.CredentialNotFoundError != nil {
			return webauthnmodels.RemoveCredentialPOSTResponse{}, fmt.Errorf("user not found")
		}
		if userFromCredential.OK == nil || userFromCredential.OK.User.ID != user.ID {
			return webauthnmodels.RemoveCredentialPOSTResponse{}, fmt.Errorf("user not found")
		}

		resp, err := (*options.RecipeImplementation.RemoveCredential)(webauthnCredentialId, userFromCredential.OK.RecipeUserId, userContext)
		if err != nil {
			return webauthnmodels.RemoveCredentialPOSTResponse{}, err
		}
		if resp.CredentialNotFoundError != nil {
			return webauthnmodels.RemoveCredentialPOSTResponse{CredentialNotFoundError: &struct{}{}}, nil
		}
		return webauthnmodels.RemoveCredentialPOSTResponse{OK: &struct{}{}}, nil
	}

	_ = strings.Contains // keep strings import used via getRecoverAccountLink

	return webauthnmodels.APIInterface{
		RegisterOptionsPOST:             &registerOptionsPOST,
		SignInOptionsPOST:               &signInOptionsPOST,
		SignUpPOST:                      &signUpPOST,
		SignInPOST:                      &signInPOST,
		ConsumeRecoverAccountTokenPOST:  &consumeRecoverAccountTokenPOST,
		RecoverAccountPOST:              &recoverAccountPOST,
		EmailExistsGET:                  &emailExistsGET,
		GenerateRecoverAccountTokenPOST: &generateRecoverAccountTokenPOST,
		ListCredentialsGET:              &listCredentialsGET,
		RegisterCredentialPOST:          &registerCredentialPOST,
		RemoveCredentialPOST:            &removeCredentialPOST,
	}
}

func getSessionRecipeUserID(session sessmodels.SessionContainer, userContext supertokens.UserContext) string {
	payload := session.GetAccessTokenPayloadWithContext(userContext)
	if payload != nil {
		if recipeUserID, ok := payload["rsub"].(string); ok && recipeUserID != "" {
			return recipeUserID
		}
	}
	return session.GetUserIDWithContext(userContext)
}

func getRecoverAccountLink(appInfo supertokens.NormalisedAppinfo, token string, tenantId string, req *http.Request, userContext supertokens.UserContext) (string, error) {
	origin, err := appInfo.GetOrigin(req, userContext)
	if err != nil {
		return "", err
	}
	websiteBasePath := appInfo.WebsiteBasePath.GetAsStringDangerous()
	return fmt.Sprintf("%s%s/webauthn/recover?token=%s&tenantId=%s",
		origin.GetAsStringDangerous(),
		websiteBasePath,
		token,
		tenantId,
	), nil
}
