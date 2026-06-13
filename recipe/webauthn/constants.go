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

package webauthn

const (
	doesEmailExistAPI              = "/webauthn/email/exists"
	registerCredentialAPI          = "/webauthn/credential"
	recoverAccountAPI              = "/user/webauthn/reset"
	generateRecoverAccountTokenAPI = "/user/webauthn/reset/token"
	registerOptionsAPI             = "/webauthn/options/register"
	signInAPI                      = "/webauthn/signin"
	signInOptionsAPI               = "/webauthn/options/signin"
	signUpAPI                      = "/webauthn/signup"
	listCredentialsAPI             = "/webauthn/credential/list"
	removeCredentialAPI            = "/webauthn/credential/remove"
)
