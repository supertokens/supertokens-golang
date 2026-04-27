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

package webauthnmodels

import (
	"github.com/supertokens/supertokens-golang/supertokens"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
)

type RecipeInterface struct {
	GetGeneratedOptions *func(
		webauthnGeneratedOptionsId string,
		tenantId string,
		userContext supertokens.UserContext,
	) (GetGeneratedOptionsResponse, error)

	RegisterOptions *func(
		email *string,
		recoverAccountToken *string,
		relyingPartyId string,
		relyingPartyName string,
		origin string,
		timeout *int,
		attestation *Attestation,
		residentKey *ResidentKey,
		userVerification *UserVerification,
		userPresence *bool,
		supportedAlgorithmIds []COSEAlgorithmIdentifier,
		tenantId string,
		userContext supertokens.UserContext,
	) (RegisterOptionsResponse, error)

	SignInOptions *func(
		relyingPartyId string,
		origin string,
		timeout *int,
		userVerification *UserVerification,
		userPresence *bool,
		tenantId string,
		userContext supertokens.UserContext,
	) (SignInOptionsResponse, error)

	SignUp *func(
		webauthnGeneratedOptionsId string,
		credential RegistrationPayload,
		tenantId string,
		userContext supertokens.UserContext,
	) (SignUpResponse, error)

	SignIn *func(
		webauthnGeneratedOptionsId string,
		credential AuthenticationPayload,
		tenantId string,
		userContext supertokens.UserContext,
	) (SignInResponse, error)

	VerifyCredentials *func(
		webauthnGeneratedOptionsId string,
		credential AuthenticationPayload,
		tenantId string,
		userContext supertokens.UserContext,
	) (VerifyCredentialsResponse, error)

	GenerateRecoverAccountToken *func(
		userId string,
		email string,
		tenantId string,
		userContext supertokens.UserContext,
	) (GenerateRecoverAccountTokenResponse, error)

	ConsumeRecoverAccountToken *func(
		token string,
		tenantId string,
		userContext supertokens.UserContext,
	) (ConsumeRecoverAccountTokenResponse, error)

	RegisterCredential *func(
		recipeUserId string,
		webauthnGeneratedOptionsId string,
		credential RegistrationPayload,
		userContext supertokens.UserContext,
	) (RegisterCredentialResponse, error)

	GetUserFromCredentialId *func(
		credentialId string,
		tenantId string,
		userContext supertokens.UserContext,
	) (GetUserFromCredentialIdResponse, error)

	GetUserByEmail *func(
		email string,
		tenantId string,
		userContext supertokens.UserContext,
	) (*supertokens.User, error)

	GetUserByID *func(
		userID string,
		userContext supertokens.UserContext,
	) (*supertokens.User, error)

	ListCredentials *func(
		recipeUserId string,
		userContext supertokens.UserContext,
	) (ListCredentialsResponse, error)

	RemoveCredential *func(
		webauthnCredentialId string,
		recipeUserId string,
		userContext supertokens.UserContext,
	) (RemoveCredentialResponse, error)

	SendEmail *func(
		input emaildelivery.EmailType,
		userContext supertokens.UserContext,
	) error
}

// --- Response types ---

type GetGeneratedOptionsResponse struct {
	OK *struct {
		WebauthnGeneratedOptionsId string
		CreatedAt                  int64
		ExpiresAt                  int64
		Email                      *string
		RelyingPartyId             string
		RelyingPartyName           string
		Origin                     string
		Challenge                  string
		Timeout                    int
		UserVerification           UserVerification
		UserPresence               bool
	}
	OptionsNotFoundError *struct{}
}

type RegisterOptionsResponse struct {
	OK *struct {
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
			Transports []AuthenticatorTransport
			Type       string
		}
		Attestation      Attestation
		PubKeyCredParams []struct {
			Alg  int
			Type string
		}
		AuthenticatorSelection struct {
			RequireResidentKey bool
			ResidentKey        ResidentKey
			UserVerification   UserVerification
		}
	}
	RecoverAccountTokenInvalidError *struct{}
	InvalidEmailError               *struct{ Err string }
	InvalidOptionsError             *struct{}
}

type SignInOptionsResponse struct {
	OK *struct {
		WebauthnGeneratedOptionsId string
		CreatedAt                  int64
		ExpiresAt                  int64
		RpId                       string
		Challenge                  string
		Timeout                    int
		UserVerification           UserVerification
	}
	InvalidOptionsError *struct{}
}

type SignUpResponse struct {
	OK *struct {
		User         supertokens.User
		RecipeUserId string
	}
	EmailAlreadyExistsError   *struct{}
	OptionsNotFoundError      *struct{}
	InvalidOptionsError       *struct{}
	InvalidCredentialsError   *struct{}
	InvalidAuthenticatorError *struct{ Reason string }
}

type SignInResponse struct {
	OK *struct {
		User         supertokens.User
		RecipeUserId string
	}
	InvalidCredentialsError   *struct{}
	InvalidOptionsError       *struct{}
	InvalidAuthenticatorError *struct{ Reason string }
	CredentialNotFoundError   *struct{}
	UnknownUserIdError        *struct{}
	OptionsNotFoundError      *struct{}
}

type VerifyCredentialsResponse struct {
	OK                        *struct{}
	InvalidCredentialsError   *struct{}
	InvalidOptionsError       *struct{}
	InvalidAuthenticatorError *struct{ Reason string }
	CredentialNotFoundError   *struct{}
	OptionsNotFoundError      *struct{}
}

type GenerateRecoverAccountTokenResponse struct {
	OK                 *struct{ Token string }
	UnknownUserIdError *struct{}
}

type ConsumeRecoverAccountTokenResponse struct {
	OK *struct {
		Email  string
		UserId string
	}
	RecoverAccountTokenInvalidError *struct{}
}

type RegisterCredentialResponse struct {
	OK                        *struct{}
	InvalidCredentialsError   *struct{}
	InvalidOptionsError       *struct{}
	InvalidAuthenticatorError *struct{ Reason string }
	OptionsNotFoundError      *struct{}
}

type GetUserFromCredentialIdResponse struct {
	OK *struct {
		User         supertokens.User
		RecipeUserId string
	}
	CredentialNotFoundError *struct{}
}

type Credential struct {
	WebauthnCredentialId string
	RelyingPartyId       string
	RecipeUserId         string
	CreatedAt            int64
}

type ListCredentialsResponse struct {
	OK *struct {
		Credentials []Credential
	}
}

type RemoveCredentialResponse struct {
	OK                      *struct{}
	CredentialNotFoundError *struct{}
}
