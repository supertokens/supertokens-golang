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
	"net/http"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// --- Primitive types ---

type Base64URLString = string
type COSEAlgorithmIdentifier = int

type ResidentKey string

const (
	ResidentKeyRequired    ResidentKey = "required"
	ResidentKeyPreferred   ResidentKey = "preferred"
	ResidentKeyDiscouraged ResidentKey = "discouraged"
)

type UserVerification string

const (
	UserVerificationRequired    UserVerification = "required"
	UserVerificationPreferred   UserVerification = "preferred"
	UserVerificationDiscouraged UserVerification = "discouraged"
)

type Attestation string

const (
	AttestationNone       Attestation = "none"
	AttestationIndirect   Attestation = "indirect"
	AttestationDirect     Attestation = "direct"
	AttestationEnterprise Attestation = "enterprise"
)

type AuthenticatorTransport string

const (
	TransportBLE       AuthenticatorTransport = "ble"
	TransportCable     AuthenticatorTransport = "cable"
	TransportHybrid    AuthenticatorTransport = "hybrid"
	TransportInternal  AuthenticatorTransport = "internal"
	TransportNFC       AuthenticatorTransport = "nfc"
	TransportSmartCard AuthenticatorTransport = "smart-card"
	TransportUSB       AuthenticatorTransport = "usb"
)

type AuthenticatorAttachment string

const (
	AttachmentPlatform      AuthenticatorAttachment = "platform"
	AttachmentCrossPlatform AuthenticatorAttachment = "cross-platform"
)

// --- Credential payload types ---

type CredentialPayloadBase struct {
	ID                      string                   `json:"id"`
	RawID                   string                   `json:"rawId"`
	AuthenticatorAttachment *AuthenticatorAttachment `json:"authenticatorAttachment,omitempty"`
	ClientExtensionResults  map[string]interface{}   `json:"clientExtensionResults"`
	Type                    string                   `json:"type"` // always "public-key"
}

type AuthenticatorAssertionResponseJSON struct {
	ClientDataJSON    Base64URLString  `json:"clientDataJSON"`
	AuthenticatorData Base64URLString  `json:"authenticatorData"`
	Signature         Base64URLString  `json:"signature"`
	UserHandle        *Base64URLString `json:"userHandle,omitempty"`
}

type AuthenticationPayload struct {
	CredentialPayloadBase
	Response AuthenticatorAssertionResponseJSON `json:"response"`
}

type AuthenticatorAttestationResponseJSON struct {
	ClientDataJSON     Base64URLString          `json:"clientDataJSON"`
	AttestationObject  Base64URLString          `json:"attestationObject"`
	AuthenticatorData  *Base64URLString         `json:"authenticatorData,omitempty"`
	Transports         []AuthenticatorTransport `json:"transports,omitempty"`
	PublicKeyAlgorithm *COSEAlgorithmIdentifier `json:"publicKeyAlgorithm,omitempty"`
	PublicKey          *Base64URLString         `json:"publicKey,omitempty"`
}

type RegistrationPayload struct {
	CredentialPayloadBase
	Response AuthenticatorAttestationResponseJSON `json:"response"`
}

type APIOptions struct {
	RecipeImplementation RecipeInterface
	AppInfo              supertokens.NormalisedAppinfo
	Config               TypeNormalisedInput
	RecipeID             string
	Req                  *http.Request
	Res                  http.ResponseWriter
	OtherHandler         http.HandlerFunc
	EmailDelivery        emaildelivery.Ingredient
}

