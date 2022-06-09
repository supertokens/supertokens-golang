/*
 * Copyright (c) 2022, VRAI Labs and/or its affiliates. All rights reserved.
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

package emaildelivery

import "github.com/supertokens/supertokens-golang/supertokens"

type EmailDeliveryInterface struct {
	SendEmail *func(input EmailType, userContext supertokens.UserContext) error
}

type TypeInput struct {
	Service  *EmailDeliveryInterface
	Override func(originalImplementation EmailDeliveryInterface) EmailDeliveryInterface
}

type TypeInputWithService struct {
	Service  EmailDeliveryInterface
	Override func(originalImplementation EmailDeliveryInterface) EmailDeliveryInterface
}

type EmailType struct {
	EmailVerification *EmailVerificationType
	PasswordReset     *PasswordResetType
	PasswordlessLogin *PasswordlessLoginType
}

type EmailVerificationType struct {
	User            User
	EmailVerifyLink string
}

type PasswordResetType struct {
	User              User
	PasswordResetLink string
}

type PasswordlessLoginType struct {
	Email            string
	UserInputCode    *string
	UrlWithLinkCode  *string
	CodeLifetime     uint64
	PreAuthSessionId string
}

type User struct {
	ID    string
	Email string
}
