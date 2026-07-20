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
	"github.com/supertokens/supertokens-golang/supertokens"
)

// TypeInput is the developer-facing config for the webauthn recipe.
type TypeInput struct {
	// GetRelyingPartyId returns the RP ID (usually the registrable domain suffix, e.g. "example.com").
	// If nil, defaults to the hostname of APIDomain.
	GetRelyingPartyId func(tenantId string, req *http.Request, userContext supertokens.UserContext) (string, error)

	// GetRelyingPartyName returns a human-readable name shown to the user by the authenticator.
	// If nil, defaults to AppName.
	GetRelyingPartyName func(tenantId string, userContext supertokens.UserContext) (string, error)

	// GetOrigin returns the origin used for WebAuthn (e.g. "https://example.com").
	// If nil, defaults to AppInfo.GetOrigin.
	GetOrigin func(tenantId string, req *http.Request, userContext supertokens.UserContext) (string, error)

	// ValidateEmailAddress optionally overrides email validation. Return an error message string or nil.
	ValidateEmailAddress func(email string, tenantId string, userContext supertokens.UserContext) *string

	Override      *OverrideStruct
	EmailDelivery *emaildelivery.TypeInput
}

// TypeNormalisedInput is the internal normalised config.
type TypeNormalisedInput struct {
	GetRelyingPartyId    func(tenantId string, req *http.Request, userContext supertokens.UserContext) (string, error)
	GetRelyingPartyName  func(tenantId string, req *http.Request, userContext supertokens.UserContext) (string, error)
	GetOrigin            func(tenantId string, req *http.Request, userContext supertokens.UserContext) (string, error)
	ValidateEmailAddress func(email string, tenantId string, userContext supertokens.UserContext) *string

	Override               OverrideStruct
	GetEmailDeliveryConfig func() emaildelivery.TypeInputWithService
}

const (
	DefaultRegisterOptionsTimeout          = 60000
	DefaultRegisterOptionsAttestation      = AttestationNone
	DefaultRegisterOptionsResidentKey      = ResidentKeyRequired
	DefaultRegisterOptionsUserVerification = UserVerificationPreferred
	DefaultRegisterOptionsUserPresence     = true
	DefaultSignInOptionsTimeout            = 60000
	DefaultSignInOptionsUserVerification   = UserVerificationPreferred
	DefaultSignInOptionsUserPresence       = true
)

var DefaultRegisterOptionsSupportedAlgorithmIDs = []COSEAlgorithmIdentifier{-8, -7, -257}

type OverrideStruct struct {
	Functions func(originalImplementation RecipeInterface) RecipeInterface
	APIs      func(originalImplementation APIInterface) APIInterface
}