type APIInterface struct {
	RegisterOptionsPOST *func(
		email *string,
		recoverAccountToken *string,
		tenantId string,
		options APIOptions,
		userContext supertokens.UserContext,
	) (RegisterOptionsPOSTResponse, error)

	SignInOptionsPOST *func(
		tenantId string,
		options APIOptions,
		userContext supertokens.UserContext,
	) (SignInOptionsPOSTResponse, error)

	SignUpPOST *func(
		webauthnGeneratedOptionsId string,
		credential RegistrationPayload,
		session *sessmodels.SessionContainer,
		tenantId string,
		options APIOptions,
		userContext supertokens.UserContext,
	) (SignUpPOSTResponse, error)

	SignInPOST *func(
		webauthnGeneratedOptionsId string,
		credential AuthenticationPayload,
		session *sessmodels.SessionContainer,
		tenantId string,
		options APIOptions,
		userContext supertokens.UserContext,
	) (SignInPOSTResponse, error)

	ConsumeRecoverAccountTokenPOST *func(
		token string,
		tenantId string,
		options APIOptions,
		userContext supertokens.UserContext,
	) (ConsumeRecoverAccountTokenPOSTResponse, error)

	RecoverAccountPOST *func(
		webauthnGeneratedOptionsId string,
		credential RegistrationPayload,
		token string,
		tenantId string,
		options APIOptions,
		userContext supertokens.UserContext,
	) (RecoverAccountPOSTResponse, error)

	EmailExistsGET *func(
		email string,
		tenantId string,
		options APIOptions,
		userContext supertokens.UserContext,
	) (EmailExistsGETResponse, error)

	GenerateRecoverAccountTokenPOST *func(
		email string,
		tenantId string,
		options APIOptions,
		userContext supertokens.UserContext,
	) (GenerateRecoverAccountTokenPOSTResponse, error)

	ListCredentialsGET *func(
		session sessmodels.SessionContainer,
		tenantId string,
		options APIOptions,
		userContext supertokens.UserContext,
	) (ListCredentialsGETResponse, error)

	RegisterCredentialPOST *func(
		webauthnGeneratedOptionsId string,
		credential RegistrationPayload,
		session sessmodels.SessionContainer,
		tenantId string,
		options APIOptions,
		userContext supertokens.UserContext,
	) (RegisterCredentialPOSTResponse, error)

	RemoveCredentialPOST *func(
		webauthnCredentialId string,
		session sessmodels.SessionContainer,
		tenantId string,
		options APIOptions,
		userContext supertokens.UserContext,
	) (RemoveCredentialPOSTResponse, error)
}

type GenerateRecoverAccountTokenPOSTResponse struct {
	OK           *struct{}
	GeneralError *supertokens.GeneralErrorResponse
}

type RegisterOptionsPOSTResponse struct {
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
	InvalidOptionsError             *struct{}
	InvalidEmailError               *struct{ Err string }
	GeneralError                    *supertokens.GeneralErrorResponse
}

type SignInOptionsPOSTResponse struct {
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
	GeneralError        *supertokens.GeneralErrorResponse
}

type SignUpPOSTResponse struct {
	OK *struct {
		User         supertokens.User
		RecipeUserId string
	}
	Session                    *sessmodels.SessionContainer
	EmailAlreadyExistsError    *struct{}
	OptionsNotFoundError       *struct{}
	InvalidOptionsError        *struct{}
	InvalidCredentialsError    *struct{}
	InvalidAuthenticatorError  *struct{ Reason string }
	LinkingToSessionUserFailed *struct{ Reason string }
	GeneralError               *supertokens.GeneralErrorResponse
}

type SignInPOSTResponse struct {
	OK *struct {
		User         supertokens.User
		Session      sessmodels.SessionContainer
		RecipeUserId string
	}
	InvalidCredentialsError    *struct{}
	InvalidOptionsError        *struct{}
	InvalidAuthenticatorError  *struct{ Reason string }
	CredentialNotFoundError    *struct{}
	UnknownUserIdError         *struct{}
	OptionsNotFoundError       *struct{}
	LinkingToSessionUserFailed *struct{ Reason string }
	GeneralError               *supertokens.GeneralErrorResponse
}

type ConsumeRecoverAccountTokenPOSTResponse struct {
	OK *struct {
		Email  string
		UserId string
	}
	RecoverAccountTokenInvalidError *struct{}
	GeneralError                    *supertokens.GeneralErrorResponse
}

type RecoverAccountPOSTResponse struct {
	OK                              *struct{}
	RecoverAccountTokenInvalidError *struct{}
	InvalidCredentialsError         *struct{}
	OptionsNotFoundError            *struct{}
	InvalidOptionsError             *struct{}
	InvalidAuthenticatorError       *struct{ Reason string }
	GeneralError                    *supertokens.GeneralErrorResponse
}

type EmailExistsGETResponse struct {
	OK           *struct{ Exists bool }
	GeneralError *supertokens.GeneralErrorResponse
}

type ListCredentialsGETResponse struct {
	OK *struct {
		Credentials []struct {
			WebauthnCredentialId string
			RelyingPartyId       string
			RecipeUserId         string
			CreatedAt            int64
		}
	}
	GeneralError *supertokens.GeneralErrorResponse
}

type RegisterCredentialPOSTResponse struct {
	OK                        *struct{}
	InvalidCredentialsError   *struct{}
	InvalidOptionsError       *struct{}
	InvalidAuthenticatorError *struct{ Reason string }
	OptionsNotFoundError      *struct{}
	GeneralError              *supertokens.GeneralErrorResponse
}

type RemoveCredentialPOSTResponse struct {
	OK                      *struct{}
	CredentialNotFoundError *struct{}
	GeneralError            *supertokens.GeneralErrorResponse
}
